package utils

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"ppt/login/db"
	"time"
)

var ctx = context.Background()

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		token := c.PostForm("token")
		if token == "" {
			c.JSON(200, gin.H{
				"code":     401,
				"response": "should sync token",
			})
			c.Abort()
			return
		}
		key := FormatTokenKey(name)
		curToken, err := db.GetRedis().Get(ctx, key).Result()
		if errors.Is(err, redis.Nil) {
			c.JSON(200, gin.H{
				"code":     401,
				"response": "token out of date, please login again",
			})
			c.Abort()
			return
		} else if err != nil {
			c.JSON(200, gin.H{
				"code":     501,
				"response": "whoops something went wrong",
			})
			c.Abort()
			return
		}
		if curToken != token {
			c.JSON(200, gin.H{
				"code":     402,
				"response": "token invalid, please check your token",
			})
			c.Abort()
			return
		}
		UpdateToken(c, token)
	}
}

func FormatTokenKey(name string) string {
	return fmt.Sprintf("token_%s", name)
}

func UpdateToken(c *gin.Context, token string) {
	db.GetRedis().Set(ctx, token, token, 1*time.Hour)
}
