package dao

import (
	MySQL "SilverBusinessServer/dao/sql"
	"SilverBusinessServer/env"
	"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
	"database/sql"
)

// QueryAllDevices :
func QueryAllDevices() (deviceList []lib.Device, err error) {
	log.HTTP.Info("获取所有设备")
	session, err := MySQL.DB.Begin()
	if err != nil {
		return nil, err
	}

	deviceList, err = MySQL.QueryAllDevices(session)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	session.Commit()
	return deviceList, nil
}

// QueryDeviceGroupListByUserID : Query device group list
//return deviceGroupList : device group list
//return err : error info
func QueryDeviceGroupList() (deviceGroupList []lib.Group, err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return nil, err
	}

	deviceGroupList, err = MySQL.QueryDeviceGroupList(session)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	session.Commit()
	return deviceGroupList, nil
}

// CreateDeviceGroup : Create device group
//param userID : user id
//param groupName : group name
//param faceGroupID : face group id
//return groupID : group id
//return err : error info
func CreateDeviceGroup(groupName string) (groupID int64, err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return -1, err
	}

	//Get device group id
	next, err := MySQL.GetSequenceIDByName(env.SequenceNameGroupID)
	if err != nil {
		session.Rollback()
		return 0, err
	}
	groupID = next

	//Insert device group info
	err = MySQL.InsertDeviceGroup(session, groupID, groupName)
	if err != nil {
		session.Rollback()
		return 0, err
	}

	session.Commit()
	return groupID, nil
}

// QueryDeviceGroup : Query device group detail | 查询设备组详情
//param groupID : group id
//return group : device group info
//return err : error info
func QueryDeviceGroup(groupID int64) (group lib.Group, err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return group, err
	}

	group, err = MySQL.QueryDeviceGroup(session, groupID)
	if err != nil {
		session.Rollback()
		return group, err
	}

	session.Commit()
	return group, nil
}

// UpdateDeviceGroupName : Update device group name | 更新重新编辑设备组的组名
//param groupID : group id
//param groupName : group name, not to update device group name if param is empty string
//return err : error info
func UpdateDeviceGroupName(groupID int64, groupName string) (err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return err
	}

	err = MySQL.UpdateDeviceGroupName(session, groupID, groupName)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

// DeleteDeviceGroup : Delete device group | 删除设备组
//param groupID : group id
//return err : error info
func DeleteDeviceGroup(groupID int64) (err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return err
	}

	err = MySQL.DeleteDeviceGroup(session, groupID)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

// DeleteDeviceFromGroup : Delete device from group | 从设备组中删除设备
func DeleteDeviceFromGroup(groupID, deviceID int64) (err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return err
	}

	err = MySQL.DeleteDeviceFromGroup(session, groupID, deviceID)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

// UpdateGroupDevices : Update group devices group | 更新设备组设备列表
//param groupID : device group id
//param addDeviceIDList : add device id list
//param deleteDeviceIDList : delete device id list
//return err : error info
func UpdateGroupDevices(groupID int64, addDeviceIDList, deleteDeviceIDList []int64) (err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return err
	}

	log.HTTP.Info("参数信息：", addDeviceIDList, deleteDeviceIDList)

	if len(deleteDeviceIDList) > 0 {
		err = MySQL.DeleteDevicesFromGroup(session, groupID, deleteDeviceIDList)
		if err != nil {
			log.HTTP.Error("删除错误")
			session.Rollback()
			return err
		}
	}

	if len(addDeviceIDList) > 0 {
		err = MySQL.AddDevicesToGroup(session, groupID, addDeviceIDList)
		if err != nil {
			log.HTTP.Error("添加错误")
			session.Rollback()
			return err
		}
	}

	session.Commit()
	return nil
}

// AddGroupDevices : Add devices to group | 向设备组中添加设备列表
//param groupID : device group id
//param deviceIDList : add device id list
//return err : error info
func AddGroupDevices(groupID int64, deviceIDList []int64) (err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return err
	}

	if len(deviceIDList) > 0 {
		err = MySQL.AddDevicesToGroup(session, groupID, deviceIDList)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	session.Commit()
	return nil
}

// DeleteGroupDevices : Delete devices to group | 从设备组中删除设备
//param groupID : device group id
//param deviceIDList : delete device id list
//return err : error info
func DeleteGroupDevices(groupID int64, deviceIDList []int64) (err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return err
	}

	if len(deviceIDList) > 0 {
		err = MySQL.DeleteDevicesFromGroup(session, groupID, deviceIDList)
		if err != nil {
			session.Rollback()
			return err
		}
	}

	session.Commit()
	return nil
}

// QueryNonGroupDeviceList : Query non group device list | 未分组设备
//return deviceList : device list
//return err : error info
func QueryNonGroupDeviceList() (deviceList []lib.Device, err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return nil, err
	}

	deviceList, err = MySQL.QueryNonGroupDeviceList(session)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	session.Commit()
	return deviceList, nil
}

//QueryNonGroupDeviceNumber : Query non group device number | 查询未分组的设备
func QueryNonGroupDeviceNumber(keyword string) (number int, err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return -1, err
	}

	number, err = MySQL.QueryNonGroupDeviceNumber(session, keyword)
	if err != nil {
		session.Rollback()
		return -1, err
	}

	session.Commit()
	return number, nil
}

