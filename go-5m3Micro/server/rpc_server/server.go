package rpc_server

import (
	"context"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"net"
	"net/url"
	"time"

	apimd "go-5m3Micro/api/metadata"
	srvints "go-5m3Micro/go-5m3Micro/server/rpc_server/server_interceptors"
	"go-5m3Micro/pkg/host"
	"go-5m3Micro/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type ServerOption func(*Server)

type Server struct {
	*grpc.Server
	address            string
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	grpcOptions        []grpc.ServerOption
	lis                net.Listener
	health             *health.Server
	endpoint           *url.URL
	metadata           *apimd.Server
	timeout            time.Duration
	enableTracing      bool
	enableMetric       bool
}

func (s *Server) Address() string {
	return s.address
}

func NewServer(opts ...ServerOption) *Server {
	// 默认配置
	srv := &Server{
		address:       ":0",
		health:        health.NewServer(),
		enableTracing: true,
		enableMetric:  true,
	}

	// 填充参数配置
	for _, opt := range opts {
		opt(srv)
	}

	// 添加必要拦截器
	interceptors := []grpc.UnaryServerInterceptor{srvints.UnaryCrashInterceptor}
	if srv.timeout > 0 {
		interceptors = append(interceptors, srvints.UnaryTimeoutInterceptor(srv.timeout))
	}
	if len(srv.unaryInterceptors) > 0 {
		interceptors = append(interceptors, srv.unaryInterceptors...)
	}
	if srv.enableMetric {
		interceptors = append(interceptors, srvints.UnaryMetricsInterceptor)
	}
	// 将多个一元拦截器（Unary Interceptor）链接成一个拦截器链
	options := []grpc.ServerOption{grpc.ChainUnaryInterceptor(interceptors...)}
	// 将参数 grpcOptions []grpc.ServerOption 传入 options
	if len(srv.grpcOptions) > 0 {
		options = append(options, srv.grpcOptions...)
	}

	if srv.enableTracing {
		options = append(options, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}

	// 创建 gRPC 服务
	srv.Server = grpc.NewServer(options...)
	// 创建 metadata 服务
	srv.metadata = apimd.NewServer(srv.Server)
	// 创建 listener
	if err := srv.ListenAndEndpoint(); err != nil {
		return nil
	}
	//注册 健康检查 服务
	grpc_health_v1.RegisterHealthServer(srv.Server, srv.health)
	//注册 metadata 服务 ， 后续通过 gRPC 单一接口 映射 全部服务实例列表
	apimd.RegisterMetadataServer(srv.Server, srv.metadata)

	return srv
}

// ListenAndEndpoint 完成 IP 与 PORT 的提取
func (s *Server) ListenAndEndpoint() error {
	if s.lis == nil {
		lis, err := net.Listen("tcp", s.address)
		if err != nil {
			return err
		}
		s.lis = lis
	}
	addr, err := host.Extract(s.address, s.lis)
	if err != nil {
		_ = s.lis.Close()
		return err
	}
	s.endpoint = &url.URL{
		Scheme: "grpc",
		Host:   addr,
	}
	return nil
}

// Start 启动服务
func (s *Server) Start(ctx context.Context) error {
	log.Infof("[gRPC] server listening on %s", s.lis.Addr().String())
	// 设置服务状态为 SERVING
	s.health.Resume()
	return s.Server.Serve(s.lis)
}

// Exit 退出服务
func (s *Server) Exit(ctx context.Context) error {
	// 设置服务状态为 NOT_SERVING ， 防止接收其它请求
	s.health.Shutdown()
	// server 服务优雅退出
	s.Server.GracefulStop()
	log.Info("[gRPC] server exit")
	return nil
}

func WithEnableTracing(enableTracing bool) ServerOption {
	return func(o *Server) {
		o.enableTracing = enableTracing
	}
}

func WithEnableMetric(enableMetric bool) ServerOption {
	return func(o *Server) {
		o.enableMetric = enableMetric
	}
}

func WithAddress(address string) ServerOption {
	return func(server *Server) {
		server.address = address
	}
}

func WithUnaryInterceptors(unaryInterceptors ...grpc.UnaryServerInterceptor) ServerOption {
	return func(server *Server) {
		server.unaryInterceptors = unaryInterceptors
	}
}

func WithStreamInterceptors(streamInterceptors ...grpc.StreamServerInterceptor) ServerOption {
	return func(server *Server) {
		server.streamInterceptors = streamInterceptors
	}
}

func WithGrpcOptions(grpcOptions ...grpc.ServerOption) ServerOption {
	return func(server *Server) {
		server.grpcOptions = grpcOptions
	}
}

func WithListener(lis net.Listener) ServerOption {
	return func(server *Server) {
		server.lis = lis
	}
}

func WithHealth(health *health.Server) ServerOption {
	return func(server *Server) {
		server.health = health
	}
}

func WithEndpoint(endpoint url.URL) ServerOption {
	return func(server *Server) {
		server.endpoint = &endpoint
	}
}

func WithTimeout(timeout time.Duration) ServerOption {
	return func(server *Server) {
		server.timeout = timeout
	}
}
