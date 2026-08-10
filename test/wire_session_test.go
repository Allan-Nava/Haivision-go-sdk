package test

import (
	"encoding/json"
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/device"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/session"
)

// Payload letterale della doc: POST /api/session
func TestInitSessionResponse(t *testing.T) {
	documented := `{
	    "response": {
	        "type": "Session",
	        "message": "Session successfully started for haiadmin",
	        "sessionID": "sess-abc123",
	        "lastLoginDate": 1536777877871,
	        "numLoginFailures": 0
	    }
	}`
	var obj session.BaseResponseInitSession
	if err := json.Unmarshal([]byte(documented), &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if obj.Response.SessionID != "sess-abc123" {
		t.Errorf("SessionID = %q, atteso sess-abc123", obj.Response.SessionID)
	}
	if obj.Response.LastLoginDate != 1536777877871 {
		t.Errorf("LastLoginDate = %d", obj.Response.LastLoginDate)
	}
}

// Payload letterale della doc: GET /api/session. Nota `lastLoginDate` e `numLoginFailures`
// a `null`: vanno in campi non-puntatore, quindi restano a zero senza errore.
func TestSessionInfoResponseWithNulls(t *testing.T) {
	documented := `{
	    "sessionID": "sess-abc123",
	    "displayName": "Administrator",
	    "email": "haiadmin@localhost",
	    "roles": ["Administrator"],
	    "startAt": 1536937838919,
	    "expireAt": 1536938857529,
	    "lastLoginDate": null,
	    "numLoginFailures": null,
	    "isLicensed": true
	}`
	var obj session.ResponseSessionInfo
	if err := json.Unmarshal([]byte(documented), &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(obj.Roles) != 1 || obj.Roles[0] != "Administrator" {
		t.Errorf("Roles = %v", obj.Roles)
	}
	if obj.ExpireAt != 1536938857529 {
		t.Errorf("ExpireAt = %d", obj.ExpireAt)
	}
	if !obj.IsLicensed {
		t.Error("IsLicensed = false, atteso true")
	}
	if obj.LastLoginDate != 0 {
		t.Errorf("LastLoginDate = %d, atteso 0 (il payload ha null)", obj.LastLoginDate)
	}
}

// Payload letterale della doc: GET /api/devices — array top-level, `serialNumber` a null.
func TestDeviceInfoResponse(t *testing.T) {
	documented := `[{
	   "_id": "wlk9FE3_sOcu_9",
	   "type": "gateway",
	   "ip": "127.0.0.1",
	   "name": "Haivision Gateway",
	   "lastConnectedAt": 1675178018888,
	   "statusCode": "ok",
	   "status": "Online",
	   "statusDetails": "Connection has been established in the last 1 minutes.",
	   "serialNumber": null,
	   "firmware": "5.5.201209.1506",
	   "hasAdminError": false,
	   "pendingSync": false,
	   "lastConnection": "<1m"
	}]`
	var obj []device.ResponseDeviceInfo
	if err := json.Unmarshal([]byte(documented), &obj); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(obj) != 1 {
		t.Fatalf("attesi 1 device, ottenuti %d", len(obj))
	}
	d := obj[0]
	if d.ID != "wlk9FE3_sOcu_9" {
		t.Errorf("ID = %q (il campo JSON è `_id`)", d.ID)
	}
	if d.Type != "gateway" {
		t.Errorf("Type = %q", d.Type)
	}
	// serialNumber è *string proprio perché la doc lo mostra a null
	if d.SerialNumber != nil {
		t.Errorf("SerialNumber = %v, atteso nil", *d.SerialNumber)
	}
}
