package dao

import (
	"SilverBusinessServer/dao/sql"
	"SilverBusinessServer/lib"
	"SilverBusinessServer/sword"
	"errors"
)

//--------------------------------------------------------------------------
// 修改设备的设备名
func QueryChangeDeviceName(deviceID int64, deviceName string) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.UpdateDeviceNameByDeviceID(session, deviceID, deviceName)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//-------------------------------------------------------------------------------------------------
// query the normal device without algorithm analysis
func QueryUnmonitorDevice() ([]lib.Device, error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return nil, err
	}

	deviceList, err := sql.QueryUnmonitorDevice(session)
	if err != nil {
		session.Rollback()
		return nil, err
	}
	
	session.Commit()
	return deviceList, nil
}

//--------------------------------------------------------------------------
//启停相机算法  start-启动，stop-停止
func StartOrStopAlg(deviceId int64, command string, alg lib.Alg, algDevice lib.AlgDevice) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}
	err = sql.StartOrStopAlg(session, deviceId, command, alg, algDevice)
	if err != nil {
		session.Rollback()
		return err
	}
	session.Commit()
	return nil
}

//根据device_id 获得alg
func GetAlg(deviceId int64) (alg lib.Alg, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return alg, err
	}
	alg, err = sql.GetAlg(session, deviceId)
	if err != nil {
		session.Rollback()
		return alg, err
	}
	session.Commit()
	return alg, nil
}

//--------------------------------------------------------------------------
// 添加摄像机
func AddConcentrateDeviceByDeviceidList(deviceidList []int64, algID string) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}
	err = sql.AddDeviceToAlg(session, deviceidList, algID)
	if err != nil {
		session.Rollback()
		return err
	}
	session.Commit()
	return nil
}

//--------------------------------------------------------------------------
// Query the algorithm configuration for a particular device
func QueryAlgConfig(deviceID int) (algID string, algConfig string, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return "", "", err
	}

	algID, algConfig, err = sql.QueryAlgConfig(session, deviceID)
	if err != nil {
		session.Rollback()
		return "", "", err
	}

	session.Commit()
	return
}

//--------------------------------------------------------------------------
// Get the device type
func GetDeviceType(deviceID int) (deviceType int, deviceUUID, algID, algPara string, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return 0, "", "", "", err
	}

	deviceType, deviceUUID, algID, algPara, err = sql.GetDeviceType(session, deviceID)
	if err != nil {
		session.Rollback()
		return 0, "", "", "", err
	}

	session.Commit()
	return
}

//--------------------------------------------------------------------------
//修改算法配置
func UpdateALGByAlgConfig(AlgConfig string, deviceId int64) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.UpdateALGBySql(session, AlgConfig, deviceId)
	if err != nil {
		session.Rollback()
		return err
	}
	session.Commit()
	return nil
}

//--------------------------------------------------------------------------
//删除相机（此时算法已经停止运行了）
func DeleteDeviceALG(deviceId int64) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.DeleteDeviceALG(session, deviceId)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//-----------------------------------------------------------------------------------------------------------
// Query the total number of alg devices
func QueryAlgDeviceNumber(keyword string) (int, error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return -1, err
	}

	num, err := sql.QueryAlgDeviceNumber(session, keyword)
	if err != nil {
		session.Rollback()
		return -1, err
	}

	session.Commit()
	return num, nil
}

//---------------------------------------------------------------------------
//查询所有算法分析设备
func QueryAllAlgs(keyword string, offset, count int64) (deviceList []lib.AlgDevice, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return nil, err
	}

	deviceList, err = sql.QueryAllAlgs(session, keyword, offset, count)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	session.Commit()
	return deviceList, nil
}

//判断此时相机的状态和操作是否冲突（比如操作码是start 但是相机状态就是start）
func ConflictOrNo(alg lib.Alg, command string) (err error) {
	var sign = 0
	if alg.TaskID == "" {
		if command == "start" {
			return nil
		} else if command == "stop" {
			return errors.New("conflict")
		} else {
			return errors.New("unknown error")
		}
	}
	var taskIDs []string
	taskIDs = append(taskIDs, alg.TaskID)
	result, err := sword.QueryTaskStatus(taskIDs)
	if err != nil || result.Err != 0 {
		return errors.New("unknown error")
	}
	//判断数据长度
	if len(result.Data.TaskStatusList) > 0 {
		if result.Data.TaskStatusList[0].TaskStatus == 1 {
			sign = 1
		}
	} else {
		return errors.New("no result.")
	}
	//继续判断
	if command == "start" {
		if sign == 0 {
			return nil
		} else {
			return errors.New("conflict")
		}
	} else if command == "stop" {
		if sign == 0 {
			return errors.New("conflict")
		} else if sign == 1 {
			return nil
		}
	}

	return errors.New("unknown cmmand")

}

//通过deviceid获取该algdevice信息
func GetAlgDeviceByDeviceID(deviceId int64, command string) (algDevice lib.AlgDevice, err error) {
	//只有command 是start时候 才获取algDevice 信息
	if command == "start" {
		session, err := sql.DB.Begin()
		if err != nil {
			return algDevice, err
		}
		algDevice, err = sql.GetAlgDeviceByDeviceID(session, deviceId)
		if err != nil {
			session.Rollback()
			return algDevice, err
		}
		session.Commit()
	}

	return algDevice, nil

}

//查询所有算法分析相机的算法运行状态
func QueryAllAlgStatus() (result []lib.AlgStatus, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return nil, err
	}
	result, err = sql.QueryAllAlgStatus(session)
	if err != nil {
		session.Rollback()
		return nil, err
	}
	session.Commit()
	return result, nil

}
