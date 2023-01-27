package route

/*
	{
	  "data": { [
	    <Array of Route Objects>
	  ]},
	  "numPages": 1,
	  "numResults": 3,
	  "numActiveOutputConnections": 1,
	  "pendingRouteCreates": 1
	}
*/
type ResponseRoutes[TS RequestSource, TD RequestDestination] struct {
	Data                       []RouteModel[TS, TD] `json:"data"`
	NumPages                   int                  `json:"numPages"`
	NumResults                 int                  `json:"numResults"`
	NumActiveOutputConnections int                  `json:"numActiveOutputConnections"`
	PendingRouteCreates        int                  `json:"pendingRouteCreates"`
}

type RouteModel[TS RequestSource, TD RequestDestination] struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Source               TS     `json:"source"`
	Destinations         []TD   `json:"destinations"`
	ElapsedTime          string `json:"elapsedTime"`
	State                string `json:"state"`
	PendingUpdates       int    `json:"pendingUpdates"`
	SummaryStatusCode    string `json:"summaryStatusCode"`
	SummaryStatusDetails string `json:"summaryStatusDetails"`
	HasPendingDelete     bool   `json:"hasPendingDelete"`
}

// SRT
//type ResponseRoutesSRT struct {}
