package user

import (
	"context"
	uproto "go-5m3Micro/api/user/v1"
	srvv1 "go-5m3Micro/app/user/srv/service/v1"
	"go-5m3Micro/pkg/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServer struct {
	Srv srvv1.UserSrv
	uproto.UnimplementedUserServer
}

func (us *UserServer) GetUserListInfo(ctx context.Context, req *uproto.PackageReq) (*uproto.UserListInfo, error) {
	//TODO implement me
	//panic("implement me")
	log.Info("GetUserListInfo is called")
	list, err := us.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (us *UserServer) GetUserByMobile(ctx context.Context, req *uproto.MobileReq) (*uproto.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserServer) GetUserById(ctx context.Context, req *uproto.IdReq) (*uproto.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserServer) CreateUser(ctx context.Context, info *uproto.CreateUserInfo) (*uproto.UserInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserServer) UpdateUser(ctx context.Context, info *uproto.UpdateUserInfo) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserServer) CheckPassword(ctx context.Context, info *uproto.CheckPasswordInfo) (*uproto.CheckResult, error) {
	//TODO implement me
	panic("implement me")
}

func (us *UserServer) mustEmbedUnimplementedUserServer() {
	//TODO implement me
	panic("implement me")
}

func NewUserServer(srv srvv1.UserSrv) *UserServer {
	return &UserServer{Srv: srv}
}
