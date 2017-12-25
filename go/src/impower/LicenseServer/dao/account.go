package dao

import (
	"impower/LicenseServer/dao/orm"
	"impower/LicenseServer/env"
	"impower/LicenseServer/lib"
	"impower/LicenseServer/log"
)

// IsAdmin : juge is or not admin
func IsAdmin(accountID int64) (isAdmin bool) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return false
	}

	isAdmin = orm.IsAdmin(session, accountID)

	session.Commit()
	return isAdmin
}

// Register : Register account
func Register(userName, userPassword, phoneNo string) (userID int64, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return -1, err
	}

	//Account is existing
	isExisting, err := orm.IsUserNameExisting(session, userName)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return -1, err
	}
	if isExisting {
		session.Commit()
		return 0, nil
	}

	//Get account id
	next, err := orm.GetSequenceIDByName(env.SequenceNameUserID)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return 0, err
	}
	userID = next

	//Insert user info
	err = orm.InsertAccount(session, userID, userName, userPassword, phoneNo)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return -1, err
	}

	session.Commit()
	return userID, nil
}

// Authenticate : User authenticate
func Authenticate(userName, userPassword string) (userID int64, userType int, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return -1, -1, err
	}

	var isExisting bool

	isExisting, userID, userType, err = orm.IsAccountPasswordExisting(session, userName, userPassword)

	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return -1, -1, err
	}

	if isExisting {
		session.Commit()
		return userID, userType, nil
	}

	//Account is not existing
	session.Commit()
	return -1, -1, nil
}

// Login : User login
func Login(userName, userPassword string) (userID int64, userType int, err error) {
	//Authenticate user account
	userID, userType, err = Authenticate(userName, userPassword)
	if err != nil {
		log.Root.Error(err)
		return -1, -1, err
	}

	return userID, userType, nil
}

// RecordLogin : Record login
func RecordLogin(userID int64, login string) (err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return err
	}

	err = orm.InsertLoginRecord(session, userID, login)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	session.Commit()
	return nil
}

// ModifyPassword : Modify password
func ModifyPassword(userName, oldPassword, newPassword string) (err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return err
	}

	err = orm.ModifyPassword(session, userName, oldPassword, newPassword)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	session.Commit()
	return nil
}

// ResetPassword : Reset password
func ResetPassword(userID int64, newPassword string) (err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return err
	}

	err = orm.ResetPassword(session, userID, newPassword)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
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
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return err
	}

	err = orm.DeleteAccount(session, UserID)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
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
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return false, err
	}

	isExisting, err = orm.IsUserNameExisting(session, userName)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return false, err
	}

	session.Commit()
	return isExisting, nil
}

// IsAccountExisting : Is account is existing
func IsAccountExisting(account, password string) (isExisting bool, accountID int64, userType int, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return false, -1, -1, err
	}

	isExisting, accountID, userType, err = orm.IsAccountPasswordExisting(session, account, password)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return false, -1, -1, err
	}

	session.Commit()
	return isExisting, accountID, userType, nil
}

// QueryUserInfo : Query user info
//param userID : user id
//return user : user info
//retuen err : error info
func QueryUserInfo(userID int64) (user lib.User, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return user, err
	}

	user, err = orm.QueryUserInfo(session, userID)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return user, err
	}

	session.Commit()
	return user, nil
}

// QueryAllUserInfo : Query all user info
func QueryAllUserInfo() (user []lib.User, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return user, err
	}

	user, err = orm.QueryAllUserInfo(session)

	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return user, err
	}

	session.Commit()
	return user, nil
}
