package config

import (
	options "go-5m3Micro/app/pkg/options"
	cliflag "go-5m3Micro/pkg/common/cli/flag"
	"go-5m3Micro/pkg/log"
)

type Config struct {
	Log           *log.Options             `json:"log" mapstructure:"log"`
	ServerOptions *options.ServerOptions   `json:"server" mapstructure:"server"`
	Registry      *options.RegistryOptions `json:"registry" mapstructure:"registry"`
}

func (c *Config) Flags() cliflag.NamedFlagSets {
	var fss cliflag.NamedFlagSets
	logsFlagSet := fss.FlagSet("logs")
	serverFlagSet := fss.FlagSet("server")
	registryFlagSet := fss.FlagSet("registry")
	c.Log.AddFlags(logsFlagSet)
	c.ServerOptions.AddFlags(serverFlagSet)
	c.Registry.AddFlags(registryFlagSet)
	return fss
}

func (c *Config) Validate() []error {
	var errs []error
	errs = append(errs, c.Log.Validate()...)
	errs = append(errs, c.ServerOptions.Validate()...)
	errs = append(errs, c.Registry.Validate()...)
	return errs
}

func New() *Config {
	cfg := &Config{
		Log:           log.NewOptions(),
		ServerOptions: options.NewUserServerOptions(),
		Registry:      options.NewRegistryOptions()}
	return cfg
}
