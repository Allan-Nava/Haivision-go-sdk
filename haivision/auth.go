package haivision

import (
	"encoding/json"
	"log"

	"gopkg.in/validator.v2"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/device"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/session"
)

/*
Requests
POST /api/session

	{
	 "username" :  "[user name]",
	 "password" :  "[password]"
	}

Response

	{
	    "response": {
	        "type": "Session",
	        "message": "Session successfully started for haiadmin",
	        "sessionID": "[Session ID]",
	        "lastLoginDate": 1536777877871,
	        "numLoginFailures": 0
	    }
	}
*/
func (o *Haivision) InitSession(username string, password string) (*session.BaseResponseInitSession, error) {
	requestBody := &session.RequestInitSession{
		Username: username,
		Password: password,
	}
	//
	if errs := validator.Validate(requestBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(SESSION, requestBody)
	if err != nil {
		return nil, err
	}
	o.debugPrint(resp)
	var obj session.BaseResponseInitSession
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

/*
Requests
GET /api/devices
cookie: sessionID: [Session ID]

Response
[

	{
	 "_id": "[Device ID]",
	 "type": "gateway",
	 "ip": "127.0.0.1",
	 "name": "Haivision Gateway",
	 "lastConnectedAt": [Date/time shown in Unix time],
	 "statusCode": "ok",
	 "status": "Online",
	 "statusDetails": "Connection has been established in the last 1 minutes.",
	 "serialNumber": null,
	 "firmware": "5.0.180611.1530",
	 "hasAdminError": false,
	 "pendingSync": false,
	 "lastConnection": "<1m"
	}

]
*/
func (o *Haivision) GetDeviceInfo() (*device.BaseResponseDeviceInfo, error) {
	resp, err := o.restyGet(DEVICE_INFO, nil)
	if err != nil {
		return nil, err
	}
	log.Println("GetDeviceInfo resp ",resp.Body())
	var obj device.BaseResponseDeviceInfo
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

//
