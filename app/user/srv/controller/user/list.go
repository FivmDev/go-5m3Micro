package user

import (
	"context"
	userpb "go-5m3Micro/api/user/v1"
	srvv1 "go-5m3Micro/app/user/srv/service/v1"
	metav1 "go-5m3Micro/pkg/common/meta/v1"
)

func userDtoToPbResponse(src srvv1.UserDto) userpb.UserInfo {
	return userpb.UserInfo{
		NickName: src.UserDo.Name,
	}
}

func (us *UserServer) List(ctx context.Context, req *userpb.PackageReq) (*userpb.UserListInfo, error) {
	packageReq := &metav1.ListMeta{
		Page:     int(req.PN),
		PageSize: int(req.PSize),
	}
	response, err := us.Srv.List(ctx, packageReq)
	if err != nil {
		return nil, err
	}
	var pbResponse userpb.UserListInfo
	pbResponse.Total = uint64(response.TotalNumber)
	if pbResponse.Total > 0 {
		for _, value := range response.Items {
			pbValue := userDtoToPbResponse(*value)
			pbResponse.Data = append(pbResponse.Data, &pbValue)
		}
	}
	return &pbResponse, nil
}
