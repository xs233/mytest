package api

import (
	"math/rand"
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

// reqLoginPostJSON :
type reqLoginPostJSON struct {
	//Username
	Username string `form:"userName" json:"userName"`
	//Password
	Password string `form:"userPassword" json:"userPassword"`
}

// HandleLoginPost : Handle user login
// url: /v1/ap1/sessions
func HandleLoginPost(c *gin.Context) {
	log.Root.Info("HandleLoginPost BEGIN")

	//Get JSON data from body
	var reqJSON reqLoginPostJSON
	if c.BindJSON(&reqJSON) != nil {
		log.Root.Error("Parse JSON error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Login
	userID, userType, err := dao.Login(reqJSON.Username, reqJSON.Password)
	if err != nil {
		log.Root.Error("User login failed")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrLogin.Code,
			"errMsg": errcode.ErrLogin.String,
		})
		return
	}

	if userID < 0 {
		log.Root.Error("User not pass authenticated")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrAuthenticate.Code,
			"errMsg": errcode.ErrAuthenticate.String,
		})
		return
	}

	//Record login
	err = dao.RecordLogin(userID, "login")
	if err != nil {
		log.Root.Warn("Record login failed. userID: %v", userID)
	}

	//Set secret cookie:userID
	cookieKey := "account"
	cookieValue := strconv.FormatInt(userID, 10)
	lib.SetSecretCookie(c.Writer, cookieKey, cookieValue)

	c.JSON(http.StatusOK, gin.H{
		"err":      errcode.ErrNoError.Code,
		"errMsg":   errcode.ErrNoError.String,
		"userID":   userID,
		"userType": userType,
	})
	return
}

type reqModifyPasswordPostJSON struct {
	UserName    string `form:"userName" json:"userName"`
	OldPassword string `form:"oldPassword" json:"oldPassword"`
	NewPassword string `form:"newPassword" json:"newPassword"`
}

