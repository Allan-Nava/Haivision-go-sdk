package haivision

import "fmt"

const (
	// AUTH
	SESSION     = "/api/session"
	DEVICE_INFO = "/api/devices"
	// ROUTES
	LIST_ROUTES  = "/api/gateway/%s/routes"
	CREATE_ROUTE = "/api/devices/%s/updates"
	// ROUTE_COMMMAND: il refuso nel nome (tre M) è mantenuto perché rinominare una costante
	// esportata è breaking — la rinomina è pianificata in v2.0.0 (`exported-naming-typos`).
	ROUTE_COMMMAND      = "/api/devices/%s/commands"
	ROUTE_CONFIGURATION = "/api/gateway/%s/routes/%s"
	// STATS
	ROUTES_STATISTICS = "/api/gateway/%s/statistics"
	// le statistiche per singolo client SRT stanno su un sotto-path dedicato
	ROUTES_CLIENT_STATISTICS = "/api/gateway/%s/statistics/client"
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
	// POST_ROUTE_COMMAND è l'endpoint dei comandi (start/stop route). Esisteva la costante ma
	// non l'helper, e StartOrStopRoute postava per errore su POST_CREATE_ROUTE (`/updates`).
	POST_ROUTE_COMMAND = func(deviceId string) string {
		return fmt.Sprintf(ROUTE_COMMMAND, deviceId)
	}
	//
	GET_ROUTES_STATISTICS = func(deviceId string) string {
		return fmt.Sprintf(ROUTES_STATISTICS, deviceId)
	}
	//
	GET_ROUTES_CLIENT_STATISTICS = func(deviceId string) string {
		return fmt.Sprintf(ROUTES_CLIENT_STATISTICS, deviceId)
	}
	//
	GET_ROUTE_CONFIGURATION = func(deviceId string, routeId string) string {
		return fmt.Sprintf(ROUTE_CONFIGURATION, deviceId, routeId)
	}
	//
)
