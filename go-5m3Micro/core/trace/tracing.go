package trace

type Options struct {
	// 名称
	Name string `json:"name"`
	// url
	Endpoint string `json:"endpoint"`
	// 采样率
	Sampler float32 `json:"sampler"`
	// 目前仅有 jaeger zapkin
	Batcher string `json:"batcher"`
}

func NewTelemetryOptions() *Options {
	return &Options{
		Name:     "go-5m3Micro",
		Sampler:  1.0,
		Batcher:  "jaeger",
		Endpoint: "http://192.168.145.128:14268/api/traces",
	}
}
