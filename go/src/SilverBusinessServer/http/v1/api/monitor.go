package api

//算法监控
import (
	"SilverBusinessServer/dao"
	"SilverBusinessServer/http/errcode"
	"SilverBusinessServer/sword"
	"SilverBusinessServer/livekeeper"
	"SilverBusinessServer/log"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//查询所有算法分析相机
/*	注：客户端定时查询所有设备的算法运行状况，业务后台接收到查询请求后，首先到数据库中查询所有
设备的算法信息（deviceID、algConfig、taskID），taskID为空的表示未通过任务分析框架sword创建任务，
其算法状态为未运行，taskID非空的任务通过任务分析框架sword提供的接口查询任务状态，或通过sword提供的
接口查询出其维护的所有任务的状态，业务层整理所有算法状态
*/
func HandleQueryAllAlgDevicesGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryAllAlgDevicesGet BEGIN")
	
	keyword := c.Query("keyword")
	offsetStr := c.Query("offset")
	countStr := c.Query("count")

	offset, _ := strconv.ParseInt(offsetStr, 10, 64)
	count, _ := strconv.ParseInt(countStr, 10, 64)

	deviceTotalNum, err := dao.QueryAlgDeviceNumber(keyword)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	deviceList, err := dao.QueryAllAlgs(keyword, offset, count)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":             errcode.ErrNoError.Code,
		"errMsg":          errcode.ErrNoError.String,
		"deviceTotalNum":  deviceTotalNum,
		"deviceList":      deviceList,
	})
	return
}

//--------------------------------------------------------------------------
//添加摄像机
func HandleAddAlgDevicePost(c *gin.Context) {
	log.HTTP.Info("HandleAddAlgDevicePost BEGIN")
	var reqJSON reqAddConcentrateDeviceJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}
	
	if err := dao.AddConcentrateDeviceByDeviceidList(reqJSON.DeviceIDList, reqJSON.AlgID); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrAddDidList.Code,
			"errMsg": errcode.ErrAddDidList.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqAddConcentrateDeviceJSON struct {
	DeviceIDList []int64 `form:"deviceIDList" json:"deviceIDList"`
	AlgID        string  `form:"algID" json:"algID"`
}

//--------------------------------------------------------------------------
// query the normal device without algorithm analysis
func HandleUnmonitorDeviceGet(c *gin.Context) {
	log.HTTP.Info("HandleUnmonitorDeviceGet BEGIN")
	deviceList, err := dao.QueryUnmonitorDevice()
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQueryUnmonitorDevice.Code,
			"errMsg": errcode.ErrQueryUnmonitorDevice.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    	   errcode.ErrNoError.Code,
		"errMsg": 	   errcode.ErrNoError.String,
		"deviceList":  deviceList,
	})
	return
}

//--------------------------------------------------------------------------
// Query the algorithm configuration for a particular device
func HandleQueryAlgConfigGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryAlgConfigGet BEGIN")
	deviceID, err := strconv.ParseInt(c.Param("did"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	algID, algConfig, err := dao.QueryAlgConfig(int(deviceID))
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQueryAlgConfig.Code,
			"errMsg": errcode.ErrQueryAlgConfig.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    	   errcode.ErrNoError.Code,
		"errMsg": 	   errcode.ErrNoError.String,
		"algID":  	   algID,
		"algConfig":   algConfig,
	})
	return
}

