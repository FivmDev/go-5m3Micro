package server_interceptors

import (
	"context"
	"google.golang.org/grpc/status"
	"strconv"
	"time"

	"go-5m3Micro/go-5m3Micro/core/metric"
	"google.golang.org/grpc"
)

const serverNameSpace = "go_5m3Micro_rpc_server"

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

func UnaryMetricsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	startTime := time.Now()
	resp, err = handler(ctx, req)
	metricServerReqDurationMs.Observe(int64(time.Since(startTime)/time.Millisecond), info.FullMethod)
	metricServerResCode.Inc(info.FullMethod, strconv.Itoa(int(status.Code(err))))
	return resp, err
}
