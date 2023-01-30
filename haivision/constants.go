package haivision

import "fmt"

const (
	// AUTH
	SESSION     = "/api/session"
	DEVICE_INFO = "/api/devices"
	// ROUTES
	LIST_ROUTES         = "/api/gateway/%s/routes"
	CREATE_ROUTE        = "/api/devices/%s/updates"
	ROUTE_COMMMAND      = "/api/devices/%s/commands"
	ROUTE_CONFIGURATION = "/api/gateway/%s/routes/%s"
	// STATS
	ROUTES_STATISTICS = "/api/gateway/%s/statistics"
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
	GET_ROUTES_STATISTICS = func(deviceId string) string {
		return fmt.Sprintf(ROUTES_STATISTICS, deviceId)
	}
	//
	GET_ROUTE_CONFIGURATION = func(deviceId string, routeId string) string {
		return fmt.Sprintf(ROUTE_CONFIGURATION, deviceId, routeId)
	}
	//
)
