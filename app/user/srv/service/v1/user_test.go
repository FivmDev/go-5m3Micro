package v1

import (
	"context"
	"go-5m3Micro/app/user/srv/data/mock"
	metav1 "go-5m3Micro/pkg/common/meta/v1"
)

func List(ctx context.Context, req *metav1.ListMeta) (*UserDtoList, error) {
	_, _ = NewUserService(mock.NewUsers()).List(ctx, req)
	//TO DO MOCK EVENT
	return nil, nil
}
