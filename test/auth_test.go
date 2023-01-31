package test

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/device"
)

func TestAuth(t *testing.T) {
	t.Log("TestAuth")
}

func TestDeviceInfo(t *testing.T) {
	var data = `[{
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
	byte_data := []byte(data)
	log.Println("TestDeviceInfo ", data, byte_data)
	//
	var obj []device.ResponseDeviceInfo
	if err := json.Unmarshal(byte_data, &obj); err != nil {
		//panic(err)
		log.Println(err)
	}
}
