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
