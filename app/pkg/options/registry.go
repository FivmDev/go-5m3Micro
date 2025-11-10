package options

import (
	"go-5m3Micro/pkg/errors"

	"github.com/spf13/pflag"
)

type RegistryOptions struct {
	Address string `json:"address"`
	Scheme  string `json:"scheme"`
}

func NewRegistryOptions() *RegistryOptions {
	return &RegistryOptions{
		Address: "127.0.0.1:8500",
		Scheme:  "http",
	}
}

func (r *RegistryOptions) Validate() []error {
	var errs []error
	if len(r.Address) == 0 || len(r.Scheme) == 0 {
		errs = append(errs, errors.New("Registry address or scheme is empty"))
	}
	return errs
}

func (r *RegistryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&r.Address, "address", r.Address, "The address and port of the server")
	fs.StringVar(&r.Scheme, "scheme", r.Scheme, "The scheme of the server")
}
