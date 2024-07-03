package indexgateway

import (
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gogo/status"
	"github.com/grafana/dskit/concurrency"
	"github.com/grafana/dskit/grpcclient"
	"github.com/grafana/dskit/instrument"
	"github.com/grafana/dskit/middleware"
	"github.com/grafana/dskit/ring"
	"github.com/grafana/dskit/ring/client"
	"github.com/grafana/dskit/services"
	"github.com/grafana/dskit/tenant"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/grafana/loki/v3/pkg/distributor/clientpool"
	"github.com/grafana/loki/v3/pkg/logproto"
	"github.com/grafana/loki/v3/pkg/storage/stores/series/index"
	"github.com/grafana/loki/v3/pkg/util/constants"
	"github.com/grafana/loki/v3/pkg/util/discovery"
	util_math "github.com/grafana/loki/v3/pkg/util/math"
)

const (
	maxQueriesPerGrpc      = 100
	maxConcurrentGrpcCalls = 10
)

// ClientConfig configures the Index Gateway client used to communicate with
// the Index Gateway server.
type ClientConfig struct {
	// Mode sets in which mode the client will operate. It is actually defined at the
	// index_gateway YAML section and reused here.
	Mode Mode `yaml:"-"`

	// PoolConfig defines the behavior of the gRPC connection pool used to communicate
	// with the Index Gateway.
	//
	// Only relevant for the ring mode.
	// It is defined at the distributors YAML section and reused here.
	PoolConfig clientpool.PoolConfig `yaml:"-"`

	// Ring is the Index Gateway ring used to find the appropriate Index Gateway instance
	// this client should talk to.
	//
	// Only relevant for the ring mode.
	Ring ring.ReadRing `yaml:"-"`

	// GRPCClientConfig configures the gRPC connection between the Index Gateway client and the server.
	//
	// Used by both, ring and simple mode.
	GRPCClientConfig grpcclient.Config `yaml:"grpc_client_config"`

	// Address of the Index Gateway instance responsible for retaining the index for all tenants.
	//
	// Only relevant for the simple mode.
	Address string `yaml:"server_address,omitempty"`

	// Forcefully disable the use of the index gateway client for the storage.
	// This is mainly useful for the index-gateway component which should always use the storage.
	Disabled bool `yaml:"-"`

	// LogGatewayRequests configures if requests sent to the gateway should be logged or not.
	// The log messages are of type debug and contain the address of the gateway and the relevant tenant.
	LogGatewayRequests bool `yaml:"log_gateway_requests"`

	GRPCUnaryClientInterceptors  []grpc.UnaryClientInterceptor  `yaml:"-"`
	GRCPStreamClientInterceptors []grpc.StreamClientInterceptor `yaml:"-"`
}

// RegisterFlagsWithPrefix register client-specific flags with the given prefix.
//
// Flags that are used by both, client and server, are defined in the indexgateway package.
func (i *ClientConfig) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	i.GRPCClientConfig.RegisterFlagsWithPrefix(prefix+".grpc", f)
	f.StringVar(&i.Address, prefix+".server-address", "", "Hostname or IP of the Index Gateway gRPC server running in simple mode. Can also be prefixed with dns+, dnssrv+, or dnssrvnoa+ to resolve a DNS A record with multiple IP's, a DNS SRV record with a followup A record lookup, or a DNS SRV record without a followup A record lookup, respectively.")
	f.BoolVar(&i.LogGatewayRequests, prefix+".log-gateway-requests", false, "Whether requests sent to the gateway should be logged or not.")
}

func (i *ClientConfig) RegisterFlags(f *flag.FlagSet) {
	i.RegisterFlagsWithPrefix("index-gateway-client", f)
}

type GatewayClient struct {
	logger                            log.Logger
	cfg                               ClientConfig
	storeGatewayClientRequestDuration *prometheus.HistogramVec
	dnsProvider                       *discovery.DNS
	pool                              *client.Pool
	ring                              ring.ReadRing
	limits                            Limits
	done                              chan struct{}
}

