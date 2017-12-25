package sql

import (
	"SilverBusinessServer/lib"
	"SilverBusinessServer/sword"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// QueryDeviceNumber : Query device number
func QueryDeviceNumber(session *sql.Tx, keyword string) (number int, err error) {
	sql := `select count(device_id) from imp_t_device where device_name like '%` + keyword + `%'`

	rows, err := session.Query(sql)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(&number); err != nil {
			return -1, err
		}
		return number, nil
	}
	return -1, nil
}

// QueryDeviceListByPage : Query device list by page
func QueryDeviceListByPage(session *sql.Tx, keyword string, offset, count int64) (devices []lib.Device, err error) {

	var sql string

	if keyword == "" {
		sql = `select t1.device_id, t1.device_uuid, t1.device_name, 
						t1.main_stream_url, t1.sub_stream_url
						from imp_t_device t1
						where t1.device_id >=
						(
						select t2.device_id	from imp_t_device t2 limit ?,1
						) limit ?`
	} else {
		sql = `select t1.device_id, t1.origin_id, t1.device_name,
						t1.main_stream_url, t1.sub_stream_ur
						from imp_t_device t1
						where  t1.device_name like '%` + keyword + `%' and t1.camera_id >=
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

	devices = []lib.Device{}
	for rows.Next() {
		device := lib.Device{}

		if err = rows.Scan(&device.DeviceID, &device.DeviceVmsID, &device.DeviceName,
			&device.DeviceIP, &device.RtspURL, &device.MainStreamURL, &device.SubStreamURL); err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}

	return devices, nil
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
//查询所有算法分析摄像机
func QueryAllAlgs(session *sql.Tx) (deviceList []lib.AlgDevice, err error) {
	sql := `select  t.device_id ,
	                d.device_name,
	                d.rtsp_url,
	                d.main_stream_url,
	                d.sub_stream_url,
					d.device_ip,
	                t.task_id,
	                t.alg_config 
	        from  imp_t_alg t INNER  JOIN imp_t_device d on t.device_id=d.device_id `
	rows, err := session.Query(sql)
	defer rows.Close()
	if err != nil {
		fmt.Println("Query data failure.")
		return nil, err
	}
	for rows.Next() {
		var (
			deviceID      interface{}
			deviceName    interface{}
			rtspURL       interface{}
			mainStreamURL interface{}
			subStreamURL  interface{}
			deviceIP      interface{}
			taskID        interface{}
			algConfig     interface{}
		)

		if err = rows.Scan(&deviceID, &deviceName, &rtspURL, &mainStreamURL, &subStreamURL, &deviceIP, &taskID, &algConfig); err != nil {
			fmt.Println("获取设备表信息失败")
			return deviceList, err
		}
		device := lib.AlgDevice{}
		if deviceID != nil {
			device.DeviceID, _ = strconv.ParseInt(string(deviceID.([]uint8)), 10, 64)
		} else {
			device.DeviceID = 0
		}

		if deviceName != nil {
			device.DeviceName = string(deviceName.([]uint8))
		} else {
			device.DeviceName = ""
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
		if deviceIP != nil {
			device.DeviceIP = string(deviceIP.([]uint8))
		} else {
			device.DeviceIP = ""
		}
		if taskID != nil {
			//发送请求获得任务状态结果
			td := string(taskID.([]uint8))
			var taskIDs []string
			taskIDs = append(taskIDs, td)
			result, err := sword.QueryTaskStatus(taskIDs)
			if result.Err != 0 || err != nil {
				//查询失败 默认其状态为0
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
		} else {
			device.AlgStatus = 0
		}
		if algConfig != nil {
			device.AlgConfig = string(algConfig.([]uint8))
		} else {
			device.AlgConfig = ""
		}

		deviceList = append(deviceList, device)
	}
	fmt.Println("查询所有算法分析摄像机")
	return deviceList, nil
}

//-----------------------------------------------------------------------------------------------
//添加摄像机
func AddDeviceToAlg(session *sql.Tx, deviceidList []int64) (err error) {
	//对deviceidList 去重
	listID := RemoveDuplicate(deviceidList)
	var TagDeviceMapString = ""
	idList, err := AlgDistinct(session)
	if err != nil {
		return err
	}
	var config string
	config = `'{"sensitivity":3,"perimeter":[{"pointNum":4,"pointList":[{"x":0.1,"y":0.1},{"x":0.9,"y":0.1},{"x":0.9,"y":0.9},{"x":0.1,"y":0.9}]}]}'`
	for _, v := range listID {
		tag := 0
		for _, l := range idList {
			if l == v {
				tag = 1
			}
		}
		if tag == 0 {
			TagDeviceMapString = TagDeviceMapString + "(" + strconv.FormatInt(v, 10) + " ," + config + ")" + ","
		}
	}
	//如果为空
	if TagDeviceMapString == "" {
		return nil
	}

	sql := "INSERT into imp_t_alg (device_id,alg_config) VALUES " + TagDeviceMapString
	strune := []rune(sql)
	sql = sql[0 : len(strune)-1]
	fmt.Println(sql)
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
//启停相机算法  start-启动，stop-停止
func StartOrStopAlg(session *sql.Tx, deviceId int64, command string, alg lib.Alg, algDevice lib.AlgDevice) (err error) {
	//一：command 是 start 启动时
	if command == "start" {
		if alg.TaskID == "" {
			//如果task_id为空时 则CreateTask后StartTask
			//1.1 createTask
			var taskPara = lib.TaskPara{DeviceID: deviceId, StreamURL: algDevice.RtspURL, BeginTime: -1, EndTime: -1, AlgConfig: algDevice.AlgConfig}
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
	var rows *sql.Rows
	sql := `select device_id,task_id,alg_config from imp_t_alg where device_id = ?`
	rows, err = session.Query(sql, deviceId)
	defer rows.Close()
	if err != nil {
		return alg, err
	}
	var (
		deviceID  interface{}
		taskID    interface{}
		algConfig interface{}
	)
	if rows.Next() {
		if err = rows.Scan(&deviceID, &taskID, &algConfig); err != nil {
			fmt.Println(err)
			return alg, err
		}
	}
	alg.DeviceID = InterfaceToInt64(deviceID)
	if taskID != nil {
		//去空格
		alg.TaskID = strings.Replace(string(taskID.([]uint8)), " ", "", -1)
	} else {
		alg.TaskID = ""
	}
	if algConfig != nil {
		alg.AlgConfig = string(algConfig.([]uint8))
	} else {
		alg.AlgConfig = ""
	}
	//fmt.Println(alg)
	return alg, nil
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
		fmt.Println("Query data failure.")
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
			fmt.Println("Query data failure.")
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
	if taskID != nil {
		algDevice.AlgStatus = -1
	} else {
		algDevice.AlgStatus = 0
	}
	if algConfig != nil {
		algDevice.AlgConfig = string(algConfig.([]uint8))
	} else {
		algDevice.AlgConfig = ""
	}
	return algDevice, nil
}

//查询alg表获得devieid以及对应的算法运行状态algstatus
func QueryAllAlgStatus(session *sql.Tx) (result []lib.AlgStatus, err error) {
	sql := `select t.device_id , t.task_id from imp_t_alg t  WHERE t.device_id in (select  device_id from imp_t_device )`
	rows, err := session.Query(sql)
	defer rows.Close()
	if err != nil {
		fmt.Println("Query data failure.")
		return result, err
	}
	for rows.Next() {
		var (
			deviceID interface{}
			taskID   interface{}
		)
		if err = rows.Scan(&deviceID, &taskID); err != nil {
			fmt.Println("获取设备表信息失败")
			return result, err
		}
		device := lib.AlgStatus{}
		if deviceID != nil {
			device.DeviceID, _ = strconv.ParseInt(string(deviceID.([]uint8)), 10, 64)
		} else {
			device.DeviceID = 0
		}
		if taskID != nil {
			//发送请求获得任务状态结果
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

		} else {
			device.AlgStatus = 0
		}
		result = append(result, device)
	}
	return result, nil
}

// Getting the algorithm configuration for the specified device
func GetAlgconfig(session *sql.Tx, deviceID int64) (algConfig string, err error) {
	sql := `SELECT alg_config FROM imp_t_alg WHERE device_id=?`
	err = session.QueryRow(sql, deviceID).Scan(&algConfig)
	return
}
