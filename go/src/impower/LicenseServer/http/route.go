package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"

	"impower/LicenseServer/env"
	v1api "impower/LicenseServer/http/v1/api"
)

// Router for gin
var (
	Router     = gin.New()
	RootRouter = Router.Group("/")
	V1Router   = RootRouter.Group("/v1")
	V2Router   = RootRouter.Group("/v2")
)

// CorsConf :
var CorsConf = cors.Config{
	Origins:         `http://localhost:8080`,
	Methods:         "GET, PUT, POST, DELETE, OPTIONS",
	RequestHeaders:  "Origin, Authorization, Content-Type",
	ExposedHeaders:  "",
	MaxAge:          60 * time.Second,
	Credentials:     true,
	ValidateHeaders: false,
}

func binds() {
	/*------------------V1.api---------------------*/
	if env.Get("httpserver.module.v1.api").(bool) {
		//acount module
		V1Router.POST("/api/sessions", v1api.HandleLoginPost)
		V1Router.PUT("/api/password/modify", v1api.HandleModifyPasswordPut)
		V1Router.PUT("/api/password/reset", MiddleWareCheckUserHasLogin(), v1api.HandleResetPasswordPut)
		V1Router.DELETE("/api/sessions/:sid", MiddleWareCheckUserHasLogin(), v1api.HandleLogoutDelete)
		V1Router.GET("/api/users/:uid", MiddleWareCheckUserHasLogin(), v1api.HandleQueryUserInfoGet)
		V1Router.DELETE("/api/users/:uid", MiddleWareCheckUserHasLogin(), v1api.HandleUserDelete)
		V1Router.GET("/api/users", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAllUserInfoGet)
		V1Router.POST("/api/users", MiddleWareCheckUserHasLogin(), v1api.HandleCreateUserPost)
		V1Router.POST("/api/phones/:phone/codes", v1api.HandleGeneratePhoneCodePost)

		//license module
		V1Router.GET("/api/mac/record", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAuthRecordByPageGet)
		V1Router.POST("/api/mac/record", MiddleWareCheckUserHasLogin(), v1api.HandleAuthRecordPost)
		V1Router.GET("/api/mac/record/batch", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAuthMacByBatchGet)
		V1Router.GET("/api/license", v1api.HandleApplyLicenseGet)
		V1Router.GET("/api/p2p/record", MiddleWareCheckUserHasLogin(), v1api.HandleQueryP2PRecordByPageGet)
		V1Router.POST("/api/p2p/record", MiddleWareCheckUserHasLogin(), v1api.HandleP2PRecordPost)
		V1Router.GET("/api/p2p/record/batch", MiddleWareCheckUserHasLogin(), v1api.HandleQueryP2PByBatchGet)
		V1Router.GET("/api/p2p", v1api.HandleApplyP2PGet)
	}

}

func init() {
	//Set gin mode
	ginModeDebug := env.Get("debug").(bool)
	if ginModeDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	Router.Use(gin.Recovery(), MiddleWareLogger(), cors.Middleware(CorsConf))
	RootRouter.Use(MiddleWareLogger(), cors.Middleware(CorsConf))
	V1Router.Use(MiddleWareLogger(), cors.Middleware(CorsConf))
	V2Router.Use(MiddleWareLogger(), cors.Middleware(CorsConf))
	binds()
}
