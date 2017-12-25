package api

import (
	"SilverBusinessServer/dao"
	"SilverBusinessServer/echo"
	"SilverBusinessServer/http/errcode"
	"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//user to login | 用户登录，向服务端会话池中增加会话
//URL:/v1/api/sessions		--------------------------------------------------------------------------------------------------
func HandleLoginPost(c *gin.Context) {
	log.HTTP.Info("HandleLoginPost BEGIN")

	//Get usr login name and possword from body
	var reqJSON reqLoginPostJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("Parse JSON error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	log.HTTP.Info("", reqJSON.Username, reqJSON.Password)
	//Login
	userID, userType, rightList, err := dao.Login(reqJSON.Username, reqJSON.Password)

	if err != nil {
		log.HTTP.Error("login error.")
		c.JSON(http.StatusOK, gin.H{ // login error
			"err":    errcode.ErrLogin.Code,
			"errMsg": errcode.ErrLogin.String,
		})
		return
	}

	log.HTTP.Info("", userID, userType)
	if userID < 0 {
		log.HTTP.Error("Username or password error.")
		c.JSON(http.StatusOK, gin.H{ // 鉴权未通过
			"err":    errcode.ErrAuthenticate.Code,
			"errMsg": errcode.ErrAuthenticate.String,
		})
		return
	}

	//record login 登录记录
	err = dao.RecordLogin(userID, "login")
	if err != nil {
	}

	// set secret cookie:userID
	cookieKey := "account"
	cookieValue := strconv.FormatInt(userID, 10) //将用户ID转换成十进制的int型
	lib.SetSecretCookie(c.Writer, cookieKey, cookieValue)

	c.JSON(http.StatusOK, gin.H{
		"err":       errcode.ErrNoError.Code,
		"errMsg":    errcode.ErrNoError.String,
		"userID":    userID,
		"userType":  userType,
		"rightList": rightList,
	})
	return
}

type reqLoginPostJSON struct {
	Username string `form:"userName" json:"userName"`
	Password string `form:"userPassword" json:"userPassword"`
}

//----------------------------------------------------------------------------------------------------------------------
// #define CMD_MODIFY_PASSWORD	  修改密码  URL: /v1/api/password/modify
func HandleModifyPasswordPut(c *gin.Context) {
	log.HTTP.Info("HandleModifyPasswordPut BEGIN.")
	// Get JSON data from body
	var reqJSON reqModifyPasswordPutJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("Form error.")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Modify password
	if dao.ModifyPassword(reqJSON.UserName, reqJSON.OldPassword, reqJSON.NewPassword) != nil {
		log.HTTP.Error("Modify password error.")
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

type reqModifyPasswordPutJSON struct {
	UserName    string `form:"userName" json:"userName"`
	OldPassword string `form:"oldPassword" json:"oldPassword"`
	NewPassword string `form:"newPassword" json:"newPassword"`
}

//------------------------------------------------------------------------------------------------------------------------
//#define CMD_RESET_PASSWORD	//重置密码 URL: "/v1/api/password/reset"
//注：PUT 管理员重置用户密码为"111111"，须校验是否为管理员，否则为伪造请求 【其实就是修改密码，只是要找到这个用户名的，然后将他的密码修改为111111
func HandleResetPasswordPut(c *gin.Context) {
	log.HTTP.Info("HandleResetPasswordPut BEGIN.")
	/*
		校验是否是管理员，否则为伪造请求
		逻辑：userID取到后取数据库中查查，查他的usertype 是不是1，不是的直接返回是伪造请求
	*/
	//Get current login user id from cookie
	currentUserID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.HTTP.Error("Get cookie failed|获取用户的userid失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	if false == dao.IsAdmin(currentUserID) {
		log.HTTP.Error("Fake request")
		//Judge user whether valid or not
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	var reqJSON reqResetPasswordPutJSON
	if c.BindJSON(&reqJSON) != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//reset password
	if dao.ResetPassword(reqJSON.UserName, reqJSON.NewPassword) != nil {
		log.HTTP.Error("Reset password error")
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

type reqResetPasswordPutJSON struct { //新密码为“111111”，需要加密
	UserName    string `form:"userName" json:"userName"`
	NewPassword string `form:"newPassword" json:"newPassword"`
}

//----------------------------------------------------------------------------------------------------------------------------
// #define CMD_RESET_RIGHT	//修改用户权限    URL: "/v1/api/right/modify"
// HandleRightPut
func HandleRightPut(c *gin.Context) {
	log.HTTP.Info("HandleRightPut BEGIN")
	var reqJSON reqResetRightPutJSON
	if c.BindJSON(&reqJSON) != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Get current login user id from cookie
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.HTTP.Error("Get cookie failed")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Judge account is admin or not, just admin has right to register device
	if false == dao.IsAdmin(accountID) { //表示不是管理员
		//Judge user whether valid or not
		if reqJSON.UserID != accountID {
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrFakeRequest.Code,
				"errMsg": errcode.ErrFakeRequest.String,
			})
			return
		}
	}

	if err := dao.ResetRight(reqJSON.UserID, reqJSON.RightList); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrResetRight.Code,
			"errMsg": errcode.ErrResetRight.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
}

type reqResetRightPutJSON struct {
	UserID 	  int64 `form:"userID" json:"userID"`
	RightList string `form:"rightList" json:"rightList"`
}

//----------------------------------------------------------------------------------------------------------------------------
// #define CMD_LOGOUT	//用户注销，从服务端会话池中删除会话    URL: "/v1/api/sessions/:sid"	//sid（会话ID）即uid（用户ID）
// MiddleWareCheckUserHasLogin : 需要先检查用户是否是登陆状态
func HandleLogoutDelete(c *gin.Context) {
	log.HTTP.Info("HandleLogoutDelete BEGIN")
	//Get session id (user id) from url
	userID, err := strconv.ParseInt(c.Param("sid"), 10, 64)
	if err != nil {
		log.HTTP.Error("Form error")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Get current login user id from cookie | 从缓存中读取用户的userid
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.HTTP.Error("Get cookie failed")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Judge user whether valid or not
	if userID != accountID {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	//Record logout | 退出登录
	err = dao.RecordLogin(userID, "logout")
	if err != nil {
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

//---------------------------------------------------------------------------------------------------------------------------------
//#define CMD_USER_INFO	//查询用户个人信息   URL: "/v1/api/users/:uid"
//注：GET 须校验是否为管理员或参数uid是否与cookie中userID一致，否则为伪造请求
func HandleGainUserInfoGet(c *gin.Context) {
	log.HTTP.Info("HandleGainUserInfoGet BEGIN")
	//Get user id from url
	userID, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		log.HTTP.Error("Params error")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	//Get current login user id from cookie
	accountID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.HTTP.Error("Get cookie failed")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//Judge account is admin or not, just admin has right to register device
	if false == dao.IsAdmin(accountID) { //表示不是管理员
		//Judge user whether valid or not
		if userID != accountID {
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

//---------------------------------------------------------------------------------------------------------------------
//删除用户		需要校验是否是管理员，如果是管理员的进行下面的操作，否则返回伪造请求，返回错误信息
func HandleUserDelete(c *gin.Context) {
	log.HTTP.Info("HandleUserDelete BEGIN")
	//Get user id from url
	willDeluserID, err := strconv.ParseInt(c.Param("uid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	//Get current login user id from cookie
	currentUserID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.HTTP.Error("Get cookie failed")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	//校验是不是管理员
	if false == dao.IsAdmin(currentUserID) {
		//Judge user whether valid or not
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	err = dao.DeleteAccount(willDeluserID)
	if err != nil {
		log.HTTP.Error("Delete error")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrDelete.Code,
			"errMsg": errcode.ErrDelete.String,
		})
		return
	}
	//通知所有订阅客户端 该用户被删除
	echo.PublishDelete(willDeluserID)

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

//---------------------------------------------------------------------------------------------------------------------------
// #define CMD_USERS	//查询所有用户
//注：GET 和 POST须校验是否为管理员，否则为伪造请求
//URL: "/v1/api/users"
////查询所有非管理员用户 (由客户端在获取的列表中过滤了类型是1的
func HandleQueryAllUserInfoGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryAllUserInfoGet BEGIN")
	//校验是不是管理员， 从cookie中获取
	userID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.HTTP.Error("Get cookie failed")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}
	//Judge account is admin or not
	if false == dao.IsAdmin(userID) {
		//Judge user whether valid or not
		log.HTTP.Error("Fake request")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	userInfos, err := dao.QueryAllUserInfo()
	if err != nil {
		log.HTTP.Error("Query error")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":      errcode.ErrNoError.Code,
		"errMsg":   errcode.ErrNoError.String,
		"userList": userInfos,
	})
	return
}

//-----------------------------------------------------------------------------------------------------------------------------
//添加用户			//校验是不是管理员
func HandleAddUsersPost(c *gin.Context) {
	log.HTTP.Info("开始添加用户")
	//Get current login user id from cookie
	currentUserID, err := lib.GetCurrentUser(c.Request)
	if err != nil {
		log.HTTP.Error("Get cookie failed|获取用户的userid失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrCookie.Code,
			"errMsg": errcode.ErrCookie.String,
		})
		return
	}

	if false == dao.IsAdmin(currentUserID) {
		log.HTTP.Error("Fake request")
		//Judge user whether valid or not
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrFakeRequest.Code,
			"errMsg": errcode.ErrFakeRequest.String,
		})
		return
	}

	/*
		向数据库中添加用户
	*/
	var reqJSON reqAddUserPostJSON
	if c.BindJSON(&reqJSON) != nil {
		log.HTTP.Error("添加失败")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}
	//在这里写一个添加到数据库的操作，返回的err来判断是不是添加成功，不是的，返回错误
	if dao.AddUser(reqJSON.Username, reqJSON.Password, reqJSON.RightList) != nil {
		log.HTTP.Error("Add user error")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrAddUser.Code,
			"errMsg": errcode.ErrAddUser.String,
		})
		return
	}

	log.HTTP.Info("添加用户成功")
	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqAddUserPostJSON struct {
	Username string `form:"userName" json:"userName"`
	Password string `form:"userPassword" json:"userPassword"`
	RightList string `form:"rightList" json:"rightList"`
}
