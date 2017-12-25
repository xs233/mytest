package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"impower/LicenseServer/http/errcode"
	"impower/LicenseServer/lib"
	"impower/LicenseServer/log"
)

// MiddleWareLogger :
func MiddleWareLogger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		log.Root.Info("[%v %v %v] %v %v", statusCode, method, path, latency, clientIP)

	}
}

// MiddleWareCheckUserHasLogin :
func MiddleWareCheckUserHasLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := lib.GetSecretCookie(c.Request, "account"); err == nil {
			c.Set("account", cookie)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err":    errcode.ErrUnLogin.Code,
				"errmsg": errcode.ErrUnLogin.String,
			})
			c.Abort()
		}
	}
}
