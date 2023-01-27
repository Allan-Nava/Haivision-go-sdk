package haivision

/*{
    "response": {
        "type": "Session",
        "message": "Session successfully started for haiadmin",
        "sessionID": "[Session ID]",
        "lastLoginDate": 1536777877871,
        "numLoginFailures": 0
    }
}*/

type ResponseInitSession struct {
	Type             string `json:"type" validate:"nonzero"`
	Message          string `json:"message" validate:"nonzero"`
	SessionID        string `json:"sessionID" validate:"nonzero"`
	LastLoginDate    int64  `json:"lastLoginDate" validate:"nonzero"`
	NumLoginFailures int    `json:"numLoginFailures" validate:"nonzero"`
}
