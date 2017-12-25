package orm

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"impower/LicenseServer/env"
	"impower/LicenseServer/lib"
)

// QueryAuthRecordTotalNumber : Query autnentication record total number
func QueryAuthRecordTotalNumber(session *sql.Tx, userID int64, isAdmin bool, orderNumber, deviceSupplier string, tmBegin, tmEnd time.Time) (totalNum int, err error) {
	matchCond := " 1 = 1 AND "
	if len(orderNumber) > 0 {
		matchCond = matchCond + `order_number LIKE '%` + orderNumber + `%' AND `
	}
	if len(deviceSupplier) > 0 {
		matchCond = matchCond + `device_supplier LIKE '%` + deviceSupplier + `%' AND `
	}

	sql := "SELECT COUNT(1) FROM lic_t_authrecord WHERE "
	if isAdmin {
		sql = sql + matchCond + "record_time >= ? AND record_time <= ?"
	} else {
		sql = sql + matchCond + "record_time >= ? AND record_time <= ? AND user_id = " + strconv.FormatInt(userID, 10)
	}

	rows, err := session.Query(sql, tmBegin, tmEnd)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&totalNum); err != nil {
			return -1, err
		}
		return totalNum, nil
	}
	return -1, nil
}

// QueryAuthRecordByPage : Query authentication record by page
func QueryAuthRecordByPage(session *sql.Tx, userID int64, isAdmin bool, orderNumber, deviceSupplier string, tmBegin, tmEnd time.Time, offset, count int) (recordList []lib.AuthRecord, err error) {
	matchCond := " 1 = 1 AND "
	if len(orderNumber) > 0 {
		matchCond = matchCond + `order_number LIKE '%` + orderNumber + `%' AND `
	}
	if len(deviceSupplier) > 0 {
		matchCond = matchCond + `device_supplier LIKE '%` + deviceSupplier + `%' AND `
	}

	sql := `SELECT t1.batch, t1.device_supplier, t1.device_type, t1.device_model, t1.device_number, t1.order_number, t1.support_p2p, t2.user_name, t1.record_time 
			FROM lic_t_authrecord t1, lic_t_user t2 WHERE `
	if isAdmin {
		sql = sql + matchCond + `t1.record_time >= ? AND t1.record_time <= ? AND t1.user_id = t2.user_id 
			ORDER BY t1.id DESC LIMIT ?,?`
	} else {
		sql = sql + matchCond + `t1.record_time >= ? AND t1.record_time <= ? AND t1.user_id = t2.user_id 
			AND t1.user_id = ` + strconv.FormatInt(userID, 10) +
			` ORDER BY t1.id DESC LIMIT ?,?`
	}

	recordList = []lib.AuthRecord{}

	rows, err := session.Query(sql, tmBegin, tmEnd, offset, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		record := lib.AuthRecord{}
		if err = rows.Scan(&record.RecordBatch, &record.DeviceSupplier, &record.DeviceType, &record.DeviceModel,
			&record.DeviceNumber, &record.OrderNumber, &record.SupportP2P, &record.RecordUser, &record.RecordTime); err != nil {
			return nil, err
		}

		recordList = append(recordList, record)
	}

	return recordList, nil
}