//--------------------------------------------------------------------------
//启停相机算法  start-启动，stop-停止
func HandleControlAlgPost(c *gin.Context) {
	log.HTTP.Info("HandleControlAlgPost BEGIN")
	//get device id from url address
	deviceId, err := strconv.ParseInt(c.Param("did"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	var reqJSON reqOpenEndALGJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	log.HTTP.Info(deviceId)
	deviceType, deviceUUID, algID, _, err := dao.GetDeviceType(int(deviceId))
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrGetDeviceType.Code,
			"errMsg": errcode.ErrGetDeviceType.String,
		})
		return
	}

	if deviceType == 1 {
		//如果是start 判断sword 计算能力是否充足
		mark, err := JudgeCapacity(reqJSON.Command)
		if err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrJudgeCapacity.Code,
				"errMsg": errcode.ErrJudgeCapacity.String,
			})
			return
		}
		if mark != 1 {
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrScarceCapacity.Code,
				"errMsg": errcode.ErrScarceCapacity.String,
			})
			return
		}

		//获得alg
		alg, err := dao.GetAlg(deviceId)
		if err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrGetALG.Code,
				"errMsg": errcode.ErrGetALG.String,
			})
			return
		}
		//判断此时相机的状态和操作是否冲突（比如操作码是start 但是相机状态就是start）
		if err := dao.ConflictOrNo(alg, reqJSON.Command); err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrConflict.Code,
				"errMsg": errcode.ErrConflict.String,
			})
			return
		}
		//获得algdevice信息（详见 lib.AlgDevice 数据结构）
		algDevice, err := dao.GetAlgDeviceByDeviceID(deviceId, reqJSON.Command)
		if err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrGetALGDevice.Code,
				"errMsg": errcode.ErrGetALGDevice.String,
			})
			return
		}
		//开始"启/停" 操作
		if err := dao.StartOrStopAlg(deviceId, reqJSON.Command, alg, algDevice); err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrOpenEndALG.Code,
				"errMsg": errcode.ErrOpenEndALG.String,
			})
			return
		}
	} else {
		var isStart bool
		if reqJSON.Command == "start" {
			isStart = true
		} else {
			isStart = false
		}

		if err := livekeeper.Notify(isStart, deviceUUID, algID); err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrOpenEndALG.Code,
				"errMsg": errcode.ErrOpenEndALG.String,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqOpenEndALGJSON struct {
	Command string `form:"command" json:"command"`
}

//--------------------------------------------------------------------------
//修改算法配置，须先停止算法运行，后台须作判断
func HandleUpdateAlgPut(c *gin.Context) {
	log.HTTP.Info("HandleUpdateAlgPut BEGIN")
	//1.获得deviceid
	didStr := c.Param("did")
	deviceId, err := strconv.ParseInt(didStr, 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	var reqJSON reqUpdateALGJSON
	if err := c.BindJSON(&reqJSON); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}

	deviceType, deviceUUID, algID, _, err := dao.GetDeviceType(int(deviceId))
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrGetDeviceType.Code,
			"errMsg": errcode.ErrGetDeviceType.String,
		})
		return
	}

	if deviceType == 1 {
		//2.判断算法运行状态
		//判断imp_t_alg 表中的task_id 是否为空
		sign, taskID, err := JudgedTask(deviceId)
		if err != nil || sign == -1 {
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrParams.Code,
				"errMsg": errcode.ErrParams.String,
			})
			return
		}
		//task_id 不为空
		if sign == 1 {
			//判断任务是否在运行
			sign1, err := JudgedRunningState(taskID)
			if sign1 == -1 || err != nil {
				c.JSON(http.StatusOK, gin.H{
					"err":    errcode.ErrJudgedRunningState.Code,
					"errMsg": errcode.ErrJudgedRunningState.String,
				})
				return
				//如果在运行则提示用户停止算法运行
			} else if sign1 == 1 {
				c.JSON(http.StatusOK, gin.H{
					"err":    errcode.ErrStopALG.Code,
					"errMsg": errcode.ErrStopALG.String,
				})
				return
				//停止状态 则先执行deleteTask
			} else if sign1 == 0 {
				result, err := sword.DeleteTask(taskID)
				if err != nil || result.Err != 0 {
					c.JSON(http.StatusOK, gin.H{
						"err":    errcode.ErrDeleteTask.Code,
						"errMsg": errcode.ErrDeleteTask.String,
					})
					return
				}
			}

		}
	} else {
		if err := livekeeper.SetAlgPara(deviceUUID, algID, reqJSON.AlgConfig); err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrUpdateAlgConfig.Code,
				"errMsg": errcode.ErrUpdateAlgConfig.String,
			})
			return
		}
	}

	//3.开始修改算法配置
	if err := dao.UpdateALGByAlgConfig(reqJSON.AlgConfig, deviceId); err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrUpdateAlgConfig.Code,
			"errMsg": errcode.ErrUpdateAlgConfig.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

type reqUpdateALGJSON struct {
	AlgID     string `form:"algID" json:"algID"`
	AlgConfig string `form:"algConfig" json:"algConfig"`
}

