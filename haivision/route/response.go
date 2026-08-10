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
type ResponseRoutes[TS ResponseSource, TD ResponseDestination] struct {
	Data                       []ResponseRouteModel[TS, TD] `json:"data"`
	NumPages                   int                          `json:"numPages"`
	NumResults                 int                          `json:"numResults"`
	NumActiveOutputConnections int                          `json:"numActiveOutputConnections"`
	PendingRouteCreates        int                          `json:"pendingRouteCreates"`
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

type ResponseRouteModel[TS ResponseSource, TD ResponseDestination] struct {
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

type ResponseCreateRoute struct {
	Status string `json:"status"`
}

/*
Response
[
  {
    "action": "command",
    "command": "start-route",
    "parameters": {
      "routeID": "[Route ID]"
    },
    "deviceID": "[Device ID]",
    "createdAt": [Date/time in Unix time],
    "completedAt": 0,
    "result": null,
    "state": "pending",
    "_id": "a5x4-7KEApdS0UuAUUCSog"
  }
]
*/

type ResponseStartOrRoute struct {
	Response []struct {
		Action     string `json:"action"`
		Command    string `json:"command"`
		Parameters struct {
			RouteID string `json:"routeID"`
		} `json:"parameters"`
		DeviceID    string `json:"deviceID"`
		CreatedAt   int64  `json:"createdAt"`
		CompletedAt int64  `json:"completedAt"`
		Result      string `json:"result"`
		State       string `json:"state"`
		ID          string `json:"_id"`
	}
}
