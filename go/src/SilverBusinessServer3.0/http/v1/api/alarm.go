package api

import (
	"SilverBusinessServer/dao"
	"SilverBusinessServer/http/errcode"
	"SilverBusinessServer/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//-------------------------------------------------------------------------------------------------------------------------------------
// 添加报警备注，通过相机ID和报警时间戳标识报警    通过id和时间戳去数据库中找报警信息，然后添加报警备注
func HandleAddAlarmRemarkPost(c *gin.Context) {
	log.HTTP.Info("创建报警备注")
	var reqJSON reqAddAlarmRemarkJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error(" JSON Err|获取信息错误")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if dao.AddAlarmRemark(reqJSON.DeviceID, reqJSON.AlarmTime, reqJSON.Remark) != nil {
		log.HTTP.Info("添加报警失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrAddAlarmRemark.Code,
			"errMsg": errcode.ErrAddAlarmRemark.String,
		})
		return
	}
	//成功了
	log.HTTP.Info("实时报警添加备注成功")
	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqAddAlarmRemarkJSON struct {
	DeviceID  int64  `form:"deviceID" json:"deviceID"`
	AlarmTime int64  `form:"alarmTime" json:"alarmTime"`
	Remark    string `form:"remark" json:"remark"`
}

//------------------------------------------------------------------------------------------------------------------------------------
//实时报警：处理或归档报警，通过相机ID和报警时间戳标识报警
func HandlePigeonholeAlarmPut(c *gin.Context) {
	log.HTTP.Info("处理或者归档报警信息")
	var reqJSON reqPigeonholeAlarmJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("JSON Err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if dao.PigeonholeAlarm(reqJSON.DeviceID, reqJSON.AlarmTime, reqJSON.ProcessStatus, reqJSON.ArchiveFlag) != nil {
		log.HTTP.Error("归档处理失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrPrOrAr.Code,
			"errMsg": errcode.ErrPrOrAr.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqPigeonholeAlarmJSON struct {
	DeviceID      int64 `form:"deviceID" json:"deviceID"`
	AlarmTime     int64 `form:"alarmTime" json:"alarmTime"`
	ProcessStatus int   `form:"processStatus" json:"processStatus"`
	ArchiveFlag   int   `form:"archiveFlag" json:"archiveFlag"`
}

//-------------------------------------------------------------------------------------------------------------------------------------
//删除报警
func HandleDeleteAlarmDelete(c *gin.Context) {
	log.HTTP.Info("HandleDeleteAlarmDelete BEGIN")
	var reqJSON reqDeleteAlarmJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("JSON Err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if dao.DeleteAlarm(reqJSON.DeviceID, reqJSON.AlarmTime) != nil {
		log.HTTP.Error("Delete Err")
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

type reqDeleteAlarmJSON struct {
	DeviceID  int64 `form:"deviceID" json:"deviceID"`
	AlarmTime int64 `form:"alarmTime" json:"alarmTime"`
}

//------------------------------------------------------------------------------------------------------------------------------------------------
//分页查询报警，查询参数中archiveFlag=1表示只查询归档的报警，其他值表示查询所有报警
func HandleQueryAlarmsByPageGet(c *gin.Context) {
	log.HTTP.Info("分页查询报警")
	deviceIDStr := c.Query("deviceID")
	beginTimeStr := c.Query("beginTime")
	endTimeStr := c.Query("endTime")
	archiveFlagStr := c.Query("archiveFlag")
	offsetStr := c.Query("offset")
	countStr := c.Query("count")

	deviceID, _ := strconv.ParseInt(deviceIDStr, 10, 64)
	beginTime, _ := strconv.ParseInt(beginTimeStr, 10, 64)
	endTime, _ := strconv.ParseInt(endTimeStr, 10, 64)
	archiveFlag, _ := strconv.Atoi(archiveFlagStr)
	offset, _ := strconv.Atoi(offsetStr)
	count, _ := strconv.Atoi(countStr)

	AlarmList, err := dao.QueryAlarms(deviceID, beginTime, endTime, archiveFlag, offset, count)
	if err != nil { //表示只查询归档的报警
		log.HTTP.Error("归档错误")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	totalNum, err := dao.QueryAlarmTotal(deviceID, beginTime, endTime, archiveFlag)
	if err != nil { //表示查询总数
		log.HTTP.Error("查询总数错误")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":       errcode.ErrNoError.Code,
		"errMsg":    errcode.ErrNoError.String,
		"totalNum":  totalNum,
		"alarmList": AlarmList,
	})
	return
}

type reqQueryAlarmByPageJSON struct {
	DeviceID    int64 `form:"deviceID" json:"deviceID"`
	BeginTime   int64 `form:"beginTime" json:"beginTime"`
	EndTime     int64 `form:"endTime" json:"endTime"`
	ArchiveFlag int   `form:"archiveFlag" json:"archiveFlag"`
	Offset      int   `form:"offset" json:"offset"`
	Count       int   `form:"count" json:"count"`
}

//--------------------------------------------------------------------------------------------------------------------------------------------------------
//按时间段批量删除报警，归档报警不删除，只删除未归档报警
func HandleDeleteAlarmsByTimeDelete(c *gin.Context) {
	log.HTTP.Info("按照时间段来删除未归档的报警")
	var reqJSON reqDeleteAlarmsByTimeJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("JSON Err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if dao.DeleteAlarmsByTime(reqJSON.BeginTime, reqJSON.EndTime) != nil {
		log.HTTP.Error("按照时间段来删除报警错误")
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

type reqDeleteAlarmsByTimeJSON struct {
	BeginTime int64 `form:"beginTime" json:"beginTime"`
	EndTime   int64 `form:"endTime" json:"endTime"`
}

//--------------------------------------------------------------------------------------------------------------------------------------------------------
//查询报警图片
func HandleQueryAlarmImageByidGet(c *gin.Context) {
	//get image id from url address
	imageId, err := strconv.ParseInt(c.Param("iid"), 10, 64)
	if err != nil {
		log.HTTP.Error("get image id err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	image, err := dao.QueryAlarmImageInfo(imageId)
	if err != nil {
		log.HTTP.Error("根据照片id来获取照片信息失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrAlarmImage.Code,
			"errMsg": errcode.ErrAlarmImage.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
		"image":  image,
	})
	return
}

//-------------------------------------------------------------------------------------------------------------------------------------------------------
//历史报警信息添加备注
func HandleAddAlarmRemarkByOldPost(c *gin.Context) {
	log.HTTP.Info("历史报警信息，添加备注，进入接口")
	//get image id from url address
	alarmid, err := strconv.ParseInt(c.Param("aid"), 10, 64)
	if err != nil {
		log.HTTP.Error("get alarm id err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	log.HTTP.Info("获取参数。。")
	var reqJSON reqAlarmRemarkByOldJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("JSON Err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}
	log.HTTP.Info("开始sql操作")
	if dao.AddAlarmRemarkByOld(alarmid, reqJSON.Remark) != nil {
		log.HTTP.Error("历史报警添加备注错误")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrAddAlarmRemark.Code,
			"errMsg": errcode.ErrAddAlarmRemark.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqAlarmRemarkByOldJSON struct {
	Remark string `form:"remark" json:"remark"`
}

//----------------------------------------------------------------------------
//历史报警：处理或归档报警，通过相机ID和报警时间戳标识报警
func HandlePigeonholeAlarmByOldPut(c *gin.Context) {
	log.HTTP.Info("历史报警信息，处理或者归档报警")
	//get image id from url address
	alarmid, err := strconv.ParseInt(c.Param("aid"), 10, 64)
	if err != nil {
		log.HTTP.Error("get alarm id err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	var reqJSON reqPigeonholeAlarmByOldJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("JSON Err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if dao.PigeonholeAlarmByOld(alarmid, reqJSON.ArchiveFlag, reqJSON.ProcessStatus) != nil {
		log.HTTP.Error("归档处理失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrPrOrAr.Code,
			"errMsg": errcode.ErrPrOrAr.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqPigeonholeAlarmByOldJSON struct {
	ProcessStatus int `form:"processStatus" json:"processStatus"`
	ArchiveFlag   int `form:"archiveFlag" json:"archiveFlag"`
}

//------------------------------------------------------------------------------
//历史报警：删除报警，通过报警id去删除报警
func HandleDeleteAlarmByOldDelete(c *gin.Context) {
	log.HTTP.Info("历史报警信息，删除报警")
	//get image id from url address
	alarmid, err := strconv.ParseInt(c.Param("aid"), 10, 64)
	if err != nil {
		log.HTTP.Error("get alarm id err")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	if dao.DeleteAlarmByAlarmId(alarmid) != nil {
		log.HTTP.Error("删除失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrPrOrAr.Code,
			"errMsg": errcode.ErrPrOrAr.String,
		})
		return
	}

	//成功了
	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}
