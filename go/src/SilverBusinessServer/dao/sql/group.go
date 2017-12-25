package sql

import (
	"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
	"SilverBusinessServer/env"
	"database/sql"
	"strconv"
	"errors"
	//"strings"
)

type DeviceInterfaceStruct struct {
	DeviceID      interface{}
	DeviceType    interface{}
	DeviceUUID    interface{}
	DeviceVmsID   interface{}
	DeviceName    interface{}
	DeviceIP      interface{}
	RtspURL       interface{}
	MainStreamURL interface{}
	SubStreamURL  interface{}
	P2PKey        interface{}
}
//-------------------------------------------------------------------------------
// assignDevice : Assign the device info according to the DeviceInterfaceStruct
func assignDevice(dis DeviceInterfaceStruct) (device lib.Device) {
	if dis.DeviceID != nil {
		device.DeviceID, _ = strconv.ParseInt(string(dis.DeviceID.([]uint8)), 10, 64)
	} else {
		device.DeviceID = 0
	}

	if dis.DeviceType != nil {
		typeValue, _ := strconv.ParseInt(string(dis.DeviceType.([]uint8)), 10, 32)
		device.DeviceType = int(typeValue)
	} else {
		device.DeviceType = 0
	}
	
	if dis.DeviceUUID != nil {
		device.DeviceUUID = string(dis.DeviceUUID.([]uint8))
	} else {
		device.DeviceUUID = ""
	}

	if dis.DeviceVmsID != nil {
		device.DeviceVmsID = string(dis.DeviceVmsID.([]uint8))
	} else {
		device.DeviceVmsID = ""
	}

	if dis.DeviceName != nil {
		device.DeviceName = string(dis.DeviceName.([]uint8))
	} else {
		device.DeviceName = ""
	}

	if dis.DeviceIP != nil {
		device.DeviceIP = string(dis.DeviceIP.([]uint8))
	} else {
		device.DeviceIP = ""
	}

	if dis.RtspURL != nil {
		device.RtspURL = string(dis.RtspURL.([]uint8))
	} else {
		device.RtspURL = ""
	}

	if dis.MainStreamURL != nil {
		device.MainStreamURL = string(dis.MainStreamURL.([]uint8))
	} else {
		device.MainStreamURL = ""
	}

	if dis.SubStreamURL != nil {
		device.SubStreamURL = string(dis.SubStreamURL.([]uint8))
	} else {
		device.SubStreamURL = ""
	}

	if dis.P2PKey != nil {
		device.P2PKey = string(dis.P2PKey.([]uint8))
	} else {
		device.P2PKey = ""
	}

	return
}

func scanDevice(row *sql.Rows) (lib.Device, error) {
	var device DeviceInterfaceStruct
	if err := row.Scan(&device.DeviceID, &device.DeviceType, &device.DeviceUUID, &device.DeviceVmsID, &device.DeviceName, 
		&device.DeviceIP, &device.RtspURL, &device.MainStreamURL, &device.SubStreamURL, &device.P2PKey); err != nil {
		return lib.Device{} , err
	}
	
	return assignDevice(device), nil
}

//-------------------------------------------------------------------------------
// QueryAllDevices : 查询所有设备
func QueryAllDevices(session *sql.Tx) (deviceList []lib.Device, err error) {
	//Query device list
	sql := `SELECT device_id, device_type, device_uuid, device_vms_id, device_name, device_ip, rtsp_url, 
			main_stream_url, sub_stream_url, p2p_key FROM imp_t_device`

	rows, err := session.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if device, err := scanDevice(rows); err != nil {
			return nil, err
		} else {
			deviceList = append(deviceList, device)
		}
	}

	return
}

