package orm

import (
	"database/sql"
	"errors"
	"time"

	"impower/LicenseServer/lib"
)

// DeletePhoneCode : Delete phone code
func DeletePhoneCode(session *sql.Tx, phoneNo string) (err error) {
	sql := "delete from lic_t_phoneverify where phone_no = ? "

	_, err = session.Exec(sql, phoneNo)
	if err != nil {
		return err
	}

	return nil
}

// InsertPhoneCode : Insert phone code
func InsertPhoneCode(session *sql.Tx, phoneNo, verifyCode string, deadline time.Time) (err error) {
	sql := "insert into lic_t_phoneverify values (?,?,?) "

	_, err = session.Exec(sql, phoneNo, verifyCode, deadline)
	if err != nil {
		return err
	}

	return nil
}

// VerifyPhoneCode : Verify phone code
func VerifyPhoneCode(session *sql.Tx, phoneNo, verifyCode string) (err error) {
	sql := "select deadline from lic_t_phoneverify where phone_no = ? and verify_code = ?"

	rows, err := session.Query(sql, phoneNo, verifyCode)
	if err != nil {
		return err
	}
	defer rows.Close()

	if false == rows.Next() {
		return errors.New("invalid phone or code")
	}

	var deadlineStr string
	rows.Scan(&deadlineStr)

	var deadline time.Time
	if deadline, err = time.Parse(lib.TimeLayout, deadlineStr); err != nil {
		return err
	}

	if deadline.Before(lib.Now()) {
		return errors.New("phone code is expire")
	}

	return nil
}
