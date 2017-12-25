package dao

import (
	"SilverBusinessServer/dao/sql"
	"SilverBusinessServer/env"
	"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
)

// IsAdmin : juge is or not admin |是不是管理员
func IsAdmin(accountID int64) (isAdmin bool) {
	session, err := sql.DB.Begin()
	if err != nil {
		return false
	}

	isAdmin = sql.IsAdmin(session, accountID)

	session.Commit()
	return isAdmin
}

// Register : Register account | 注册账号
func Register(username, password string) (UserID int64, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return -1, err
	}

	//Account is existing
	isExisting, err := sql.IsUserNameExisting(session, username)
	if err != nil {
		session.Rollback()
		return -1, err
	}
	if isExisting {
		session.Commit()
		return 0, nil
	}

	//Get account id
	next, err := sql.GetSequenceIDByName(env.SequenceNameUserID)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	UserID = next
	//Insert user info
	err = sql.InsertAccount(session, UserID, username, password)
	if err != nil {
		session.Rollback()
		return -1, err
	}
	session.Commit()
	return UserID, nil
}

// Authenticate : User authenticate | 用户鉴定
func Authenticate(username, password string) (userID int64, usertype int, rightList string, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		log.HTTP.Error("sql err")
		return -1, -1, "", err
	}

	var isExisting bool

	isExisting, userID, usertype, rightList, err = sql.IsAccountAndPasswordExisting(session, username, password)

	if err != nil {
		log.HTTP.Error("预计是这里")
		session.Rollback()
		return -1, -1, "", err
	}

	if isExisting {
		session.Commit()
		return userID, usertype, rightList, nil
	}

	//Account is not existing
	session.Commit()
	return -1, -1, rightList, nil
}

// Login : User login 用户登录
func Login(username, password string) (userID int64, usertype int, rightList string, err error) {

	//Authenticate user account | 验证用户账号存在
	userID, usertype, rightList, err = Authenticate(username, password)
	if err != nil {
		log.HTTP.Error("验证账号不存在，或者有问题")
		return -1, -1, "", err
	}

	return userID, usertype, rightList, nil
}

// RecordLogin : Record login
func RecordLogin(userID int64, login string) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.InsertLoginRecord(session, userID, login)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

// ModifyPassword : 修改密码
func ModifyPassword(username, oldPassword, newPassword string) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.ModifyPassword(session, username, oldPassword, newPassword)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//reset Password | 重置用户的密码，设置为111111 【重置密码就是相当于修改密码，只是重置的密码是固定的，是“111111”的md5加密
func ResetPassword(username, newPassword string) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		log.HTTP.Error("db err")
		return err
	}

	err = sql.ResetPasswordSQL(session, username, newPassword)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

// DeleteAccount : Delete account
//param account : installer account
//param masterUserID : master user id
//retuen err : error info
func DeleteAccount(UserID int64) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.DelecteAccount(session, UserID)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

// IsUserNameExisting : Judge username is existing or not
//param userName : user name
//retuen isExisting : user name is existing or not, 0-not existing,1-existing
//retuen err : error info
func IsUserNameExisting(userName string) (isExisting bool, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return false, err
	}

	isExisting, err = sql.IsUserNameExisting(session, userName)
	if err != nil {
		session.Rollback()
		return false, err
	}

	session.Commit()
	return isExisting, nil
}

// IsAccountExisting : Is account is existing
func IsAccountExisting(account, password string) (isExisting bool, accountID int64, userType int, rightList string, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return false, -1, -1, "", err
	}

	isExisting, accountID, userType, rightList, err = sql.IsAccountAndPasswordExisting(session, account, password)
	if err != nil {
		session.Rollback()
		return false, -1, -1, "", err
	}

	session.Commit()
	return isExisting, accountID, userType, rightList, nil
}

// QueryUserInfo : Query user info
//param userID : user id
//return user : user info
//retuen err : error info
func QueryUserInfo(userID int64) (user lib.User, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return user, err
	}

	user, err = sql.QueryUserInfo(session, userID)
	if err != nil {
		session.Rollback()
		return user, err
	}

	session.Commit()
	return user, nil
}

// QueryAllUserInfo : Query all user info | 是所有非管理员的用户
func QueryAllUserInfo() (user []lib.UserAccount, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return user, err
	}

	user, err = sql.QueryAllUserInfo(session)

	if err != nil {
		session.Rollback()
		return user, err
	}

	session.Commit()
	return user, nil
}

//add User
func AddUser(username, password, rightList string) (err error) {
	log.HTTP.Info("在数据中添加用户")
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	// 创建userID的序列
	UserID, err := sql.GetSequenceIDByName(env.SequenceNameUserID)
	if err != nil {
		log.HTTP.Error("创建userID序列的时候错误")
		session.Rollback()
		return err
	}

	err = sql.AddUserInfoSQL(session, UserID, username, password, rightList)
	if err != nil {
		log.HTTP.Error("数据库中添加用户失败")
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

// Reset user permissions
func ResetRight(userID int64, rightList string) error {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	if err := sql.ResetRight(session, userID, rightList); err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}