// NewGatewayClient instantiates a new client used to communicate with an Index Gateway instance.
//
// If it is configured to be in ring mode, a pool of GRPC connections to all Index Gateway instances is created using a ring.
// Otherwise, it creates a GRPC connection pool to as many addresses as can be resolved from the given address.
func NewGatewayClient(cfg ClientConfig, r prometheus.Registerer, limits Limits, logger log.Logger, metricsNamespace string) (*GatewayClient, error) {
	latency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: constants.Loki,
		Name:      "index_gateway_request_duration_seconds",
		Help:      "Time (in seconds) spent serving requests when using the index gateway",
		Buckets:   instrument.DefBuckets,
	}, []string{"operation", "status_code"})
	if r != nil {
		err := r.Register(latency)
		if err != nil {
			alreadyErr, ok := err.(prometheus.AlreadyRegisteredError)
			if !ok {
				return nil, err
			}
			latency = alreadyErr.ExistingCollector.(*prometheus.HistogramVec)
		}
	}

	sgClient := &GatewayClient{
		logger:                            logger,
		cfg:                               cfg,
		storeGatewayClientRequestDuration: latency,
		ring:                              cfg.Ring,
		limits:                            limits,
		done:                              make(chan struct{}),
	}

	dialOpts, err := cfg.GRPCClientConfig.DialOption(instrumentation(cfg, sgClient.storeGatewayClientRequestDuration))
	if err != nil {
		return nil, errors.Wrap(err, "index gateway grpc dial option")
	}
	factory := func(addr string) (client.PoolClient, error) {
		igPool, err := NewClientPool(addr, dialOpts)
		if err != nil {
			return nil, errors.Wrap(err, "new index gateway grpc pool")
		}

		return igPool, nil
	}

	//FIXME(ewelch) we don't expose the pool configs nor set defaults, and register flags is kind of messed up with remote config being defined somewhere else
	//make a separate PR to make the pool config generic so it can be used with proper names in multiple places.
	sgClient.cfg.PoolConfig.RemoteTimeout = 2 * time.Second
	sgClient.cfg.PoolConfig.ClientCleanupPeriod = 5 * time.Second
	sgClient.cfg.PoolConfig.HealthCheckIngesters = true

	if sgClient.cfg.Mode == RingMode {
		sgClient.pool = clientpool.NewPool("index-gateway", sgClient.cfg.PoolConfig, sgClient.ring, client.PoolAddrFunc(factory), logger, metricsNamespace)
	} else {
		// Note we don't use clientpool.NewPool because we want to provide our own discovery function
		poolCfg := client.PoolConfig{
			CheckInterval:      sgClient.cfg.PoolConfig.ClientCleanupPeriod,
			HealthCheckEnabled: sgClient.cfg.PoolConfig.HealthCheckIngesters,
			HealthCheckTimeout: sgClient.cfg.PoolConfig.RemoteTimeout,
		}
		clients := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: constants.Loki,
			Name:      "index_gateway_clients",
			Help:      "The current number of index gateway clients.",
		})
		if r != nil {
			err := r.Register(clients)
			if err != nil {
				alreadyErr, ok := err.(prometheus.AlreadyRegisteredError)
				if !ok {
					return nil, err
				}
				clients = alreadyErr.ExistingCollector.(prometheus.Gauge)
			}
		}
		//TODO(ewelch) we can't use metrics in the provider because of duplicate registration errors
		dnsProvider := discovery.NewDNS(logger, sgClient.cfg.PoolConfig.ClientCleanupPeriod, sgClient.cfg.Address, nil)
		sgClient.dnsProvider = dnsProvider

		// Make an attempt to do one DNS lookup so we can start with addresses
		dnsProvider.RunOnce()

		discovery := func() ([]string, error) {
			return dnsProvider.Addresses(), nil
		}
		sgClient.pool = client.NewPool("index gateway", poolCfg, discovery, client.PoolAddrFunc(factory), clients, logger)

	}

	// We have to start the pool service, it will handle removing stale clients in the background
	err = services.StartAndAwaitRunning(context.Background(), sgClient.pool)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start index gateway connection pool")
	}

	return sgClient, nil
}

// Stop stops the execution of this gateway client.
func (s *GatewayClient) Stop() {
	ctx, cancel := context.WithTimeoutCause(context.Background(), 10*time.Second, errors.New("service shutdown timeout expired"))
	defer cancel()
	err := services.StopAndAwaitTerminated(ctx, s.pool)
	if err != nil {
		level.Error(s.logger).Log("msg", "failed to stop index gateway connection pool", "err", err)
	}
	s.dnsProvider.Stop()
}

