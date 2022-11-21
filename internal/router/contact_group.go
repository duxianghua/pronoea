package router

import (
	"github.com/duxianghua/pronoea/internal/handler"
	"github.com/gin-gonic/gin"
)

func contactGroupApi(r *gin.RouterGroup) {
	health := r.Group("contactgroup") // {} must new line
	{
		health.GET("/", (&handler.ContactGroupAPI{}).List)
		health.GET("/:name", (&handler.ContactGroupAPI{}).Get)
		health.POST("/:name", (&handler.ContactGroupAPI{}).Create)
		health.DELETE("/:name", (&handler.ContactGroupAPI{}).Delete)
		health.PUT("/:name", (&handler.ContactGroupAPI{}).Update)
	}
}
