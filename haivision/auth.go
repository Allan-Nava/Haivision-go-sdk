package haivision

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

// POST 
func (o *Haivision) InitSession(username string, password string) error {
	return nil
}

// GET 
func (o *Haivision) GetDeviceInfo() error {
    return nil
}