package app

import (
	"context"
	"golang.org/x/sync/errgroup"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go-5m3Micro/go-5m3Micro/registry"
	"go-5m3Micro/go-5m3Micro/server"
	"go-5m3Micro/pkg/log"

	"github.com/google/uuid"
)

type App struct {
	opt    options
	mutex  sync.Mutex
	sr     *registry.ServiceInstance
	cancel func()
}

func New(opts ...Option) *App {
	o := options{
		sigs:              []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL},
		registerTimeout:   5 * time.Second,
		deregisterTimeout: 5 * time.Second,
	}

	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	for _, opt := range opts {
		opt(&o)
	}
	return &App{opt: o}
}

func (a *App) Run() error {
	//获取注册服务实例信息
	sr, err := a.BuildInstance()
	if err != nil {
		log.Errorf("创建服务注册实例的信息 : %v", err)
		return err
	}
	a.mutex.Lock()
	a.sr = sr
	a.mutex.Unlock()

	var srvs []server.Server
	if a.opt.rpcServer != nil {
		srvs = append(srvs, a.opt.rpcServer)
	}
	if a.opt.restServer != nil {
		srvs = append(srvs, a.opt.restServer)
	}
	// 启动 rpc server 与 rest server
	tempCtx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	eg, ctx := errgroup.WithContext(tempCtx)
	var wg sync.WaitGroup
	for _, srv := range srvs {
		wg.Add(1)
		eg.Go(func() error {
			<-ctx.Done()
			// 如果 某服务 启动失败 ，自动执行 优雅退出 ，但需考虑 退出超时
			c, tempCancel := context.WithTimeout(context.Background(), a.opt.srvTimeout)
			defer tempCancel()
			return srv.Exit(c)

		})
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}
	wg.Wait()

	//注册服务
	if a.opt.registrar != nil {
		ctx, cCancel := context.WithTimeout(context.Background(), a.opt.registerTimeout)
		defer cCancel()
		err = a.opt.registrar.Register(ctx, a.sr)
		if err != nil {
			log.Errorf("注册服务失败 : %v", err)
			return err
		}
	}

	//监听退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opt.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return a.Exit()
		}
	})
	if err = eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (a *App) Exit() error {
	a.mutex.Lock()
	sr := a.sr
	a.mutex.Unlock()

	// 注销 registrarServer
	if a.opt.registrar != nil {
		log.Infof("start to deregister")
		ctx, cCancel := context.WithTimeout(context.Background(), a.opt.deregisterTimeout)
		defer cCancel()
		err := a.opt.registrar.Deregister(ctx, sr)
		if err != nil {
			log.Errorf("注销服务失败 : %v", err)
			return err
		}
	}

	// 注销 rpcServer restServer
	a.cancel()

	return nil
}

// BuildInstance 创建服务注册实例的信息
func (a *App) BuildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	// 从 rpcServer 主动获取信息
	if a.opt.rpcServer != nil {
		endpoint := url.URL{
			Scheme: "grpc",
			Host:   a.opt.rpcServer.Address(),
		}
		endpoints = append(endpoints, endpoint.String())
	}
	//// 从 restServer 主动获取信息
	//if a.opt.restServer != nil {
	//	endpoint := url.URL{
	//		Scheme: "http",
	//		Host:   a.opt.restServer.,
	//	}
	//	endpoints = append(endpoints, endpoint.String())
	//}

	for _, e := range a.opt.endpoints {
		endpoints = append(endpoints, e.String())
	}

	return &registry.ServiceInstance{
		Id:        a.opt.id,
		Name:      a.opt.name,
		Endpoints: endpoints,
	}, nil
}
