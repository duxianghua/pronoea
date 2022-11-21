package internal

import (
	"github.com/duxianghua/pronoea/internal/router"
	"github.com/gin-gonic/gin"
)

func Service() *gin.Engine {
	r := gin.Default()
	router.Register(r)
	return r
}
