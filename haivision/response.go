package haivision

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

type BaseResponseDeviceInfo struct {
	Responses []ResponseDeviceInfo
}

/*
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
*/
type ResponseDeviceInfo struct {
	ID               string `json:"_id" validate:"nonzero"`
	Type             string `json:"type" validate:"nonzero"`
	IP               string `json:"ip" validate:"nonzero"`
	Name             string `json:"name" validate:"nonzero"`
	LastConnectedAt  int64  `json:"lastConnectedAt" validate:"nonzero"`
	StatusCode       string `json:"statusCode" validate:"nonzero"`
	Status           string `json:"status" validate:"nonzero"`
	StatusDetails    string `json:"statusDetails" validate:"nonzero"`
	SerialNumber     string `json:"serialNumber" validate:"nonzero"`
	Firmware         string `json:"firmware" validate:"nonzero"`
	HasAdminError    bool   `json:"hasAdminError" validate:"nonzero"`
	PendingSync      bool   `json:"pendingSync" validate:"nonzero"`
	LastConnection   string `json:"lastConnection" validate:"nonzero"`
	LastConnectionAt string `json:"lastConnectionAt" validate:"nonzero"`
}
