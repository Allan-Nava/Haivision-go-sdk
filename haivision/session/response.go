package session

/*
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
type BaseResponseInitSession struct {
	Response ResponseInitSession `json:"response"`
}

type ResponseInitSession struct {
	Type             string `json:"type"`
	Message          string `json:"message"`
	SessionID        string `json:"sessionID"`
	LastLoginDate    int64  `json:"lastLoginDate"`
	NumLoginFailures int    `json:"numLoginFailures"`
}

/*
	{
	    "sessionID": "[Session ID]",
	    "displayName": "Administrator",
	    "email": "haiadmin@localhost",
	    "roles": [
	        "Administrator"
	    ],
	    "startAt": 1536937838919,
	    "expireAt": 1536938857529,
	    "lastLoginDate": null,
	    "numLoginFailures": null,
	    "isLicensed": true
	}
*/
type ResponseSessionInfo struct {
	SessionID        string   `json:"sessionID"`
	DisplayName      string   `json:"displayName"`
	Email            string   `json:"email"`
	Roles            []string `json:"roles"`
	StartAt          int64    `json:"startAt"`
	ExpireAt         int64    `json:"expireAt"`
	LastLoginDate    int64    `json:"lastLoginDate"`
	NumLoginFailures int      `json:"numLoginFailures"`
	IsLicensed       bool     `json:"isLicensed"`
}
