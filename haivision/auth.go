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


// POST http://1.2.3.4:8081/v1/vhosts/default/apps/app:startPush
func (o *Haivision) InitSession(username string, password string)  error {
    
}