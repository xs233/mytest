package lib

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"impower/LicenseServer/env"
	"impower/LicenseServer/log"
)

var (
	apikey  = env.Get("phoneverify.apikey").(string)
	tplid   = env.Get("phoneverify.tplid").(string)
	company = env.Get("phoneverify.company").(string)
)

// SendShortMessage :
func SendShortMessage(phone string, code int) (err error) {
	resp, err := http.Post("http://yunpian.com/v1/sms/tpl_send.json",
		"application/x-www-form-urlencoded",
		strings.NewReader("tpl_id="+tplid+"&apikey="+apikey+"&tpl_value=%23code%23%3D"+strconv.Itoa(code)+"%26%23company%23%3D"+company+"&mobile="+phone))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	u := map[string]interface{}{}
	json.Unmarshal(body, &u)
	if u["code"].(float64) != 0 {
		log.Root.Info(string(body))
		return errors.New(string(body))
	}
	return nil
}
