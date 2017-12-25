package sword

import (
	"SilverBusinessServer/env"
	"SilverBusinessServer/lib"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//Sword计算框架 ip/port配置
var SwordIP = map[string]string{
	"ip":   env.Get("swordserver.host").(string),
	"port": env.Get("swordserver.port").(string),
}

//向算法监控模块Monitor 的 查询任务状态接口：http://[ip]:[port]/queryTaskStatus（GET）输入taskID列表，返回taskID列表对应的任务状态
func QueryTaskStatus(taskIDs []string) (result lib.TaskResult, err error) {
	url := "http://" + SwordIP["ip"] + ":" + SwordIP["port"] + "/queryTaskStatus?taskID="
	for _, id := range taskIDs {
		url = url + id + "&taskID="
	}
	index := strings.LastIndex(url, "&")
	urlNew := url[0:index]
	resp, err := http.Get(urlNew)
	if err != nil {
		return result, errors.New("QueryTaskStatus http  get error.")
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, errors.New("Unmarshal JSON error.")
	}
	return result, nil
}

//向算法监控模块Monitor 的 查询任务状态接口：http://[ip]:[port]/queryAllTaskStatus 获得所有任务状态
func QueryAllTaskStatus() (result lib.TaskResult, err error) {
	url := "http://" + SwordIP["ip"] + ":" + SwordIP["port"] + "/queryAllTaskStatus"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return result, errors.New("Query all task status http  get error.")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, errors.New("Unmarshal JSON error.")
	}
	return result, nil
}

//创建一个分析任务
func CreateTask(taskPara lib.TaskPara) (result lib.CreateTask, err error) {
	url := "http://" + SwordIP["ip"] + ":" + SwordIP["port"] + "/createTask"
	MarTP, _ := json.Marshal(taskPara)
	mapTaskInfo := map[string]string{
		"taskPara": string(MarTP),
	}
	taskJSON, err := json.Marshal(mapTaskInfo)
	if err != nil {
		return result, errors.New("Marshal JSON error.")
	}
	resp, err := http.Post(url, "text/plain", strings.NewReader(string(taskJSON)))
	defer resp.Body.Close()
	if err != nil {
		return result, errors.New("CreateTask http  post error.")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, errors.New("Unmarshal JSON error.")
	}

	return result, nil
}

//启动一个分析任务
func StartTask(taskID string) (result lib.StartTask, err error) {

	url := "http://" + SwordIP["ip"] + ":" + SwordIP["port"] + "/startTask"
	mapTaskInfo := map[string]string{
		"taskID": taskID,
	}
	taskJSON, err := json.Marshal(mapTaskInfo)
	if err != nil {
		return result, errors.New("Marshal JSON error.")
	}
	resp, err := http.Post(url, "text/plain", strings.NewReader(string(taskJSON)))
	defer resp.Body.Close()
	if err != nil {
		return result, errors.New("StartTask http  post error.")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, errors.New("Unmarshal JSON error.")
	}
	return result, nil
}

//停止分析任务
func StopTask(taskID string) (result lib.StartTask, err error) {
	//测试使用
	url := "http://" + SwordIP["ip"] + ":" + SwordIP["port"] + "/stopTask"
	mapTaskInfo := map[string]string{
		"taskID": taskID,
	}
	taskJSON, err := json.Marshal(mapTaskInfo)
	if err != nil {
		return result, errors.New("Marshal JSON error.")
	}
	resp, err := http.Post(url, "text/plain", strings.NewReader(string(taskJSON)))
	defer resp.Body.Close()
	if err != nil {
		return result, errors.New("StopTask http post error.")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, errors.New("Unmarshal JSON error.")
	}
	return result, nil
}

//删除分析任务
func DeleteTask(taskID string) (result lib.DeleteTask, err error) {
	url := "http://" + SwordIP["ip"] + ":" + SwordIP["port"] + "/deleteTask"
	mapTaskInfo := map[string]string{
		"taskID": taskID,
	}
	taskJSON, err := json.Marshal(mapTaskInfo)
	if err != nil {
		return result, errors.New("Marshal JSON error.")
	}
	resp, err := http.Post(url, "text/plain", strings.NewReader(string(taskJSON)))
	defer resp.Body.Close()
	if err != nil {
		return result, errors.New("DeleteTask http post error.")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, errors.New("Unmarshal JSON error.")
	}
	return result, nil
}

//查询可用的任务计算能力
func QueryTaskCapacity() (result lib.TaskCapacity, err error) {
	url := "http://" + SwordIP["ip"] + ":" + SwordIP["port"] + "/queryTaskCapacity"
	resp, err := http.Get(url)
	if err != nil {
		return result, errors.New("QueryTaskCapacity error.")
	}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, errors.New("Unmarshal JSON error.")
	}
	fmt.Println(result)
	return result, nil
}
