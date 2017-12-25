package sql

import (
	"SilverBusinessServer/lib"
	"SilverBusinessServer/sword"
	"SilverBusinessServer/livekeeper"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"errors"
	"strconv"
	"strings"
)

//-------------------------------------------------------------------------------
// assignDevice : Assign the device info according to the DeviceInterfaceStruct
func assignAlgDevice(dis DeviceInterfaceStruct) (device lib.AlgDevice) {
	if dis.DeviceID != nil {
		device.DeviceID = dis.DeviceID.(int64)
	} else {
		device.DeviceID = 0
	}

	if dis.DeviceType != nil {
		device.DeviceType = int(dis.DeviceType.(int64))
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

func scanAlgDevice(row *sql.Rows) (lib.AlgDevice, error) {
	var device DeviceInterfaceStruct
	var algID, algConfig interface{}
	if err := row.Scan(&device.DeviceID, &device.DeviceType, &device.DeviceUUID, &device.DeviceVmsID, &device.DeviceName, 
		&device.DeviceIP, &device.RtspURL, &device.MainStreamURL, &device.SubStreamURL, &device.P2PKey, &algID, &algConfig); err != nil {
		return lib.AlgDevice{} , err
	}
	
	algDevice := assignAlgDevice(device)
	if algID != nil {
		algDevice.AlgID = string(algID.([]uint8))
	} else {
		algDevice.AlgID = ""
	}

	if algConfig != nil {
		algDevice.AlgConfig = string(algConfig.([]uint8))
	} else {
		algDevice.AlgConfig = ""
	}

	return algDevice, nil
}

//---------------------------------------------------------------------------------------------------------------
//修改算法配置
func UpdateALGBySql(session *sql.Tx, algCon string, deviceId int64) (err error) {
	sql := "update imp_t_alg set alg_config = ? , task_id = ? where device_id = ?"
	var taskID = ""
	_, err = session.Exec(sql, algCon, taskID, deviceId)
	if err != nil {
		return err
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------
//删除设备，之前已经停止算法
func DeleteDeviceALG(session *sql.Tx, deviceID int64) (err error) {
	//删除tag记录
	sql := "delete from imp_t_alg where device_id = ?"
	_, err = session.Exec(sql, deviceID)
	if err != nil {
		return err
	}
	return nil
}

//-----------------------------------------------------------------------------------------------------------
// Query the total number of alg devices
func QueryAlgDeviceNumber(session *sql.Tx, keyword string) (int, error) {
	sql := `select count(t.device_id) from imp_t_alg t INNER JOIN imp_t_device d on t.device_id=d.device_id
	and d.device_name like '%` + keyword + `%'`

	rows, err := session.Query(sql)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if false == rows.Next() {
		return -1, errors.New("QueryAlgDeviceNumber error")
	}

	var number int
	rows.Scan(&number)
	return number, nil
}

//-----------------------------------------------------------------------------------------------------------
//查询所有算法分析摄像机
func QueryAllAlgs(session *sql.Tx, keyword string, offset, count int64) (deviceList []lib.AlgDevice, err error) {
	var sql string
	if keyword == "" {
		sql = `select  t.device_id,
		d.device_type,
		d.device_uuid,
		d.device_vms_id,
		d.device_name,
		d.device_ip,
		d.rtsp_url,
		d.main_stream_url,
		d.sub_stream_url,
		d.p2p_key,
		t.alg_id,
		t.alg_config 
		from  imp_t_alg t INNER  JOIN imp_t_device d on t.device_id=d.device_id ORDER BY device_id LIMIT ?, ?`
	} else {
		sql = `select  t.device_id,
		d.device_type,
		d.device_uuid,
		d.device_vms_id,
		d.device_name,
		d.device_ip,
		d.rtsp_url,
		d.main_stream_url,
		d.sub_stream_url,
		d.p2p_key,
		t.alg_id,
		t.alg_config 
		from  imp_t_alg t INNER  JOIN imp_t_device d on t.device_id=d.device_id and d.device_name like '%` + keyword + `%' 
		ORDER BY device_id LIMIT ?, ?`
	}
	
	rows, err := session.Query(sql, offset, count)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if algDevice, err := scanAlgDevice(rows); err != nil {
			return nil, err
		} else {
			deviceList = append(deviceList, algDevice)
		}	
	}

	if result, err := QueryAllAlgStatus(session); err != nil {
		return nil, err
	} else {
		for _, item1 := range result {
			for i, item2 := range deviceList {
				if item1.DeviceID == item2.DeviceID {
					deviceList[i].AlgStatus = item1.AlgStatus
				}
			}
		}
	}

	return deviceList, nil
}

//-----------------------------------------------------------------------------------------------
//添加摄像机
func AddDeviceToAlg(session *sql.Tx, deviceidList []int64, algID string) (err error) {
	//对deviceidList 去重
	listID := RemoveDuplicate(deviceidList)
	var TagDeviceMapString = ""
	idList, err := AlgDistinct(session)
	if err != nil {
		return err
	}
	var config string
	if algID == "aod" {
		config = `'{"sensitivity":3,"perimeter":[{"pointNum":4,"pointList":[{"x":0.1,"y":0.1},{"x":0.9,"y":0.1},{"x":0.9,"y":0.9},{"x":0.1,"y":0.9}]}]}'`
	} else if algID == "opa" {
		config = `'{"devHeight":3.0,"devDistance":8.0,"alarmCompression":0,"rule":{"perimeter":[{"pointNum":4,"pointList":[{"x":0.01,"y":0.01},{"x":0.01,"y":0.99},{"x":0.99,"y":0.99},{"x":0.99,"y":0.01}]}],"tripwire":[],"doubleTripwire":[]}}'`
	}
	
	for _, v := range listID {
		tag := 0
		for _, l := range idList {
			if l == v {
				tag = 1
			}
		}
		if tag == 0 {
			TagDeviceMapString = TagDeviceMapString + "(" + strconv.FormatInt(v, 10) + "," + "1,'" + algID + "'," + config + ")" + ","
		}
	}
	//如果为空
	if TagDeviceMapString == "" {
		return nil
	}

	sql := "INSERT into imp_t_alg (device_id, device_type, alg_id, alg_config) VALUES " + TagDeviceMapString
	strune := []rune(sql)
	sql = sql[0 : len(strune)-1]
	_, err = session.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

//数组去重
func RemoveDuplicate(list []int64) []int64 {
	var x []int64 = []int64{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

//摄像机添加中的“去重”
func AlgDistinct(session *sql.Tx) (deviceidList []int64, err error) {
	var rows *sql.Rows
	sql := `select device_id  from imp_t_alg`
	rows, err = session.Query(sql)
	if err != nil {
		rows.Close()
		return deviceidList, err
	}
	for rows.Next() {
		var deviceID interface{}
		if err = rows.Scan(&deviceID); err != nil {
			return deviceidList, err
		}
		if deviceID != nil {
			deviceIDi, _ := strconv.ParseInt(string(deviceID.([]uint8)), 10, 64)
			deviceidList = append(deviceidList, deviceIDi)
		}
	}
	return deviceidList, nil
}

//-------------------------------------------------------------------------------------------------
// query the normal device without algorithm analysis
func QueryUnmonitorDevice(session *sql.Tx) ([]lib.Device, error) {
	sql := `SELECT device_id, device_type, device_uuid, device_vms_id, device_name, rtsp_url, 
				   main_stream_url,sub_stream_url, device_ip, p2p_key
			FROM  imp_t_device WHERE device_id NOT IN (SELECT device_id FROM imp_t_alg) AND device_type=1`

	rows, err := session.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var deviceList []lib.Device
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
// Query the algorithm configuration for a particular device
func QueryAlgConfig(session *sql.Tx, deviceID int) (algID string, algConfig string, err error) {
	sql := `SELECT alg_id, alg_config FROM imp_t_alg WHERE device_id=?`
	err = session.QueryRow(sql, deviceID).Scan(&algID, &algConfig)
	return
}

//--------------------------------------------------------------------------
// Get the device type
func GetDeviceType(session *sql.Tx, deviceID int) (deviceType int, deviceUUID, algID, algPara string, err error) {
	sql := `SELECT t.device_type, t.device_uuid, d.alg_id, d.alg_config FROM imp_t_device t 
		INNER JOIN imp_t_alg d ON t.device_id=d.device_id AND t.device_id=?`
	var uuidInterface interface{}
	if err = session.QueryRow(sql, deviceID).Scan(&deviceType, &uuidInterface, &algID, &algPara); err != nil {
		return
	}
	
	if uuidInterface != nil {
		deviceUUID = string(uuidInterface.([]uint8))
	} else {
		deviceUUID = ""
	}

	return
}

//-------------------------------------------------------------------------------------------------
//启停相机算法  start-启动，stop-停止
func StartOrStopAlg(session *sql.Tx, deviceId int64, command string, alg lib.Alg, algDevice lib.AlgDevice) (err error) {
	//一：command 是 start 启动时
	if command == "start" {
		if alg.TaskID == "" {
			//如果task_id为空时 则CreateTask后StartTask
			//1.1 createTask
			var taskPara = lib.TaskPara{DeviceID: deviceId, StreamURL: algDevice.RtspURL, BeginTime: -1, EndTime: -1, AlgID: alg.AlgID, AlgConfig: algDevice.AlgConfig}
			result, err := sword.CreateTask(taskPara)
			if err != nil {
				return err
			}
			if result.Err != 0 {
				return errors.New(result.ErrMsg)
			}
			//1.2 startTask
			reuslt1, err := sword.StartTask(result.Data.TaskID)
			if err != nil {
				return err
			}
			if reuslt1.Err != 0 {
				return errors.New(reuslt1.ErrMsg)
			}
			//更新数据
			task := result.Data.TaskID
			sql1 := "update imp_t_alg set task_id = ?  where device_id = ? "
			_, err = session.Exec(sql1, task, deviceId)
			if err != nil {
				return err
			}
			return nil
		} else {
			//如果task_id不为空 则直接StartTask
			result, err := sword.StartTask(alg.TaskID)
			if err != nil {
				return err
			}
			if result.Err != 0 {
				return errors.New(result.ErrMsg)
			}
			return nil
		}
		//二：command是stop 停止时
	} else if command == "stop" {
		result, err := sword.StopTask(alg.TaskID)
		if err != nil {
			return err
		}
		if result.Err != 0 {
			return errors.New(result.ErrMsg)
		}
		return nil
	}
	return nil
}

//根据device_id 获得alg
func GetAlg(session *sql.Tx, deviceId int64) (alg lib.Alg, err error) {
	sqlStr := `SELECT device_id, device_type, alg_id, alg_config, task_id FROM imp_t_alg where device_id = ?`
	var TaskID interface{}
	if err := session.QueryRow(sqlStr, deviceId).Scan(&alg.DeviceID, &alg.DeviceType, &alg.AlgID, &alg.AlgConfig, &TaskID); err != nil {
		if err == sql.ErrNoRows {
			return alg, nil
		}
		return alg, err
	}

	if TaskID != nil {
		alg.TaskID = string(TaskID.([]uint8))
	} else {
		alg.TaskID = ""
	}

	return
}

//断言 inteface to int64
func InterfaceToInt64(inter interface{}) int64 {
	var tempInt64 int64
	if inter == nil {
		tempInt64 = 0
		return tempInt64
	}
	switch inter.(type) {
	case string:
		tempInt64, _ = strconv.ParseInt(inter.(string), 10, 64)
		break
	case int64:
		tempInt64 = inter.(int64)
		break
	case int:
		tempInt64 = inter.(int64)
		break
	}
	return tempInt64
}

//断言interface to string
func InterfaceTostring(inter interface{}) string {
	if inter == nil {
		return ""
	}
	var temp string
	switch inter.(type) {
	case string:
		temp = inter.(string)
		break
	case float64:
		temp = strconv.FormatFloat(inter.(float64), 'f', -1, 64)
		break
	case int64:
		temp = strconv.FormatInt(inter.(int64), 10)
		break
	case int:
		temp = strconv.Itoa(inter.(int))
		break
	}
	return temp
}

//根据deviceid获得algdevice信息（详见 lib.AlgDevice 数据结构）
func GetAlgDeviceByDeviceID(session *sql.Tx, deviceId int64) (algDevice lib.AlgDevice, err error) {
	sql := `select  t.device_id ,
	                d.device_name,
	                d.rtsp_url,
	                d.main_stream_url,
	                d.sub_stream_url,
	                t.task_id,
	                t.alg_config 
	        from  imp_t_alg t LEFT  JOIN imp_t_device d 
			      on t.device_id=d.device_id 
	              where t.device_id=? `
	rows, err := session.Query(sql, deviceId)
	defer rows.Close()
	if err != nil {
		return algDevice, err
	}
	var (
		deviceID      interface{}
		deviceName    interface{}
		rtspURL       interface{}
		mainStreamURL interface{}
		subStreamURL  interface{}
		taskID        interface{}
		algConfig     interface{}
	)
	if rows.Next() {
		if err = rows.Scan(&deviceID, &deviceName, &rtspURL, &mainStreamURL, &subStreamURL, &taskID, &algConfig); err != nil {
			return algDevice, err
		}
	}
	if deviceID != nil {
		//algDevice.DeviceID, _ = strconv.ParseInt(string(deviceID.([]uint8)), 10, 64)
		algDevice.DeviceID = InterfaceToInt64(deviceID)
	} else {
		algDevice.DeviceID = 0
	}
	if deviceName != nil {
		algDevice.DeviceName = string(deviceName.([]uint8))
	} else {
		algDevice.DeviceName = ""
	}
	if rtspURL != nil {
		algDevice.RtspURL = string(rtspURL.([]uint8))
	} else {
		algDevice.RtspURL = ""
	}
	if mainStreamURL != nil {
		algDevice.MainStreamURL = string(mainStreamURL.([]uint8))
	} else {
		algDevice.MainStreamURL = ""
	}
	if subStreamURL != nil {
		algDevice.SubStreamURL = string(subStreamURL.([]uint8))
	} else {
		algDevice.SubStreamURL = ""
	}
	// if taskID != nil {
	// 	algDevice.AlgStatus = -1
	// } else {
	// 	algDevice.AlgStatus = 0
	// }
	if algConfig != nil {
		algDevice.AlgConfig = string(algConfig.([]uint8))
	} else {
		algDevice.AlgConfig = ""
	}
	return algDevice, nil
}

//查询alg表获得devieid以及对应的算法运行状态algstatus
func QueryAllAlgStatus(session *sql.Tx) (result []lib.AlgStatus, err error) {
	// sql := `select t.device_id, t.device_type, t.device_uuid, t.alg_id, t.task_id from imp_t_alg t
	// 	WHERE t.device_id in (select  device_id from imp_t_device )`

	sql := `SELECT t.device_id, t.device_type, d.device_uuid, t.alg_id, t.task_id from imp_t_alg t 
		INNER JOIN imp_t_device d
		ON t.device_id =  d.device_id`
	rows, err := session.Query(sql)
	defer rows.Close()
	if err != nil {
		return result, err
	}
	for rows.Next() {
		var (
			deviceID    interface{}
			deviceType  interface{}
			deviceUUID  interface{}
			algID       interface{}
			taskID      interface{}
		)
		if err = rows.Scan(&deviceID, &deviceType, &deviceUUID, &algID, &taskID); err != nil {
			return result, err
		}
		device := lib.AlgStatus{}
		if deviceID != nil {
			device.DeviceID, _ = strconv.ParseInt(string(deviceID.([]uint8)), 10, 64)
		} else {
			device.DeviceID = 0
		}

		if algID != nil {
			device.AlgID = string(algID.([]uint8))
		} else {
			device.AlgID = ""
		}

		typeValue, _ := strconv.ParseInt(string(deviceType.([]uint8)), 10, 32)
		if taskID != nil && typeValue == 1 {
			// sword
			td := string(taskID.([]uint8))
			td = strings.Replace(td, " ", "", -1)
			if td == "" {
				device.AlgStatus = 0
			} else {
				var taskIDs []string
				taskIDs = append(taskIDs, td)
				result, err := sword.QueryTaskStatus(taskIDs)
				if result.Err != 0 || err != nil {
					//如果查询失败则赋值0
					device.AlgStatus = 0
				} else {
					var sign = 0
					for _, tas := range result.Data.TaskStatusList {
						if tas.TaskID == td {
							if tas.TaskStatus == 0 {
								device.AlgStatus = 0
							} else if tas.TaskStatus == 1 {
								device.AlgStatus = 1
							}
							sign = 1
							break
						}
					}
					//如果查不到结果则赋值为0
					if sign == 0 {
						device.AlgStatus = 0
					}
				}
			}
		} else if typeValue == 2 {
			// livekeeper
			var deviceList []string
			deviceList = append(deviceList, string(deviceUUID.([]uint8)))
			if deviceAlgStatus, err := livekeeper.QueryDeviceAlgStatus(deviceList); err != nil {
				device.AlgStatus = 0
			} else {
				for _, v := range deviceAlgStatus.Data.DevAlgStatusList {
					device.AlgStatus = v.AlgStatus
				}
			}
		} else {
			device.AlgStatus = 0
		}
		result = append(result, device)
	}

	return result, nil
}
