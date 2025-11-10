package direct

import "google.golang.org/grpc/resolver"

type directResolver struct{}

func newDirectResolver() *directResolver {
	return &directResolver{}
}

func (r *directResolver) ResolveNow(resolver.ResolveNowOptions) {

}

func (r *directResolver) Close() {

}

var _ resolver.Resolver = (*directResolver)(nil)
