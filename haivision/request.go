package haivision

type RequestSourceModelUDPandRTP struct {
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
// need to finish with all type: udp rtp, srt, hls
type RequestCreateRoute struct {
	Action      string `json:"action" required:"true" validate:"nonnil,min=1"`
	DeviceID    string `json:"deviceID" required:"true" validate:"nonnil,min=1"`
	ElementType string `json:"elementType" required:"true" validate:"nonnil,min=1"`
	Fields      struct {
		Name       string `json:"name" required:"true" validate:"nonnil,min=1"`
		StartRoute bool   `json:"startRoute" required:"true" validate:"nonnil,min=1"`
		Source     struct {
			Name     string `json:"name" required:"true" validate:"nonnil,min=1"`
			ID       string `json:"id" required:"true" validate:"nonnil,min=1"`
			Address  string `json:"address" required:"true" validate:"nonnil,min=1"`
			Protocol string `json:"protocol" required:"true" validate:"nonnil,min=1"`
		} `json:"source" required:"true" validate:"nonnil,min=1"`
		Destinations []struct {
		} `json:"destinations" required:"true" validate:"nonnil,min=1"`
	}
}

//
