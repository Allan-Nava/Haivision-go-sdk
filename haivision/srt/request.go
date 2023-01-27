package srt

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
POST API Requests
Use the following destinations model when issuing the Create a Route, Update a Route, and Start or Stop a Route's Destination API requests. Definition of each destination depends on the protocol.

UDP and RTP
SRT
HLS
Name	Type	Description
name	string	Name of destination for route. (Unique name with length 1-60.)
id	string	Optional when creating destination, required when updating. Destination ID.
address	string	
IP address of route destination. (0.0.0.0 for listener.)

If SRT path redundancy is used, this value is not needed. The srtNetworkBondingParams object contains the individual address assignments.

protocol	string	Destination protocol: srt.
port	number	
Port number of route destination: 1–65535.

If SRT path redundancy is used, this value is not needed. The srtNetworkBondingParams object contains the individual port assignments.

networkInterface	string	
Optional. Network interface name. (Empty string if auto.)

If SRT path redundancy is used, this value is not needed. The srtNetworkBondingParams object contains the individual network interface names.

retainHeader	Boolean	Optional. To retain headers for RTP tunneling through SRT.
action	string	Optional. Destination action: start, stop, or none.
mtu	number	
Optional. Destination MTU. Range = 280–1500.

If SRT path redundancy is used, this value is not needed. The srtNetworkBondingParams object contains the individual MTU values.

ttl	number	
Optional. Destination TTL. Range = 1–255.

If SRT path redundancy is used, this value is not needed. The srtNetworkBondingParams object contains the individual TTL values.

tos	number	
Optional. Destination ToS. Range = 0–255.

If SRT path redundancy is used, this value is not needed. The srtNetworkBondingParams object contains the individual ToS values.

srtEncryption	string	Optional. Encryption mode: AES128, AES256, or None.
srtPassPhrase	string	Optional. SRT Passphrase.
srtLatency

number	Optional. SRT latency.
srtMode	string	Optional. SRT mode: caller, listener, or rendezvous.
srtOverhead	string	Optional. Overhead used for SRT.
srtStreamID	string	
Optional. SRT Caller mode only. Stream ID string to identify the listener.

Note

Only available for Haivision SRT Gateway.

useFec	Boolean	SRT Caller mode with path redundancy disabled only. Enable FEC: true or false
srtFecCols	number	If useFEC is enabled, number of columns.
srtFecRows	number	If useFEC is enabled, number of rows.
srtFecLayout	string	If useFEC is enabled, FEC layout: saircase or even.
srtFecArq	string	If useFEC is enabled, FEC ARQ: never or onreq.
srtConnectionLimit	string	SRT Listener mode only. SRT Caller connection limit.
srtGroupMode	string	
Optional. SRT path redundancy mode: none, broadcast, backup, or any.

srtNetworkBondingParams	object list	If srtGroupMode is not set to none, array of SRT network bonding parameters.

*/
type RequestDestinationModelSrt struct {
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
}
