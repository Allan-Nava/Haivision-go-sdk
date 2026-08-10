package route

import (
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/hls"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtmp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtsp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
	udprtp "github.com/Allan-Nava/Haivision-go-sdk/haivision/udp_rtp"
)

// Le quattro operazioni sulle route passano tutte dallo stesso endpoint
// `POST /api/devices/{deviceID}/updates` e si distinguono per il campo `action`.
// Valori presi dalla REST API Integrator's Reference (HMG 3.7.x), non inventati.
const (
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"

	// ElementTypeRoute è l'unico elementType usato da questo SDK.
	ElementTypeRoute = "route"
)

// Comandi accettati da POST /api/devices/{deviceID}/commands.
const (
	START_ROUTE = "start-route"
	STOP_ROUTE  = "stop-route"
)

// Comandi accettati nel campo `action` di una singola destinazione, dentro una update:
// è così che si avvia/ferma una destinazione, non con un endpoint dedicato.
const (
	DestinationActionStart = "start"
	DestinationActionStop  = "stop"
)

// --- type set dei protocolli --------------------------------------------------------------

// RequestSource sono i modelli di sorgente accettati in una create/update.
type RequestSource interface {
	srt.RequestSourceModelSRT | udprtp.RequestSourceModelUdpRtp | rtmp.RequestSourceModelRTMP | rtsp.RequestSourceModelRTSP
}

// RequestDestination sono i modelli di destinazione accettati in una create/update.
type RequestDestination interface {
	srt.RequestDestinationModelSrt | udprtp.RequestDestinationModelUdpRtp | rtmp.RequestDestinationModelRtmp | rtsp.RequestDestinationModelRtsp
}

// ResponseSource sono i modelli di sorgente restituiti dal gateway.
type ResponseSource interface {
	udprtp.ResponseSourceUdpRtp | srt.ResponseSourceSrt | rtmp.ResponseSourceRtmp
}

// ResponseDestination sono i modelli di destinazione restituiti dal gateway.
type ResponseDestination interface {
	udprtp.ResponseDestinationUdpRtp | srt.ResponseDestinationSrt | hls.ResponseDestinationHls
}

// --- corpo delle richieste ---------------------------------------------------------------

// RouteFields è il contenuto di `fields` in una create o update.
//
// StartRoute è un puntatore perché la doc lo prevede **solo nella create**: in una update va
// omesso, e con un `bool` non-puntatore `omitempty` non saprebbe distinguere "false" da "assente".
type RouteFields[TS RequestSource, TD RequestDestination] struct {
	Name         string `json:"name" validate:"required"`
	StartRoute   *bool  `json:"startRoute,omitempty"`
	Source       TS     `json:"source" validate:"required"`
	Destinations []TD   `json:"destinations" validate:"required,min=1,dive"`
}

/*
RequestCreateRoute — POST /api/devices/[Device ID]/updates

	{
	 "action": "create",
	 "deviceID": "[Device ID]",
	 "elementType": "route",
	 "fields": {
	  "name": "[Route name]",
	  "startRoute": [true,false],
	  "source": { <Source object> },
	  "destinations": [ <Destination object list> ]
	 }
	}

Nota sui tag: senza `json:"fields"` la struct serializzerebbe col nome del campo Go — `"Fields"`
maiuscolo — che il gateway non riconosce. È il bug corretto in v1.1.0; i test di
`test/wire_route_test.go` falliscono se ricompare.
*/
type RequestCreateRoute[TS RequestSource, TD RequestDestination] struct {
	Action      string              `json:"action" validate:"required"`
	DeviceID    string              `json:"deviceID" validate:"required"`
	ElementType string              `json:"elementType" validate:"required"`
	Fields      RouteFields[TS, TD] `json:"fields" validate:"required"`
}

/*
RequestUpdateRoute — POST /api/devices/[Device ID]/updates

	{
	 "action": "update",
	 "deviceID": "[Device ID]",
	 "elementType": "route",
	 "elementID": "[Route ID]",
	 "fields": { "name": …, "source": {…}, "destinations": […] }
	}

Rispetto alla create ha `elementID` e **non** ha `startRoute`.

È anche il modo documentato per avviare o fermare una **singola destinazione**: si manda una
update con `action` valorizzata sulla destinazione interessata (`start`/`stop`) e le altre
destinazioni incluse senza `action`. Ometterle le rimuoverebbe dalla route.
*/
type RequestUpdateRoute[TS RequestSource, TD RequestDestination] struct {
	Action      string              `json:"action" validate:"required"`
	DeviceID    string              `json:"deviceID" validate:"required"`
	ElementType string              `json:"elementType" validate:"required"`
	ElementID   string              `json:"elementID" validate:"required"`
	Fields      RouteFields[TS, TD] `json:"fields" validate:"required"`
}

/*
RequestDeleteRoute — POST /api/devices/[Device ID]/updates

	{
	 "action": "delete",
	 "deviceID": "[Device ID]",
	 "elementType": "route",
	 "elementID": "[Route ID]"
	}

Non ha `fields`: la delete non è generica, quindi non usa i type parameter.
*/
type RequestDeleteRoute struct {
	Action      string `json:"action" validate:"required"`
	DeviceID    string `json:"deviceID" validate:"required"`
	ElementType string `json:"elementType" validate:"required"`
	ElementID   string `json:"elementID" validate:"required"`
}

// CommandParameters sono i parametri di un comando su route.
// È un tipo con nome, non una struct anonima: una struct anonima porta i suoi tag dentro il
// tipo, quindi un consumer che volesse costruire il literal dovrebbe riscriverli identici.
type CommandParameters struct {
	RouteID string `json:"routeID" validate:"required"`
}

/*
RequestStartOrStopRoutes — POST /api/devices/[Device ID]/commands

	{
	 "deviceID": "[Device ID]",
	 "command": "[start-route|stop-route]",
	 "parameters": { "routeID": "[Route ID]" }
	}
*/
type RequestStartOrStopRoutes struct {
	DeviceID   string            `json:"deviceID" validate:"required"`
	Command    string            `json:"command" validate:"required,oneof=start-route stop-route"`
	Parameters CommandParameters `json:"parameters" validate:"required"`
}
