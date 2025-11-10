package middlewares

import "github.com/gin-gonic/gin"

var Middlewares = defaultMiddleware()

func defaultMiddleware() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery": gin.Recovery(),
		"cors":     Cors(),
	}
}
