package haivision

import (
	"context"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/device"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/session"
)

/*
InitSession — POST /api/session

	{ "username": "[user name]", "password": "[password]" }

Risposta:

	{
	  "response": {
	    "type": "Session",
	    "message": "Session successfully started for haiadmin",
	    "sessionID": "[Session ID]",
	    "lastLoginDate": 1536777877871,
	    "numLoginFailures": 0
	  }
	}

Di norma non serve chiamarla: lo fa Connect, che imposta anche il cookie di sessione.
*/
func (c *Client) InitSession(ctx context.Context, username, password string) (*session.BaseResponseInitSession, error) {
	body := &session.RequestInitSession{Username: username, Password: password}
	if err := validateRequest("InitSession", body); err != nil {
		return nil, err
	}
	resp, err := c.post(ctx, SESSION, body)
	if err != nil {
		return nil, err
	}
	c.debugResponse("InitSession", resp)
	var obj session.BaseResponseInitSession
	if err := decode(resp, "InitSession", &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

/*
GetDeviceInfo — GET /api/devices (cookie sessionID)

Risposta: array top-level di device.

	[ { "_id": "[Device ID]", "type": "gateway", "ip": "127.0.0.1", … } ]
*/
func (c *Client) GetDeviceInfo(ctx context.Context) ([]device.ResponseDeviceInfo, error) {
	resp, err := c.get(ctx, DEVICE_INFO, nil)
	if err != nil {
		return nil, err
	}
	c.debugResponse("GetDeviceInfo", resp)
	var obj []device.ResponseDeviceInfo
	if err := decode(resp, "GetDeviceInfo", &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// GetSessionInfo — GET /api/session. `ExpireAt` dice quando la sessione scade: il gateway non
// la rinnova da solo e l'SDK non lo fa per te, quindi è il campo da guardare per decidere
// quando richiamare Connect.
func (c *Client) GetSessionInfo(ctx context.Context) (*session.ResponseSessionInfo, error) {
	resp, err := c.get(ctx, SESSION, nil)
	if err != nil {
		return nil, err
	}
	c.debugResponse("GetSessionInfo", resp)
	var obj session.ResponseSessionInfo
	if err := decode(resp, "GetSessionInfo", &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
