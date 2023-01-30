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

/*
type Response interface {
	GetData() []Route
	GetNumPages() int
	GetNumResults() int
	GetNumActiveOutputConnections() int
	GetPendingRouteCreates() int
}

type Route interface {
	GetID() string
	GetName() string
	GetSource() RequestSource
	GetDestinations() []RequestDestination
	GetElapsedTime() string
	GetState() string
	GetPendingUpdates() int
	GetSummaryStatusCode() string
	GetSummaryStatusDetails() string
	GetHasPendingDelete() bool
}

func (r ResponseRoutes[TS RequestSource, TD RequestDestination]) GetData() []Route {
	// convert []RouteModel[TS, TD] to []Route and return
}

func (r ResponseRoutes[TS RequestSource, TD RequestDestination]) GetNumPages() int {
	return r.NumPages
}

// implement the rest of the methods for ResponseRoutes

func (r RouteModel[TS RequestSource, TD RequestDestination]) GetID() string {
	return r.ID
}
*/

type ResponseCreateRoute struct {
	Status string `json:"status"`
}
