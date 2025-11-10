package srv

import (
	"github.com/hashicorp/consul/api"
	"go-5m3Micro/app/pkg/options"
	"go-5m3Micro/app/user/srv/config"
	gapp "go-5m3Micro/go-5m3Micro/app"
	"go-5m3Micro/go-5m3Micro/registry"
	"go-5m3Micro/go-5m3Micro/registry/consul"
	"go-5m3Micro/pkg/app"
	"go-5m3Micro/pkg/log"
)

func NewApp(basename string) *app.App {
	cfg := config.New()
	application := app.NewApp("go-5m3Micro", basename,
		app.WithOptions(cfg),
		app.WithRunFunc(run(cfg)),
		app.WithNoConfig())
	return application
}

func NewRegistrar(registry *options.RegistryOptions) (registry.Registrar, error) {
	cfg := api.DefaultConfig()
	cfg.Address = registry.Address
	cfg.Scheme = registry.Scheme
	cli, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return consul.New(cli, consul.WithHealthCheck(false)), nil
}

func NewUserApp(cfg *config.Config) (*gapp.App, error) {
	log.Init(cfg.Log)
	defer log.Flush()
	register, err := NewRegistrar(cfg.Registry)
	if err != nil {
		return nil, err
	}
	rpcServer, err := NewUserRPCServer(cfg)
	if err != nil {
		return nil, err
	}
	return gapp.New(gapp.WithName("user-server"),
		gapp.WithRPCServer(rpcServer),
		gapp.WithRegistrar(register)), nil
}

func run(cfg *config.Config) app.RunFunc {
	return func(basename string) error {
		userApp, err := NewUserApp(cfg)
		if err != nil {
			return err
		}
		if err = userApp.Run(); err != nil {
			log.Errorf("run user app err:%v", err)
			return err
		}
		return nil
	}
}
