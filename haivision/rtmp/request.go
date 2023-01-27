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
	StreamName 	     string `json:"streamName" required:"true" validate:"nonnil,min=1"`
	RtmpMode		 string `json:"RtmpMode" required:"true" validate:"nonnil,min=1"`
}