// QueryNonGroupDeviceListByPage : Query non group device list by page | 查询未分组的设备
//return deviceList : device list
//return err : error info
func QueryNonGroupDeviceListByPage(keyword string, offset, count int64) (deviceList []lib.Device, err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return nil, err
	}

	deviceList, err = MySQL.QueryNonGroupDeviceListByPage(session, keyword, offset, count)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	session.Commit()
	return deviceList, nil
}

//查询喜欢的数量
func QueryFavoritesNumber(session *sql.Tx, deviceIDListStr string, beginTime, endTime int64) (number int, err error) {
	var rows *sql.Rows
	var sqlStr string

	if "" == deviceIDListStr {
		sqlStr = `select count(id) from cp_t_favorite where image_time between ? and ?`
		rows, err = session.Query(sqlStr, beginTime, endTime)
	} else {
		sqlStr := `select count(id) from cp_t_favorite where device_id in ` + deviceIDListStr + ` and image_time between ? and ?`
		rows, err = session.Query(sqlStr, beginTime, endTime)
	}
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

// QueryFavoritesByPage :
func QueryFavoritesByPage(session *sql.Tx, deviceIDListStr string, beginTime, endTime int64, offset, count int) (favoriteList []lib.Favorite, err error) {

	var rows *sql.Rows
	var sqlStr string

	if "" == deviceIDListStr {
		sqlStr = `select t1.id, t1.device_id, t1.image_url, t1.image_time, t1.search_rule_id
						from cp_t_favorite t1 
						where t1.image_time between ? and ? and t1.image_time <=
						(
						select t2.image_time from cp_t_favorite t2 
						where t2.image_time between ? and ? order by t2.image_time desc limit ?,1
						)
						order by t1.image_time desc limit ?`
		rows, err = session.Query(sqlStr, beginTime, endTime, beginTime, endTime, offset, count)
	} else {
		sqlStr = `select t1.id, t1.device_id, t1.image_url, t1.image_time, t1.search_rule_id
						from cp_t_favorite t1 
						where t1.device_id in ` + deviceIDListStr + ` and t1.image_time between ? and ? and t1.image_time <=
						(
						select t2.image_time from cp_t_favorite t2 
						where t2.device_id in ` + deviceIDListStr + ` and t2.image_time between ? and ? order by t2.image_time desc limit ?,1
						)
						order by t1.image_time desc limit ?`

		rows, err = session.Query(sqlStr, beginTime, endTime, beginTime, endTime, offset, count)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	favoriteList = []lib.Favorite{}
	for rows.Next() {
		favorite := lib.Favorite{}
		if err = rows.Scan(&favorite.FavoriteID, &favorite.DeviceID, &favorite.ImageURL, &favorite.ImageTime, &favorite.SearchRuleID); err != nil {
			return nil, err
		}
		favoriteList = append(favoriteList, favorite)
	}

	return favoriteList, nil
}

// QueryFavoriteDevice :
func QueryFavoriteDevice(session *sql.Tx) (groupList []lib.FavoriteGroup, err error) {

	var rows *sql.Rows

	sqlStr := `select t1.group_id, t1.group_name, t3.device_id, t3.device_name
						from cp_t_group t1, cp_t_groupdevice t2, cp_t_device t3, cp_t_favorite t4
						where t1.group_id = t2.group_id and t2.device_id = t3.device_id and t2.device_id = t4.device_id`
	rows, err = session.Query(sqlStr)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mapIDGroup := make(map[int64]lib.FavoriteGroup)

	for rows.Next() {
		var (
			groupID    int64
			groupName  string
			deviceID   int64
			deviceName string
		)
		if err = rows.Scan(&groupID, &groupName, &deviceID, &deviceName); err != nil {
			return nil, err
		}

		if v, ok := mapIDGroup[groupID]; ok {
			deviceGroup := lib.DeviceGroup{deviceID, deviceName}
			v.DeviceList = append(v.DeviceList, deviceGroup)
		} else {
			favoriteGroup := lib.FavoriteGroup{
				Id:   groupID,
				Name: groupName,
				DeviceList: []lib.DeviceGroup{
					lib.DeviceGroup{
						Id:   deviceID,
						Name: deviceName,
					},
				},
			}
			mapIDGroup[groupID] = favoriteGroup
		}
	}

	groupList = []lib.FavoriteGroup{}
	for _, v := range mapIDGroup {
		groupList = append(groupList, v)
	}
	return groupList, nil
}

//DeleteFavorite : delete favorite
func DeleteFavorite(session *sql.Tx, fid int64) (err error) {

	sql := `delete from cp_t_favorite  where id=?`

	_, err = session.Exec(sql, fid)
	if err != nil {
		return err
	}
	return nil
}

//AddFavorite : Add favorite
func AddFavorite(session *sql.Tx, favorite lib.Favorite) (err error) {

	sql := `insert into cp_t_favorite values (?, ?, ?, ?, ?)`

	_, err = session.Exec(sql, favorite.FavoriteID, favorite.DeviceID, favorite.ImageURL,
		favorite.ImageTime, favorite.SearchRuleID)
	if err != nil {
		return err
	}

	return nil
}

//------------------------------------------------------------------------------------------
//导入设备列表到数据库中
func ImportDeviceToList(deviceList []lib.Device) (err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return err
	}

	err = MySQL.ImportDeviceListSql(session, deviceList)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//----------------------------------------------------------------------------
//删除设备
func DeleteDevice(deviceID int64) (err error) {
	session, err := MySQL.DB.Begin()
	if err != nil {
		return err
	}

	err = MySQL.DeleteDeviceByMySQL(session, deviceID)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}
