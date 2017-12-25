package dao

import (
	"time"

	"impower/LicenseServer/dao/orm"
	"impower/LicenseServer/log"
)

// VerifyPhoneCode : Verify phone code
func VerifyPhoneCode(phoneNo, verifyCode string) (err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return err
	}

	err = orm.VerifyPhoneCode(session, phoneNo, verifyCode)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	err = orm.DeletePhoneCode(session, phoneNo)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	session.Commit()
	return nil
}

// SavePhoneCode : Save phone code
func SavePhoneCode(phoneNo, verifyCode string, deadline time.Time) (err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return err
	}

	err = orm.DeletePhoneCode(session, phoneNo)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	err = orm.InsertPhoneCode(session, phoneNo, verifyCode, deadline)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	session.Commit()
	return nil
}