func (s *GatewayClient) QueryPages(ctx context.Context, queries []index.Query, callback index.QueryPagesCallback) error {
	if len(queries) <= maxQueriesPerGrpc {
		return s.doQueries(ctx, queries, callback)
	}

	jobsCount := len(queries) / maxQueriesPerGrpc
	if len(queries)%maxQueriesPerGrpc != 0 {
		jobsCount++
	}
	return concurrency.ForEachJob(ctx, jobsCount, maxConcurrentGrpcCalls, func(ctx context.Context, idx int) error {
		return s.doQueries(ctx, queries[idx*maxQueriesPerGrpc:util_math.Min((idx+1)*maxQueriesPerGrpc, len(queries))], callback)
	})
}

func (s *GatewayClient) QueryIndex(_ context.Context, _ *logproto.QueryIndexRequest, _ ...grpc.CallOption) (logproto.IndexGateway_QueryIndexClient, error) {
	panic("not implemented")
}

func (s *GatewayClient) GetChunkRef(ctx context.Context, in *logproto.GetChunkRefRequest) (*logproto.GetChunkRefResponse, error) {
	return poolDo(ctx, s, true, func(innerCtx context.Context, client logproto.IndexGatewayClient) (*logproto.GetChunkRefResponse, error) {
		return client.GetChunkRef(innerCtx, in)
	})
}

func (s *GatewayClient) GetSeries(ctx context.Context, in *logproto.GetSeriesRequest) (*logproto.GetSeriesResponse, error) {
	return poolDo(ctx, s, true, func(innerCtx context.Context, client logproto.IndexGatewayClient) (*logproto.GetSeriesResponse, error) {
		return client.GetSeries(innerCtx, in)
	})
}

func (s *GatewayClient) LabelNamesForMetricName(ctx context.Context, in *logproto.LabelNamesForMetricNameRequest) (*logproto.LabelResponse, error) {
	return poolDo(ctx, s, true, func(innerCtx context.Context, client logproto.IndexGatewayClient) (*logproto.LabelResponse, error) {
		return client.LabelNamesForMetricName(innerCtx, in)
	})
}

func (s *GatewayClient) LabelValuesForMetricName(ctx context.Context, in *logproto.LabelValuesForMetricNameRequest) (*logproto.LabelResponse, error) {
	return poolDo(ctx, s, true, func(innerCtx context.Context, client logproto.IndexGatewayClient) (*logproto.LabelResponse, error) {
		return client.LabelValuesForMetricName(innerCtx, in)
	})
}

func (s *GatewayClient) GetStats(ctx context.Context, in *logproto.IndexStatsRequest) (*logproto.IndexStatsResponse, error) {
	return poolDo(ctx, s, true, func(ctx context.Context, client logproto.IndexGatewayClient) (*logproto.IndexStatsResponse, error) {
		return client.GetStats(ctx, in)
	})
}

func (s *GatewayClient) GetVolume(ctx context.Context, in *logproto.VolumeRequest) (*logproto.VolumeResponse, error) {
	return poolDo(ctx, s, true, func(innerCtx context.Context, client logproto.IndexGatewayClient) (*logproto.VolumeResponse, error) {
		return client.GetVolume(innerCtx, in)
	})
}

func (s *GatewayClient) GetShards(
	ctx context.Context,
	in *logproto.ShardsRequest,
) (res *logproto.ShardsResponse, err error) {

	return poolDo(
		ctx,
		s,
		true,
		func(innerCtx context.Context, client logproto.IndexGatewayClient) (*logproto.ShardsResponse, error) {
			perReplicaResult := &logproto.ShardsResponse{}
			streamer, err := client.GetShards(innerCtx, in)
			if err != nil {
				return nil, errors.Wrap(err, "get shards")
			}

			// TODO(owen-d): stream currently unused (buffered) because query planning doesn't expect a streamed response,
			// but can be improved easily in the future by using a stream here.
			for {
				resp, err := streamer.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return nil, errors.WithStack(err)
				}
				perReplicaResult.Merge(resp)
			}

			// Since `poolDo` retries on error, we only want to set the response if we got a successful response.
			// This avoids cases where we add duplicates to the response on retries.
			res = perReplicaResult

			return res, nil
		})
}

// TODO(owen-d): this was copied from ingester_querier.go -- move it to a shared pkg
// isUnimplementedCallError tells if the GRPC error is a gRPC error with code Unimplemented.
func isUnimplementedCallError(err error) bool {
	if err == nil {
		return false
	}

	s, ok := status.FromError(err)
	if !ok {
		return false
	}
	return (s.Code() == codes.Unimplemented)
}

