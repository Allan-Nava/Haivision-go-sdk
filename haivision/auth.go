package haivision

import (
	"encoding/json"

	"gopkg.in/validator.v2"
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
func (o *Haivision) InitSession(username string, password string) (*BaseResponseInitSession, error) {
	var requestBody RequestInitSession
	//
	if errs := validator.Validate(requestBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(SESSION, requestBody)
	if err != nil {
		return nil, err
	}
	var obj BaseResponseInitSession
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// GET
func (o *Haivision) GetDeviceInfo() error {
	return nil
}
