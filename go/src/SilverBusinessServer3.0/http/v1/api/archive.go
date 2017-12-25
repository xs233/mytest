package api

import (
	"SilverBusinessServer/dao"
	"SilverBusinessServer/http/errcode"
	"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

/////////////////////////////////////////////////////设备管理device/////////////////////////////////////////////////////
//#define CMD_DEVICES	//从CMS系统中导入设备，查询系统中所有设备
//URL: "/v1/api/devices"
//用户查询所有设备列表，用户ID从cookie中获取
func HandleQueryAllCamerasListGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryAllCamerasListGet BEGIN")
	//Get current login user id from cookie	 这个是表示从cookie中获取用户ID
	//currentUserID, err := lib.GetCurrentUser(c.Request)
	//if err != nil {
	//	log.HTTP.Error("获取用户的user id失败")
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrCookie.Code,
	//		"errMsg": errcode.ErrCookie.String,
	//	})
	//	return
	//}
	//if false == dao.IsAdmin(currentUserID) {
	//	log.HTTP.Error("验证是不是管理员")
	//	//Judge user whether valid or not
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrFakeRequest.Code,
	//		"errMsg": errcode.ErrFakeRequest.String,
	//	})
	//	return
	//}

	deviceList, err := dao.QueryAllDevices()
	if err != nil {
		log.HTTP.Error("获取所有设备错误")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":        errcode.ErrNoError.Code,
		"errMsg":     errcode.ErrNoError.String,
		"deviceList": deviceList,
	})
	return
}

//------------------------------------------------------------------------------------------
//从配置文件中导入摄像  这是需要将数据保存到数据库中的
func HandleImportCameraFromPlistPost(c *gin.Context) {
	log.HTTP.Info("HandleImportCameraFromPlistPost BEGIN")
	var reqJSON reqDeviceListIntoJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("JSON Err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if dao.ImportDeviceToList(reqJSON.DeviceList) != nil {
		log.HTTP.Error("导入设备列表失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrAddDidList.Code,
			"errMsg": errcode.ErrAddDidList.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqDeviceListIntoJSON struct {
	DeviceList []lib.Device `form:"deviceList" json:"deviceList"`
}

//------------------------------------------------------------------------------------------
//从CMS中导入相机
func HandleImportCameraFromCMSPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

//---------------------------------------------------------------------------------------------------------------------------
//#define CMD_DEVICES	//修改设备名		从url中获取设备的id 在数据库dao中把设备名改了
//URL: "/v1/api/devices/:did"
func HandleChangeCameraNamePut(c *gin.Context) {
	log.HTTP.Info("修改设备名")
	DeviceID, err := strconv.ParseInt(c.Param("did"), 10, 64) //表示从url中获取设备id
	if err != nil {
		log.HTTP.Error("get device id err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	var reqJSON reqChangeCameraNameJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("JSON Err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrUpdate.Code,
			"errMsg": errcode.ErrUpdate.String,
		})
		return
	}
	if dao.QueryChangeDeviceName(DeviceID, reqJSON.DeviceName) != nil {
		log.HTTP.Error("修改设备名失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrUpdate.Code,
			"errMsg": errcode.ErrUpdate.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqChangeCameraNameJSON struct {
	DeviceName string `form:"deviceName" json:"deviceName"`
}

//删除设备
func HandleDeleteCameraDelete(c *gin.Context) {
	//1.从url中获取deviceid
	DeviceID, err := strconv.ParseInt(c.Param("did"), 10, 64) //表示从url中获取设备id
	if err != nil {
		log.HTTP.Error("get device id err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//2.校验是不是管理员
	//Get current login user id from cookie	 这个是表示从cookie中获取用户ID
	currentUserID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.HTTP.Error("获取用户的user id失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}
	if false == dao.IsAdmin(currentUserID) {
		log.HTTP.Error("验证是不是管理员")
		//Judge user whether valid or not
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	//3判断alg表中是否存在该deviceid
	alg, err := dao.GetAlg(DeviceID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrGetALG.Code,
			"errMsg": errcode.ErrGetALG.String,
		})
		return
	}
	//能查到则告诉操作人员
	if alg.DeviceID != 0 {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrExistAlg.Code,
			"errMsg": errcode.ErrExistAlg.String,
		})
		return
	}

	//	//3.判断算法运行状态
	//	//判断imp_t_alg 表中的task_id 是否为空
	//	sign, taskID, err := JudgedTask(DeviceID)
	//	if err != nil || sign == -1 {
	//		c.JSON(http.StatusOK, gin.H{
	//			"err":    errcode.ErrParams.Code,
	//			"errMsg": errcode.ErrParams.String,
	//		})
	//		return
	//	}
	//	//task_id 不为空
	//	if sign == 1 {
	//		//判断任务是否在运行
	//		sign1, err := JudgedRunningState(taskID)
	//		if sign1 == -1 || err != nil {
	//			c.JSON(http.StatusOK, gin.H{
	//				"err":    errcode.ErrJudgedRunningState.Code,
	//				"errMsg": errcode.ErrJudgedRunningState.String,
	//			})
	//			return
	//			//如果在运行则提示用户停止算法运行
	//		} else if sign1 == 1 {
	//			c.JSON(http.StatusOK, gin.H{
	//				"err":    errcode.ErrStopALG.Code,
	//				"errMsg": errcode.ErrStopALG.String,
	//			})
	//			return
	//		}
	//	}

	//4.进入数据库中操作，获取是不是设备组中的，不是话，从设备表中直接删除
	if dao.DeleteDevice(DeviceID) != nil {
		log.HTTP.Error("删除设备失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrDelete.Code,
			"errMsg": errcode.ErrDelete.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}
