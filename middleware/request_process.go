package middleware

import (
	"bytes"
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"ppt/logger"
	"ppt/util"
)

const (
	RequestCryptAesKey = "ppt01Dg28eN38Eqd12=="
	SHA256SignKey      = "d2f8A9kB3ew58DnG"
)

func RequestParse() gin.HandlerFunc {
	return func(c *gin.Context) {
		var plainBody []byte
		var originParams, sign string
		var err error
		err = c.Request.ParseForm()
		if err != nil {
			logger.Error("RequestParse ParseForm error", zap.Error(err))
			c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10001, "error": err.Error()})
			c.Abort()
			return
		}
		switch c.Request.Method {
		case http.MethodGet:
			originParams = c.Request.FormValue("params")
			plainBody, err = util.EcbDecrypt(originParams, RequestCryptAesKey)
			if err != nil {
				logger.Error("RequestParse EcbDecrypt error", zap.Error(err), zap.String("originParams", originParams))
				c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10001, "error": err.Error()})
				c.Abort()
				return
			}
		case http.MethodPost, http.MethodPut:
			payload, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Error("RequestParse read request body error", zap.Error(err))
				c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10001, "error": err.Error()})
				c.Abort()
				return
			}
			originParams = string(payload)
			plainBody, err = util.EcbDecrypt(originParams, RequestCryptAesKey)
			if err != nil {
				logger.Error("RequestParse EcbDecrypt originParams error", zap.Error(err), zap.String("originParams", originParams))
				c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10002, "error": err.Error()})
				c.Abort()
				return
			}
		}
		sign = c.Request.FormValue("sign")
		c.Set("OriginParams", originParams)
		c.Set("PlainBody", plainBody)
		c.Set("Sign", sign)
		c.Next()
	}
}

func RequestCheckSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		sign := c.MustGet("Sign").(string)
		originParams := c.MustGet("OriginParams").(string)
		calcSign := util.SignSHA256WithKey(originParams, SHA256SignKey)
		if sign != calcSign {
			logger.Error("RequestCheckSign sign not match", zap.String("req_sign", sign), zap.String("origin_params", originParams), zap.String("calc_sign", calcSign))
			c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10003, "error": "sign not match"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func RequestDecompressGzip() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				logger.Error("RequestDecompressGzip NewReader error", zap.Error(err))
				c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10001, "error": err.Error()})
				c.Abort()
				return
			}
			defer gzipReader.Close()

			body, err := io.ReadAll(gzipReader)
			if err != nil {
				logger.Error("RequestDecompressGzip ReadAll error", zap.Error(err))
				c.JSON(http.StatusBadRequest, map[string]interface{}{"code": 10002, "error": err.Error()})
				c.Abort()
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		c.Next()
	}
}
