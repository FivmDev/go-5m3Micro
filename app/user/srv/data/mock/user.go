package mock

import (
	"context"
	"go-5m3Micro/app/user/srv/data"
	metav1 "go-5m3Micro/pkg/common/meta/v1"
)

type users struct {
}

func NewUsers() *users {
	return &users{}
}

func (u users) List(ctx context.Context, req *metav1.ListMeta) (*data.UserDoList, error) {
	//TO DO MOCK EVENT
	return &data.UserDoList{
		TotalNumber: 1,
		Items: []*data.UserDo{
			{
				Name: "test1",
			},
		},
	}, nil
}

// 接口验证
var _ data.UserStore = (*users)(nil)
