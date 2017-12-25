package livekeeper

import (
	"SilverBusinessServer/env"
	"SilverBusinessServer/lib"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var host string = env.Get("livekeeper.host").(string) + ":" + env.Get("livekeeper.port").(string)
var c = &http.Client{  
    Timeout: 3 * time.Second,
}

func QueryDeviceAlgStatus(deviceList []string) (*lib.DeviceAlgStatus, error) {
	url := "http://" + host + "/queryDevAlgStatus"
	for i, v := range deviceList {
		if i == 0 {
			url += "?devID="
		} else {
			url += "&devID="
		}

		url += v
	}

	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	
	deviceAlgStatus := new(lib.DeviceAlgStatus)
	if err := json.Unmarshal([]byte(body), deviceAlgStatus); err != nil {
		return nil, err
	}

	if deviceAlgStatus.Err != 0 {
		return nil, errors.New(deviceAlgStatus.ErrMsg)
	}

	return deviceAlgStatus, nil
}

func isExpectStatus(deviceUUID string, expectStatus int) (res bool, err error) {
	var deviceList []string
	deviceList = append(deviceList, deviceUUID)
	if deviceAlgStatus, err := QueryDeviceAlgStatus(deviceList); err != nil {
		return false, err
	} else {
		var realStatus int
		for _, v := range deviceAlgStatus.Data.DevAlgStatusList {
			realStatus = v.AlgStatus
		}

		if expectStatus == realStatus {
			res = true
		} else {
			res = false
		}
	}

	return res, nil
}


func Notify(isStart bool, deviceUUID, algID string) (error) {
	var url string
	var expectStatus int
	if isStart == true {
		url = "http://" + host + "/startDevAlg"
		expectStatus = 0
	} else {
		url = "http://" + host + "/stopDevAlg"
		expectStatus = 1
	}
	
	if res, err := isExpectStatus(deviceUUID, expectStatus); err != nil {
		return err
	} else {
		if res == false {
			return errors.New("Status Conflict")
		}
	}

	requestBody := map[string]string {
		"devID": deviceUUID,
		"algID": algID,
	}

	if err := postRequest(url, requestBody); err != nil {
		return err
	}

	return nil
}

type RespondMsg struct {
	Err  	int     `from:"err" json:"err"`
	ErrMsg  string  `from:"errMsg" json:"errMsg"`
}

func SetAlgPara(deviceUUID, algID, algPara string) error {
	if res, err := isExpectStatus(deviceUUID, 0); err != nil {
		return err
	} else {
		if res == false {
			return errors.New("The alg is running!")
		}
	}

	url := "http://" + host + "/setDevAlgPara"
	requestBody := map[string]string {
		"devID": deviceUUID,
		"algID": algID,
		"algPara": algPara,
	}

	if err := postRequest(url, requestBody); err != nil {
		return err
	}

	return nil
}

func postRequest(url string, requestBody map[string]string) error {
	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	res, err := c.Post(url, "text/plain", strings.NewReader(string(requestBodyJson)))
	if err != nil {
		return err
	}

	respondBody, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	var respondMsg RespondMsg
	if err := json.Unmarshal([]byte(respondBody), &respondMsg); err != nil {
		return err
	}

	if respondMsg.Err != 0 {
		return errors.New(respondMsg.ErrMsg)
	}

	return nil
}