//-----------------------------------------------------------------------------------------------------------------------------
//删除相机，须先停止算法运行，后台须作判断
func HandleDeleteAlgDeviceDelete(c *gin.Context) {
	log.HTTP.Info("HandleDeleteAlgDeviceDelete BEGIN")
	deviceId, err := strconv.ParseInt(c.Param("did"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	deviceType, _, _, _, err := dao.GetDeviceType(int(deviceId))
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrGetDeviceType.Code,
			"errMsg": errcode.ErrGetDeviceType.String,
		})
		return
	}

	if deviceType == 1 {
		//2.判断算法运行状态
		//判断imp_t_alg 表中的task_id 是否为空
		sign, taskID, err := JudgedTask(deviceId)
		if err != nil || sign == -1 {
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrParams.Code,
				"errMsg": errcode.ErrParams.String,
			})
			return
		}
		//task_id 不为空
		if sign == 1 {
			//判断任务是否在运行
			sign1, err := JudgedRunningState(taskID)
			if sign1 == -1 || err != nil {
				c.JSON(http.StatusOK, gin.H{
					"err":    errcode.ErrJudgedRunningState.Code,
					"errMsg": errcode.ErrJudgedRunningState.String,
				})
				return
				//如果在运行则提示用户停止算法运行
			} else if sign1 == 1 {
				c.JSON(http.StatusOK, gin.H{
					"err":    errcode.ErrStopALG.Code,
					"errMsg": errcode.ErrStopALG.String,
				})
				return
				//停止状态 则先执行deleteTask
			} else if sign1 == 0 {
				result, err := sword.DeleteTask(taskID)
				if err != nil || result.Err != 0 {
					c.JSON(http.StatusOK, gin.H{
						"err":    errcode.ErrDeleteTask.Code,
						"errMsg": errcode.ErrDeleteTask.String,
					})
					return
				}
			}

		}
		//开始删除相机
		if err := dao.DeleteDeviceALG(deviceId); err != nil {
			log.HTTP.Error(err)
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrDelete.Code,
				"errMsg": errcode.ErrDelete.String,
			})
			return
		}
	} else {
		log.HTTP.Error("Smart device must not be deleted")
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrDeviceType.Code,
			"errMsg": errcode.ErrDeviceType.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
	})
	return
}

//---------------------------------------------------
//客户端定时查询所有设备的算法运行状况，定时15秒
func HandleQueryAllAlgDevicesStatusGet(c *gin.Context) {
	log.HTTP.Info("HandleQueryAllAlgDevicesStatusGet BEGIN")
	StatusList, err := dao.QueryAllAlgStatus()
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"err":        errcode.ErrNoError.Code,
		"errMsg":     errcode.ErrNoError.String,
		"deviceList": StatusList,
	})
	return
}

//通过deviceid查询imp_t_alg 的task_id 是否为空
func JudgedTask(deviceId int64) (sign int, taskID string, err error) {
	sign = -1
	alg, err := dao.GetAlg(deviceId)
	if err != nil {
		return sign, "", err
	}
	if alg.TaskID != "" {
		sign = 1
	} else {
		sign = 0
	}
	return sign, alg.TaskID, nil
}

//判断算法运行状态(sign=-1:请求失败；sign=0：已停止运行；sign=1：正在运行)
func JudgedRunningState(taskID string) (sign int, err error) {
	sign = -1
	var taskIDs []string
	taskIDs = append(taskIDs, taskID)
	result, err := sword.QueryTaskStatus(taskIDs)
	if result.Err != 0 || err != nil {
		return sign, errors.New("http error.")
	} else {
		for _, tas := range result.Data.TaskStatusList {
			if tas.TaskID == taskID {
				if tas.TaskStatus == 0 {
					sign = 0
				} else if tas.TaskStatus == 1 {
					sign = 1
				}
				break
			}
		}
	}
	return sign, nil
}

//(算法的启停)如果是start 判断sword 计算能力是否充足
func JudgeCapacity(command string) (int, error) {
	if command == "stop" {
		return 1, nil
	} else {
		result, err := sword.QueryTaskCapacity()
		if err != nil || result.Err != 0 {
			return 0, errors.New("Judge zhe Sword capacity error.")
		}
		if result.Data.TaskCapacity != 0 {
			return 1, nil
		} else {
			return 0, nil
		}
	}
}
