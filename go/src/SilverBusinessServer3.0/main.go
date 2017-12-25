// SilverBusinessServer project main.go
package main

import (
	"SilverBusinessServer/env"
	"SilverBusinessServer/http"
	"SilverBusinessServer/log"
	"runtime"
)

var (
	cHTTPServerPort = env.Get("httpserver.port").(string)
)

//程序入口
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.HTTP.Info("Listening and serving HTTP on :%v", cHTTPServerPort)
	http.Router.Run("0.0.0.0:" + cHTTPServerPort) //, cHTTPServerPort正式环境的时候添加，测试环境是0.0.0.0:8081
}
