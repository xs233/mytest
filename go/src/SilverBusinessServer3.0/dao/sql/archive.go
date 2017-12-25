package sql

import (
	"SilverBusinessServer/lib"
	"database/sql"
)

// query favorites number
func QueryFavoritesNumber(session *sql.Tx, deviceIDListStr string, beginTime, endTime int64) (number int, err error) {
	var rows *sql.Rows
	var sqlStr string

	if "" == deviceIDListStr {
		sqlStr = `select count(id) from imp_t_favorite where image_time between ? and ?`
		rows, err = session.Query(sqlStr, beginTime, endTime)
	} else {
		sqlStr := `select count(id) from imp_t_favorite where device_id in ` + deviceIDListStr + ` and image_time between ? and ?`
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

// query favorites by page
func QueryFavoritesByPage(session *sql.Tx, deviceIDListStr string, beginTime, endTime int64, offset, count int) (favoriteList []lib.Favorite, err error) {

	var rows *sql.Rows
	var sqlStr string

	if "" == deviceIDListStr {
		sqlStr = `select t1.id, t1.device_id, t1.image_url, t1.image_time, t1.search_rule_id
						from imp_t_favorite t1 
						where t1.image_time between ? and ? and t1.image_time <=
						(
						select t2.image_time from imp_t_favorite t2 
						where t2.image_time between ? and ? order by t2.image_time desc limit ?,1
						)
						order by t1.image_time desc limit ?`
		rows, err = session.Query(sqlStr, beginTime, endTime, beginTime, endTime, offset, count)
	} else {
		sqlStr = `select t1.id, t1.device_id, t1.image_url, t1.image_time, t1.search_rule_id
						from imp_t_favorite t1 
						where t1.device_id in ` + deviceIDListStr + ` and t1.image_time between ? and ? and t1.image_time <=
						(
						select t2.image_time from imp_t_favorite t2 
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

// query favorite device
func QueryFavoriteDevice(session *sql.Tx) (groupList []lib.FavoriteGroup, err error) {

	var rows *sql.Rows

	sqlStr := `select t1.group_id, t1.group_name, t3.device_id, t3.device_name
						from imp_t_group t1, imp_t_groupdevice t2, imp_t_device t3, imp_t_favorite t4
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

	sql := `delete from imp_t_favorite  where id=?`

	_, err = session.Exec(sql, fid)
	if err != nil {
		return err
	}
	return nil
}

//AddFavorite : Add favorite
func AddFavorite(session *sql.Tx, favorite lib.Favorite) (err error) {

	sql := `insert into imp_t_favorite values (?, ?, ?, ?, ?)`

	_, err = session.Exec(sql, favorite.FavoriteID, favorite.DeviceID, favorite.ImageURL,
		favorite.ImageTime, favorite.SearchRuleID)
	if err != nil {
		return err
	}

	return nil
}
