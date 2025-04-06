package middleware

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"ppt/util"
)

const (
	RequestCryptAesKey = "ppt01Dg28eN38Eqd12=="
)

func RequestParse() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10001, "error": err.Error()})
			c.Abort()
			return
		}
		plainBody, err := util.EcbDecrypt(payload, RequestCryptAesKey)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10002, "error": err.Error()})
			c.Abort()
			return
		}
		c.Set("PlainBody", plainBody)
		c.Next()
	}
}
