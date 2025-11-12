package server_interceptors

import (
	"context"
	"strconv"
	"time"

	"google.golang.org/grpc/status"
	"google.golang.org/grpc"

	"go-5m3Micro/go-5m3Micro/core/metric"
)

const serverNameSpace = "go_5m3Micro_rpc_client"

var (
	metricServerReqDurationMs = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNameSpace,
		Subsystem: "request",
		Name:      "request_duration_ms",
		Help:      "rpc_server_request_duration_ms",
		Labels:    []string{"method"},
		Buckets:   []float64{5, 15, 30, 50, 80, 120},
	})
	metricServerResCode = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNameSpace,
		Subsystem: "response",
		Name:      "response_code",
		Help:      "rpc_server_response_code",
		Labels:    []string{"method", "code"},
	})
)

func UnaryMetricsInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		startTime := time.Now()
		err = invoker(ctx, method, req, reply, cc, opts...)
		metricServerReqDurationMs.Observe(int64(time.Since(startTime)/time.Millisecond), method)
		metricServerResCode.Inc(method, strconv.Itoa(int(status.Code(err))))
		return err
	}
}
