package route

import "encoding/json"

/*
ResponseRoutes — GET /api/gateway/[Device ID]/routes

	{
	  "data": [ <Array of Route Objects> ],
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

// ResponseRouteModel è una route come la restituisce il gateway: porta i campi di stato
// (`state`, `elapsedTime`, `summaryStatus*`) che NON vanno mai inviati in una create/update.
// Per costruire il corpo di una richiesta si usa RouteFields.
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

// ResponseCreateRoute è la risposta di una create/update/delete: {"status": "..."}.
type ResponseCreateRoute struct {
	Status string `json:"status"`
}

/*
RouteCommand è un elemento della risposta di start/stop route.

	[
	  {
	    "action": "command",
	    "command": "start-route",
	    "parameters": { "routeID": "[Route ID]" },
	    "deviceID": "[Device ID]",
	    "createdAt": [Date/time in Unix time],
	    "completedAt": 0,
	    "result": null,
	    "state": "pending",
	    "_id": "a5x4-7KEApdS0UuAUUCSog"
	  }
	]

`Result` è json.RawMessage perché la doc lo mostra a `null` e il gateway può metterci un
oggetto: un `string` fallirebbe l'unmarshal nel secondo caso.
*/
type RouteCommand struct {
	Action      string            `json:"action"`
	Command     string            `json:"command"`
	Parameters  CommandParameters `json:"parameters"`
	DeviceID    string            `json:"deviceID"`
	CreatedAt   int64             `json:"createdAt"`
	CompletedAt int64             `json:"completedAt"`
	Result      json.RawMessage   `json:"result"`
	State       string            `json:"state"`
	ID          string            `json:"_id"`
}

// ResponseRouteCommand è la risposta di start/stop route: un ARRAY top-level.
//
// Fino alla v1.x era una struct che wrappava in un campo `Response` senza tag JSON, quindi
// l'unmarshal della risposta reale falliva sempre con "cannot unmarshal array into Go value".
// Il nome corretto sostituisce anche `ResponseStartOrRoute`, a cui mancava "Stop".
type ResponseRouteCommand []RouteCommand

// Pending indica se almeno un comando è ancora in stato `pending`: i comandi sul gateway sono
// asincroni, un 200 non significa che la route sia già partita.
func (r ResponseRouteCommand) Pending() bool {
	for _, c := range r {
		if c.State == "pending" {
			return true
		}
	}
	return false
}
