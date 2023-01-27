package rtmp

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

type RequestSourceModelRTMP struct {
	Name             string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID               string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address          string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol         string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port             int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	StreamName       string `json:"streamName" required:"true" validate:"nonnil,min=1"`
	RtmpMode         string `json:"RtmpMode" required:"true" validate:"nonnil,min=1"`
}

/*
Use the following destinations model when issuing the Create a Route, Update a Route, and Start or Stop a Route's Destination API requests. Definition of each destination depends on the protocol.
*/
type RequestDestinationModelRtmp struct {
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
