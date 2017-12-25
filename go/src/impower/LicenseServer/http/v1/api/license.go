package api

import (
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"impower/LicenseServer/dao"
	"impower/LicenseServer/env"
	"impower/LicenseServer/http/errcode"
	"impower/LicenseServer/lib"
	"impower/LicenseServer/log"
)

// HandleQueryAuthRecordByPageGet : Query authentication record by page
// url : /v1/api/mac/record
// params: {"orderNumber":string,"deviceSupplier":string,"beginTime":int64,"endTime":int64,"offset":int,"count":int}
func HandleQueryAuthRecordByPageGet(c *gin.Context) {
	log.Root.Info("HandleQueryAuthRecordByPageGet BEGIN")

	//Get params from url
	orderNumber := c.Query("orderNumber")
	deviceSupplier := c.Query("deviceSupplier")

	beginTimeStamp, _ := strconv.ParseInt(c.Query("beginTime"), 10, 64)
	endTimeStamp, _ := strconv.ParseInt(c.Query("endTime"), 10, 64)
	tmBegin := time.Unix(beginTimeStamp/1000, 0)
	tmEnd := time.Unix(endTimeStamp/1000, 0)
	tmBegin = lib.LocalTime(tmBegin)
	tmEnd = lib.LocalTime(tmEnd)

	offset, _ := strconv.Atoi(c.Query("offset"))
	count, _ := strconv.Atoi(c.Query("count"))

	//Get current login user id from cookie
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user id failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Query auth record total number
	totalNum, err := dao.QueryAuthRecordTotalNumber(accountID, orderNumber, deviceSupplier, tmBegin, tmEnd)
	if err != nil {
		log.Root.Error("Query authentication record total number failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	//Query auth record by page
	recordList, err := dao.QueryAuthRecordByPage(accountID, orderNumber, deviceSupplier, tmBegin, tmEnd, offset, count)
	if err != nil {
		log.Root.Error("Query authentication record by page failed.")
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
		"recordList": recordList,
	})
	return
}

type reqAuthRecordPostJSON struct {
	DeviceSupplier string   `form:"deviceSupplier" json:"deviceSupplier"`
	DeviceType     string   `form:"deviceType" json:"deviceType"`
	DeviceModel    string   `form:"deviceModel" json:"deviceModel"`
	DeviceNumber   int      `form:"deviceNumber" json:"deviceNumber"`
	OrderNumber    string   `form:"orderNumber" json:"orderNumber"`
	MacList        []string `form:"macList" json:"macList"`
	SupportP2P     int      `form:"supportP2P" json:"supportP2P"`
	PhoneNo        string   `form:"phoneNo" json:"phoneNo"`
	VerifyCode     string   `form:"verifyCode" json:"verifyCode"`
}

// HandleAuthRecordPost : Record authentication information
// url : /v1/api/mac/record
func HandleAuthRecordPost(c *gin.Context) {
	log.Root.Info("HandleAuthRecordPost BEGIN")

	//Get JSON data from body
	var reqJSON reqAuthRecordPostJSON
	if c.BindJSON(&reqJSON) != nil {
		log.Root.Error("Parse JSON error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Get current login user id from cookie
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user id failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Verify phone code
	err = dao.VerifyPhoneCode(reqJSON.PhoneNo, reqJSON.VerifyCode)
	if err != nil {
		log.Root.Error("Verify phone code error. PhoneNo: %v, Code: %v", reqJSON.PhoneNo, reqJSON.VerifyCode)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrVerifyPhoneCode.Code,
			"errMsg": errcode.ErrVerifyPhoneCode.String,
		})
		return
	}

	//Generate record batch
	batch := generateRecordBatch()
	log.Root.Info("Auth record batch: %v, device supplier: %v, order number: %v, device number: %v", batch, reqJSON.DeviceSupplier, reqJSON.OrderNumber, reqJSON.DeviceNumber)

	record := lib.AuthRecord{
		RecordBatch:    batch,
		DeviceSupplier: reqJSON.DeviceSupplier,
		DeviceType:     reqJSON.DeviceType,
		DeviceModel:    reqJSON.DeviceModel,
		DeviceNumber:   reqJSON.DeviceNumber,
		OrderNumber:    reqJSON.OrderNumber,
		SupportP2P:     reqJSON.SupportP2P,
		RecordUser:     "",
		RecordTime:     "",
	}

	//Record authentication
	err = dao.RecordAuthBatch(accountID, record, reqJSON.MacList)
	if err != nil {
		log.Root.Error("Record authentication failed. Record batch: %v", batch)
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

func generateRecordBatch() string {
	now := time.Now()
	batch := ""
	batch = batch + strconv.Itoa(now.Year())
	if now.Month() < 10 {
		batch = batch + "0" + strconv.Itoa(int(now.Month()))
	} else {
		batch = batch + strconv.Itoa(int(now.Month()))
	}
	if now.Day() < 10 {
		batch = batch + "0" + strconv.Itoa(now.Day())
	} else {
		batch = batch + strconv.Itoa(now.Day())
	}
	if now.Hour() < 10 {
		batch = batch + "0" + strconv.Itoa(now.Hour())
	} else {
		batch = batch + strconv.Itoa(now.Hour())
	}
	if now.Minute() < 10 {
		batch = batch + "0" + strconv.Itoa(now.Minute())
	} else {
		batch = batch + strconv.Itoa(now.Minute())
	}
	if now.Second() < 10 {
		batch = batch + "0" + strconv.Itoa(now.Second())
	} else {
		batch = batch + strconv.Itoa(now.Second())
	}
	millsec := now.Nanosecond() / 1000000
	if millsec < 10 {
		batch = batch + "00" + strconv.Itoa(millsec)
	} else if millsec < 100 {
		batch = batch + "0" + strconv.Itoa(millsec)
	} else {
		batch = batch + strconv.Itoa(millsec)
	}

	return batch
}

// HandleQueryAuthMacByBatchGet : Query authenticated mac list by batch
// url : /v1/api/mac/record/batch
// params: {"batch":string}
func HandleQueryAuthMacByBatchGet(c *gin.Context) {
	log.Root.Info("HandleQueryAuthMacByBatchGet BEGIN")

	//Get params from url
	batch := c.Query("batch")

	//Get current login user id from cookie
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user id failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Query authentication mac address by batch
	macList, err := dao.QueryAuthMacByBatch(accountID, batch)
	if err != nil {
		log.Root.Error("Query authentication mac address by batch failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":     errcode.ErrNoError.Code,
		"errMsg":  errcode.ErrNoError.String,
		"macList": macList,
	})
	return
}

// HandleApplyLicenseGet : Apply license
// url : /v1/api/license
func HandleApplyLicenseGet(c *gin.Context) {
	log.Root.Info("HandleApplyLicenseGet BEGIN")

	macBase64 := c.Query("mac")
	mac, err := base64.StdEncoding.DecodeString(macBase64)
	if err != nil {
		log.Root.Error("Decode mac address base64 failed. MAC address: %v", macBase64)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	//Authorise mac address
	isValid, err := dao.AuthMacAddress(string(mac))
	if err != nil {
		log.Root.Error("Authorise mac address failed. MAC address: %v", string(mac))
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	if false == isValid {
		log.Root.Error("Authenticated unpass. MAC address: %v", string(mac))
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrPermissionDenied.Code,
			"errMsg": errcode.ErrPermissionDenied.String,
		})
		return
	}

	//Encrypt mac
	macEncrypted, err := lib.EncryptByAES([]byte(mac), []byte(env.AESSecretKey), []byte(env.AESSecretIV))
	if err != nil {
		log.Root.Error("Encrypt mac failed. Mac address: %v", mac)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrServer.Code,
			"errMsg": errcode.ErrServer.String,
		})
		return
	}

	macEncryptedBase64 := base64.StdEncoding.EncodeToString(macEncrypted)

	c.JSON(http.StatusOK, gin.H{
		"err":     errcode.ErrNoError.Code,
		"errMsg":  errcode.ErrNoError.String,
		"license": macEncryptedBase64,
	})
	return
}

// HandleQueryP2PRecordByPageGet : Query p2p record by page
// url : /v1/api/p2p/record
// params: {"orderNumber":string,"beginTime":int64,"endTime":int64,"offset":int,"count":int}
func HandleQueryP2PRecordByPageGet(c *gin.Context) {
	log.Root.Info("HandleQueryP2PRecordByPageGet BEGIN")

	//Get params from url
	orderNumber := c.Query("orderNumber")

	beginTimeStamp, _ := strconv.ParseInt(c.Query("beginTime"), 10, 64)
	endTimeStamp, _ := strconv.ParseInt(c.Query("endTime"), 10, 64)
	tmBegin := time.Unix(beginTimeStamp/1000, 0)
	tmEnd := time.Unix(endTimeStamp/1000, 0)
	tmBegin = lib.LocalTime(tmBegin)
	tmEnd = lib.LocalTime(tmEnd)

	offset, _ := strconv.Atoi(c.Query("offset"))
	count, _ := strconv.Atoi(c.Query("count"))

	//Get current login user id from cookie
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user id failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Query p2p record total number
	totalNum, err := dao.QueryP2PRecordTotalNumber(accountID, orderNumber, tmBegin, tmEnd)
	if err != nil {
		log.Root.Error("Query authentication record total number failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	//Query p2p record by page
	recordList, err := dao.QueryP2PRecordByPage(accountID, orderNumber, tmBegin, tmEnd, offset, count)
	if err != nil {
		log.Root.Error("Query authentication record by page failed.")
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
		"recordList": recordList,
	})
	return
}

type reqP2PRecordPostJSON struct {
	OrderNumber string   `form:"orderNumber" json:"orderNumber"`
	P2PNumber   int      `form:"p2pNumber" json:"p2pNumber"`
	P2PList     []string `form:"p2pList" json:"p2pList"`
	PhoneNo     string   `form:"phoneNo" json:"phoneNo"`
	VerifyCode  string   `form:"verifyCode" json:"verifyCode"`
}

// HandleP2PRecordPost : Record p2p information
// url : /v1/api/p2p/record
func HandleP2PRecordPost(c *gin.Context) {
	log.Root.Info("HandleP2PRecordPost BEGIN")

	//Get JSON data from body
	var reqJSON reqP2PRecordPostJSON
	if c.BindJSON(&reqJSON) != nil {
		log.Root.Error("Parse JSON error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Get current login user id from cookie
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user id failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Verify phone code
	err = dao.VerifyPhoneCode(reqJSON.PhoneNo, reqJSON.VerifyCode)
	if err != nil {
		log.Root.Error("Verify phone code error. PhoneNo: %v, Code: %v", reqJSON.PhoneNo, reqJSON.VerifyCode)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrVerifyPhoneCode.Code,
			"errMsg": errcode.ErrVerifyPhoneCode.String,
		})
		return
	}

	//Generate record batch
	batch := generateRecordBatch()
	log.Root.Info("P2P record batch: %v, order number: %v, p2p number: %v", batch, reqJSON.OrderNumber, reqJSON.P2PNumber)

	record := lib.P2PRecord{
		RecordBatch: batch,
		P2PNumber:   reqJSON.P2PNumber,
		OrderNumber: reqJSON.OrderNumber,
		RecordUser:  "",
		RecordTime:  "",
	}

	//Record p2p
	err = dao.RecordP2PBatch(accountID, record, reqJSON.P2PList)
	if err != nil {
		log.Root.Error("Record p2p failed. Record batch: %v", batch)
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

// HandleQueryP2PByBatchGet : Query p2p list by batch
// url : /v1/api/p2p/record/batch
// params: {"batch":string}
func HandleQueryP2PByBatchGet(c *gin.Context) {
	log.Root.Info("HandleQueryP2PByBatchGet BEGIN")

	//Get params from url
	batch := c.Query("batch")

	//Get current login user id from cookie
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user id failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Query p2p list by batch
	p2pList, err := dao.QueryP2PByBatch(accountID, batch)
	if err != nil {
		log.Root.Error("Query p2p list by batch failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":     errcode.ErrNoError.Code,
		"errMsg":  errcode.ErrNoError.String,
		"p2pList": p2pList,
	})
	return
}

// HandleApplyP2PGet : Apply p2p
// url : /v1/api/p2p
func HandleApplyP2PGet(c *gin.Context) {
	log.Root.Info("HandleApplyP2PGet BEGIN")

	macBase64 := c.Query("mac")
	mac, err := base64.StdEncoding.DecodeString(macBase64)
	if err != nil {
		log.Root.Error("Decode mac address base64 failed. MAC address: %v", macBase64)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	//Authorise mac address is support p2p
	isValid, err := dao.AuthMacSupportP2P(string(mac))
	if err != nil {
		log.Root.Error("Authorise mac address failed. MAC address: %v", string(mac))
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	if false == isValid {
		log.Root.Error("Authenticated unpass. MAC address: %v", string(mac))
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrPermissionDenied.Code,
			"errMsg": errcode.ErrPermissionDenied.String,
		})
		return
	}

	//Get p2p id by mac address
	p2pID, err := dao.GetP2P(string(mac))
	if err != nil {
		log.Root.Error("Get p2p id failed. MAC address: %v", string(mac))
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	//Encrypt p2p id
	p2pEncrypted, err := lib.EncryptByAES([]byte(p2pID), []byte(env.AESSecretKey), []byte(env.AESSecretIV))
	if err != nil {
		log.Root.Error("Encrypt p2p id failed. P2P ID: %v", p2pID)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrServer.Code,
			"errMsg": errcode.ErrServer.String,
		})
		return
	}

	p2pEncryptedBase64 := base64.StdEncoding.EncodeToString(p2pEncrypted)

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
		"p2pID":  p2pEncryptedBase64,
	})
	return
}
