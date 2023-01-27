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
	Response ResponseInitSession `json:"response" validate:"nonzero"`
}

type ResponseInitSession struct {
	Type             string `json:"type" validate:"nonzero"`
	Message          string `json:"message" validate:"nonzero"`
	SessionID        string `json:"sessionID" validate:"nonzero"`
	LastLoginDate    int64  `json:"lastLoginDate" validate:"nonzero"`
	NumLoginFailures int    `json:"numLoginFailures" validate:"nonzero"`
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
	SessionID        string   `json:"sessionID" validate:"nonzero"`
	DisplayName      string   `json:"displayName" validate:"nonzero"`
	Email            string   `json:"email" validate:"nonzero"`
	Roles            []string `json:"roles" validate:"nonzero"`
	StartAt          int64    `json:"startAt" validate:"nonzero"`
	ExpireAt         int64    `json:"expireAt" validate:"nonzero"`
	LastLoginDate    int64    `json:"lastLoginDate" validate:"nonzero"`
	NumLoginFailures int      `json:"numLoginFailures" validate:"nonzero"`
	IsLicensed       bool     `json:"isLicensed" validate:"nonzero"`
}