// RecordAuthBatch : Record authentication batch
func RecordAuthBatch(session *sql.Tx, userID int64, record lib.AuthRecord) (err error) {
	sql := `INSERT INTO lic_t_authrecord VALUES (null, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := lib.Now()
	res, err := session.Exec(sql, record.RecordBatch, record.DeviceSupplier, record.DeviceType, record.DeviceModel,
		record.DeviceNumber, record.OrderNumber, record.SupportP2P, userID, now)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("insert authentication record information error")
	}

	return nil
}

// RecordMac : Record authentication
func RecordMac(session *sql.Tx, batch string, macList []string) (err error) {
	var valuesList []string
	for _, v := range macList {

		var values string
		values = "(null, '" + batch + "','" + v + "')"

		valuesList = append(valuesList, values)
	}
	valuesListStr := strings.Join(valuesList, ",")

	sql := "INSERT INTO lic_t_mac VALUES " + valuesListStr

	_, err = session.Exec(sql)

	return err
}

// QueryAuthMacByBatch : Query authentication mac address by batch
func QueryAuthMacByBatch(session *sql.Tx, userID int64, isAdmin bool, batch string) (macList []string, err error) {
	var sql string
	if isAdmin {
		sql = "SELECT mac FROM lic_t_mac WHERE batch = ?"
	} else {
		sql = `SELECT t1.mac FROM lic_t_mac t1,  lic_t_authrecord t2 
			WHERE t1.batch = ? AND t1.batch = t2.batch AND t2.user_id = ` + strconv.FormatInt(userID, 10)
	}

	macList = []string{}

	rows, err := session.Query(sql, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mac string
		if err = rows.Scan(&mac); err != nil {
			return nil, err
		}

		macList = append(macList, mac)
	}

	return macList, nil
}

// AuthMacAddress : Authenticate mac address
func AuthMacAddress(session *sql.Tx, macAddress string) (isValid bool, err error) {
	sql := "SELECT id FROM lic_t_mac WHERE mac = ?"

	rows, err := session.Query(sql, macAddress)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

// QueryP2PRecordTotalNumber : Query p2p record total number
func QueryP2PRecordTotalNumber(session *sql.Tx, userID int64, isAdmin bool, orderNumber string, tmBegin, tmEnd time.Time) (totalNum int, err error) {
	matchCond := " 1 = 1 AND "
	if len(orderNumber) > 0 {
		matchCond = matchCond + `order_number LIKE '%` + orderNumber + `%' AND `
	}

	sql := "SELECT COUNT(1) FROM lic_t_p2precord WHERE "
	if isAdmin {
		sql = sql + matchCond + "record_time >= ? AND record_time <= ?"
	} else {
		sql = sql + matchCond + "record_time >= ? AND record_time <= ? AND user_id = " + strconv.FormatInt(userID, 10)
	}

	rows, err := session.Query(sql, tmBegin, tmEnd)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&totalNum); err != nil {
			return -1, err
		}
		return totalNum, nil
	}
	return -1, nil
}

// QueryP2PRecordByPage : Query p2p record by page
func QueryP2PRecordByPage(session *sql.Tx, userID int64, isAdmin bool, orderNumber string, tmBegin, tmEnd time.Time, offset, count int) (recordList []lib.P2PRecord, err error) {
	matchCond := " 1 = 1 AND "
	if len(orderNumber) > 0 {
		matchCond = matchCond + `order_number LIKE '%` + orderNumber + `%' AND `
	}

	sql := `SELECT t1.batch, t1.p2p_number, t1.order_number, t2.user_name, t1.record_time 
			FROM lic_t_p2precord t1, lic_t_user t2 WHERE `
	if isAdmin {
		sql = sql + matchCond + `t1.record_time >= ? AND t1.record_time <= ? AND t1.user_id = t2.user_id 
			ORDER BY t1.id DESC LIMIT ?,?`
	} else {
		sql = sql + matchCond + `t1.record_time >= ? AND t1.record_time <= ? AND t1.user_id = t2.user_id 
			AND t1.user_id = ` + strconv.FormatInt(userID, 10) +
			` ORDER BY t1.id DESC LIMIT ?,?`
	}

	recordList = []lib.P2PRecord{}

	rows, err := session.Query(sql, tmBegin, tmEnd, offset, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		record := lib.P2PRecord{}
		if err = rows.Scan(&record.RecordBatch, &record.P2PNumber, &record.OrderNumber, &record.RecordUser, &record.RecordTime); err != nil {
			return nil, err
		}

		recordList = append(recordList, record)
	}

	return recordList, nil
}

// RecordP2PBatch : Record p2p batch
func RecordP2PBatch(session *sql.Tx, userID int64, record lib.P2PRecord) (err error) {
	sql := `INSERT INTO lic_t_p2precord VALUES (null, ?, ?, ?, ?, ?)`

	now := lib.Now()
	res, err := session.Exec(sql, record.RecordBatch, record.P2PNumber, record.OrderNumber, userID, now)
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("insert p2p record information error")
	}

	return nil
}

// RecordP2P : Record p2p
func RecordP2P(session *sql.Tx, batch string, p2pList []string) (err error) {
	var valuesList []string
	for _, v := range p2pList {

		var values string
		values = "(null, '" + batch + "','" + v + "','" + strconv.Itoa(env.P2P_STATUS_UNUSED) + "','')"

		valuesList = append(valuesList, values)
	}
	valuesListStr := strings.Join(valuesList, ",")

	sql := "INSERT INTO lic_t_p2p VALUES " + valuesListStr

	_, err = session.Exec(sql)

	return err
}

// QueryP2PByBatch : Query p2p list by batch
func QueryP2PByBatch(session *sql.Tx, userID int64, isAdmin bool, batch string) (p2pList []lib.P2P, err error) {
	var sql string
	if isAdmin {
		sql = "SELECT p2p_id, status, mac FROM lic_t_p2p WHERE batch = ?"
	} else {
		sql = `SELECT t1.p2p_id, t1.status, t1.mac FROM lic_t_p2p t1,  lic_t_p2precord t2 
			WHERE t1.batch = ? AND t1.batch = t2.batch AND t2.user_id = ` + strconv.FormatInt(userID, 10)
	}

	p2pList = []lib.P2P{}

	rows, err := session.Query(sql, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p2p lib.P2P
		if err = rows.Scan(&p2p.P2PID, &p2p.Status, &p2p.MAC); err != nil {
			return nil, err
		}

		p2pList = append(p2pList, p2p)
	}

	return p2pList, nil
}

// AuthMacSupportP2P : Authenticate mac address is support p2p
func AuthMacSupportP2P(session *sql.Tx, macAddress string) (isValid bool, err error) {
	sql := "SELECT t1.id FROM lic_t_mac t1, lic_t_authrecord t2 WHERE t1.mac = ? AND t1.batch = t2.batch AND t2.support_p2p = ?"

	rows, err := session.Query(sql, macAddress, env.P2P_SUPPORT)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}

	return false, nil
}

// HasAppliedP2P : Has applied p2p id or not
func HasAppliedP2P(session *sql.Tx, macAddress string) (p2pID string, err error) {
	sql := "SELECT p2p_id from lic_t_p2p WHERE mac = ? AND status = ?"

	rows, err := session.Query(sql, macAddress, env.P2P_STATUS_USED)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if false == rows.Next() {
		return "", nil
	}

	if err = rows.Scan(&p2pID); err != nil {
		return "", err
	}

	return p2pID, nil
}

// ApplyNewP2P : Apply new p2p id
func ApplyNewP2P(session *sql.Tx) (p2pID string, err error) {
	sql := "SELECT p2p_id FROM lic_t_p2p WHERE status = ?  ORDER BY id LIMIT 0,1"

	rows, err := session.Query(sql, env.P2P_STATUS_UNUSED)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if false == rows.Next() {
		return "", nil
	}

	if err = rows.Scan(&p2pID); err != nil {
		return "", err
	}

	return p2pID, nil
}

// UpdateP2P : Update p2p info
func UpdateP2P(session *sql.Tx, macAddress, p2pID string) (err error) {
	sql := "UPDATE lic_t_p2p SET status = ?, mac = ? WHERE p2p_id = ?"

	_, err = session.Exec(sql, env.P2P_STATUS_USED, macAddress, p2pID)
	if err != nil {
		return err
	}

	return nil
}
