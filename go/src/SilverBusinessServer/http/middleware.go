package http

import (
	"SilverBusinessServer/http/errcode"
	"SilverBusinessServer/lib"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//MiddleWareLogger :
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

		fmt.Println("", statusCode, method, path, latency, clientIP)

	}
}

// MiddleWareCheckUserHasLogin : 验证是否已经登录了账号
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
