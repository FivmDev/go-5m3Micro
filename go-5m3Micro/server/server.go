package server

import "context"

type Server interface {
	Start(ctx context.Context) error
	Exit(ctx context.Context) error
}
