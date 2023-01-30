package route

import (
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtmp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtsp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
	udprtp "github.com/Allan-Nava/Haivision-go-sdk/haivision/udp_rtp"
)

/*
	{
		"action": "create",
		"deviceID": "[Device ID]",
		"elementType": "route",
		"fields":
		{
		"name": "[Route name]",
		"startRoute": [true,false],
		"source":
			{
			<Source object>
			},
			"destinations": [
			<Destination object list>
			// destinations udp and rtp
			// destinations srt
			// destinations hls
			]
		}
	}
*/
type RequestCreateRoute[TS RequestSource, TD RequestDestination] struct {
	Action      string `json:"action" required:"true" validate:"nonnil,min=1"`
	DeviceID    string `json:"deviceID" required:"true" validate:"nonnil,min=1"`
	ElementType string `json:"elementType" required:"true" validate:"nonnil,min=1"`
	Fields      struct {
		Name         string `json:"name" required:"true" validate:"nonnil,min=1"`
		StartRoute   bool   `json:"startRoute" required:"true" validate:"nonnil,min=1"`
		Source       TS     `json:"source" required:"true"`
		Destinations []TD   `json:"destinations" required:"true"`
	}
}

// request source
type RequestSource interface {
	srt.RequestSourceModelSRT | udprtp.RequestSourceModelUdpRtp | rtmp.RequestSourceModelRTMP | rtsp.RequestSourceModelRTSP
}

// request destination
type RequestDestination interface {
	srt.RequestDestinationModelSrt | udprtp.RequestDestinationModelUdpRtp | rtmp.RequestDestinationModelRtmp | rtsp.RequestDestinationModelRtsp
}

//

/*
Requests
POST /api/devices/[Device ID]/commands
cookie: sessionID: [Session ID]

{
 "deviceID": "[Device ID]",
 "command": "[command]",
 "parameters": {
  "routeID": "[Route ID]"
 }
}
*/

const (
	START_ROUTE = "start-route"
	STOP_ROUTE  = "stop-route"
)

type RequestStartOrStopRoutes struct {
	DeviceID   string `json:"deviceID" required:"true" validate:"nonnil,min=1"`
	Command    string `json:"command" required:"true" validate:"nonnil,min=1"`
	Parameters struct {
		RouteID string `json:"routeID" required:"true" validate:"nonnil,min=1"`
	}
}
