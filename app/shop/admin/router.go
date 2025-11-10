package srv

import (
	user "go-5m3Micro/app/shop/admin/controller/user"
	restserver "go-5m3Micro/go-5m3Micro/server/rest_server"
)

func initRouter(r *restserver.Server) {
	v1 := r.Group("/v1")
	u := v1.Group("/users")
	controller := user.NewUserServer()
	u.GET("/list", controller.List)
}
