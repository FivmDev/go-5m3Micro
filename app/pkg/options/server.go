package options

import "github.com/spf13/pflag"

type ServerOptions struct {
	// 是否开启pprof
	EnableProfiling bool `json:"enable-profiling" mapstructure:"enable-profiling"`

	// 是否开启metrics
	EnableMetrics bool `json:"enable-metrics" mapstructure:"enable-metrics"`

	// 是否开启health check
	EnableHealthCheck bool `json:"enable-health-check" mapstructure:"enable-health-check"`

	// host
	Host string `json:"host" mapstructure:"host"`

	// port
	Port int `json:"port" mapstructure:"port"`

	// http port
	HttpPort int `json:"http-port" mapstructure:"http-port"`

	// 名称
	Name string `json:"name" mapstructure:"name"`

	// 中间件
	Middlewares []string `json:"middlewares" mapstructure:"middlewares"`
}

func NewUserServerOptions() *ServerOptions {
	return &ServerOptions{
		EnableProfiling:   true,
		EnableMetrics:     true,
		EnableHealthCheck: true,
		Host:              "127.0.0.1",
		Port:              8078,
		Name:              "user-srv",
	}
}

func (s *ServerOptions) Validate() []error {
	var errs []error
	return errs
}

func (s *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.BoolVar(&s.EnableProfiling, "server.enable-profiling", s.EnableProfiling,
		"enable-profiling, if true, will add <host>:<port>/debug/pprof/, default is true")

	fs.BoolVar(&s.EnableMetrics, "server.enable-metrics", s.EnableMetrics,
		"enable-metrics, if true, will add /metrics, default is true")

	fs.BoolVar(&s.EnableHealthCheck, "server.enable-health-check", s.EnableHealthCheck,
		"enable-health-check, if true, will add health check route, default is true")

	fs.StringVar(&s.Host, "server.host", s.Host, "server host default is 127.0.0.1")

	fs.IntVar(&s.Port, "server.port", s.Port, "server port default is 8078")
	fs.IntVar(&s.HttpPort, "server.http-port", s.HttpPort, "server http port default is 8079")
}