func (s *GatewayClient) doQueries(ctx context.Context, queries []index.Query, callback index.QueryPagesCallback) error {
	queryKeyQueryMap := make(map[string]index.Query, len(queries))
	gatewayQueries := make([]*logproto.IndexQuery, 0, len(queries))

	for _, query := range queries {
		queryKeyQueryMap[index.QueryKey(query)] = query
		gatewayQueries = append(gatewayQueries, &logproto.IndexQuery{
			TableName:        query.TableName,
			HashValue:        query.HashValue,
			RangeValuePrefix: query.RangeValuePrefix,
			RangeValueStart:  query.RangeValueStart,
			ValueEqual:       query.ValueEqual,
		})
	}

	// Don't hedge these requests because we can't make this threadsafe as it is
	// mutating state in this function via the callback function being passed.
	_, err := poolDo(ctx, s, false, func(innerCtx context.Context, client logproto.IndexGatewayClient) (*any, error) {
		return nil, s.clientDoQueries(innerCtx, gatewayQueries, queryKeyQueryMap, callback, client)
	})
	return err
}

// clientDoQueries send a query request to an Index Gateway instance using the given gRPC client.
//
// It is used by both, simple and ring mode.
func (s *GatewayClient) clientDoQueries(ctx context.Context, gatewayQueries []*logproto.IndexQuery,
	queryKeyQueryMap map[string]index.Query, callback index.QueryPagesCallback, client logproto.IndexGatewayClient,
) error {
	streamer, err := client.QueryIndex(ctx, &logproto.QueryIndexRequest{Queries: gatewayQueries})
	if err != nil {
		return errors.Wrap(err, "query index")
	}

	for {
		resp, err := streamer.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return errors.WithStack(err)
		}
		query, ok := queryKeyQueryMap[resp.QueryKey]
		if !ok {
			level.Error(s.logger).Log("msg", fmt.Sprintf("unexpected %s QueryKey received, expected queries %s", resp.QueryKey, fmt.Sprint(queryKeyQueryMap)))
			return fmt.Errorf("unexpected %s QueryKey received", resp.QueryKey)
		}
		if !callback(query, &readBatch{resp}) {
			return nil
		}
	}

	return nil
}

type callbackResponse[T any] struct {
	resp   *T
	hedged bool
	err    error
}

