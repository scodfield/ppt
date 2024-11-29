package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ppt/login/utils"
)

func AdminHandler(r *gin.Engine) {
	admin := r.Group("/admin")
	{
		admin.Use(utils.AuthMiddleware())
		admin.POST("/", AdminMain)
	}
}

func AdminMain(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Data": "This is admin page.",
	})
}
