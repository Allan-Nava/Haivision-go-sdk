package haivision

import "fmt"

const (
	// AUTH
	SESSION     = "/api/session"
	DEVICE_INFO = "/api/devices"
	// ROUTES
	LIST_ROUTES  = "/api/gateway/%s/routes"
	CREATE_ROUTE = "/api/devices/%s/updates"
	//
)

var (
	//
	GET_LIST_OF_ROUTES = func(deviceId string) string {
		return fmt.Sprintf(LIST_ROUTES, deviceId)
	}
	//
	POST_CREATE_ROUTE = func(deviceId string) string {
		return fmt.Sprintf(CREATE_ROUTE, deviceId)
	}
	//
)
