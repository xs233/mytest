package main

import (
	"impower/LicenseServer/env"
	"impower/LicenseServer/http"
	"impower/LicenseServer/log"
	"runtime"
)

var (
	cHTTPServerPort = env.Get("httpserver.port").(string)
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	log.Root.Info("Listening and serving HTTP on :%v", cHTTPServerPort)
	// Http
	http.Router.Run("0.0.0.0:" + cHTTPServerPort)
	// Https
	// http.ListenAndServeTLS(":8080", "./res/cert.pem", "./res/key.pem", httpserver.Router)
}
