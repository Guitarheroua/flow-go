package backend

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/onflow/flow/protobuf/go/flow/access"
	"github.com/onflow/flow/protobuf/go/flow/execution"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"github.com/onflow/flow-go/module"
	"github.com/onflow/flow-go/utils/grpcutils"
)

// the default timeout used when making a GRPC request to a collection node or an execution node
const defaultClientTimeout = 3 * time.Second

// ConnectionFactory is used to create an access api client
type ConnectionFactory interface {
	GetAccessAPIClient(address string) (access.AccessAPIClient, error)
	InvalidateAccessAPIClient(address string)
	GetExecutionAPIClient(address string) (execution.ExecutionAPIClient, error)
	InvalidateExecutionAPIClient(address string)
}

type ProxyConnectionFactory struct {
	ConnectionFactory
	targetAddress string
}

func (p *ProxyConnectionFactory) GetAccessAPIClient(address string) (access.AccessAPIClient, error) {
	return p.ConnectionFactory.GetAccessAPIClient(p.targetAddress)
}

func (p *ProxyConnectionFactory) GetExecutionAPIClient(address string) (execution.ExecutionAPIClient, error) {
	return p.ConnectionFactory.GetExecutionAPIClient(p.targetAddress)
}

type ConnectionFactoryImpl struct {
	CollectionGRPCPort        uint
	ExecutionGRPCPort         uint
	CollectionNodeGRPCTimeout time.Duration
	ExecutionNodeGRPCTimeout  time.Duration
	ConnectionsCache          *lru.Cache
	CacheSize                 uint
	lock                      sync.Mutex
	AccessMetrics             module.AccessMetrics
}

type ConnectionCacheStore struct {
	ClientConn *grpc.ClientConn
	mutex      *sync.Mutex
}

// createConnection creates new gRPC connections to remote node
func (cf *ConnectionFactoryImpl) createConnection(address string, timeout time.Duration) (*grpc.ClientConn, error) {

	if timeout == 0 {
		timeout = defaultClientTimeout
	}

	keepaliveParams := keepalive.ClientParameters{
		Time:    10 * time.Second,
		Timeout: timeout,
	}

	// ClientConn's default KeepAlive on connections is indefinite, assuming the timeout isn't reached
	// The connections should be safe to be persisted and reused
	// https://pkg.go.dev/google.golang.org/grpc#WithKeepaliveParams
	// https://grpc.io/blog/grpc-on-http2/#keeping-connections-alive
	conn, err := grpc.Dial(
		address,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcutils.DefaultMaxMsgSize)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepaliveParams),
		WithClientUnaryInterceptor(timeout))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to address %s: %w", address, err)
	}
	return conn, nil
}

