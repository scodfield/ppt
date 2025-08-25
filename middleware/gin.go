package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"ppt/log"
	"runtime/debug"
	"strings"
)

func GinRecover(log *log.LoggerV2, printStack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// check whether the connection is broken
				isBroken := false
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							isBroken = true
						}
					}
				}

				request, _ := httputil.DumpRequest(c.Request, false)
				if isBroken {
					log.Error("http connection is broken", zap.String("req_path", c.Request.URL.Path),
						zap.String("http_request", string(request)), zap.Any("error", err))
					return
				}

				if printStack {
					log.Error("http recover panic", zap.String("req_path", c.Request.URL.Path),
						zap.String("http_request", string(request)), zap.Any("error", err),
						zap.Any("stack", string(debug.Stack())))
				} else {
					log.Error("http recover panic", zap.String("req_path", c.Request.URL.Path),
						zap.String("http_request", string(request)), zap.Any("error", err))
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

func HttpCheckIP(whiteIP []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Request forbidden", http.StatusForbidden)
			return
		}
		ipAddr := net.ParseIP(clientIP)
		if (ipAddr != nil && ipAddr.IsPrivate()) || isWhiteIP(clientIP, whiteIP) {
			next.ServeHTTP(w, r)
			return
		}
	})
}

func isWhiteIP(ip string, whiteIP []string) bool {
	if len(whiteIP) == 0 {
		return true
	}
	for _, v := range whiteIP {
		if v == ip {
			return true
		}
	}
	return false
}
