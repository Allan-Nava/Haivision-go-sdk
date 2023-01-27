package srt

/*
POST API Requests
Use the following destinations model when issuing the Create a Route, Update a Route, and Start or Stop a Route's Destination API requests. Definition of each destination depends on the protocol.

UDP and RTP
SRT
HLS
*/
/*
type RequestUdpRtpCreateRoute struct {
	Action      string `json:"action" required:"true" validate:"nonnil,min=1"`
	DeviceID    string `json:"deviceID" required:"true" validate:"nonnil,min=1"`
	ElementType string `json:"elementType" required:"true" validate:"nonnil,min=1"`
	Fields      struct {
		Name         string                       `json:"name" required:"true" validate:"nonnil,min=1"`
		StartRoute   bool                         `json:"startRoute" required:"true" validate:"nonnil,min=1"`
		Source       RequestSourceModelSRT        `json:"source" required:"true"`
		Destinations []RequestDestinationModelSrt `json:"destinations" required:"true"`
	}
}
*/

// https://doc.haivision.com/HMG3.7.5/rest-api-integrator-s-reference/rest-api-reference/object-model-reference/source-object-model?activetab=SRT%7EUDPandRTP

type RequestSourceModelSRT struct {
	Name             string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID               string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address          string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol         string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port             int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
}

/*
Use the following destinations model when issuing the Create a Route, Update a Route, and Start or Stop a Route's Destination API requests. Definition of each destination depends on the protocol.
*/
type RequestDestinationModelSrt struct {
	Name             string `json:"name" validate:"nonnil,min=1" required:"true"`
	ID               string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address          string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol         string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port             int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	RetainHeader     string `json:"retainHeader" required:"true" validate:"nonnil,min=1"`
	Action           string `json:"action" required:"true" validate:"nonnil,min=1"`
	Ttl              string `json:"ttl" required:"true" validate:"nonnil,min=1"`
	Tos              string `json:"tos" required:"true" validate:"nonnil,min=1"`
}
