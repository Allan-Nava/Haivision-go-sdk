package udprtp

// Modelli di richiesta UDP/RTP.
//
// Esempio letterale dalla doc, che fissa i tipi:
//
//	{ "name": "Destination1Name", "id": "…", "action": "stop", "protocol": "udp",
//	  "port": 1111, "networkInterface": "", "address": "10.0.65.10",
//	  "ttl": 64, "mtu": 1496, "tos": 136, "encryption": "none", "fec": "none",
//	  "shaping": false, "maxBitrate": 10000 }
//
// Fino alla v1.x `shaping` e `maxBitrate` erano `*string`: il gateway li vuole booleano e
// numero.

type RequestSourceModelUdpRtp struct {
	Name             string  `json:"name" validate:"required"`
	ID               string  `json:"id,omitempty"`
	Address          string  `json:"address" validate:"required"`
	Protocol         string  `json:"protocol" validate:"required"`
	Port             int     `json:"port" validate:"required,min=1,max=65535"`
	NetworkInterface string  `json:"networkInterface"`
	Mode             *string `json:"mode,omitempty"`
	SourceAddress    *string `json:"sourceAddress,omitempty"`
	Fec              *string `json:"fec,omitempty"`
	RetainHeader     *bool   `json:"retainHeader,omitempty"`
}

type RequestDestinationModelUdpRtp struct {
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
	Fec              *string `json:"fec,omitempty"`
	Encryption       *string `json:"encryption,omitempty"`
	// Prompeg FEC: numerici, non stringhe
	PrompegFecLevel          *string `json:"prompegFecLevel,omitempty"`
	PrompegFecIsBlockAligned *bool   `json:"prompegFecIsBlockAligned,omitempty"`
	PrompegFecColumns        *int    `json:"prompegFecColumns,omitempty"`
	PrompegFecRows           *int    `json:"prompegFecRows,omitempty"`
	Shaping                  *bool   `json:"shaping,omitempty"`
	MaxBitrate               *int    `json:"maxBitrate,omitempty"`
}
