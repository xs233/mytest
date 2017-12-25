package orm

import (
	"database/sql"
	"errors"
	"impower/LicenseServer/env"
	"impower/LicenseServer/lib"
)

// IsAdmin : Is admin
func IsAdmin(session *sql.Tx, accountID int64) (isAdmin bool) {
	sql := "select * from lic_t_user where user_id = ? and user_type = ?"

	rows, err := session.Query(sql, accountID, int(env.AccountTypeAdmin))
	if err != nil {
		return false
	}
	defer rows.Close()

	if rows.Next() {
		return true
	}

	return false
}

// IsAccountPasswordExisting : Query is account/password is existing
func IsAccountPasswordExisting(session *sql.Tx, userName, userPassword string) (existing bool, userID int64, userType int, err error) {
	sql := "select user_id, user_type from lic_t_user where user_name = ? and password = ?"

	rows, err := session.Query(sql, userName, userPassword)
	if err != nil {
		return false, -1, -1, err
	}
	defer rows.Close()

	if false == rows.Next() {
		return false, -1, -1, nil
	}

	if err = rows.Scan(&userID, &userType); err != nil {
		return true, -1, -1, err
	}

	return true, userID, userType, nil
}

//DeleteAccount : delete account
func DeleteAccount(session *sql.Tx, userID int64) (err error) {
	sql := `delete from lic_t_user  where user_id=?`

	_, err = session.Exec(sql, userID)
	if err != nil {
		return err
	}
	return nil
}

// InsertAccount : Insert account, account could be username or phoneno or email
func InsertAccount(session *sql.Tx, userID int64, userName, userPassword, phoneNo string) (err error) {
	sql := `insert into lic_t_user values (?, ?, ?, ?, ?, ?)`

	now := lib.Now()
	res, err := session.Exec(sql, userID, env.AccountTypeUser, userName, userPassword, phoneNo, now)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("insert user error")
	}

	return nil
}

// InsertLoginRecord : Insert login Record
func InsertLoginRecord(session *sql.Tx, userID int64, login string) (err error) {
	var loginType int
	if login == "login" {
		loginType = 1
	} else {
		loginType = 0
	}
	sql := "insert into lic_t_login set user_id=?, login=?, log_time=?"

	now := lib.Now()

	_, err = session.Exec(sql, userID, loginType, now)
	if err != nil {
		return err
	}

	return nil
}

// ModifyPassword : Modify password
func ModifyPassword(session *sql.Tx, userName, oldPassword, newPassword string) (err error) {
	sql := "update lic_t_user set password = ? where user_name = ? and password = ?"

	res, err := session.Exec(sql, newPassword, userName, oldPassword)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("account and password are not right")
	}

	return nil
}

// ResetPassword : Reset password
func ResetPassword(session *sql.Tx, userID int64, newPassword string) (err error) {
	sql := "update lic_t_user set password = ? where user_id = ?"

	res, err := session.Exec(sql, newPassword, userID)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("account is not existing")
	}

	return nil
}

// IsUserNameExisting : Judge username is existing or not
func IsUserNameExisting(session *sql.Tx, userName string) (isExisting bool, err error) {
	sql := "select user_id from lic_t_user where user_name = ?"

	rows, err := session.Query(sql, userName)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if false == rows.Next() {
		return false, nil
	}

	return true, nil
}

//QueryUserInfo : Query user info
func QueryUserInfo(session *sql.Tx, userID int64) (user lib.User, err error) {
	sql := `select user_id, user_name, user_type, phone_no from lic_t_user where user_id = ?`

	rows, err := session.Query(sql, userID)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&user.UserID, &user.UserName, &user.UserType, &user.PhoneNo); err != nil {
			return user, err
		}
	}

	return user, nil
}

//QueryAllUserInfo ： Query all user info
func QueryAllUserInfo(session *sql.Tx) (userList []lib.User, err error) {
	sql := `select user_id, user_name, user_type, phone_no from lic_t_user where user_type <> ?`

	rows, err := session.Query(sql, env.AccountTypeAdmin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userList = []lib.User{}
	for rows.Next() {
		user := lib.User{}
		if err = rows.Scan(&user.UserID, &user.UserName, &user.UserType, &user.PhoneNo); err != nil {
			return nil, err
		}
		userList = append(userList, user)
	}
	return userList, nil
}
