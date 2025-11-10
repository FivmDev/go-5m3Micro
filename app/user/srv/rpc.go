package srv

import (
	uproto "go-5m3Micro/api/user/v1"
	"go-5m3Micro/app/user/srv/config"
	"go-5m3Micro/app/user/srv/controller/user"
	"go-5m3Micro/app/user/srv/data/mock"
	srvv1 "go-5m3Micro/app/user/srv/service/v1"
	rpc "go-5m3Micro/go-5m3Micro/server/rpc_server"
	"strconv"
	"strings"
)

func NewUserRPCServer(cfg *config.Config) (*rpc.Server, error) {
	data := mock.NewUsers()
	srv := srvv1.NewUserService(data)
	userSrv := user.NewUserServer(srv)

	addr := strings.Builder{}
	addr.WriteString(cfg.ServerOptions.Host)
	addr.WriteString(":")
	addr.WriteString(strconv.Itoa(cfg.ServerOptions.Port))
	rpcSrv := rpc.NewServer(rpc.WithAddress(addr.String()))
	uproto.RegisterUserServer(rpcSrv, userSrv)

	//g := gin.Default()
	//uproto.RegisterUserHTTPServer(g, userSrv)
	//g.Run()
	return rpcSrv, nil
}
