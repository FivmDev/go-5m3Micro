package direct

import (
	"google.golang.org/grpc/resolver"
	"strings"
)

func init() {
	resolver.Register(NewDirectBuilder())
}

// NewBuilder creates a directBuilder which is used to factory direct resolvers.
//
// example:
// direct://<authority>/127.0.0.1:9000,/
type directBuilder struct{}

func NewDirectBuilder() *directBuilder {
	return &directBuilder{}
}

func (d *directBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	addrs := make([]resolver.Address, 0)
	for _, addr := range strings.Split(strings.TrimPrefix(target.URL.Path, "/"), ",") {
		addrs = append(addrs, resolver.Address{Addr: addr})
	}
	// gRPC 连接的所有逻辑都从这里走
	err := cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		return nil, err
	}
	return newDirectResolver(), nil
}

func (d *directBuilder) Scheme() string {
	return "direct"
}

var _ resolver.Builder = &directBuilder{}
