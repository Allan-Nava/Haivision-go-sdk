package haivision

import "fmt"

// Path della REST API del gateway. Riferimento: REST API Integrator's Reference (HMG 3.7.x).
const (
	// AUTH
	SESSION     = "/api/session"
	DEVICE_INFO = "/api/devices"
	// ROUTES
	LIST_ROUTES = "/api/gateway/%s/routes"
	// ROUTE_UPDATES è l'endpoint di create, update e delete di una route: le tre operazioni si
	// distinguono per il campo `action` del body, non per il path.
	ROUTE_UPDATES = "/api/devices/%s/updates"
	// ROUTE_COMMANDS è l'endpoint dei comandi (start-route / stop-route).
	ROUTE_COMMANDS      = "/api/devices/%s/commands"
	ROUTE_CONFIGURATION = "/api/gateway/%s/routes/%s"
	// STATS
	ROUTES_STATISTICS        = "/api/gateway/%s/statistics"
	ROUTES_CLIENT_STATISTICS = "/api/gateway/%s/statistics/client"
)

var (
	GET_LIST_OF_ROUTES = func(deviceID string) string {
		return fmt.Sprintf(LIST_ROUTES, deviceID)
	}
	// POST_ROUTE_UPDATES: create/update/delete di una route.
	POST_ROUTE_UPDATES = func(deviceID string) string {
		return fmt.Sprintf(ROUTE_UPDATES, deviceID)
	}
	// POST_ROUTE_COMMAND: start-route / stop-route.
	POST_ROUTE_COMMAND = func(deviceID string) string {
		return fmt.Sprintf(ROUTE_COMMANDS, deviceID)
	}
	GET_ROUTES_STATISTICS = func(deviceID string) string {
		return fmt.Sprintf(ROUTES_STATISTICS, deviceID)
	}
	GET_ROUTES_CLIENT_STATISTICS = func(deviceID string) string {
		return fmt.Sprintf(ROUTES_CLIENT_STATISTICS, deviceID)
	}
	GET_ROUTE_CONFIGURATION = func(deviceID, routeID string) string {
		return fmt.Sprintf(ROUTE_CONFIGURATION, deviceID, routeID)
	}
)
