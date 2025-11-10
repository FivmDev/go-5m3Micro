package app

import (
	"net/url"
	"os"
	"time"

	"go-5m3Micro/go-5m3Micro/registry"
	restserver "go-5m3Micro/go-5m3Micro/server/rest_server"
	rpcserver "go-5m3Micro/go-5m3Micro/server/rpc_server"
)

type Option func(*options)

type options struct {
	// 服务实例信息
	id        string
	name      string
	endpoints []*url.URL

	// 信号量
	sigs []os.Signal

	// 服务注册注销接口
	registrar         registry.Registrar
	registerTimeout   time.Duration
	deregisterTimeout time.Duration

	// 服务
	rpcServer  *rpcserver.Server
	restServer *restserver.Server
	srvTimeout time.Duration
}

func WithRestServer(srv *restserver.Server) Option {
	return func(o *options) {
		o.restServer = srv
	}
}

func WithRPCServer(srv *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = srv
	}
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithEndpoints(endpoints ...*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithSignal(sigs ...os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}
