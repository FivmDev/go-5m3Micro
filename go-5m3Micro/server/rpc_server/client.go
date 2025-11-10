package rpc_server

import (
	"context"
	"go-5m3Micro/go-5m3Micro/registry"
	"go-5m3Micro/go-5m3Micro/server/rpc_server/client_interceptors"
	"go-5m3Micro/go-5m3Micro/server/rpc_server/resolver/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ClientOption func(*client)

type client struct {
	address string
	timeout time.Duration
	// discovery
	discovery          registry.Discovery
	unaryInterceptors  []grpc.UnaryClientInterceptor
	streamInterceptors []grpc.StreamClientInterceptor
	grpcOptions        []grpc.DialOption
	balanceName        string
}

// WithClientAddress 客户端连接地址
func WithClientAddress(address string) ClientOption {
	return func(c *client) {
		c.address = address
	}
}

// WithClientTimeout 超时时间
func WithClientTimeout(timeout time.Duration) ClientOption {
	return func(c *client) {
		c.timeout = timeout
	}
}

// WithDiscovery 服务发现接口
func WithDiscovery(discovery registry.Discovery) ClientOption {
	return func(c *client) {
		c.discovery = discovery
	}
}

// WithClientUnaryInterceptors 一元拦截器
func WithClientUnaryInterceptors(unaryInterceptors ...grpc.UnaryClientInterceptor) ClientOption {
	return func(c *client) {
		c.unaryInterceptors = unaryInterceptors
	}
}

// WithClientStreamInterceptors 流拦截器
func WithClientStreamInterceptors(streamInterceptors ...grpc.StreamClientInterceptor) ClientOption {
	return func(c *client) {
		c.streamInterceptors = streamInterceptors
	}
}

// WithClientGrpcOptions 配置
func WithClientGrpcOptions(grpcOptions ...grpc.DialOption) ClientOption {
	return func(c *client) {
		c.grpcOptions = grpcOptions
	}
}

// WithClientBalance 负载均衡
func WithClientBalance(balanceName string) ClientOption {
	return func(c *client) {
		c.balanceName = balanceName
	}
}

func DialInsecure(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, true, opts...)
}

func Dial(ctx context.Context, opts ...ClientOption) (*grpc.ClientConn, error) {
	return dial(ctx, false, opts...)
}

func dial(ctx context.Context, isInsecure bool, opts ...ClientOption) (*grpc.ClientConn, error) {
	cli := &client{
		timeout:     2 * time.Second,
		balanceName: "round_robin",
	}
	for _, opt := range opts {
		opt(cli)
	}
	//TODO 客户端默认拦截器
	unaryInterceptors := []grpc.UnaryClientInterceptor{
		client_interceptors.TimeoutUnaryClientInterceptor(cli.timeout)}
	streamInterceptors := []grpc.StreamClientInterceptor{}

	if len(cli.unaryInterceptors) > 0 {
		unaryInterceptors = append(unaryInterceptors, cli.unaryInterceptors...)
	}
	if len(cli.streamInterceptors) > 0 {
		streamInterceptors = append(streamInterceptors, cli.streamInterceptors...)
	}
	grpcOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
		grpc.WithChainStreamInterceptor(streamInterceptors...),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"` + cli.balanceName + `"}`),
	}
	if len(cli.grpcOptions) > 0 {
		grpcOpts = append(grpcOpts, cli.grpcOptions...)
	}
	//TODO 服务发现选项
	if cli.discovery != nil {
		grpcOpts = append(grpcOpts,
			grpc.WithResolvers(discovery.NewBuilder(cli.discovery, discovery.WithInsecure(isInsecure))))
	}
	if isInsecure {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	return grpc.NewClient(cli.address, grpcOpts...)
}