// HandleModifyPasswordPut :
// url : /v1/ap1/password
func HandleModifyPasswordPut(c *gin.Context) {
	log.Root.Info("HandleModifyPasswordPut BEGIN")

	//Get JSON data from body
	var reqJSON reqModifyPasswordPostJSON
	if c.BindJSON(&reqJSON) != nil {
		log.Root.Error("Parse JSON error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Modify password
	if dao.ModifyPassword(reqJSON.UserName, reqJSON.OldPassword, reqJSON.NewPassword) != nil {
		log.Root.Error("Modify password error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrModifyPassword.Code,
			"errMsg": errcode.ErrModifyPassword.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqResetPasswordPutJSON struct {
	UserID      int64  `form:"userID" json:"userID"`
	NewPassword string `form:"newPassword" json:"newPassword"`
}

// HandleResetPasswordPut :
// url : /v1/ap1/password/reset
func HandleResetPasswordPut(c *gin.Context) {
	log.Root.Info("HandleResetPasswordPut BEGIN")

	//Get JSON data from body
	var reqJSON reqResetPasswordPutJSON
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

	if false == dao.IsAdmin(accountID) {
		//Judge user whether valid or not
		log.Root.Error("User is not admin. Fake request.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	//Reset password
	if dao.ResetPassword(reqJSON.UserID, reqJSON.NewPassword) != nil {
		log.Root.Error("Reset password error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrResetPassword.Code,
			"errMsg": errcode.ErrResetPassword.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

// HandleLogoutDelete : User logout
// url: /v1/ap1/session/:sid
func HandleLogoutDelete(c *gin.Context) {
	log.Root.Info("HandleLogoutDelete BEGIN")

	//Get session id (user id) from url
	userID, err := strconv.ParseInt(c.Param("sid"), 10, 64)
	if err != nil {
		log.Root.Error("No session ID in URL or Session ID error.")
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

	//Judge user whether valid or not
	if userID != accountID {
		log.Root.Error("User IDs are not equal from URL and Cookie. Fake request.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	//Record logout
	err = dao.RecordLogin(userID, "logout")
	if err != nil {
		log.Root.Warn("Record login failed. userID: %v", userID)
	}

	//Delete secret cookie:accountID
	cookie := http.Cookie{
		Name:   "account",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

// HandleQueryUserInfoGet : Query user info
// url : /api/users/:uid
func HandleQueryUserInfoGet(c *gin.Context) {
	log.Root.Info("HandleQueryUserInfoGet BEGIN")

	//Get user id  from url
	userID, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		log.Root.Error("No User ID in URL or User ID error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
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

	//Judge account is admin or not, just admin has right to register device
	if false == dao.IsAdmin(accountID) {
		//Judge user whether valid or not
		if userID != accountID {
			log.Root.Error("User IDs are not equal from URL and Cookie. Fake request.")
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrFakeRequest.Code,
				"errMsg": errcode.ErrFakeRequest.String,
			})
			return
		}
	}

	//Query user info by user ID
	userInfo, err := dao.QueryUserInfo(userID)
	if err != nil {
		log.Root.Error("Query user info failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
		"user":   userInfo,
	})
	return
}

// HandleUserDelete : Delete user
// url : /v1/api/users/:uid
func HandleUserDelete(c *gin.Context) {
	log.Root.Info("HandleUserDelete BEGIN")

	willDeluserID, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		log.Root.Error("No session ID in URL or Session ID error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Get current login user id from cookie
	currentUserID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	if false == dao.IsAdmin(currentUserID) {
		//Judge user whether valid or not
		log.Root.Error("User is not admin. Fake request.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	err = dao.DeleteAccount(willDeluserID)
	if err != nil {
		log.Root.Warn("Delete User failed. will delete userID: %v", willDeluserID)
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})

}

// HandleQueryAllUserInfoGet : Query all user info
// url : /v1/api/users
func HandleQueryAllUserInfoGet(c *gin.Context) {
	log.Root.Info("HandleQueryAllUserInfoGet BEGIN")

	userID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user id failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}
	//Judge account is admin or not
	if false == dao.IsAdmin(userID) {
		//Judge user whether valid or not
		log.Root.Error("User ID is not admin . Fake request.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	userList, err := dao.QueryAllUserInfo()
	if err != nil {
		log.Root.Error("Query all user info failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":      errcode.ErrNoError.Code,
		"errMsg":   errcode.ErrNoError.String,
		"userList": userList,
	})
	return
}

// reqCreateUserPostJSON :
type reqCreateUserPostJSON struct {
	UserName     string `form:"userName" json:"userName"`
	Password     string `form:"userPassword" json:"userPassword"`
	UserPhoneNo  string `form:"userPhoneNo" json:"userPhoneNo"`
	AdminPhoneNo string `form:"adminPhoneNo" json:"adminPhoneNo"`
	VerifyCode   string `form:"verifyCode" json:"verifyCode"`
}

// HandleCreateUserPost : regist user
// url: /v1/api/users
func HandleCreateUserPost(c *gin.Context) {
	log.Root.Info("HandleCreateUserPost BEGIN")

	//Get JSON data from body
	var reqJSON reqCreateUserPostJSON
	if c.BindJSON(&reqJSON) != nil {
		log.Root.Error("Parse JSON error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})

		return
	}

	//Get current login user id from cookie
	currentUserID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.Root.Error("Get current login user failed.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	if false == dao.IsAdmin(currentUserID) {
		//Judge user whether valid or not
		log.Root.Error("User is not admin. Fake request.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	//Verify phone code
	err = dao.VerifyPhoneCode(reqJSON.AdminPhoneNo, reqJSON.VerifyCode)
	if err != nil {
		log.Root.Error("Verify phone code error. PhoneNo: %v, Code: %v", reqJSON.AdminPhoneNo, reqJSON.VerifyCode)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrVerifyPhoneCode.Code,
			"errMsg": errcode.ErrVerifyPhoneCode.String,
		})
		return
	}

	//Register
	accountID, err := dao.Register(reqJSON.UserName, reqJSON.Password, reqJSON.UserPhoneNo)
	if err != nil || accountID < 0 {
		log.Root.Error("Register user failed")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrRegister.Code,
			"errMsg": errcode.ErrRegister.String,
		})
		return
	}

	//Account is existing
	if accountID == 0 {
		log.Root.Error("User register failed, account is existing")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrAccountExisting.Code,
			"errMsg": errcode.ErrAccountExisting.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

// HandleGeneratePhoneCodePost : Generate phone code
// url : /v1/ap1/phone/:phone/codes
func HandleGeneratePhoneCodePost(c *gin.Context) {
	log.Root.Info("HandleGeneratePhoneCodePost BEGIN")

	//Get session id (user id) from url
	phone := c.Param("phone")
	_, err := strconv.ParseInt(phone, 0, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrPhoneNoFormat.Code,
			"errMsg": errcode.ErrPhoneNoFormat.String,
		})
		return
	}

	// Generate phone code
	rand.Seed(time.Now().Unix())
	code := rand.Intn(899999) + 100000
	liveTime := env.Get("phoneverify.livetime").(int64)
	deadline := lib.Now().Add(time.Second * time.Duration(liveTime))

	err = dao.SavePhoneCode(phone, strconv.Itoa(code), deadline)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrVerifyPhoneCode.Code,
			"errMsg": errcode.ErrVerifyPhoneCode.String,
		})
		return
	}

	// Send phone code
	err = lib.SendShortMessage(phone, int(code))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrSendShortMessage.Code,
			"errMsg": errcode.ErrSendShortMessage.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}
