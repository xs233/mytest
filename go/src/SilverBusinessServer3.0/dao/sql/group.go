package sql

import (
	"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
	"database/sql"
	"strconv"
	//"strings"
)

//-------------------------------------------------------------------------------
// QueryAllDevices : 查询所有设备
func QueryAllDevices(session *sql.Tx) (deviceList []lib.Device, err error) {
	//Query device list
	sql := `select device_id, device_vms_id, device_name, device_ip, rtsp_url, main_stream_url,
					sub_stream_url from imp_t_device`

	rows, err := session.Query(sql)
	if err != nil {
		log.HTTP.Error("创建sql失败")
		return nil, err
	}
	defer rows.Close()

	deviceList = []lib.Device{}
	for rows.Next() {

		var (
			deviceID      interface{}
			deviceVmsID   interface{}
			deviceName    interface{}
			deviceIP      interface{}
			rtspURL       interface{}
			mainStreamURL interface{}
			subStreamURL  interface{}
		)

		device := lib.Device{}
		if err = rows.Scan(&deviceID, &deviceVmsID, &deviceName, &deviceIP, &rtspURL, &mainStreamURL, &subStreamURL); err != nil {
			log.HTTP.Error("获取设备表信息失败")
			return nil, err
		}

		if deviceID != nil {
			device.DeviceID, _ = strconv.ParseInt(string(deviceID.([]uint8)), 10, 64)
		} else {
			device.DeviceID = 0
		}

		if deviceVmsID != nil {
			device.DeviceVmsID = string(deviceVmsID.([]uint8))
		} else {
			device.DeviceVmsID = ""
		}

		if deviceName != nil {
			device.DeviceName = string(deviceName.([]uint8))
		} else {
			device.DeviceName = ""
		}

		if deviceIP != nil {
			device.DeviceIP = string(deviceIP.([]uint8))
		} else {
			device.DeviceIP = ""
		}

		if rtspURL != nil {
			device.RtspURL = string(rtspURL.([]uint8))
		} else {
			device.RtspURL = ""
		}

		if mainStreamURL != nil {
			device.MainStreamURL = string(mainStreamURL.([]uint8))
		} else {
			device.MainStreamURL = ""
		}

		if subStreamURL != nil {
			device.SubStreamURL = string(subStreamURL.([]uint8))
		} else {
			device.SubStreamURL = ""
		}

		deviceList = append(deviceList, device)
	}

	return deviceList, nil
}

