package server

import (
	"net/http"
	"time"
	"log"
	"encoding/json"
	"encoding/base64"
	"strings"
	"errors"
	"io/ioutil"
)

func Server(addr string) (err error){
	var hander msgHandler
	http.Handle("/push", hander)
	log.Printf("LivePushProxyServer start...")
	log.Fatal(http.ListenAndServe(addr, nil))

	return
}

func (handler msgHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Verify authorization
	if err := verifyAuthorization(r.Header.Get("Authorization")); err != nil {
		log.Print(err)
		msg, _ := json.Marshal(responseMsg{AuthorizationError, "Authorization error"})
		w.Write(msg)
		return
	}

	// Read body message
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		msg, _ := json.Marshal(responseMsg{ReadBodyError, "Read body error"})
		w.Write(msg)
		return
	}

	// Unmarshal body
	if err := json.Unmarshal(body, &handler); err != nil {
		log.Print(err)
		msg, _ := json.Marshal(responseMsg{ReadBodyError, "Read body error"})
		w.Write(msg)
		return
	}

	// Construct http body
	requestBody, _ := constructBody(handler)

	// Push notification
	if err := push(url, devKey, devSecret, string(requestBody)); err != nil {
		log.Print(err)
		msg, _ := json.Marshal(responseMsg{PushError, "Push notification error"})
		w.Write(msg)
		return
	}
	
	msg, _ := json.Marshal(responseMsg{Success, "Handle success"})
	w.Write(msg)
	return
}

func verifyAuthorization(authorization string) (err error) {	
	// Determine whether the data is empty
	if authorization == "" {
		err = errors.New("The authorization is empty")
		return
	}

	// Parse authorization
	authInfo := strings.Split(string(authorization), " ")
	if len(authInfo) != 2 {
		err = errors.New("Authorization format error")
		return
	}

	// Decode info
	keySecret, err := base64.StdEncoding.DecodeString(authInfo[1])
	if err != nil {
		err = errors.New("Authorization coding error")
		return
	}

	// Verify the content
	if string(keySecret) != "impower:taiheming" {
		err = errors.New("devKey or devSecret error")
		return
	}
	return
}

func push(url, devKey, devSecret, requestBody string) error {
	c := &http.Client{  
		Timeout: 3 * time.Second,
	}

	// Construct request
	log.Printf(requestBody)
	request, err := http.NewRequest("POST", url, strings.NewReader(requestBody))
	if err != nil {
		return err
	}

	// Set request header
	authorization := "Basic " +  base64.StdEncoding.EncodeToString([]byte(devKey + ":" + devSecret))
	request.Header.Add("Authorization", authorization)

	// Post request
	response, err := c.Do(request)
	if err != nil {
		return err
	}

	// Parse reponse body
	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	log.Print(string(responseBody))
	return nil
}


func constructBody(handler msgHandler) ([]byte, error) {
	audience := audience{}
	for _, item := range handler.Subscriber {
		if len(item) > 10 {
			audience.RegistrationId = append(audience.RegistrationId, item)
		}
	}
	notificationInfo := notification{Alert: handler.MsgTitle, Extras: handler.MsgBody}
	mapNotification := make(map[string]notification)
	mapNotification["android"] = notificationInfo
	mapNotification["ios"] = notificationInfo
	options := options{TimeToLive: 60, ApnsProduction: false}
	pushBody := jPush{
		Platform: "all", 
		Audience: audience, 
		Notification: mapNotification, 
		Options: options,
	}

	return json.Marshal(pushBody)
}