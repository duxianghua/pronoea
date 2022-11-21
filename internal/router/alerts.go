package router

import (
	"github.com/duxianghua/pronoea/internal/handler"
	"github.com/gin-gonic/gin"
)

func alertsApi(r *gin.RouterGroup) {
	health := r.Group("alerts") // {} must new line
	{
		health.POST("/alertmanager", (&handler.AlertmanagerWebhook{}).Post)
	}
}