//----------------------------------------------------------------------------
// QueryDeviceGroupList : Query device group list, if group has not any devices, group will still be returned
func QueryDeviceGroupList(session *sql.Tx) (deviceGroupList []lib.Group, err error) {
	//Query group list 改
	sql := `select t3.group_id, t3.group_name, t4.device_id, t4.device_type, t4.device_uuid, t4.device_vms_id, t4.device_name, t4.device_ip,
					t4.rtsp_url, t4.main_stream_url, t4.sub_stream_url, t4.p2p_key from (select t1.group_id, t1.group_name, t2.device_id from imp_t_group t1 
					left join imp_t_groupdevice t2 on t1.group_id = t2.group_id) t3 left join imp_t_device t4 on  
					t3.device_id = t4.device_id`

	rows, err := session.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deviceGroupList = []lib.Group{}
	mapGroupIDDeviceGroup := make(map[int64]lib.Group)
	for rows.Next() {

		var groupID int64
		var groupName string
		var device DeviceInterfaceStruct
		if err = rows.Scan(&groupID, &groupName, &device.DeviceID, &device.DeviceType, &device.DeviceUUID, &device.DeviceVmsID, 
			&device.DeviceName, &device.DeviceIP, &device.RtspURL, &device.MainStreamURL, &device.SubStreamURL, &device.P2PKey); err != nil {
			return nil, err
		}

		if elem, ok := mapGroupIDDeviceGroup[groupID]; ok {
			if device.DeviceID != nil {
				elem.DeviceList = append(elem.DeviceList, assignDevice(device))
			}
			mapGroupIDDeviceGroup[groupID] = elem
		} else {
			var deviceGroup lib.Group
			deviceGroup.GroupID = groupID
			deviceGroup.GroupName = groupName
			deviceGroup.DeviceList = []lib.Device{}
			if device.DeviceID != nil {
				deviceGroup.DeviceList = append(deviceGroup.DeviceList, assignDevice(device))
			}
			mapGroupIDDeviceGroup[groupID] = deviceGroup
		}
	}

	for _, v := range mapGroupIDDeviceGroup {
		deviceGroupList = append(deviceGroupList, v)
	}
	return deviceGroupList, nil
}

//----------------------------------------------------------------------------
// InsertDeviceGroup : Insert device group info
func InsertDeviceGroup(session *sql.Tx, groupID int64, groupName string) (err error) {
	sql := "insert into imp_t_group values(?,?)"

	_, err = session.Exec(sql, groupID, groupName)
	if err != nil {
		return err
	}

	return nil
}

//-----------------------------------------------------------------------------
// QueryDeviceGroup : Query device group detail | 在mysql里面来查询设备组详情
func QueryDeviceGroup(session *sql.Tx, groupID int64) (group lib.Group, err error) {

	sql := `select t3.group_id, t3.group_name, t4.device_id, t4.device_type, t4.device_uuid, t4.device_vms_id, t4.device_name, t4.device_ip,
					t4.rtsp_url, t4.main_stream_url, t4.sub_stream_url, t4.p2p_key from (select t1.group_id, t1.group_name, t2.device_id from imp_t_group t1 
					left join imp_t_groupdevice t2 on t1.group_id = t2.group_id) t3 left join imp_t_device t4 on  
					t3.device_id = t4.device_id and t3.group_id = ?`

	rows, err := session.Query(sql, groupID)
	if err != nil {
		log.HTTP.Error(err)
		return group, err
	}
	defer rows.Close()

	group.DeviceList = []lib.Device{}
	for rows.Next() {
		var device DeviceInterfaceStruct

		if err = rows.Scan(&group.GroupID, &group.GroupName, &device.DeviceID, &device.DeviceType, &device.DeviceUUID, &device.DeviceVmsID, 
			&device.DeviceName, &device.DeviceIP, &device.RtspURL, &device.MainStreamURL, &device.SubStreamURL, &device.P2PKey); err != nil {
			return group, err
		}

		if device.DeviceID != nil {
			group.DeviceList = append(group.DeviceList, assignDevice(device))
		}
	}

	return group, nil
}

