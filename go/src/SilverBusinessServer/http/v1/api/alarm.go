package api

import (
	"SilverBusinessServer/dao"
	"SilverBusinessServer/http/errcode"
	"SilverBusinessServer/log"
	"SilverBusinessServer/lib"
	"SilverBusinessServer/echo"
	"net/http"
	"strconv"
	b64 "encoding/base64"
	"github.com/gin-gonic/gin"
)

//-------------------------------------------------------------------------------------------------------------------------------------
// 添加报警备注，通过相机ID和报警时间戳标识报警    通过id和时间戳去数据库中找报警信息，然后添加报警备注
func HandleAddAlarmRemarkPost(c *gin.Context) {
	log.HTTP.Info("HandleAddAlarmRemarkPost BEGIN")
	var reqJSON reqAddAlarmRemarkJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if err := dao.AddAlarmRemark(reqJSON.DeviceID, reqJSON.AlarmTime, reqJSON.Remark); err != nil {
		log.HTTP.Error(err)
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

type reqAddAlarmRemarkJSON struct {
	DeviceID  int64  `form:"deviceID" json:"deviceID"`
	AlarmTime int64  `form:"alarmTime" json:"alarmTime"`
	Remark    string `form:"remark" json:"remark"`
}

//------------------------------------------------------------------------------------------------------------------------------------
//实时报警：处理或归档报警，通过相机ID和报警时间戳标识报警
func HandlePigeonholeAlarmPut(c *gin.Context) {
	log.HTTP.Info("HandlePigeonholeAlarmPut BEGIN")
	var reqJSON reqPigeonholeAlarmJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if err := dao.PigeonholeAlarm(reqJSON.DeviceID, reqJSON.AlarmTime, reqJSON.ProcessStatus, reqJSON.ArchiveFlag); err != nil {
		log.HTTP.Error(err)
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
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if err := dao.DeleteAlarm(reqJSON.DeviceID, reqJSON.AlarmTime); err != nil {
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

type reqDeleteAlarmJSON struct {
	DeviceID  int64 `form:"deviceID" json:"deviceID"`
	AlarmTime int64 `form:"alarmTime" json:"alarmTime"`
}

//------------------------------------------------------------------------------------------------------------------------------------------------
//分页查询报警，查询参数中archiveFlag=1表示只查询归档的报警，其他值表示查询所有报警
func HandleQueryAlarmsByPageGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryAlarmsByPageGet BEGIN")
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
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	totalNum, err := dao.QueryAlarmTotal(deviceID, beginTime, endTime, archiveFlag)
	if err != nil { //表示查询总数
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
//push alarm to specified platforms (client and database)
func HandlePushAlarmsPost(c *gin.Context) {
	log.HTTP.Info("HandlePushAlarmsPost BEGIN")
	var reqJSON reqAlarmInfoJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if reqJSON.DeviceType == 2 {
		if deviceID, err := dao.GetDeviceID(reqJSON.DeviceUUID); err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrPushAlarm.Code,
				"errMsg": errcode.ErrPushAlarm.String,
			})
			return
		} else {
			reqJSON.DeviceID = deviceID
		}
	}

	log.HTTP.Info("AlgID: ", reqJSON.AlgID)
	textData := lib.AlarmValue{DeviceID : reqJSON.DeviceID, AlgID : reqJSON.AlgID, AlarmTime : reqJSON.AlarmTime, AlarmInfo : reqJSON.AlarmInfo}
	binData, err := b64.StdEncoding.DecodeString(reqJSON.AlarmImage)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrPushAlarm.Code,
			"errMsg": errcode.ErrPushAlarm.String,
		})
		return
	}

	if err := echo.PushAlarm(&textData, binData); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrPushAlarm.Code,
			"errMsg": errcode.ErrPushAlarm.String,
		})
		return
	}

	if err := dao.SaveAlarm(&textData, binData); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrSaveAlarm.Code,
			"errMsg": errcode.ErrSaveAlarm.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":       errcode.ErrNoError.Code,
		"errMsg":    errcode.ErrNoError.String,
	})
	return
}

type reqAlarmInfoJSON struct {
	DeviceType  int    `form:"deviceType" json:"deviceType"`
	DeviceID    int64  `form:"deviceID" json:"deviceID"`
	DeviceUUID  string `form:"deviceUUID" json:"deviceUUID"`
	AlgID	    string `form:"algID" json:"algID"`
	AlarmTime   int64  `form:"alarmTime" json:"alarmTime"`
	AlarmInfo   string `form:"alarmInfo" json:"alarmInfo"`
	AlarmImage  string `form:"alarmImage" json:"alarmImage"`
}

//--------------------------------------------------------------------------------------------------------------------------------------------------------
//按时间段批量删除报警，归档报警不删除，只删除未归档报警
func HandleDeleteAlarmsByTimeDelete(c *gin.Context) {
	log.HTTP.Info("HandleDeleteAlarmsByTimeDelete BEGIN")
	var reqJSON reqDeleteAlarmsByTimeJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if err := dao.DeleteAlarmsByTime(reqJSON.BeginTime, reqJSON.EndTime); err != nil {
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

type reqDeleteAlarmsByTimeJSON struct {
	BeginTime int64 `form:"beginTime" json:"beginTime"`
	EndTime   int64 `form:"endTime" json:"endTime"`
}

//--------------------------------------------------------------------------------------------------------------------------------------------------------
//查询报警图片
func HandleQueryAlarmImageByidGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryAlarmImageByidGet BEGIN")
	//get image id from url address
	imageId, err := strconv.ParseInt(c.Param("iid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	image, err := dao.QueryAlarmImageInfo(imageId)
	if err != nil {
		log.HTTP.Error(err)
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
	log.HTTP.Info("HandleAddAlarmRemarkByOldPost BEGIN")
	//get image id from url address
	alarmid, err := strconv.ParseInt(c.Param("aid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	var reqJSON reqAlarmRemarkByOldJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if err := dao.AddAlarmRemarkByOld(alarmid, reqJSON.Remark); err != nil {
		log.HTTP.Error(err)
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
	log.HTTP.Info("HandlePigeonholeAlarmByOldPut BEGIN")
	//get image id from url address
	alarmid, err := strconv.ParseInt(c.Param("aid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	var reqJSON reqPigeonholeAlarmByOldJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	if err := dao.PigeonholeAlarmByOld(alarmid, reqJSON.ArchiveFlag, reqJSON.ProcessStatus); err != nil {
		log.HTTP.Error(err)
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
	log.HTTP.Info("HandleDeleteAlarmByOldDelete BEGIN")
	//get image id from url address
	alarmid, err := strconv.ParseInt(c.Param("aid"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	if err := dao.DeleteAlarmByAlarmId(alarmid); err != nil {
		log.HTTP.Error(err)
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
