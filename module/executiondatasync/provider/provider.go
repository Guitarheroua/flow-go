package provider

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/module/blobs"
	"github.com/onflow/flow-go/module/executiondatasync/execution_data"
	"github.com/onflow/flow-go/module/executiondatasync/tracker"
	"github.com/onflow/flow-go/network"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type ProviderOption func(*Provider)

func WithBlobSizeLimit(size int) ProviderOption {
	return func(p *Provider) {
		p.maxBlobSize = size
	}
}

type Provider struct {
	logger      zerolog.Logger
	metrics     module.ExecutionDataProviderMetrics
	maxBlobSize int
	serializer  execution_data.Serializer
	blobService network.BlobService
	storage     *tracker.Storage
}

type ProvideJob struct {
	ExecutionDataID flow.Identifier
	Done            <-chan error
}

func NewProvider(
	logger zerolog.Logger,
	metrics module.ExecutionDataProviderMetrics,
	serializer execution_data.Serializer,
	blobService network.BlobService,
	opts ...ProviderOption,
) *Provider {
	p := &Provider{
		logger:      logger.With().Str("component", "execution_data_provider").Logger(),
		metrics:     metrics,
		maxBlobSize: execution_data.DefaultMaxBlobSize,
		serializer:  serializer,
		blobService: blobService,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

func (p *Provider) storeBlobs(parent context.Context, blockHeight uint64, blobCh <-chan blobs.Blob) <-chan error {
	ch := make(chan error, 1)
	go func() {
		defer close(ch)

		start := time.Now()

		var blobs []blobs.Blob
		var cids []cid.Cid
		var totalSize uint64
		for blob := range blobCh {
			blobs = append(blobs, blob)
			cids = append(cids, blob.Cid())
			totalSize += uint64(len(blob.RawData()))
		}

		if zerolog.GlobalLevel() <= zerolog.DebugLevel && p.logger.GetLevel() <= zerolog.DebugLevel {
			cidArr := zerolog.Arr()
			for _, cid := range cids {
				cidArr = cidArr.Str(cid.String())
			}
			p.logger.Debug().Array("cids", cidArr).Uint64("height", blockHeight).Msg("storing blobs")
		}

		err := p.storage.Update(func(trackBlobs func(uint64, ...cid.Cid) error) error {
			ctx, cancel := context.WithCancel(parent)
			defer cancel()

			if err := trackBlobs(blockHeight, cids...); err != nil {
				return fmt.Errorf("failed to track blobs: %w", err)
			}

			if err := p.blobService.AddBlobs(ctx, blobs); err != nil {
				return fmt.Errorf("failed to add blobs: %w", err)
			}

			return nil
		})
		duration := time.Since(start)

		if err != nil {
			ch <- err
			p.metrics.AddBlobsFailed()
		} else {
			p.metrics.AddBlobsSucceeded(duration, totalSize)
		}
	}()

	return ch
}

func (p *Provider) Provide(ctx context.Context, blockHeight uint64, executionData *execution_data.BlockExecutionData) (*ProvideJob, error) {
	logger := p.logger.With().Uint64("height", blockHeight).Str("block_id", executionData.BlockID.String()).Logger()
	logger.Debug().Msg("providing execution data")

	start := time.Now()

	blobCh := make(chan blobs.Blob)
	defer close(blobCh)

	errCh := p.storeBlobs(ctx, blockHeight, blobCh)
	g, gCtx := errgroup.WithContext(ctx)

	chunkDataIDs := make([]cid.Cid, len(executionData.ChunkExecutionDatas))
	for i, chunkExecutionData := range executionData.ChunkExecutionDatas {
		i := i
		chunkExecutionData := chunkExecutionData

		g.Go(func() error {
			logger.Debug().Int("chunk_index", i).Msg("adding chunk execution data")
			cedID, err := p.addChunkExecutionData(gCtx, chunkExecutionData, blobCh)
			if err != nil {
				return fmt.Errorf("failed to add chunk execution data at index %d: %w", i, err)
			}
			logger.Debug().Int("chunk_index", i).Str("chunk_execution_data_id", cedID.String()).Msg("chunk execution data added")

			chunkDataIDs[i] = cedID
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	edRoot := &execution_data.BlockExecutionDataRoot{
		BlockID:               executionData.BlockID,
		ChunkExecutionDataIDs: chunkDataIDs,
	}
	rootID, err := p.addExecutionDataRoot(ctx, edRoot, blobCh)
	if err != nil {
		return nil, fmt.Errorf("failed to add execution data root: %w", err)
	}
	logger.Debug().Str("root_id", rootID.String()).Msg("root ID computed")

	duration := time.Since(start)
	p.metrics.RootIDComputed(duration, len(executionData.ChunkExecutionDatas))

	return &ProvideJob{rootID, errCh}, nil
}

func (p *Provider) addExecutionDataRoot(
	ctx context.Context,
	edRoot *execution_data.BlockExecutionDataRoot,
	blobCh chan<- blobs.Blob,
) (flow.Identifier, error) {
	buf := new(bytes.Buffer)
	if err := p.serializer.Serialize(buf, edRoot); err != nil {
		return flow.ZeroID, fmt.Errorf("failed to serialize execution data root: %w", err)
	}

	rootBlob := blobs.NewBlob(buf.Bytes())
	blobCh <- rootBlob

	rootID, err := flow.CidToId(rootBlob.Cid())
	if err != nil {
		return flow.ZeroID, fmt.Errorf("failed to convert root blob cid to id: %w", err)
	}

	return rootID, nil
}

func (p *Provider) addChunkExecutionData(
	ctx context.Context,
	ced *execution_data.ChunkExecutionData,
	blobCh chan<- blobs.Blob,
) (cid.Cid, error) {
	cids, err := p.addBlobs(ctx, ced, blobCh)
	if err != nil {
		return cid.Undef, fmt.Errorf("failed to add chunk execution data blobs: %w", err)
	}

	for {
		if len(cids) == 1 {
			return cids[0], nil
		}

		if cids, err = p.addBlobs(ctx, cids, blobCh); err != nil {
			return cid.Undef, fmt.Errorf("failed to add cid blobs: %w", err)
		}
	}
}

// addBlobs serializes the given object, splits the serialized data into blobs, and sends them to the given channel.
func (p *Provider) addBlobs(ctx context.Context, v interface{}, blobCh chan<- blobs.Blob) ([]cid.Cid, error) {
	bcw := blobs.NewBlobChannelWriter(blobCh, p.maxBlobSize)
	defer bcw.Close()

	if err := p.serializer.Serialize(bcw, v); err != nil {
		return nil, fmt.Errorf("failed to serialize object: %w", err)
	}

	if err := bcw.Flush(); err != nil {
		return nil, fmt.Errorf("failed to flush blob channel writer: %w", err)
	}

	return bcw.CidsSent(), nil
}
