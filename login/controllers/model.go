package controllers

import "github.com/gin-gonic/gin"

func RegModelHandler(r *gin.Engine) {
	LoginHandler(r)
	AdminHandler(r)
}
