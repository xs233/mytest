package http

import (
	"SilverBusinessServer/env"
	v1api "SilverBusinessServer/http/v1/api"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

//Router for gin
var (
	Router     = gin.New()
	RootRouter = Router.Group("/")
	V1Router   = RootRouter.Group("/v1")
	V2Router   = RootRouter.Group("/v2")
)

// CorsConf :
var CorsConf = cors.Config{
	Origins:         `http://localhost:8081`, //localhost
	Methods:         "GET, PUT, POST, DELETE, OPTIONS",
	RequestHeaders:  "Origin, Authorization, Content-Type",
	ExposedHeaders:  "",
	MaxAge:          60 * time.Second,
	Credentials:     true,
	ValidateHeaders: false,
}

func binds() {
	/*-------------------V1.API定义接口-------------------*/
	if env.Get("httpserver.module.v1.api").(bool) { // 注释后面的OK，标示的是测试代码正常返回成功

		//account 账号
		V1Router.POST("/api/sessions", v1api.HandleLoginPost)                                          //用户登陆 ok ok
		V1Router.PUT("/api/password/modify", v1api.HandleModifyPasswordPut)                            //修改密码 ok ok
		V1Router.PUT("/api/password/reset", v1api.HandleResetPasswordPut)                              //重制密码 ok ok
		V1Router.PUT("/api/right/modify", v1api.HandleRightPut)										   //重置权限
		V1Router.DELETE("/api/sessions/:sid", MiddleWareCheckUserHasLogin(), v1api.HandleLogoutDelete) //用户注销 ok ok
		V1Router.GET("/api/users/:uid", MiddleWareCheckUserHasLogin(), v1api.HandleGainUserInfoGet)    //查询用户个人信息 ok ok
		V1Router.DELETE("/api/users/:uid", MiddleWareCheckUserHasLogin(), v1api.HandleUserDelete)      //删除用户 ok ok
		V1Router.GET("/api/users", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAllUserInfoGet)     //查询所有用户信息 ok ok
		V1Router.POST("/api/users", MiddleWareCheckUserHasLogin(), v1api.HandleAddUsersPost)           //添加用户 ok ok

		//Device 设备
		V1Router.GET("/api/devices", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAllCamerasListGet)     //用户查询所有设备列表，用户ID从cookie中获取 ok ok
		V1Router.POST("/api/devices", MiddleWareCheckUserHasLogin(), v1api.HandleImportCameraFromPlistPost) //从配置文件中导入摄像机 ok ok
		//V1Router.POST("/api/devices", MiddleWareCheckUserHasLogin(), v1api.HandleImportCameraFromCMSPost)                        //从CMS中导入摄像机
		V1Router.PUT("/api/devices/:did", MiddleWareCheckUserHasLogin(), v1api.HandleChangeCameraNamePut)                        //修改设备名 ok ok
		V1Router.DELETE("/api/devices/:did", MiddleWareCheckUserHasLogin(), v1api.HandleDeleteCameraDelete)                      //删除设备			---10月11日更新添加接口
		V1Router.POST("/api/groups", MiddleWareCheckUserHasLogin(), v1api.HandleCreateCamerasGroupPost)                          //创建设备组列表 ok ok
		V1Router.GET("/api/groups", MiddleWareCheckUserHasLogin(), v1api.HandleQueryDeviceGroupsGet)                             //用户查询设备组列表 ok ok
		V1Router.GET("/api/groups/:gid", MiddleWareCheckUserHasLogin(), v1api.HandleQueryDeviceGroupDetailGet)                   //用户查询设备组详情 ok ok
		V1Router.PUT("/api/groups/:gid", MiddleWareCheckUserHasLogin(), v1api.HandleUpdateDeviceGroupInfoPut)                    //用户修改设备组名 ok ok
		V1Router.DELETE("/api/groups/:gid", MiddleWareCheckUserHasLogin(), v1api.HandleDeleteDeviceGroupDelete)                  //用户删除相机组 ok ok
		V1Router.GET("/api/nongroup/devices", MiddleWareCheckUserHasLogin(), v1api.HandleQueryNonGroupDeviceListGet)             //查询所有未分组设备 ok ok
		V1Router.GET("/api/nongroup/page/devices", MiddleWareCheckUserHasLogin(), v1api.HandleQueryNonGroupDeviceListByPageGet)  //分页查询未分组的设备 ok 暂时不用，以后需要的可以扩展（未测试）
		V1Router.DELETE("/api/groups/:gid/devices/:did", MiddleWareCheckUserHasLogin(), v1api.HandleDeleteDeviceFromGroupDelete) //从设备组中删除一个设备 ok ok
		V1Router.POST("/api/groups/:gid/devices", MiddleWareCheckUserHasLogin(), v1api.HandleAddGroupDevicesPost)                //向设备组中添加设备 ok ok
		V1Router.PUT("/api/groups/:gid/devices", MiddleWareCheckUserHasLogin(), v1api.HandleUpdateGroupDevicesPut)               //更新设备组设备列表 ok ok
		V1Router.DELETE("/api/groups/:gid/devices", MiddleWareCheckUserHasLogin(), v1api.HandleDeleteGroupDevicesDelete)         //从设备组中删除设备列表 ok ok

		//alarm 实时报警
		V1Router.POST("/api/alarm/real", MiddleWareCheckUserHasLogin(), v1api.HandleAddAlarmRemarkPost)  //添加报警备注，通过相机ID和报警时间戳标识报警
		V1Router.PUT("/api/alarm/real", MiddleWareCheckUserHasLogin(), v1api.HandlePigeonholeAlarmPut)   //处理或归档报警，通过相机ID和报警时间戳标识报警
		V1Router.DELETE("/api/alarm/real", MiddleWareCheckUserHasLogin(), v1api.HandleDeleteAlarmDelete) //删除报警，通过相机ID和报警时间戳标识报警

		//alarm 历史报警
		V1Router.GET("/api/alarms", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAlarmsByPageGet)           //分页查询报警，查询参数中archiveFlag=1表示只查询归档的报警，其他值表示查询所有报警ok ok
		V1Router.DELETE("/api/alarms", MiddleWareCheckUserHasLogin(), v1api.HandleDeleteAlarmsByTimeDelete)    //按时间段批量删除报警，归档报警不删除，只删除未归档报警ok ok
		V1Router.GET("/api/images/:iid", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAlarmImageByidGet)    //查询报警图片 ok ok
		V1Router.POST("/api/alarms/:aid", MiddleWareCheckUserHasLogin(), v1api.HandleAddAlarmRemarkByOldPost)  //历史报警：，添加报警备注，aid表示的是报警id ok ok
		V1Router.PUT("/api/alarms/:aid", MiddleWareCheckUserHasLogin(), v1api.HandlePigeonholeAlarmByOldPut)   //历史报警：处理或归档报警，通过相机ID和报警时间戳标识报警 ok ok
		V1Router.DELETE("/api/alarms/:aid", MiddleWareCheckUserHasLogin(), v1api.HandleDeleteAlarmByOldDelete) //历史报警：删除报警，通过相机ID和报警时间戳标识报警 ok ok

		//monitor 算法监控
		V1Router.GET("/api/alg/devices", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAllAlgDevicesGet)         //查询所有算法分析相机
		V1Router.POST("/api/alg/devices", MiddleWareCheckUserHasLogin(), v1api.HandleAddAlgDevicePost)             //添加摄像机
		V1Router.GET("/api/alg/devices/:did", MiddleWareCheckUserHasLogin(), v1api.AlgConfigGet)				   //获取算法配置
		V1Router.POST("/api/alg/devices/:did", MiddleWareCheckUserHasLogin(), v1api.HandleControlAlgPost)          //启停相机算法
		V1Router.PUT("/api/alg/devices/:did", MiddleWareCheckUserHasLogin(), v1api.HandleUpdateAlgPut)             //修改算法配置，须先停止算法运行，后台须作判断
		V1Router.DELETE("/api/alg/devices/:did", MiddleWareCheckUserHasLogin(), v1api.HandleDeleteAlgDeviceDelete) //删除相机，须先停止算法运行，后台须作判断
		V1Router.GET("/api/alg/status", MiddleWareCheckUserHasLogin(), v1api.HandleQueryAllAlgDevicesStatus)       //客户端定时查询所有设备的算法运行状况，定时15秒

	}
}

func init() {
	//Set gin mode  debug
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