func (cf *ConnectionFactoryImpl) retrieveConnection(grpcAddress string, timeout time.Duration) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	clientMutex := new(sync.Mutex)
	store := ConnectionCacheStore{
		ClientConn: nil,
		mutex:      clientMutex,
	}
	if prev, ok, _ := cf.ConnectionsCache.PeekOrAdd(grpcAddress, store); ok {
		clientMutex = prev.(ConnectionCacheStore).mutex
	}
	// we lock this mutex to prevent a scenario where the connection is not good, which will result in
	// re-establishing the connection for this address. if the mutex is not locked, we may attempt to re-establish
	// the connection multiple times which would result in cache thrashing.
	clientMutex.Lock()
	defer clientMutex.Unlock()

	// get again to update the LRU cache
	if res, ok := cf.ConnectionsCache.Get(grpcAddress); ok {
		conn = res.(ConnectionCacheStore).ClientConn
		if cf.AccessMetrics != nil {
			cf.AccessMetrics.ConnectionFromPoolRetrieved()
		}
	}
	if conn == nil || conn.GetState() == connectivity.Shutdown {
		// this lock prevents a memory leak where a race condition may occur if 2 requests to a new connection at the
		// same address occur. the second add would overwrite the first without closing the connection
		cf.lock.Lock()
		defer cf.lock.Unlock()
		// Check if connection was created/refreshed by another thread
		if res, ok := cf.ConnectionsCache.Get(grpcAddress); ok {
			conn = res.(ConnectionCacheStore).ClientConn
			if conn != nil && conn.GetState() != connectivity.Shutdown {
				return conn, nil
			}
		}

		// updates to the cache don't trigger evictions; this line closes connections before re-establishing new ones
		if conn != nil {
			conn.Close()
		}
		var err error
		conn, err = cf.addConnection(grpcAddress, timeout, clientMutex)
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (cf *ConnectionFactoryImpl) addConnection(grpcAddress string, timeout time.Duration, mutex *sync.Mutex) (*grpc.ClientConn, error) {
	conn, err := cf.createConnection(grpcAddress, timeout)
	if err != nil {
		return nil, err
	}

	store := ConnectionCacheStore{
		ClientConn: conn,
		mutex:      mutex,
	}

	cf.ConnectionsCache.Add(grpcAddress, store)
	if cf.AccessMetrics != nil {
		cf.AccessMetrics.TotalConnectionsInPool(uint(cf.ConnectionsCache.Len()), cf.CacheSize)
	}
	return conn, nil
}

func (cf *ConnectionFactoryImpl) GetAccessAPIClient(address string) (access.AccessAPIClient, error) {

	grpcAddress, err := getGRPCAddress(address, cf.CollectionGRPCPort)
	if err != nil {
		return nil, err
	}

	conn, err := cf.retrieveConnection(grpcAddress, cf.CollectionNodeGRPCTimeout)
	if err != nil {
		return nil, err
	}

	accessAPIClient := access.NewAccessAPIClient(conn)
	return accessAPIClient, nil
}

func (cf *ConnectionFactoryImpl) InvalidateAccessAPIClient(address string) {
	cf.invalidateAPIClient(address, cf.CollectionGRPCPort)
}

func (cf *ConnectionFactoryImpl) GetExecutionAPIClient(address string) (execution.ExecutionAPIClient, error) {

	grpcAddress, err := getGRPCAddress(address, cf.ExecutionGRPCPort)
	if err != nil {
		return nil, err
	}

	conn, err := cf.retrieveConnection(grpcAddress, cf.ExecutionNodeGRPCTimeout)
	if err != nil {
		return nil, err
	}

	executionAPIClient := execution.NewExecutionAPIClient(conn)
	return executionAPIClient, nil
}

func (cf *ConnectionFactoryImpl) InvalidateExecutionAPIClient(address string) {
	cf.invalidateAPIClient(address, cf.ExecutionGRPCPort)
}

func (cf *ConnectionFactoryImpl) invalidateAPIClient(address string, port uint) {
	grpcAddress, _ := getGRPCAddress(address, port)
	if res, ok := cf.ConnectionsCache.Get(grpcAddress); ok {
		store := res.(ConnectionCacheStore)
		store.mutex.Lock()
		if store.ClientConn != nil {
			conn := store.ClientConn
			conn.Close()
			store.ClientConn = nil
		}
		store.mutex.Unlock()
	}
}

// getExecutionNodeAddress translates flow.Identity address to the GRPC address of the node by switching the port to the
// GRPC port from the libp2p port
func getGRPCAddress(address string, grpcPort uint) (string, error) {
	// split hostname and port
	hostnameOrIP, _, err := net.SplitHostPort(address)
	if err != nil {
		return "", err
	}
	// use the hostname from identity list and port number as the one passed in as argument
	grpcAddress := fmt.Sprintf("%s:%d", hostnameOrIP, grpcPort)

	return grpcAddress, nil
}

func WithClientUnaryInterceptor(timeout time.Duration) grpc.DialOption {

	clientTimeoutInterceptor := func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {

		// create a context that expires after timeout
		ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)

		defer cancel()

		// call the remote GRPC using the short context
		err := invoker(ctxWithTimeout, method, req, reply, cc, opts...)

		return err
	}

	return grpc.WithUnaryInterceptor(clientTimeoutInterceptor)
}
