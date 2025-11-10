package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
	uproto "go-5m3Micro/api/user/v1"
	"go-5m3Micro/go-5m3Micro/registry/consul"
	rpc "go-5m3Micro/go-5m3Micro/server/rpc_server"
	_ "go-5m3Micro/go-5m3Micro/server/rpc_server/resolver/direct"
	"go-5m3Micro/go-5m3Micro/server/rpc_server/selector"
	"go-5m3Micro/go-5m3Micro/server/rpc_server/selector/random"
	"time"
)

func main() {
	selector.SetGlobalSelector(random.NewBuilder())
	rpc.InitBuilder()
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.145.128:8500"
	cfg.Scheme = "http"
	n, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	r := consul.New(n, consul.WithHealthCheck(false))
	con, err := rpc.DialInsecure(context.Background(),
		rpc.WithClientAddress("discovery:///user-server"),
		//rpc.WithClientTimeout(3*time.Second),
		rpc.WithDiscovery(r))
	if err != nil {
		panic(err)
	}
	cli := uproto.NewUserClient(con)
	for {
		res, err := cli.GetUserListInfo(context.Background(), &uproto.PackageReq{
			PN:    0,
			PSize: 0,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
		time.Sleep(3 * time.Second)
	}
}
