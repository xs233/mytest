package api

import (
	"SilverBusinessServer/dao"
	"SilverBusinessServer/http/errcode"
	//"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//----------------------------------------------------------------------------------------------------------------------------
// #define CMD_GROUPS	//设备组
// 用户查询设备组列表，用户ID从cookie中获取
// URL: "/v1/api/groups"		首先，从cookie中获取user id
func HandleQueryDeviceGroupsGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryDeviceGroupsGet BEGIN")
	//Get current login user id from cookie
	//currentUserID, err := lib.GetCurrentUser(c.Request)
	//if err != nil {
	//	log.HTTP.Error("获取登录用户的userid失败")
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrCookie.Code,
	//		"errMsg": errcode.ErrCookie.String,
	//	})
	//	return
	//}
	//if false == dao.IsAdmin(currentUserID) {
	//	log.HTTP.Error("验证是不是管理员错误")
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrFakeRequest.Code,
	//		"errMsg": errcode.ErrFakeRequest.String,
	//	})
	//	return
	//}
	//Query all device group list | 获取所有设备的列表
	deviceGroupList, err := dao.QueryDeviceGroupList()
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":       errcode.ErrNoError.Code,
		"errMsg":    errcode.ErrNoError.String,
		"groupList": deviceGroupList,
	})
	return
}

//---------------------------------------------------------------------------------------------------------------------------
//用户创建设备组
func HandleCreateCamerasGroupPost(c *gin.Context) {
	log.HTTP.Info("HandleCreateCamerasGroupPost BEGIN")
	var reqJSON reqCamerasGroupJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Create device group
	groupID, err := dao.CreateDeviceGroup(reqJSON.GroupName)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCreate.Code,
			"errMsg": errcode.ErrCreate.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":     errcode.ErrNoError.Code,
		"errMsg":  errcode.ErrNoError.String,
		"groupID": groupID,
	})
	return
}

type reqCamerasGroupJSON struct {
	GroupName string `form:"groupName" json:"groupName"`
}

//----------------------------------------------------------------------------------------------------------------------------
//用户查询设备组详情
//URL: "/v1/api/groups/:gid"
func HandleQueryDeviceGroupDetailGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryDeviceGroupDetailGet BEGIN")
	//Get group id  from url
	groupID, err := strconv.ParseInt(c.Param("gid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	group, err := dao.QueryDeviceGroup(groupID)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
		"group":  group,
	})
	return
}

