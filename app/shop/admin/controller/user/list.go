package user

import (
	"github.com/gin-gonic/gin"
	"go-5m3Micro/pkg/log"
)

func (us *userServer) List(c *gin.Context) {
	log.Info("List is called")
}
