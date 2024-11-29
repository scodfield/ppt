package controllers

import (
	"actor2/util"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"ppt/login/db"
	"ppt/login/utils"
)

var ctx = context.Background()

func LoginHandler(r *gin.Engine) {
	acc := r.Group("/account")
	{
		acc.POST("/registry", AccRegistryHandler)
		acc.POST("/login", AccLoginHandler)
	}
}

type PlayerReq struct {
	Name string `form:"name" json:"name" binding:"required"`
	Pass string `form:"pass" json:"pass" binding:"required"`
}

func AccRegistryHandler(c *gin.Context) {
	var playerReq PlayerReq
	if err := c.ShouldBind(&playerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	existed, err := db.GetRedis().HSetNX(ctx, "table_acc", playerReq.Name, playerReq.Pass).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !existed {
		c.JSON(http.StatusOK, gin.H{"regResponse": "Acc has already registed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"loginResult": "Welcome aboard " + playerReq.Name,
	})
}

func AccLoginHandler(c *gin.Context) {
	var playerReq PlayerReq
	if err := c.ShouldBind(&playerReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	_, err := db.GetRedis().HGet(ctx, "table_acc", playerReq.Name).Result()
	if errors.Is(err, redis.Nil) {
		c.JSON(http.StatusOK, gin.H{
			"loginResult": "Not registry",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token := util.FormatTokenKey(playerReq.Name)
	utils.UpdateToken(c, token)
	c.JSON(http.StatusOK, gin.H{
		"loginResult": "Welcome back " + playerReq.Name,
		"token":       token,
	})
}
