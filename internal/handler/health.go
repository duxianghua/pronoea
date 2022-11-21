package handler

import "github.com/gin-gonic/gin"

type HealthHandler struct{}

func (p *HealthHandler) Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
