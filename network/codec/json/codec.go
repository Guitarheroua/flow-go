package json

import (
	"bytes"
	"fmt"
	"github.com/onflow/flow-go/model/encoding/json"
	"github.com/onflow/flow-go/network/codec"
	"github.com/rs/zerolog"
	"io"
	"math/rand"
	"time"

	"github.com/onflow/flow-go/network"
)

// Codec represents a CBOR codec for our network.
type Codec struct {
	codec      json.Codec
	marshaller json.Marshaler
	l          *zerolog.Logger
}

var _ network.Codec = (*Codec)(nil)

// NewCodec creates a new CBOR codec.
func NewCodec() *Codec {
	c := &Codec{}
	return c
}

func (c *Codec) WithLogger(l *zerolog.Logger) *Codec {
	c.l = l
	return c
}

func (c *Codec) NewEncoder(w io.Writer) network.Encoder {
	return c.codec.NewEncoder(w)
}

func (c *Codec) NewDecoder(r io.Reader) network.Decoder {
	return &Decoder{decoder: c.codec.NewDecoder(r), marshaller: &c.marshaller}
}

func (c *Codec) Encode(v interface{}) ([]byte, error) {
	var start time.Time
	if c.l != nil {
		start = time.Now()
	}

	// encode the value
	code, what, err := codec.MessageCodeFromInterface(v)
	if err != nil {
		return nil, fmt.Errorf("could not determine envelope code: %w", err)
	}

	var data bytes.Buffer
	data.WriteByte(code.Uint8())

	encoder := c.NewEncoder(&data)
	err = encoder.Encode(v)
	if err != nil {
		return nil, fmt.Errorf("could not encode CBOR payload with envelope code %decoder AKA %s: %w", code, what, err) // e.g. 2, "CodeBlockProposal", <CBOR error>
	}

	dataBytes := data.Bytes()

	var duration time.Duration
	if c.l != nil {
		duration = time.Since(start)
	}

	byteUUID := generateOneByteUUID()
	dataBytes = append(dataBytes, byteUUID)

	if c.l != nil {
		c.l.Debug().Msg(fmt.Sprintf("JSON: Execution of Encode took: %s for %s with size: %d bytes.", duration, fmt.Sprintf("1-byte UUID: 0x%02X\n", byteUUID), data.Len()))
	}

	return dataBytes, nil
}

func (c *Codec) Decode(data []byte) (interface{}, error) {
	var start time.Time
	if c.l != nil {
		start = time.Now()
	}

	msgCode, err := codec.MessageCodeFromPayload(data)
	if err != nil {
		return nil, err
	}

	msgInterface, what, err := codec.InterfaceFromMessageCode(msgCode)
	if err != nil {
		return nil, err
	}

	// unmarshal the payload
	err = c.marshaller.Unmarshal(data[1:len(data)-1], msgInterface)
	if err != nil {
		return nil, codec.NewMsgUnmarshalErr(data[0], what, err)
	}

	var duration time.Duration
	if c.l != nil {
		duration = time.Since(start)
	}

	byteUUID := data[len(data)-1]
	data = data[:len(data)-1]

	if c.l != nil {
		c.l.Debug().Msg(fmt.Sprintf("JSON: Execution of Decode took: %s for %s with size: %d bytes.", duration, fmt.Sprintf("1-byte UUID: 0x%02X\n", byteUUID), len(data)))
	}

	return msgInterface, nil
}

func generateOneByteUUID() byte {
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random byte value
	return byte(r.Intn(256))
}
