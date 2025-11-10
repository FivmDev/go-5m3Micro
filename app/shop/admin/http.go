package srv

import (
	"go-5m3Micro/app/user/srv/config"
	restserver "go-5m3Micro/go-5m3Micro/server/rest_server"
)

func NewUserHttpServer(cfg *config.Config) (*restserver.Server, error) {

	srv := restserver.NewServer(restserver.WithPort(cfg.ServerOptions.HttpPort),
		restserver.WithMiddlewares(cfg.ServerOptions.Middlewares))

	initRouter(srv)
	//addr := strings.Builder{}
	//addr.WriteString(cfg.ServerOptions.Host)
	//addr.WriteString(":")
	//addr.WriteString(strconv.Itoa(cfg.ServerOptions.HttpPort))
	//rpcSrv := restserver.NewServer(restserver.WithAddress(addr.String()))
	return srv, nil
}
