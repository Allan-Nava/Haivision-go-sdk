package rtmp

// Modelli di richiesta RTMP. Vedi le note sui tipi in srt/request.go: i campi opzionali sono
// puntatori con omitempty, `ttl`/`tos`/`mtu` sono numerici come negli esempi della doc.

type RequestSourceModelRTMP struct {
	Name             string `json:"name" validate:"required"`
	ID               string `json:"id,omitempty"`
	Address          string `json:"address" validate:"required"`
	Protocol         string `json:"protocol" validate:"required"`
	Port             int    `json:"port" validate:"required,min=1,max=65535"`
	NetworkInterface string `json:"networkInterface"`
	StreamName       string `json:"streamName" validate:"required"`
	// RtmpMode: il campo JSON è `rtmpMode` minuscolo — fino alla v1.x era `RtmpMode`, quindi
	// il gateway non lo riconosceva.
	RtmpMode *string `json:"rtmpMode,omitempty"`
	Mode     *string `json:"mode,omitempty"`
}

type RequestDestinationModelRtmp struct {
	Name             string  `json:"name" validate:"required"`
	ID               string  `json:"id,omitempty"`
	Address          string  `json:"address" validate:"required"`
	Protocol         string  `json:"protocol" validate:"required"`
	Port             int     `json:"port" validate:"required,min=1,max=65535"`
	NetworkInterface string  `json:"networkInterface"`
	StreamName       *string `json:"streamName,omitempty"`
	Action           *string `json:"action,omitempty"`
	Ttl              *int    `json:"ttl,omitempty"`
	Tos              *int    `json:"tos,omitempty"`
	Mtu              *int    `json:"mtu,omitempty"`
	RetainHeader     *bool   `json:"retainHeader,omitempty"`
}