// poolDo will call callback providing a random index gateway client, if hedge is true, it will dispatch
// a second call to callback after the hedge timeout providing a different random index gateway client.
// NOTE: if hedge is true, callback must be safe to call concurrently from different threads.
func poolDo[T any](ctx context.Context, gatewayClient *GatewayClient, hedge bool, callback func(ctx context.Context, client logproto.IndexGatewayClient) (*T, error)) (*T, error) {
	userID, err := tenant.TenantID(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "index gateway client get tenant ID")
	}

	// This function returns the gateways in a random order
	addrs, err := gatewayClient.getServerAddresses(userID)
	if err != nil {
		return nil, err
	}

	if len(addrs) == 0 {
		level.Error(gatewayClient.logger).Log("msg", fmt.Sprintf("no index gateway instances found for tenant %s", userID))
		return nil, fmt.Errorf("no index gateway instances found for tenant %s", userID)
	}

	respChan := make(chan callbackResponse[T])
	firstAddr := addrs[0]

	if gatewayClient.cfg.LogGatewayRequests {
		level.Debug(gatewayClient.logger).Log("msg", "sending request to gateway", "gateway", firstAddr, "tenant", userID)
	}

	genericClient, err := gatewayClient.pool.GetClientFor(firstAddr)
	if err != nil {
		level.Error(gatewayClient.logger).Log("msg", fmt.Sprintf("failed to get client for instance %s", firstAddr), "err", err)
		return nil, fmt.Errorf("failed to get index gateway client for instance %s", firstAddr)
	}

	// Send first request
	client1 := genericClient.(logproto.IndexGatewayClient)
	ctx1, cancel1 := context.WithCancelCause(ctx)
	defer cancel1(errors.New("canceled from defer statement"))
	go func() {
		req1Start := time.Now()
		r, err := callback(ctx1, client1)
		level.Info(gatewayClient.logger).Log("msg", "initial request returned", "duration", time.Since(req1Start), "err", err)
		// Because we only listen for once response to make sure we don't leak a goroutine
		// here also watch for the context to be canceled
		select {
		case respChan <- callbackResponse[T]{resp: r, hedged: false, err: err}:
		case <-ctx1.Done():
		}
	}()

	select {
	case r := <-respChan:
		return r.resp, r.err
	case <-time.After(1 * time.Second):
		// Send hedged request if we have another address and we are able to get a client from it.
		// Ignore errors here and just don't send the hedged request
		if hedge && len(addrs) > 1 {
			gc2, err := gatewayClient.pool.GetClientFor(addrs[1])
			client2 := gc2.(logproto.IndexGatewayClient)
			if err == nil {
				ctx2, cancel2 := context.WithCancelCause(ctx)
				defer cancel2(errors.New("canceled from defer statement"))
				go func() {
					level.Info(gatewayClient.logger).Log("msg", "sending hedged request", "err", err)
					req2start := time.Now()
					r, err := callback(ctx2, client2)
					level.Info(gatewayClient.logger).Log("msg", "hedged request returned", "duration", time.Since(req2start), "err", err)
					// Because we only listen for once response to make sure we don't leak a goroutine
					// here also watch for the context to be canceled
					select {
					case respChan <- callbackResponse[T]{resp: r, hedged: true, err: err}:
					case <-ctx2.Done():
					}
				}()
			}
		}
	}

	// Wait for either request to return or on the underlying context to finish
	select {
	case r := <-respChan:
		level.Info(gatewayClient.logger).Log("msg", "returning index gateway response", "hedged", r.hedged, "err", r.err)
		return r.resp, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (s *GatewayClient) getServerAddresses(tenantID string) ([]string, error) {
	var addrs []string
	// The GRPC pool we use only does discovery calls when cleaning up already existing connections,
	// so the list of addresses should always be provided from the external provider (ring or DNS)
	// and not from the RegisteredAddresses method as this list is only populated after a call to GetClientFor
	if s.cfg.Mode == RingMode {
		r := GetShuffleShardingSubring(s.ring, tenantID, s.limits)
		rs, err := r.GetReplicationSetForOperation(IndexesRead)
		if err != nil {
			return nil, errors.Wrap(err, "index gateway get ring")
		}
		addrs = rs.GetAddresses()
	} else {
		addrs = s.dnsProvider.Addresses()
	}

	// shuffle addresses to make sure we don't always access the same Index Gateway instances in sequence for same tenant.
	rand.Shuffle(len(addrs), func(i, j int) {
		addrs[i], addrs[j] = addrs[j], addrs[i]
	})

	return addrs, nil
}

func (s *GatewayClient) NewWriteBatch() index.WriteBatch {
	panic("unsupported")
}

func (s *GatewayClient) BatchWrite(_ context.Context, _ index.WriteBatch) error {
	panic("unsupported")
}

type readBatch struct {
	*logproto.QueryIndexResponse
}

func (r *readBatch) Iterator() index.ReadBatchIterator {
	return &grpcIter{
		i:                  -1,
		QueryIndexResponse: r.QueryIndexResponse,
	}
}

type grpcIter struct {
	i int
	*logproto.QueryIndexResponse
}

func (b *grpcIter) Next() bool {
	b.i++
	return b.i < len(b.Rows)
}

func (b *grpcIter) RangeValue() []byte {
	return b.Rows[b.i].RangeValue
}

func (b *grpcIter) Value() []byte {
	return b.Rows[b.i].Value
}

func instrumentation(cfg ClientConfig, clientRequestDuration *prometheus.HistogramVec) ([]grpc.UnaryClientInterceptor, []grpc.StreamClientInterceptor) {
	var unaryInterceptors []grpc.UnaryClientInterceptor
	unaryInterceptors = append(unaryInterceptors, cfg.GRPCUnaryClientInterceptors...)
	unaryInterceptors = append(unaryInterceptors, otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()))
	unaryInterceptors = append(unaryInterceptors, middleware.ClientUserHeaderInterceptor)
	unaryInterceptors = append(unaryInterceptors, middleware.UnaryClientInstrumentInterceptor(clientRequestDuration))

	var streamInterceptors []grpc.StreamClientInterceptor
	streamInterceptors = append(streamInterceptors, cfg.GRCPStreamClientInterceptors...)
	streamInterceptors = append(streamInterceptors, otgrpc.OpenTracingStreamClientInterceptor(opentracing.GlobalTracer()))
	streamInterceptors = append(streamInterceptors, middleware.StreamClientUserHeaderInterceptor)
	streamInterceptors = append(streamInterceptors, middleware.StreamClientInstrumentInterceptor(clientRequestDuration))

	return unaryInterceptors, streamInterceptors
}
