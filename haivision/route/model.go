package route

import (
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

type RequestSource interface {
	srt.RequestSourceModelSRT | udprtp.RequestSourceModelUdpRtp
}

type RequestDestination interface {
	srt.RequestDestinationModelSrt | udprtp.RequestDestinationModelUdpRtp
}