//-----------------------------------------------------------------------------------------------------------------------------
//用户修改设备组名
func HandleUpdateDeviceGroupInfoPut(c *gin.Context) {
	log.HTTP.Info("HandleUpdateDeviceGroupInfoPut BEGIN")
	//Get group id  from url
	groupID, err := strconv.ParseInt(c.Param("gid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}
	//Get JSON data from body
	var reqJSON reqUpdateDeviceGroupInfoPutJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//校验是不是管理员
	//userID, err := lib.GetCurrentUser(c.Request)
	//if err != nil {
	//	log.HTTP.Error(" 判断是不是管理员")
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrCookie.Code,
	//		"errMsg": errcode.ErrCookie.String,
	//	})
	//	return
	//}
	//
	//if false == dao.IsAdmin(userID) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrFakeRequest.Code,
	//		"errMsg": errcode.ErrFakeRequest.String,
	//	})
	//	return
	//}

	err = dao.UpdateDeviceGroupName(groupID, reqJSON.GroupName)
	if err != nil {
		log.HTTP.Error(err)
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

type reqUpdateDeviceGroupInfoPutJSON struct {
	GroupName string `form:"groupName" json:"groupName"`
}

//-----------------------------------------------------------------------------------------------------------------------------
//用户删除相机组
// HandleDeleteDeviceGroupDelete : Delete device group
func HandleDeleteDeviceGroupDelete(c *gin.Context) {
	log.HTTP.Info("HandleDeleteDeviceGroupDelete BEGIN")
	//Get group id  from url
	groupID, err := strconv.ParseInt(c.Param("gid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}
	//是不是管理员
	//userID, err := lib.GetCurrentUser(c.Request)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrCookie.Code,
	//		"errMsg": errcode.ErrCookie.String,
	//	})
	//	return
	//}
	//
	//if false == dao.IsAdmin(userID) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrFakeRequest.Code,
	//		"errMsg": errcode.ErrFakeRequest.String,
	//	})
	//	return
	//}
	//Delete device group
	err = dao.DeleteDeviceGroup(groupID)
	if err != nil {
		log.HTTP.Error(err)
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

//-----------------------------------------------------------------------------------------------------------------------------
//查询所有未分组的设备
//URL: "/v1/api/nongroup/devices"
func HandleQueryNonGroupDeviceListGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryNonGroupDeviceListGet BEGIN")
	//Get device list Not at group from mysql
	deviceList, err := dao.QueryNonGroupDeviceList()
	if err != nil {
		log.HTTP.Error(err)
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

//----------------------------------------------------------------------------------------------------------------------------
//URL: "/v1/api/nongroup/page/devices"
//分页查询未分组的设备
func HandleQueryNonGroupDeviceListByPageGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryNonGroupDeviceListByPageGet BEGIN")

	keyword := c.Query("keyword")
	offsetStr := c.Query("offset")
	countStr := c.Query("count")

	offset, _ := strconv.ParseInt(offsetStr, 10, 64)
	count, _ := strconv.ParseInt(countStr, 10, 64)

	totalNum, err := dao.QueryNonGroupDeviceNumber(keyword)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	deviceList, err := dao.QueryNonGroupDeviceListByPage(keyword, offset, count)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":        errcode.ErrNoError.Code,
		"errMsg":     errcode.ErrNoError.String,
		"totalNum":   totalNum,
		"deviceList": deviceList,
	})
	return
}

//-----------------------------------------------------------------------------------------------------------------------------
//从设备组中删除一个设备
//URL: "/v1/api/groups/:gid/devices/:did"
func HandleDeleteDeviceFromGroupDelete(c *gin.Context) {
	log.HTTP.Info("HandleDeleteDeviceFromGroupDelete BEGIN")
	//Get group id  from url
	groupID, err := strconv.ParseInt(c.Param("gid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	//Get group id  from url
	deviceID, err := strconv.ParseInt(c.Param("did"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	//userID, err := lib.GetCurrentUser(c.Request)
	//if err != nil {
	//	log.HTTP.Error("校验是不是管理员")
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrCookie.Code,
	//		"errMsg": errcode.ErrCookie.String,
	//	})
	//	return
	//}
	//
	//if false == dao.IsAdmin(userID) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrFakeRequest.Code,
	//		"errMsg": errcode.ErrFakeRequest.String,
	//	})
	//	return
	//}
	//Delete device group
	err = dao.DeleteDeviceFromGroup(groupID, deviceID)
	if err != nil {
		log.HTTP.Error(err)
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

//----------------------------------------------------------------------------------------------------------------------------
//向设备组中添加设备列表
//URL: "/v1/api/groups/:gid/devices"
func HandleAddGroupDevicesPost(c *gin.Context) {
	log.HTTP.Info("HandleAddGroupDevicesPost BEGIN")
	//Get group id  from url
	groupID, err := strconv.ParseInt(c.Param("gid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	//Get JSON data from body
	var reqJSON reqAddGroupDevicesJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//userID, err := lib.GetCurrentUser(c.Request)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrCookie.Code,
	//		"errMsg": errcode.ErrCookie.String,
	//	})
	//	return
	//}
	//
	//if false == dao.IsAdmin(userID) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrFakeRequest.Code,
	//		"errMsg": errcode.ErrFakeRequest.String,
	//	})
	//	return
	//}

	//Add devices to group
	err = dao.AddGroupDevices(groupID, reqJSON.DeviceList)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCreate.Code,
			"errMsg": errcode.ErrCreate.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqAddGroupDevicesJSON struct {
	DeviceList []int64 `form:"deviceIDList" json:"deviceIDList"`
}

//------------------------------------------------------------------------------------------------------------------------------------
//更新设备组设备列表
func HandleUpdateGroupDevicesPut(c *gin.Context) {
	log.HTTP.Info("HandleUpdateGroupDevicesPut BEGIN")
	//Get group id  from url
	groupID, err := strconv.ParseInt(c.Param("gid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	//Get JSON data from body
	var reqJSON reqUpdateGroupDevicesJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//userID, err := lib.GetCurrentUser(c.Request)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrCookie.Code,
	//		"errMsg": errcode.ErrCookie.String,
	//	})
	//	return
	//}
	//
	//if false == dao.IsAdmin(userID) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrFakeRequest.Code,
	//		"errMsg": errcode.ErrFakeRequest.String,
	//	})
	//	return
	//}

	//Delete devices from group
	err = dao.UpdateGroupDevices(groupID, reqJSON.AddDeviceList, reqJSON.DeleteDeviceList)
	if err != nil {
		log.HTTP.Error(err)
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

type reqUpdateGroupDevicesJSON struct {
	AddDeviceList    []int64 `form:"addDeviceList" json:"addDeviceList"`
	DeleteDeviceList []int64 `form:"deleteDeviceList" json:"deleteDeviceList"`
}

//------------------------------------------------------------------------------------------------------------------------------------
//从设备组中删除设备列表
func HandleDeleteGroupDevicesDelete(c *gin.Context) {
	log.HTTP.Info("HandleDeleteGroupDevicesDelete BEGIN")
	//Get group id  from url
	groupID, err := strconv.ParseInt(c.Param("gid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	//Get JSON data from body
	var reqJSON reqDeleteGroupDevicesJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//userID, err := lib.GetCurrentUser(c.Request)
	//if err != nil {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrCookie.Code,
	//		"errMsg": errcode.ErrCookie.String,
	//	})
	//	return
	//}
	//
	//if false == dao.IsAdmin(userID) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"err":    errcode.ErrFakeRequest.Code,
	//		"errMsg": errcode.ErrFakeRequest.String,
	//	})
	//	return
	//}

	//Delete devices to group
	err = dao.DeleteGroupDevices(groupID, reqJSON.DeviceIDList)
	if err != nil {
		log.HTTP.Error(err)
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

type reqDeleteGroupDevicesJSON struct {
	DeviceIDList []int64 `form:"deviceIDList" json:"deviceIDList"`
}
