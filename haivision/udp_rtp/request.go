package udprtp

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
	}*/
//
/*type RequestUdpRtpCreateRoute struct {
	Action      string `json:"action" required:"true" validate:"nonnil,min=1"`
	DeviceID    string `json:"deviceID" required:"true" validate:"nonnil,min=1"`
	ElementType string `json:"elementType" required:"true" validate:"nonnil,min=1"`
	Fields      struct {
		Name         string                          `json:"name" required:"true" validate:"nonnil,min=1"`
		StartRoute   bool                            `json:"startRoute" required:"true" validate:"nonnil,min=1"`
		Source       RequestSourceModelUdpRtp        `json:"source" required:"true"`
		Destinations []RequestDestinationModelUdpRtp `json:"destinations" required:"true"`
	}
}*/

type RequestSourceModelUdpRtp struct {
	Name             string `json:"name" validate:"nonnil,min=1" required:"true"`
	ID               string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address          string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol         string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port             int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	RetainHeader     string `json:"retainHeader" required:"true" validate:"nonnil,min=1"`
	SourceAddress    string `json:"sourceAddress" required:"true" validate:"nonnil,min=1"`
	Fec              string `json:"fec" required:"true" validate:"nonnil,min=1"`
}

type RequestDestinationModelUdpRtp struct {
	Name                     string  `json:"name" validate:"nonnil,min=1" required:"true"`
	ID                       string  `json:"id" required:"true" validate:"nonnil,min=1"`
	Address                  string  `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol                 string  `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port                     int     `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface         string  `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	RetainHeader             string  `json:"retainHeader" required:"true" validate:"nonnil,min=1"`
	Action                   string  `json:"action" required:"true" validate:"nonnil,min=1"`
	Ttl                      string  `json:"ttl" required:"true" validate:"nonnil,min=1"`
	Tos                      string  `json:"tos" required:"true" validate:"nonnil,min=1"`
	Fec                      string  `json:"fec" required:"true" validate:"nonnil,min=1"`
	PrompegFecLevel          *string `json:"prompegFecLevel" `
	PrompegFecIsBlockAligned *string `json:"prompegFecIsBlockAligned" `
	PrompegFecColumns        *string `json:"prompegFecColumns" `
	PrompegFecRows           *string `json:"prompegFecRows" `
	Shaping                  *string `json:"shaping" `
	MaxBitrate               *string `json:"maxBitrate" `
}
