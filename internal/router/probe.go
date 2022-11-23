package router

import (
	"github.com/duxianghua/pronoea/internal/handler"
	"github.com/gin-gonic/gin"
)

func probeApi(r *gin.RouterGroup) {
	health := r.Group("probe") // {} must new line
	{
		health.GET("/", (&handler.ProbeAPI{}).List)
		health.GET("/:name", (&handler.ProbeAPI{}).Get)
		health.POST("/:name", (&handler.ProbeAPI{}).Create)
		health.DELETE("/:name", (&handler.ProbeAPI{}).Delete)
		health.PUT("/:name", (&handler.ProbeAPI{}).Update)
		health.GET("/:name/status", (&handler.ProbeAPI{}).Status)
	}
}
