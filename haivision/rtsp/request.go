package rtsp

// Modelli di richiesta RTSP. Vedi le note sui tipi in srt/request.go.

type RequestSourceModelRTSP struct {
	Name             string  `json:"name" validate:"required"`
	ID               string  `json:"id,omitempty"`
	Address          string  `json:"address" validate:"required"`
	Protocol         string  `json:"protocol" validate:"required"`
	Port             int     `json:"port" validate:"required,min=1,max=65535"`
	NetworkInterface string  `json:"networkInterface"`
	Mode             *string `json:"mode,omitempty"`
}

type RequestDestinationModelRtsp struct {
	Name             string  `json:"name" validate:"required"`
	ID               string  `json:"id,omitempty"`
	Address          string  `json:"address" validate:"required"`
	Protocol         string  `json:"protocol" validate:"required"`
	Port             int     `json:"port" validate:"required,min=1,max=65535"`
	NetworkInterface string  `json:"networkInterface"`
	Action           *string `json:"action,omitempty"`
	Ttl              *int    `json:"ttl,omitempty"`
	Tos              *int    `json:"tos,omitempty"`
	Mtu              *int    `json:"mtu,omitempty"`
	RetainHeader     *bool   `json:"retainHeader,omitempty"`
}
