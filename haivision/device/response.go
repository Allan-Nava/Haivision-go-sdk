package device

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
	ID               string `json:"_id" `
	Type             string `json:"type" `
	IP               string `json:"ip" `
	Name             string `json:"name" `
	LastConnectedAt  int64  `json:"lastConnectedAt" `
	StatusCode       string `json:"statusCode" `
	Status           string `json:"status" `
	StatusDetails    string `json:"statusDetails" `
	SerialNumber     *string `json:"serialNumber" `
	Firmware         string `json:"firmware" `
	HasAdminError    bool   `json:"hasAdminError" `
	PendingSync      bool   `json:"pendingSync" `
	LastConnection   string `json:"lastConnection" `
	LastConnectionAt string `json:"lastConnectionAt" `
}