//----------------------------------------------------------------------------
// QueryDeviceGroupList : Query device group list, if group has not any devices, group will still be returned
func QueryDeviceGroupList(session *sql.Tx) (deviceGroupList []lib.Group, err error) {
	//Query group list 改
	sql := `select t3.group_id, t3.group_name, t4.device_id, t4.device_vms_id, t4.device_name, t4.device_ip,
					t4.rtsp_url, t4.main_stream_url, t4.sub_stream_url from (select t1.group_id, t1.group_name, t2.device_id from imp_t_group t1 
					left join imp_t_groupdevice t2 on t1.group_id = t2.group_id) t3 left join imp_t_device t4 on  
					t3.device_id = t4.device_id`

	rows, err := session.Query(sql)
	if err != nil {
		log.HTTP.Error("sql语句错误")
		return nil, err
	}
	defer rows.Close()

	deviceGroupList = []lib.Group{}
	mapGroupIDDeviceGroup := make(map[int64]lib.Group)
	for rows.Next() {

		var groupID int64
		var groupName string
		var (
			deviceID      interface{}
			deviceVmsID   interface{}
			deviceName    interface{}
			deviceIP      interface{}
			rtspURL       interface{}
			mainStreamURL interface{}
			subStreamURL  interface{}
		)

		if err = rows.Scan(&groupID, &groupName, &deviceID, &deviceVmsID, &deviceName,
			&deviceIP, &rtspURL, &mainStreamURL, &subStreamURL); err != nil {
			return nil, err
		}

		device := lib.Device{}
		if deviceID != nil {
			device.DeviceID, _ = strconv.ParseInt(string(deviceID.([]uint8)), 10, 64)

			if deviceVmsID != nil {
				device.DeviceVmsID = string(deviceVmsID.([]uint8))
			} else {
				device.DeviceVmsID = ""
			}

			if deviceName != nil {
				device.DeviceName = string(deviceName.([]uint8))
			} else {
				device.DeviceName = ""
			}

			if deviceIP != nil {
				device.DeviceIP = string(deviceIP.([]uint8))
			} else {
				device.DeviceIP = ""
			}

			if rtspURL != nil {
				device.RtspURL = string(rtspURL.([]uint8))
			} else {
				device.RtspURL = ""
			}

			if mainStreamURL != nil {
				device.MainStreamURL = string(mainStreamURL.([]uint8))
			} else {
				device.MainStreamURL = ""
			}

			if subStreamURL != nil {
				device.SubStreamURL = string(subStreamURL.([]uint8))
			} else {
				device.SubStreamURL = ""
			}
		}

		if elem, ok := mapGroupIDDeviceGroup[groupID]; ok {
			if deviceID != nil {
				elem.DeviceList = append(elem.DeviceList, device)
			}
			mapGroupIDDeviceGroup[groupID] = elem
		} else {
			var deviceGroup lib.Group
			deviceGroup.GroupID = groupID
			deviceGroup.GroupName = groupName
			deviceGroup.DeviceList = []lib.Device{}
			if deviceID != nil {
				deviceGroup.DeviceList = append(deviceGroup.DeviceList, device)
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

	sql := `select t3.group_id, t3.group_name, t4.device_id, t4.device_vms_id, t4.device_name, t4.device_ip,
					t4.rtsp_url, t4.main_stream_url, t4.sub_stream_url from (select t1.group_id, t1.group_name, t2.device_id from imp_t_group t1 
					left join imp_t_groupdevice t2 on t1.group_id = t2.group_id) t3 left join imp_t_device t4 on  
					t3.device_id = t4.device_id and t3.group_id = ?`

	rows, err := session.Query(sql, groupID)
	if err != nil {
		log.HTTP.Error("sql语句错误")
		return group, err
	}
	defer rows.Close()

	group.DeviceList = []lib.Device{}
	for rows.Next() {
		var (
			deviceID      interface{}
			deviceVmsID   interface{}
			deviceName    interface{}
			deviceIP      interface{}
			rtspURL       interface{}
			mainStreamURL interface{}
			subStreamURL  interface{}
		)

		if err = rows.Scan(&group.GroupID, &group.GroupName, &deviceID, &deviceVmsID, &deviceName,
			&deviceIP, &rtspURL, &mainStreamURL, &subStreamURL); err != nil {
			return group, err
		}

		if deviceID != nil {
			device := lib.Device{}
			device.DeviceID = deviceID.(int64)
			if deviceVmsID != nil {
				device.DeviceVmsID = string(deviceVmsID.([]uint8))
			} else {
				device.DeviceVmsID = ""
			}

			if deviceName != nil {
				device.DeviceName = string(deviceName.([]uint8))
			} else {
				device.DeviceName = ""
			}

			if deviceIP != nil {
				device.DeviceIP = string(deviceIP.([]uint8))
			} else {
				device.DeviceIP = ""
			}

			if rtspURL != nil {
				device.RtspURL = string(rtspURL.([]uint8))
			} else {
				device.RtspURL = ""
			}

			if mainStreamURL != nil {
				device.MainStreamURL = string(mainStreamURL.([]uint8))
			} else {
				device.MainStreamURL = ""
			}

			if subStreamURL != nil {
				device.SubStreamURL = string(subStreamURL.([]uint8))
			} else {
				device.SubStreamURL = ""
			}

			group.DeviceList = append(group.DeviceList, device)
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
	log.HTTP.Info("测试删除数据：", groupID, deviceID)
	sql := "delete from imp_t_groupdevice where group_id = ? and device_id = ?"
	_, err = session.Exec(sql, groupID, deviceID)
	log.HTTP.Info("error:", err)
	if err != nil {
		log.HTTP.Error("从设备组中删除一个设备失败")
		return err
	}

	////需要将设备同时从设备表中一起删除（李阳说的）(10月10，目前给注释了)
	//sql2 := "delete from imp_t_device where device_id = ?"
	//_, err = session.Exec(sql2, deviceID)
	//if err != nil {
	//	log.HTTP.Error("从设备列表中删除一个设备失败")
	//	return err
	//}

	return nil
}

// DeleteDevicesFromGroup : Delete devices from group | 更新设备组设备列表之删除设备
func DeleteDevicesFromGroup(session *sql.Tx, groupID int64, deviceIDList []int64) (err error) {
	deviceIDListString := lib.IntSliceToString(deviceIDList)
	sql := "delete from imp_t_groupdevice where group_id = ? and device_id in (" + deviceIDListString + ")"

	_, err = session.Exec(sql, groupID)
	if err != nil {
		log.HTTP.Error("删除sql语句错误")
		return err
	}

	//strings.Split(deviceIDListString, ",")

	log.HTTP.Info("将id列表打印出来:", deviceIDListString, deviceIDList)
	//需要将设备同时从设备表中一起删除
	sql2 := "delete from imp_t_device where device_id = ?"

	_, err = session.Exec(sql2, deviceIDListString)
	if err != nil {
		log.HTTP.Error("从设备列表中删除一个设备失败")
		return err
	}
	log.HTTP.Info("测试从组中删除相机成功啦")
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
		log.HTTP.Error("添加sql语句错误")
		return err
	}

	return nil
}

// QueryNonGroupDeviceList : Query non group device list | 查询非分组的设备
func QueryNonGroupDeviceList(session *sql.Tx) (deviceList []lib.Device, err error) {

	sql := `select device_id, device_vms_id, device_name, device_ip, rtsp_url, main_stream_url, sub_stream_url  
		    from imp_t_device where device_id not in (select device_id from imp_t_groupdevice)`

	rows, err := session.Query(sql)
	if err != nil {
		log.HTTP.Error("sql语句错误")
		return nil, err
	}
	defer rows.Close()

	deviceList = []lib.Device{}
	for rows.Next() {
		var (
			deviceID      interface{}
			deviceVmsID   interface{}
			deviceName    interface{}
			deviceIP      interface{}
			rtspURL       interface{}
			mainStreamURL interface{}
			subStreamURL  interface{}
		)

		device := lib.Device{}
		if err = rows.Scan(&deviceID, &deviceVmsID, &deviceName,
			&deviceIP, &rtspURL, &mainStreamURL, &subStreamURL); err != nil {
			return nil, err
		}
		device.DeviceID, _ = strconv.ParseInt(string(deviceID.([]uint8)), 10, 64)

		if deviceVmsID != nil {
			device.DeviceVmsID = string(deviceVmsID.([]uint8))
		} else {
			device.DeviceVmsID = ""
		}

		if deviceName != nil {
			device.DeviceName = string(deviceName.([]uint8))
		} else {
			device.DeviceName = ""
		}

		if deviceIP != nil {
			device.DeviceIP = string(deviceIP.([]uint8))
		} else {
			device.DeviceIP = ""
		}

		if rtspURL != nil {
			device.RtspURL = string(rtspURL.([]uint8))
		} else {
			device.RtspURL = ""
		}

		if mainStreamURL != nil {
			device.MainStreamURL = string(mainStreamURL.([]uint8))
		} else {
			device.MainStreamURL = ""
		}

		if subStreamURL != nil {
			device.SubStreamURL = string(subStreamURL.([]uint8))
		} else {
			device.SubStreamURL = ""
		}

		deviceList = append(deviceList, device)
	}

	return deviceList, nil
}

//QueryNonGroupDeviceNumber : Query non group device number | 查询未分组的设备
func QueryNonGroupDeviceNumber(session *sql.Tx, keyword string) (number int, err error) {

	sql := `select count(device_id) from imp_t_device where device_id not in (select device_id from imp_t_groupdevice)
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
		sql = `select t1.device_id, t1.device_vms_id, t1.device_name,
						t1.device_ip, t1.rtsp_url, t1.main_stream_url, t1.sub_stream_url
						from imp_t_device t1
						where t1.device_id not in (select device_id from imp_t_groupdevice) and t1.device_id >=
						(
						select t2.device_id	from imp_t_device t2 limit ?,1
						) limit ?`
	} else {
		sql = `select t1.device_id, t1.device_vms_id, t1.device_name,
						t1.device_ip, t1.rtsp_url, t1.main_stream_url, t1.sub_stream_url
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
		log.HTTP.Error("sql语句错误")
		return nil, err
	}
	defer rows.Close()

	deviceList = []lib.Device{}
	for rows.Next() {
		var (
			deviceID      interface{}
			deviceVmsID   interface{}
			deviceName    interface{}
			deviceIP      interface{}
			rtspURL       interface{}
			mainStreamURL interface{}
			subStreamURL  interface{}
		)

		device := lib.Device{}
		if err = rows.Scan(&deviceID, &deviceVmsID, &deviceName,
			&deviceIP, &rtspURL, &mainStreamURL, &subStreamURL); err != nil {
			return nil, err
		}

		if deviceVmsID != nil {
			device.DeviceVmsID = string(deviceVmsID.([]uint8))
		} else {
			device.DeviceVmsID = ""
		}

		if deviceName != nil {
			device.DeviceName = string(deviceName.([]uint8))
		} else {
			device.DeviceName = ""
		}

		if deviceIP != nil {
			device.DeviceIP = string(deviceIP.([]uint8))
		} else {
			device.DeviceIP = ""
		}

		if rtspURL != nil {
			device.RtspURL = string(rtspURL.([]uint8))
		} else {
			device.RtspURL = ""
		}

		if mainStreamURL != nil {
			device.MainStreamURL = string(mainStreamURL.([]uint8))
		} else {
			device.MainStreamURL = ""
		}

		if subStreamURL != nil {
			device.SubStreamURL = string(subStreamURL.([]uint8))
		} else {
			device.SubStreamURL = ""
		}

		deviceList = append(deviceList, device)
	}
	log.HTTP.Info("查询未分组设备成功")
	return deviceList, nil
}

//--------------------------------------------------------------------------
// 根据设备的id来修改设备组名
func UpdateDeviceNameByDeviceID(session *sql.Tx, deviceId int64, deviceName string) (err error) {
	sql := "update imp_t_device set device_name = ? where device_id = ?"

	_, err = session.Exec(sql, deviceName, deviceId)
	if err != nil {
		return err
	}
	log.HTTP.Info("根据设备的id来修改设备组名成功")
	return nil
}

//------------------------------------------------------------------------------------------
//在数据库中导入设备列表
func ImportDeviceListSql(session *sql.Tx, deviceList []lib.Device) (err error) {

	sql := "insert into imp_t_device(device_vms_id,device_name,device_ip,rtsp_url,main_stream_url,sub_stream_url,access_account,access_password) values"

	sqlSel := "select device_ip, rtsp_url from imp_t_device"
	rows, err := session.Query(sqlSel)
	if err != nil {
		log.HTTP.Error("创建sql失败")
		return err
	}
	defer rows.Close()

	var deviceListH = []lib.Device{}
	for rows.Next() {
		var (
			deviceIP interface{}
			rtspURL  interface{}
		)
		if err = rows.Scan(&deviceIP, &rtspURL); err != nil {
			log.HTTP.Error("获取设备表信息失败")
			return err
		}
		device := lib.Device{}
		if deviceIP != nil {
			device.DeviceIP = string(deviceIP.([]uint8))
		} else {
			device.DeviceIP = ""
		}
		if rtspURL != nil {
			device.RtspURL = string(rtspURL.([]uint8))
		} else {
			device.RtspURL = ""
		}

		deviceListH = append(deviceListH, device)
	}

	var str string
	for _, v := range deviceList {
		var fl bool
		for _, vv := range deviceListH {
			if v.DeviceIP == vv.DeviceIP && v.RtspURL == v.RtspURL {
				fl = true
				break
			}
		}
		if fl {
			continue
		}
		str += "('" + "','" + v.DeviceName + "','" + v.DeviceIP + "','" + v.RtspURL + "','" + "','" + "','" + "','" + "'),"

	}
	var err2 error
	if len(str) < 1 {
		log.HTTP.Error("str是否为空：", len(str))
		return err2
	}
	str = str[0 : len(str)-1]
	sql += str
	_, err = session.Exec(sql)
	if err != nil {
		return err
	}
	log.HTTP.Info("从数据库中导入设备列表成功")
	return nil
}

//----------------------------------------------------------------------------
func DeleteDeviceByMySQL(session *sql.Tx, deviceid int64) (err error) {
	//1.查询所有设备组的表，循环从每一个表中便利，查询是否有这个deviceID的设备，有的话，获取设备组的groupid，然后从设备组中删除，调用DeleteDeviceFromGroup里面的内容
	//sql := "select * from imp_t_groupdevice"
	//rows, err := session.Query(sql)
	//if err != nil {
	//	log.HTTP.Info("查询设备组失败")
	//	return err
	//}
	//defer rows.Close()
	//
	//for rows.Next() {
	//}
	//sql := "delete from imp_t_groupdevice where device_id = ?"

	//2.从设备表中，将该设备删除
	sql1 := "delete from imp_t_device where device_id = ?"
	_, err = session.Exec(sql1, deviceid)
	if err != nil {
		log.HTTP.Error("删除一个设备失败")
		return err
	}
	log.HTTP.Info("删除设备成功")
	return nil
}
