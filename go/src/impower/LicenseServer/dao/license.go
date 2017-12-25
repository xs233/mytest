package dao

import (
	"errors"
	"time"

	"impower/LicenseServer/dao/orm"
	"impower/LicenseServer/lib"
	"impower/LicenseServer/log"
)

// QueryAuthRecordTotalNumber : Query authentication record total number
//param userID : user id
//param orderNumber : order number
//param deviceSupplier : device supplier
//param tmBegin : query begin time
//param tmEnd : query end time
//return totalNum : total number
//retuen err : error info
func QueryAuthRecordTotalNumber(userID int64, orderNumber, deviceSupplier string, tmBegin, tmEnd time.Time) (totalNum int, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return totalNum, err
	}

	isAdmin := orm.IsAdmin(session, userID)

	totalNum, err = orm.QueryAuthRecordTotalNumber(session, userID, isAdmin, orderNumber, deviceSupplier, tmBegin, tmEnd)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return totalNum, err
	}

	session.Commit()
	return totalNum, nil
}

// QueryAuthRecordByPage : Query authentication record by page
//param userID : user id
//param orderNumber : order number
//param deviceSupplier : device supplier
//param tmBegin : query begin time
//param tmEnd : query end time
//param offset : offset
//param count : count
//return recordList : record list
//retuen err : error info
func QueryAuthRecordByPage(userID int64, orderNumber, deviceSupplier string, tmBegin, tmEnd time.Time, offset, count int) (recordList []lib.AuthRecord, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return nil, err
	}

	isAdmin := orm.IsAdmin(session, userID)

	recordList, err = orm.QueryAuthRecordByPage(session, userID, isAdmin, orderNumber, deviceSupplier, tmBegin, tmEnd, offset, count)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return nil, err
	}

	session.Commit()
	return recordList, nil
}

// RecordAuth : Record authentication
//param userID : user id
//param record : authentication record information
//param macList : mac address list
//retuen err : error info
func RecordAuthBatch(userID int64, record lib.AuthRecord, macList []string) (err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return err
	}

	//Record authentication
	err = orm.RecordAuthBatch(session, userID, record)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	//Record mac address list
	err = orm.RecordMac(session, record.RecordBatch, macList)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	session.Commit()
	return nil
}

// QueryAuthMacByBatch : Query authentication mac address by batch
//param userID : user id
//param batch : authentication record batch
//retuen macList : mac address list
//retuen err : error info
func QueryAuthMacByBatch(userID int64, batch string) (macList []string, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return nil, err
	}

	isAdmin := orm.IsAdmin(session, userID)

	macList, err = orm.QueryAuthMacByBatch(session, userID, isAdmin, batch)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return nil, err
	}

	session.Commit()
	return macList, nil
}

// AuthMacAddress : Authenticate mac address
//param macAddress :  mac address
//retuen isValid : mac address is in authenticated list or not
//retuen err : error info
func AuthMacAddress(macAddress string) (isValid bool, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return false, err
	}

	isValid, err = orm.AuthMacAddress(session, macAddress)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return false, err
	}

	session.Commit()
	return isValid, nil
}

// QueryP2PRecordTotalNumber : Query p2p record total number
//param userID : user id
//param orderNumber : order number
//param tmBegin : query begin time
//param tmEnd : query end time
//return totalNum : total number
//retuen err : error info
func QueryP2PRecordTotalNumber(userID int64, orderNumber string, tmBegin, tmEnd time.Time) (totalNum int, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return totalNum, err
	}

	isAdmin := orm.IsAdmin(session, userID)

	totalNum, err = orm.QueryP2PRecordTotalNumber(session, userID, isAdmin, orderNumber, tmBegin, tmEnd)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return totalNum, err
	}

	session.Commit()
	return totalNum, nil
}

// QueryP2PRecordByPage : Query p2p record by page
//param userID : user id
//param orderNumber : order number
//param tmBegin : query begin time
//param tmEnd : query end time
//param offset : offset
//param count : count
//return recordList : record list
//retuen err : error info
func QueryP2PRecordByPage(userID int64, orderNumber string, tmBegin, tmEnd time.Time, offset, count int) (recordList []lib.P2PRecord, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return nil, err
	}

	isAdmin := orm.IsAdmin(session, userID)

	recordList, err = orm.QueryP2PRecordByPage(session, userID, isAdmin, orderNumber, tmBegin, tmEnd, offset, count)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return nil, err
	}

	session.Commit()
	return recordList, nil
}

// RecordP2PBatch : Record p2p
//param userID : user id
//param record : p2p record information
//param p2pList : p2p list
//retuen err : error info
func RecordP2PBatch(userID int64, record lib.P2PRecord, p2pList []string) (err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return err
	}

	//Record p2p batch
	err = orm.RecordP2PBatch(session, userID, record)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	//Record p2p list
	err = orm.RecordP2P(session, record.RecordBatch, p2pList)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return err
	}

	session.Commit()
	return nil
}

// QueryP2PByBatch : Query p2p list by batch
//param userID : user id
//param batch : p2p record batch
//retuen p2pList : p2p list
//retuen err : error info
func QueryP2PByBatch(userID int64, batch string) (p2pList []lib.P2P, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return nil, err
	}

	isAdmin := orm.IsAdmin(session, userID)

	p2pList, err = orm.QueryP2PByBatch(session, userID, isAdmin, batch)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return nil, err
	}

	session.Commit()
	return p2pList, nil
}

// AuthMacSupportP2P : Authenticate mac address is support p2p
//param macAddress :  mac address
//retuen isValid : mac address is in authenticated list and support p2p or not
//retuen err : error info
func AuthMacSupportP2P(macAddress string) (isValid bool, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return false, err
	}

	isValid, err = orm.AuthMacSupportP2P(session, macAddress)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return false, err
	}

	session.Commit()
	return isValid, nil
}

// GetP2P : Get p2p id by mac address
//param macAddress :  mac address
//retuen p2pID : p2p id
//retuen err : error info
func GetP2P(macAddress string) (p2pID string, err error) {
	session, err := orm.DB.Begin()
	if err != nil {
		log.Root.Error(err)
		return "", err
	}

	p2pID, err = orm.HasAppliedP2P(session, macAddress)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return "", err
	}

	//The mac address has applied p2p id already
	if p2pID != "" {
		session.Commit()
		return p2pID, nil
	}

	//Apply new p2p id
	p2pID, err = orm.ApplyNewP2P(session)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return "", err
	}

	if p2pID == "" {
		session.Rollback()
		log.Root.Error("No usable p2p id.")
		return "", errors.New("no usable p2p id")
	}

	err = orm.UpdateP2P(session, macAddress, p2pID)
	if err != nil {
		session.Rollback()
		log.Root.Error(err)
		return "", err
	}

	session.Commit()
	return p2pID, nil
}