// UpdateDeviceGroupName : Update device group name | 更新设备组名
func UpdateDeviceGroupName(session *sql.Tx, groupID int64, groupName string) (err error) {
	sql := "update imp_t_group set group_name = ? where group_id = ?"

	_, err = session.Exec(sql, groupName, groupID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDeviceGroup : Delete device group | 删除设备组
func DeleteDeviceGroup(session *sql.Tx, groupID int64) (err error) {
	sql := "delete from imp_t_group where group_id = ?"

	_, err = session.Exec(sql, groupID)
	if err != nil {
		return err
	}

	sql = "delete from imp_t_groupdevice where group_id = ?"

	_, err = session.Exec(sql, groupID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDeviceFromGroup : Delete device from group | 从设备组中删除设备
func DeleteDeviceFromGroup(session *sql.Tx, groupID, deviceID int64) (err error) {
	sql := "delete from imp_t_groupdevice where group_id = ? and device_id = ?"
	_, err = session.Exec(sql, groupID, deviceID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteDevicesFromGroup : Delete devices from group | 更新设备组设备列表之删除设备
func DeleteDevicesFromGroup(session *sql.Tx, groupID int64, deviceIDList []int64) (err error) {
	deviceIDListString := lib.IntSliceToString(deviceIDList)
	sql := "delete from imp_t_groupdevice where group_id = ? and device_id in (" + deviceIDListString + ")"

	_, err = session.Exec(sql, groupID)
	if err != nil {
		return err
	}

	//strings.Split(deviceIDListString, ",")

	//需要将设备同时从设备表中一起删除
	sql2 := "delete from imp_t_device where device_id = ?"

	_, err = session.Exec(sql2, deviceIDListString)
	if err != nil {
		return err
	}
	
	return nil
}

// AddDevicesToGroup : Add devices to group | 向设备组中添加设备列表
func AddDevicesToGroup(session *sql.Tx, groupID int64, deviceIDList []int64) (err error) {

	var groupDeviceMapString string
	for _, v := range deviceIDList {
		groupDeviceMapString = groupDeviceMapString + "(" + strconv.FormatInt(groupID, 10) + "," + strconv.FormatInt(v, 10) + ")" + ","
	}
	if len(groupDeviceMapString) > 0 {
		groupDeviceMapString = groupDeviceMapString[0 : len(groupDeviceMapString)-1]
	}
	sql := "insert into imp_t_groupdevice values " + groupDeviceMapString

	_, err = session.Exec(sql)
	if err != nil {
		return err
	}

	return nil
}

// QueryNonGroupDeviceList : Query non group device list | 查询非分组的设备
func QueryNonGroupDeviceList(session *sql.Tx) (deviceList []lib.Device, err error) {

	sql := `SELECT device_id, device_type, device_uuid, device_vms_id, device_name, device_ip, rtsp_url, main_stream_url, 
			sub_stream_url, p2p_key FROM imp_t_device WHERE device_id NOT IN (select device_id from imp_t_groupdevice)`

	rows, err := session.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deviceList = []lib.Device{}
	for rows.Next() {
		if device, err := scanDevice(rows); err != nil {
			return nil, err
		} else {
			deviceList = append(deviceList, device)
		}
	}

	return deviceList, nil
}

//QueryNonGroupDeviceNumber : Query non group device number | 查询未分组的设备
func QueryNonGroupDeviceNumber(session *sql.Tx, keyword string) (number int, err error) {

	sql := `SELECT count(device_id) FROM imp_t_device where device_id NOT IN (select device_id from imp_t_groupdevice)
		and device_name like '%` + keyword + `%'`
	//like '%` + keyword + `%'  这个表示的是模糊查询
	rows, err := session.Query(sql)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if false == rows.Next() {
		return -1, nil
	}

	rows.Scan(&number)
	return number, nil
}

// QueryNonGroupDeviceListByPage : Query non group device list by page | 查询未分组的设备
func QueryNonGroupDeviceListByPage(session *sql.Tx, keyword string, offset, count int64) (deviceList []lib.Device, err error) {

	var sql string

	if keyword == "" {
		// TODO :: there is a problem
		sql = `select t1.device_id, t1.device_type, t1.device_uuid, t1.device_vms_id, t1.device_name,
						t1.device_ip, t1.rtsp_url, t1.main_stream_url, t1.sub_stream_url, t1.p2p_key
						from imp_t_device t1
						where t1.device_id not in (select device_id from imp_t_groupdevice) and t1.device_id >=
						(
						select t2.device_id	from imp_t_device t2 limit ?,1
						) limit ?`
	} else {
		sql = `select t1.device_id, t1.device_type, t1.device_uuid, t1.device_vms_id, t1.device_name,
						t1.device_ip, t1.rtsp_url, t1.main_stream_url, t1.sub_stream_url, t1.p2p_key
						from imp_t_device t1
						where  t1.device_id not in (select device_id from imp_t_groupdevice) 
						and t1.device_name like '%` + keyword + `%' and t1.device_id >=
						(
						select t2.device_id from imp_t_device t2
						where t2.device_name like '%` + keyword + `%' limit ?,1
						) limit ?`
	}

	rows, err := session.Query(sql, offset, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	deviceList = []lib.Device{}
	for rows.Next() {
		if device, err := scanDevice(rows); err != nil {
			return nil, err
		} else {
			deviceList = append(deviceList, device)
		}
	}

	return deviceList, nil
}

//--------------------------------------------------------------------------
// 根据设备的id来修改设备组名
func UpdateDeviceNameByDeviceID(session *sql.Tx, deviceId int64, deviceName string) error {
	sql := "update imp_t_device set device_name = ? where device_id = ?"

	if result, err := session.Exec(sql, deviceName, deviceId); err != nil {
		return err
	} else {
		if num, err := result.RowsAffected(); err != nil {
			return err
		} else {
			if num == 0 {
				return errors.New("SQL statement execution success, but no rows was affected")
			}
		}
	}

	return nil
}

//------------------------------------------------------------------------------------------
//在数据库中导入设备列表
func ImportDeviceListSql(session *sql.Tx, deviceList []lib.Device) error {
	if len(deviceList) == 0 {
		return errors.New("the deviceList is empty")
	}

	sql := `INSERT IGNORE INTO imp_t_device (device_type, device_name, device_ip, rtsp_url, access_account, 
			access_password) VALUES `

	for _, v := range deviceList {
		sql = sql + "(" + "1" + ",  '" +
					      v.DeviceName + "', '" +
						  v.DeviceIP + "', '" +
						  v.RtspURL + "', '" +
						  v.Account + "', '" +
						  v.Password +
					"'),"
	}

	sql = sql[:len(sql)-1]
	log.HTTP.Info(sql)

	if _, err := session.Exec(sql); err != nil {
		return err
	}

	return nil
} 

//------------------------------------------------------------------------------------------
//register information to t_device table
func RegisterInfoToDeviceTable(session *sql.Tx, info lib.RegisterLiveDevices) (status int, deviceID int64, err error) {
	sql := `INSERT IGNORE INTO imp_t_device (device_type, device_uuid, device_name, device_ip, rtsp_url, access_account, 
			access_password) VALUES (2, ?, ?, ?, ?, ?, ?)`
	account, _ := env.Get("rtsp.account").(string)
	password, _ := env.Get("rtsp.password").(string)
	rtspURL := "rtsp://" + account + ":" + 
				password + "@" + 
				info.DeviceIP + env.Get("rtsp.suffix").(string)

	if result, err := session.Exec(sql, info.DeviceUUID, info.DeviceIP, info.DeviceIP, rtspURL, account, password); err != nil {
		return lib.RegisterFailed, 0, err
	} else {
		if deviceID, err = result.LastInsertId(); err != nil {
			return lib.RegisterFailed, 0, err
		} else {
			if deviceID == 0 {
				return lib.RegisterMultiple, 0, errors.New("you had registered the device")
			}
		}
	}

	return lib.RegisterSuccess, deviceID, nil
}

//------------------------------------------------------------------------------------------
//register information to t_alg table
func RegisterInfoToAlgTable(session *sql.Tx, info lib.RegisterLiveDevices, deviceID int64) error {
	sql := `INSERT INTO imp_t_alg (device_id, device_type, alg_id, alg_config) VALUES (?, ?, ?, ?)`
	if _, err := session.Exec(sql, deviceID, 2, info.AlgID, info.AlgConfig); err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------------
func DeleteDeviceByMySQL(session *sql.Tx, deviceID int64) error {
	sql := "DELETE FROM imp_t_device WHERE device_id = ? AND device_type = 1"
	if result, err := session.Exec(sql, deviceID); err != nil {
		return err
	} else {
		if num, err := result.RowsAffected(); err != nil {
			return err
		} else {
			if num == 0 {
				return errors.New("SQL statement execution success, but no rows was affected")
			}
		}
	}

	return nil
}
