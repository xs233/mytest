package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func exec(session *sql.Tx, sqlStr string) (err error) {
	result, _ := session.Exec(sqlStr)
	var num int64
	if num, err = result.RowsAffected(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(num)
		if num, err = result.LastInsertId(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(num)
		}
	}

	return
}

func queryRow(session *sql.Tx, sqlStr string, alg ...interface{}) error {
	return session.QueryRow(sqlStr).Scan(alg...)
}

func main() {
	if db, err := sql.Open("mysql", "impower:impower2017@tcp(192.168.0.182:3306)/silverdb?charset=utf8"); err != nil {
		fmt.Println(err)
	} else {
		session, _ := db.Begin()
		sqlStr := `select t3.group_id, t3.group_name, t4.device_id, t4.device_type, t4.device_uuid, t4.device_vms_id, t4.device_name, t4.device_ip,
		t4.rtsp_url, t4.main_stream_url, t4.sub_stream_url, t4.p2p_key from (select t1.group_id, t1.group_name, t2.device_id from imp_t_group t1 
		left join imp_t_groupdevice t2 on t1.group_id = t2.group_id) t3 left join imp_t_device t4 on  
		t3.device_id = t4.device_id`
		var groupID int64
		var groupName sql.NullString
		var device Device
		if err := queryRow(session, sqlStr, &groupID, &groupName, &device.DeviceID, &device.DeviceType, &device.DeviceUUID, &device.DeviceVmsID, 
			&device.DeviceName, &device.DeviceIP, &device.RtspURL, &device.MainStreamURL, &device.SubStreamURL, &device.P2PKey); err != nil {
			session.Rollback()
			fmt.Println(err)
			return
		}
		
		err := sql.ErrNoRows
		fmt.Println(device)
		session.Commit()
	}
}

type Device struct {
	DeviceID      sql.NullInt64  `form:"deviceID" json:"deviceID"`
	DeviceType	  sql.NullInt64    `form:"deviceType" json:"deviceType"`
	DeviceUUID	  sql.NullString `form:"deviceUUID" json:"deviceUUID"`
	DeviceVmsID   sql.NullString `form:"deviceVmsID" json:"deviceVmsID"`
	DeviceName    sql.NullString `form:"deviceName" json:"deviceName"`
	DeviceIP      sql.NullString `form:"deviceIP" json:"deviceIP"`
	RtspURL       sql.NullString `form:"rtspURL" json:"rtspURL"`
	MainStreamURL sql.NullString `form:"mainStreamURL" json:"mainStreamURL"`
	SubStreamURL  sql.NullString `form:"subStreamURL" json:"subStreamURL"`
	Account		  sql.NullString `form:"account" json:"account"`
	Password	  sql.NullString `form:"password" json:"password"`
	P2PKey		  sql.NullString `form:"p2pKey" json:"p2pKey"`
}