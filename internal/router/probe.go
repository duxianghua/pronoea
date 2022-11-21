package router

import (
	"github.com/duxianghua/pronoea/internal/handler"
	"github.com/gin-gonic/gin"
)

func probeApi(r *gin.RouterGroup) {
	health := r.Group("probe") // {} must new line
	{
		health.GET("/", (&handler.ProbeAPI{}).List)
		health.GET("/:namespace/:name", (&handler.ProbeAPI{}).Get)
		health.POST("/:namespace/:name", (&handler.ProbeAPI{}).Create)
		health.DELETE("/:namespace/:name", (&handler.ProbeAPI{}).Delete)
		health.PUT("/:namespace/:name", (&handler.ProbeAPI{}).Update)
		health.GET("/:namespace/:name/status", (&handler.ProbeAPI{}).Status)
	}
}
