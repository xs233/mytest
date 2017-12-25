package dao

import (
	"SilverBusinessServer/dao/sql"
	"SilverBusinessServer/env"
	"SilverBusinessServer/lib"
)

//--------------------------------------------------------------------------------------------------------------------------------------
//Add Alarm Remark | 添加报警备注
func AddAlarmRemark(deviceID int64, alarmTime int64, remark string) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.AddRemarkByAlarmAndTime(session, deviceID, alarmTime, remark) //在数据库里操作，根据id和时间戳来添加报警备注
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//------------------------------------------------------------------------------------------------------------
//处理或归档报警，通过相机ID和报警时间戳标识报警
func PigeonholeAlarm(deviceID int64, alarmTime int64, processStatus int, archiveFlag int) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	//处理或归档报警
	err = sql.PigeonholeByAlarmAndTime(session, deviceID, alarmTime, processStatus, archiveFlag)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//------------------------------------------------------------------------------------------------------------
//删除报警，通过相机ID和报警时间戳标识报警
func DeleteAlarm(deviceID int64, alarmTime int64) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.DeleteAlarmByIdAndtime(session, deviceID, alarmTime)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//------------------------------------------------------------------------------------------------------------
//分页查询报警，查询参数中archiveFlag=1表示只查询归档的报警，其他值表示查询所有报警
func QueryAlarms(deviceID int64, beginTime int64, endTime int64, archiveFlag int, offset int, count int) (alarmList []lib.Alarm, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return nil, err
	}

	alarmList, err = sql.QueryAlarmByArOrNO(session, deviceID, beginTime, endTime, archiveFlag, offset, count)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	session.Commit()
	return alarmList, nil
}

//
func QueryAlarmTotal(deviceID int64, beginTime int64, endTime int64, archiveFlag int) (totalNum int64, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return -1, err
	}

	totalNum, err = sql.QueryAlarmTotalSql(session, deviceID, beginTime, endTime, archiveFlag)
	if err != nil {
		session.Rollback()
		return -1, err
	}

	session.Commit()
	return totalNum, nil
}

//--------------------------------------------------------------------------
// Get deviceId according to deviceUUID
func GetDeviceID(deviceUUID string) (deviceID int64, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return 0, err
	}

	if deviceID, err = sql.GetDeviceID(session, deviceUUID); err != nil {
		session.Rollback()
		return
	}

	session.Commit()
	return
}

//--------------------------------------------------------------------------
//Save the alarm information to the database
func SaveAlarm(value *lib.AlarmValue, binData []byte) error {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	//Get image id
	imageID, err := sql.GetSequenceIDByName(env.SequenceNameImageID)
	if err != nil {
		session.Rollback()
		return err
	}

	if err := sql.SaveAlarmToTableAlarm(session, value, imageID); err != nil {
		session.Rollback()
		return err
	} else {
		if err := sql.SaveAlarmToTableImage(session, binData, imageID); err != nil {
			session.Rollback()
			return err
		}
	}
	

	session.Commit()
	return nil
}

//--------------------------------------------------------------------------
//根据时间段来删除报警
func DeleteAlarmsByTime(beginTime int64, endTime int64) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.DeleteAlarmsByBeginAndEndTime(session, beginTime, endTime)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//--------------------------------------------------------------------------
//查询报警图片信息
func QueryAlarmImageInfo(imageId int64) (image []lib.AlarmImage, err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return nil, err
	}

	imageinfo, err := sql.QueryAlarmImageInfoSQL(session, imageId)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	session.Commit()
	return imageinfo, nil
}

//----------------------------------------------------------------------------
//历史报警消息添加备注
func AddAlarmRemarkByOld(alarmid int64, remark string) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.AddAlarmRemarkByOldSQL(session, alarmid, remark)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//----------------------------------------------------------------------------
//历史报警消息归档
func PigeonholeAlarmByOld(alarmid int64, archiveFlag, processStatus int) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.PigeonholeAlarmByOldSQL(session, alarmid, archiveFlag, processStatus)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}

//----------------------------------------------------------------------------
//历史报警消息删除
func DeleteAlarmByAlarmId(alarmid int64) (err error) {
	session, err := sql.DB.Begin()
	if err != nil {
		return err
	}

	err = sql.DeleteAlarmByAlarmIdSql(session, alarmid)
	if err != nil {
		session.Rollback()
		return err
	}

	session.Commit()
	return nil
}
