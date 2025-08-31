package controllers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
	"ppt/log"
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

type UserRegistration struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Pass  string `form:"pass" json:"pass" binding:"required"`
	Email string `form:"email" json:"email" binding:"required"`
}

func AccRegistryHandler(c *gin.Context) {
	var userReg UserRegistration
	if err := c.ShouldBind(&userReg); err != nil {
		log.Error("AccRegistryHandler UserRegistration bind error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existed, err := db.WhetherUserNameRegistered(userReg.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !existed {
		log.Info("AccRegistryHandler user name", zap.String("name", userReg.Name), zap.String("email", userReg.Email))
		c.JSON(http.StatusOK, gin.H{"regResponse": "Acc has already registered"})
		return
	}
	err = db.RegUserName(userReg.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if UserID := db.GenerateUserID(); UserID > 0 {
		c.JSON(http.StatusOK, gin.H{
			"loginResult": "Welcome aboard " + userReg.Name,
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
