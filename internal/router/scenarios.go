package router

import (
	"github.com/duxianghua/pronoea/internal/handler"
	"github.com/gin-gonic/gin"
)

func ScenariosAPI(r *gin.RouterGroup) {
	scenarios := r.Group("scenarios") // {} must new line
	{
		scenarios.GET("/", (&handler.ScenariosAPI{}).List)
		scenarios.GET("/:name", (&handler.ScenariosAPI{}).Get)
		scenarios.POST("/:name", (&handler.ScenariosAPI{}).Create)
		scenarios.DELETE("/:name", (&handler.ScenariosAPI{}).Delete)
		scenarios.PUT("/:name", (&handler.ScenariosAPI{}).Update)
		scenarios.GET("/:name/status", (&handler.ScenariosAPI{}).Status)
		scenarios.PATCH("/:name", (&handler.ScenariosAPI{}).Patch)
	}
}
