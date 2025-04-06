package controllers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"ppt/login/db"
	"ppt/util"
)

var ctx = context.Background()

func LoginHandler(r *gin.Engine) {
	acc := r.Group("/account")
	{
		acc.GET("/login", LoginGetHandler)
		acc.POST("/registry", AccRegistryHandler)
		acc.POST("/login", AccLoginHandler)
	}
}

func LoginGetHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

type PlayerReg struct {
	Name string `form:"name" json:"name" binding:"required"`
	Pass string `form:"pass" json:"pass" binding:"required"`
}

func AccRegistryHandler(c *gin.Context) {
	var playerReg PlayerReg
	if err := c.ShouldBind(&playerReg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	existed, err := db.GetRedis().HSetNX(ctx, "table_acc", playerReg.Name, playerReg.Pass).Result()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !existed {
		c.JSON(http.StatusOK, gin.H{"regResponse": "Acc has already registed"})
		return
	}
	if UserID := db.GenerateUserID(); UserID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"loginResult": "Welcome aboard " + playerReg.Name,
			"userID":      UserID,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"loginResult": "Sorry, something went wrong ",
	})
}

type PlayerLogin struct {
	Name string `form:"name" json:"name" binding:"required"`
	ID   int64  `form:"id" json:"id" binding:"required"`
}

func AccLoginHandler(c *gin.Context) {
	var playerLogin PlayerLogin
	if err := c.ShouldBind(&playerLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	_, err := db.GetRedis().HGet(ctx, "table_acc", playerLogin.Name).Result()
	if errors.Is(err, redis.Nil) {
		c.JSON(http.StatusOK, gin.H{
			"loginResult": "Not registry",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//token := util.FormatTokenKey(playerLogin.Name)
	token, _ := util.GenerateToken(playerLogin.ID, playerLogin.Name)
	db.UpdateToken(playerLogin.ID, token)
	c.JSON(http.StatusOK, gin.H{
		"loginResult": "Welcome back " + playerLogin.Name,
		"token":       token,
	})
}
