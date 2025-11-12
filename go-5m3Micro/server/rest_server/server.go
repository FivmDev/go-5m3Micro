package rest_server

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"go-5m3Micro/go-5m3Micro/server/rest_server/middlewares"
	"go-5m3Micro/go-5m3Micro/server/rest_server/pprof"
	"go-5m3Micro/go-5m3Micro/server/rest_server/validation"
	"go-5m3Micro/pkg/errors"
	"go-5m3Micro/pkg/log"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

type JwtInfo struct {
	// 名称 ，默认 go-5m3Micro-jwt
	Realm string
	// 加密
	Key string
	// 超时时间 ，默认 5*day
	Timeout time.Duration
	// 自动刷新时间 ，默认 5*day
	MaxRefresh time.Duration
}

// Server 基于 Gin 封装
type Server struct {
	*gin.Engine
	// 端口号 ，默认 8080
	port int
	// 开发模式 ，默认 debug
	mode string
	// 健康检查 ， 默认 开启 ，如果开启添加 /health 接口
	healthy bool
	// 性能调试 ，默认 开启 ，如果开启添加 /debug/pprof 接口
	enableProfiling bool
	//中间件
	middlewares []string
	// JWT 配置
	jwt *JwtInfo
	// 翻译器
	transName string
	trans     ut.Translator
	// http.server
	httpServer *http.Server
	// tracer name
	tracerName   string
	enableMetric bool
}

func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		Engine:          gin.Default(),
		port:            8080,
		mode:            "debug",
		healthy:         true,
		enableProfiling: true,
		jwt: &JwtInfo{
			Realm:      "go-5m3Micro-jwt",
			Key:        "DgHFLxp7aZ1Pi6hWpRukG8uP",
			Timeout:    5 * 24 * time.Hour,
			MaxRefresh: 5 * 24 * time.Hour,
		},
		transName:    "zh",
		tracerName:   "go-5m3Micro-tracer",
		enableMetric: true,
	}

	for _, opt := range opts {
		opt(srv)
	}

	if srv.enableMetric {
		// get global Monitor object
		m := ginmetrics.GetMonitor()
		// +optional set metric path, default /debug/metrics
		m.SetMetricPath("/metrics")
		// +optional set slow time, default 5s
		m.SetSlowTime(10)
		// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
		// used to p95, p99
		m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
		// set middleware for gin
		m.Use(srv)
	}

	for _, m := range srv.middlewares {
		if mw, ok := middlewares.Middlewares[m]; ok {
			srv.Use(mw)
			log.Infof("middlewares: %v init [success]", m)
		} else {
			log.Warnf("can not find middlewares: %v ", m)
		}
	}

	srv.Use(middlewares.TracingHandler(srv.tracerName))
	return srv
}

func (srv *Server) Start(ctx context.Context) error {
	// 检查 参数 错误
	if srv.mode != gin.DebugMode && srv.mode != gin.TestMode && srv.mode != gin.ReleaseMode {
		return errors.New("must be in one of gin-debug, gin-release, gin-test")
	}
	// 设置 开发环境 ， 打印 日志 格式
	gin.SetMode(srv.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s(%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	// TODO 初始化翻译器
	err := srv.initTranslator(srv.transName)
	if err != nil {
		return err
	}
	// 注册 手机号 验证器
	validation.RegisterMobile(srv.trans)
	// 设置可信代理 设置为 nil
	err = srv.SetTrustedProxies(nil)
	if err != nil {
		return err
	}
	// 根据配置初始化 pprof 路由
	pprof.Register(srv.Engine)
	// 自定义 Run ，目的是为了取到 http.server ，确保优雅退出
	addr := ":" + strconv.Itoa(srv.port)
	srv.httpServer = &http.Server{
		Addr:    addr,
		Handler: srv.Engine,
	}
	if err = srv.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	log.Infof("http_server is running on port %d", srv.port)
	return nil
}

func (srv *Server) Exit(ctx context.Context) error {
	log.Infof("http_server is shutting down on port %d", srv.port)
	if err := srv.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	log.Infof("http_server stopped")
	return nil
}
