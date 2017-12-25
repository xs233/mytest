package api

//算法监控
import (
	"SilverBusinessServer/dao"
	"SilverBusinessServer/http/errcode"
	"SilverBusinessServer/sword"
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
	deviceList, err := dao.QueryAllAlgs()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrQuery.Code,
			"errMsg": errcode.ErrQuery.String,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"err":        errcode.ErrNoError.Code,
		"errMsg":     errcode.ErrNoError.String,
		"deviceList": deviceList,
	})
	return
}

//--------------------------------------------------------------------------
//添加摄像机
func HandleAddAlgDevicePost(c *gin.Context) {
	var reqJSON reqAddConcentrateDeviceJSON
	if c.BindJSON(&reqJSON) != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}
	if dao.AddConcentrateDeviceByDeviceidList(reqJSON.DeviceIDList) != nil {
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
}

//--------------------------------------------------------------------------
// Getting the algorithm configuration for the specified device
func AlgConfigGet(c *gin.Context) {
	//get device id from url address
	deviceID, err := strconv.ParseInt(c.Param("did"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}

	algConfig, err := dao.GetAlgconfig(deviceID)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrGetAlgConfig.Code,
			"errMsg": errcode.ErrGetAlgConfig.String,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"err":    errcode.ErrNoError.Code,
		"errMsg": errcode.ErrNoError.String,
		"algConfig": algConfig,
	})
	return
}


//--------------------------------------------------------------------------
//启停相机算法  start-启动，stop-停止
func HandleControlAlgPost(c *gin.Context) {
	//get device id from url address
	deviceId, err := strconv.ParseInt(c.Param("did"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}
	var reqJSON reqOpenEndALGJSON
	if c.BindJSON(&reqJSON) != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}
	//如果是start 判断sword 计算能力是否充足
	mark, err := JudgeCapacity(reqJSON.Command)
	if err != nil {
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
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrGetALG.Code,
			"errMsg": errcode.ErrGetALG.String,
		})
		return
	}
	//判断此时相机的状态和操作是否冲突（比如操作码是start 但是相机状态就是start）
	if dao.ConflictOrNo(alg, reqJSON.Command) != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrConflict.Code,
			"errMsg": errcode.ErrConflict.String,
		})
		return
	}
	//获得algdevice信息（详见 lib.AlgDevice 数据结构）
	algDevice, err := dao.GetAlgDeviceByDeviceID(deviceId, reqJSON.Command)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrGetALGDevice.Code,
			"errMsg": errcode.ErrGetALGDevice.String,
		})
		return
	}
	//开始"启/停" 操作
	if dao.StartOrStopAlg(deviceId, reqJSON.Command, alg, algDevice) != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrOpenEndALG.Code,
			"errMsg": errcode.ErrOpenEndALG.String,
		})
		return
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
	//1.获得deviceid
	didStr := c.Param("did")
	deviceId, err := strconv.ParseInt(didStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}
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
	//3.开始修改算法配置
	var reqJSON reqUpdateALGJSON
	if c.BindJSON(&reqJSON) != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrForm.Code,
			"errMsg": errcode.ErrForm.String,
		})
		return
	}
	if dao.UpdateALGByAlgConfig(reqJSON.AlgConfig, deviceId) != nil {
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
	AlgConfig string `form:"algConfig" json:"algConfig"`
}

//-----------------------------------------------------------------------------------------------------------------------------
//删除相机，须先停止算法运行，后台须作判断
func HandleDeleteAlgDeviceDelete(c *gin.Context) {
	deviceId, err := strconv.ParseInt(c.Param("did"), 10, 64)
	if err != nil {
		log.HTTP.Error(err)
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrParams.Code,
			"errMsg": errcode.ErrParams.String,
		})
		return
	}
	//2.判断算法运行状态
	//判断imp_t_alg 表中的task_id 是否为空
	sign, taskID, err := JudgedTask(deviceId)
	if err != nil || sign == -1 {
		log.HTTP.Error("JudgedTask failed")
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
			log.HTTP.Error("JudgedRunningState failed")
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrJudgedRunningState.Code,
				"errMsg": errcode.ErrJudgedRunningState.String,
			})
			return
			//如果在运行则提示用户停止算法运行
		} else if sign1 == 1 {
			log.HTTP.Error("Alg running")
			c.JSON(http.StatusOK, gin.H{
				"err":    errcode.ErrStopALG.Code,
				"errMsg": errcode.ErrStopALG.String,
			})
			return
			//停止状态 则先执行deleteTask
		} else if sign1 == 0 {
			result, err := sword.DeleteTask(taskID)
			if err != nil || result.Err != 0 {
				log.HTTP.Error("DeleteTask failed")
				c.JSON(http.StatusOK, gin.H{
					"err":    errcode.ErrDeleteTask.Code,
					"errMsg": errcode.ErrDeleteTask.String,
				})
				return
			}
		}

	}
	//开始删除相机
	if dao.DeleteDeviceALG(deviceId) != nil {
		c.JSON(http.StatusOK, gin.H{
			"err":    errcode.ErrDelete.Code,
			"errMsg": errcode.ErrDelete.String,
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
func HandleQueryAllAlgDevicesStatus(c *gin.Context) {
	StatusList, err := dao.QueryAllAlgStatus()
	if err != nil {
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
