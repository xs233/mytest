package sql

import (
	"SilverBusinessServer/lib"
	"SilverBusinessServer/log"
	MYSQL "database/sql"
	"reflect"
	//"strconv"
	"encoding/base64"
	//"github.com/gin-gonic/gin"
	//"gopkg.in/mgo.v2"
	"fmt"
)

//---------------------------------------------------------------------------------------------------------------------------------------
// add alarm remark | 实时报警：在数据库里操作，根据id和时间戳来添加报警备注 【就是相当于更新报警表中的报警备注，类似于修改密码】
func AddRemarkByAlarmAndTime(session *MYSQL.Tx, deviceID int64, alarmTime int64, remark string) (err error) {
	sql := "update imp_t_alarm set remark = ? where device_id = ? and alarm_time = ?"
	_, err = session.Exec(sql, remark, deviceID, alarmTime)
	if err != nil {
		log.HTTP.Error("sql操作失败了，没有添加备注")
		return err
	}

	return nil
}

//---------------------------------------------------------------------------------------------------------------------------------------
// 实时报警：处理或归档报警
//"processStatus":int 	   	更新处理状态，0-未处理，1-已处理，-1-无效
//"archiveFlag":  int      更新归档标识，0-未归档，1-已归档，-1-无效
func PigeonholeByAlarmAndTime(session *MYSQL.Tx, deviceID int64, alarmTime int64, processStatus int, archiveFlag int) (err error) {
	var sql string

	if archiveFlag == -1 {
		sql = "update imp_t_alarm set process_status = ? where device_id = ? and alarm_time = ?"
		_, err := session.Exec(sql, processStatus, deviceID, alarmTime)
		if err != nil {
			log.HTTP.Error("", err)
			return err
		}
	} else if processStatus == -1 {
		sql = "update imp_t_alarm set archive_flag = ? where device_id = ? and alarm_time = ?"
		_, err := session.Exec(sql, archiveFlag, deviceID, alarmTime)

		if err != nil {
			log.HTTP.Error("", err)
			return err
		}
	} else {
		sql = "update imp_t_alarm set archive_flag = ?,process_status = ?  where device_id = ? and alarm_time = ?"
		_, err := session.Exec(sql, processStatus, archiveFlag, deviceID, alarmTime)

		if err != nil {
			log.HTTP.Error("", err)
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

//----------------------------------------------------------------------------------------------------------------------------------------
//实时报警：删除报警，通过相机ID和报警时间戳标识报警
func DeleteAlarmByIdAndtime(session *MYSQL.Tx, deviceID int64, alarmTime int64) (err error) {
	sql := `delete from imp_t_alarm where device_id = ? and alarm_time = ?` //标准MySQL语句
	_, err = session.Exec(sql, deviceID, alarmTime)
	if err != nil {
		return err
	}
	return nil
}

//------------------------------------------------------------------------------------------------------------
//按时间段批量删除报警，归档报警不删除，只删除未归档报警 【需要先从数据库中获取报警是否归档了，对于没有归档的，可以删除，已经归档的，就不好删除了】
func DeleteAlarmsByBeginAndEndTime(session *MYSQL.Tx, beginTime int64, endTime int64) (err error) {
	//先查询时间段内的所有报警
	sql := "delete from imp_t_alarm where alarm_time <= ? and alarm_time >= ? and archive_flag != 1"

	_, err = session.Exec(sql, endTime, beginTime)
	if err != nil {
		return err
	}
	return nil
}

//------------------------------------------------------------------------------------------------------------
//分页查询报警，查询参数中archiveFlag=1表示只查询归档的报警，其他值表示查询所有报警
func QueryAlarmByArOrNO(session *MYSQL.Tx, deviceID int64, beginTime int64, endTime int64, archiveFlag int, offset int, count int) (alarmlist []lib.Alarm, err error) {
	var sql string
	var rows *MYSQL.Rows
	var errs error
	fmt.Println("参数:", endTime, beginTime, offset, count)
	if archiveFlag == 1 { // 表示查询的是归档的
		if deviceID == 0 { //查询的是所有的
			sql = "select t1.* from imp_t_alarm t1, imp_t_device t2 where t1.device_id = t2.device_id and t1.alarm_time <= ? and t1.alarm_time >= ? and t1.archive_flag = 1 limit ?,?"
			rows, errs = session.Query(sql, endTime, beginTime, offset, count)
		} else { //查询device_id的
			sql = "select * from imp_t_alarm where alarm_time <= ? and alarm_time >= ? and device_id = ? and archive_flag = 1 limit ?,?"
			rows, errs = session.Query(sql, endTime, beginTime, deviceID, offset, count)
		}

	} else { //查询所有报警
		if deviceID == 0 {
			sql = "select t1.* from imp_t_alarm t1, imp_t_device t2 where t1.device_id = t2.device_id and t1.alarm_time <= ? and t1.alarm_time >= ? limit ?,?"
			rows, errs = session.Query(sql, endTime, beginTime, offset, count)
		} else {
			sql = "select * from imp_t_alarm where alarm_time <= ? and alarm_time >= ? and device_id = ? limit ?,?"
			rows, errs = session.Query(sql, endTime, beginTime, deviceID, offset, count)
		}

	}

	if errs != nil {
		log.HTTP.Error("sql失败")
		return nil, errs
	}
	defer rows.Close()

	alarmlist = []lib.Alarm{}
	for rows.Next() {

		var (
			alarm_id       interface{}
			device_id      interface{}
			alarm_time     interface{}
			image_id       interface{}
			alarm_info     interface{}
			process_status interface{}
			archive_flag   interface{}
			remark         interface{}
		)
		alarm := lib.Alarm{}
		if err = rows.Scan(&alarm_id, &device_id, &alarm_time,
			&image_id, &alarm_info, &process_status, &archive_flag, &remark); err != nil {
			return nil, err
		}
		alarm.AlarmID = alarm_id.(int64)
		alarm.DeviceID = device_id.(int64)
		alarm.AlarmTime = alarm_time.(int64)
		alarm.ImageID = image_id.(int64)

		if alarm_info != nil {
			alarm.AlarmInfo = string(alarm_info.([]uint8))
		} else {
			alarm.AlarmInfo = ""
		}

		alarm.ProcessStatus = process_status.(int64)
		alarm.ArchiveFlag = archive_flag.(int64)

		if remark != nil {
			alarm.Remark = string(remark.([]uint8))
		} else {
			alarm.Remark = ""
		}

		alarmlist = append(alarmlist, alarm)
	}
	fmt.Println("获取报警列表:", alarmlist)
	return alarmlist, nil
}

//---------------------------------------------------------------------------
//查询报警总数
func QueryAlarmTotalSql(session *MYSQL.Tx, deviceID int64, beginTime int64, endTime int64, archiveFlag int) (totalNum int64, err error) {
	log.HTTP.Info("查询报警总数的参数：", deviceID, beginTime, endTime, archiveFlag)
	var sql string
	var rows *MYSQL.Rows
	if archiveFlag == 1 { // 表示查询的是归档的
		if deviceID == 0 { //查询的是所有的
			sql = "select count(t1.alarm_id) from imp_t_alarm t1, imp_t_device t2 where t1.device_id = t2.device_id and t1.alarm_time <= ? and t1.alarm_time >= ? and t1.archive_flag = 1"
			rows, _ = session.Query(sql, endTime, beginTime)

		} else { //查询device_id的
			sql = "select count(alarm_id) from imp_t_alarm where alarm_time <= ? and alarm_time >= ? and device_id = ? and archive_flag = 1"
			rows, _ = session.Query(sql, endTime, beginTime, deviceID)
		}

	} else { //查询所有报警
		if deviceID == 0 {

			sql = "select count(t1.alarm_id) from imp_t_alarm t1, imp_t_device t2 where t1.device_id = t2.device_id and t1.alarm_time <= ? and t1.alarm_time >= ?"
			rows, _ = session.Query(sql, endTime, beginTime)
		} else {
			sql = "select count(alarm_id) from imp_t_alarm where alarm_time <= ? and alarm_time >= ? and device_id = ?"
			rows, _ = session.Query(sql, endTime, beginTime, deviceID)
		}
	}
	defer rows.Close()

	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return -1, err
		}
		totalNum = count
	}
	log.HTTP.Info("----------_>total:", totalNum)

	return totalNum, nil

}

//--------------------------------------------------------------------------
//查询报警图片信息
func QueryAlarmImageInfoSQL(session *MYSQL.Tx, imageID int64) (image []lib.AlarmImage, err error) {
	sql := "select image_id, image_width, image_height, image_format, image_data from imp_t_image where image_id = ?"
	rows, err := session.Query(sql, imageID)
	if err != nil {
		log.HTTP.Error("sql失败")
		return nil, err
	}
	defer rows.Close()

	imageinfo := []lib.AlarmImage{}
	for rows.Next() {
		var (
			sImageID     interface{}
			sImageWith   interface{}
			sImageHeight interface{}
			sImageFormat interface{}
			sImageData   interface{}
		)
		imageIn := lib.AlarmImage{}
		if err = rows.Scan(&sImageID, &sImageWith, &sImageHeight, &sImageFormat, &sImageData); err != nil {
			return nil, err
		}
		imageIn.ImageID = sImageID.(int64)
		imageIn.ImageWidth = sImageWith.(int64)
		imageIn.ImageHeight = sImageHeight.(int64)
		if sImageFormat != nil {
			imageIn.ImageFormat = string(sImageFormat.([]uint8))
		} else {
			imageIn.ImageFormat = ""
		}

		if sImageData != nil {
			imageIn.ImageData = base64.StdEncoding.EncodeToString(sImageData.([]uint8))
		} else {
			imageIn.ImageData = ""
		}

		log.HTTP.Info("打印imageData图片类型：", reflect.TypeOf(sImageData), len(imageIn.ImageData))
		imageinfo = append(imageinfo, imageIn)
	}
	return imageinfo, nil
}

//---------------------------------------------------------------------------
//历史报警消息添加报警备注sql处理
func AddAlarmRemarkByOldSQL(session *MYSQL.Tx, alarmid int64, remark string) (err error) {
	log.HTTP.Info("begin sql operation:", remark, alarmid)
	sql := "update imp_t_alarm set remark = ? where alarm_id = ?"
	_, err = session.Exec(sql, remark, alarmid)
	if err != nil {
		log.HTTP.Error("sql operation error,no remark to add", err)
		return err
	}
	return nil
}

//---------------------------------------------------------------------------
//历史报警消息归档sql处理
func PigeonholeAlarmByOldSQL(session *MYSQL.Tx, alarmid int64, archiveFlag, processStatus int) (err error) {
	log.HTTP.Info("参数：", alarmid, archiveFlag, processStatus)
	var sql string
	if archiveFlag == -1 {
		sql = "update imp_t_alarm set process_status = ? where alarm_id = ?"
		_, err := session.Exec(sql, processStatus, alarmid)
		if err != nil {
			log.HTTP.Error("", err)
			return err
		}
	} else if processStatus == -1 {
		sql = "update imp_t_alarm set archive_flag = ? where alarm_id = ?"
		_, err := session.Exec(sql, archiveFlag, alarmid)

		if err != nil {
			log.HTTP.Error("", err)
			return err
		}
	} else {
		sql = "update imp_t_alarm set archive_flag = ?,process_status = ?  where alarm_id = ?"
		_, err := session.Exec(sql, processStatus, archiveFlag, alarmid)
		if err != nil {
			log.HTTP.Error("", err)
			return err
		}
	}

	return nil
}

//---------------------------------------------------------------------------
//历史报警消息删除sql处理
func DeleteAlarmByAlarmIdSql(session *MYSQL.Tx, alarmid int64) (err error) {
	sql := `delete from imp_t_alarm where alarm_id = ?` //标准MySQL语句
	_, err = session.Exec(sql, alarmid)
	if err != nil {
		return err
	}
	return nil
}
