package sql

import (
	"SilverBusinessServer/env"
	"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
	"database/sql"
	"errors"
)

// Is Admin | 是否是管理员
func IsAdmin(session *sql.Tx, accountID int64) (isAdmin bool) {
	sql := "select * from imp_t_user where user_id = ? and user_type = ?" // 标准MySQL语句

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

// Is Account And Password Existing | 验证账号和密码是否存在
func IsAccountAndPasswordExisting(session *sql.Tx, username, password string) (isExisting bool, userID int64, userType int, rightList string, err error) {
	sql := `SELECT t1.user_id, t1.user_type, t2.right_list FROM imp_t_user t1 
			LEFT JOIN imp_t_right t2 ON t1.user_id = t2.user_id
			WHERE t1.user_name = ? AND t1.password = ?`
	rows, err := session.Query(sql, username, password)
	if err != nil {
		log.HTTP.Error("sql验证错误")
		return false, -1, -1, "", err
	}

	defer rows.Close()
	if false == rows.Next() {
		log.HTTP.Error("rows错误")
		return false, -1, -1, "", nil
	}

	var rightInterface interface{}
	rows.Scan(&userID, &userType, &rightInterface)
	if rightInterface != nil {
		rightList = string(rightInterface.([]uint8))
	} else {
		rightList = ""
	}
	
	return true, userID, userType, rightList, nil
}

// Delete Account | 删除账号
func DelecteAccount(session *sql.Tx, userID int64) (err error) {
	sql := `delete from imp_t_user where user_id=?`
	_, err = session.Exec(sql, userID)
	if err != nil {
		return err
	}
	return nil
}

// Insert Account | 添加账号
func InsertAccount(session *sql.Tx, userID int64, userName, passWord string) (err error) {
	now := lib.Now() //用户获取当前的时间
	sql := `insert into imp_t_user values (?, ?, ?, ?,  ?)`

	_, err = session.Exec(sql, userID, 2, userName, passWord, now)
	if err != nil {
		return err
	}
	return nil
}

//插入登陆
func InsertLoginRecord(session *sql.Tx, userID int64, login string) (err error) {
	var loginType int
	if login == "login" {
		loginType = 1
	} else {
		loginType = 0
	}
	sql := "insert into imp_t_login set user_id=?, login=?, log_time=?"

	now := lib.Now()

	_, err = session.Exec(sql, userID, loginType, now)
	if err != nil {
		return err
	}

	return nil
}

// Modify Password | 修改密码
func ModifyPassword(session *sql.Tx, username, oldPassword, newPassword string) (err error) {
	sql := "update imp_t_user set password = ? where user_name = ? and password = ?"
	res, err := session.Exec(sql, newPassword, username, oldPassword)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("account and password are not right")
	}
	return nil
}

func ResetPasswordSQL(session *sql.Tx, username, newPassword string) (err error) {
	sql := "update imp_t_user set password = ? where user_name = ?"
	_, err = session.Exec(sql, newPassword, username)
	if err != nil {
		return err
	}

	return nil
}

// Is UserName Existing | 是否是现有的用户名称
func IsUserNameExisting(session *sql.Tx, userName string) (isExisting bool, err error) {
	sql := "select user_id from imp_t_user where user_name = ?"

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

// Query User Info | 查询用户信息
func QueryUserInfo(session *sql.Tx, userID int64) (user lib.User, err error) {

	sql := `SELECT t1.user_id, t1.user_name, t1.user_type, t2.right_list FROM imp_t_user t1 
			INNER JOIN imp_t_right t2 ON t1.user_id = t2.user_id WHERE t2.user_id=?`

	rows, err := session.Query(sql, userID)
	if err != nil {
		return user, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&user.UserID, &user.UserName, &user.UserType, &user.RightList); err != nil {
			return user, err
		}
	}

	return user, nil
}

// QueryAllUserInfo | 查询所有用户的信息
func QueryAllUserInfo(session *sql.Tx) (users []lib.UserAccount, err error) {

	sql := `SELECT t1.user_id, t1.user_name, t2.right_list FROM imp_t_user t1 INNER JOIN imp_t_right t2 
			ON t1.user_id = t2.user_id`

	rows, err := session.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users = []lib.UserAccount{}
	for rows.Next() {
		user := lib.UserAccount{}
		if err = rows.Scan(&user.UserID, &user.UserName, &user.RightList); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

//AddUserInfoSQL | 向表中添加新的用户，添加用户的账号和密码 insert into imp_t_useruser
func AddUserInfoSQL(session *sql.Tx, userID int64, username, password, rightList string) (err error) {
	now := lib.Now()
	sql := `insert into imp_t_user values (?, ?, ?, ?,  ?)`
	_, err = session.Exec(sql, userID, 2, username, password, now)
	if err != nil {
		log.HTTP.Error("数据库中添加的时候失败")
		return err
	}

	sql = `INSERT INTO imp_t_right VALUES (?, ?)`
	_, err = session.Exec(sql, userID, rightList)
	if err != nil {
		log.HTTP.Error(err)
		return err
	}

	log.HTTP.Info("数据库添加成功")
	return nil
}

// Reset user permissions
func ResetRight(session *sql.Tx, userID int64, rightList string) error {
	sql := `UPDATE imp_t_right SET right_list = ? WHERE user_id = ?`
	if _, err := session.Exec(sql, rightList, userID); err != nil {
		return err
	}

	return nil
}