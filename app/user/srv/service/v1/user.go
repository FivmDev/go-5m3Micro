package v1

import (
	"context"
	"go-5m3Micro/app/user/srv/data"
	metav1 "go-5m3Micro/pkg/common/meta/v1"
)

type UserDto struct {
	//Name string `json:"name"`
	UserDo data.UserDo
}

type UserDtoList struct {
	TotalNumber int32      `json:"total_number"`
	Items       []*UserDto `json:"items,omitempty"`
}

type UserSrv interface {
	List(ctx context.Context, req *metav1.ListMeta) (*UserDtoList, error)
}

type UserService struct {
	userStore data.UserStore
}

func NewUserService(us data.UserStore) *UserService {
	return &UserService{userStore: us}
}

func (u *UserService) List(ctx context.Context, req *metav1.ListMeta) (*UserDtoList, error) {
	response, err := u.userStore.List(ctx, req)
	if err != nil {
		return nil, err
	}
	var userDtoListRes UserDtoList
	userDtoListRes.TotalNumber = response.TotalNumber
	if userDtoListRes.TotalNumber > 0 {
		for _, value := range response.Items {
			userDtoListRes.Items = append(userDtoListRes.Items, &UserDto{*value})
		}
	}
	return &userDtoListRes, nil
}

var _ UserSrv = (*UserService)(nil)
