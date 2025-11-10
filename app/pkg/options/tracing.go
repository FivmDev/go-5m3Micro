package options

import (
	"go-5m3Micro/pkg/errors"

	"github.com/spf13/pflag"
)

type TelemetryOptions struct {
	// 名称
	Name string `json:"name"`
	// url
	Endpoint string `json:"endpoint"`
	// 采样率
	Sampler float32 `json:"sampler"`
	// 目前仅有 jaeger zipkin
	Batcher string `json:"batcher"`
}

func NewTelemetryOptions() *TelemetryOptions {
	return &TelemetryOptions{
		Name:     "go-5m3Micro",
		Sampler:  1.0,
		Batcher:  "jaeger",
		Endpoint: "http://192.168.145.128:14268/api/traces",
	}
}

func (to *TelemetryOptions) Validate() []error {
	errs := []error{}
	if len(to.Endpoint) == 0 {
		errs = append(errs, errors.New("missing endpoint"))
	}
	if to.Batcher != "jaeger" && to.Batcher != "zipkin" {
		errs = append(errs, errors.New("batcher must be one of 'jaeger', 'zipkin'"))
	}
	return errs
}

func (to *TelemetryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&to.Name, "telemetry-name", to.Name, "telemetry-name")
	fs.StringVar(&to.Endpoint, "telemetry-endpoint", to.Endpoint, "The endpoint to send traces to")
	fs.Float32Var(&to.Sampler, "telemetry-sampler", to.Sampler, "The sampler to use")
	fs.StringVar(&to.Batcher, "telemetry-batcher", to.Batcher, "batcher must be one of 'jaeger', 'zipkin'")
}
