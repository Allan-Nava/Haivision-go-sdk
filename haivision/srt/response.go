package srt

// need to be finished
type ResponseSourceSrt struct {
	Name             string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID               string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address          string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol         string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port             int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	SrtPassPhrase    string `json:"srtPassPhrase" required:"true" validate:"nonnil,min=1"`
	SrtLatency       int    `json:"srtLatency" required:"true" validate:"nonnil,min=1"`
	SrtRcvBuf        int    `json:"srtRcvBuf" required:"true" validate:"nonnil,min=1"`
	SrtStreamId      string `json:"srtStreamId" required:"true" validate:"nonnil,min=1"`
	UseFec           bool   `json:"useFec" required:"true" validate:"nonnil,min=1"`
}

type ResponseDestinationSrt struct {
	Name             string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID               string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address          string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol         string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port             int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	RetainHeader     string `json:"retainHeader" required:"true" validate:"nonnil,min=1"`
	Mtu              int    `json:"mtu" required:"true" validate:"nonnil,min=1"`
	Ttl              string `json:"ttl" required:"true" validate:"nonnil,min=1"`
	Tos              string `json:"tos" required:"true" validate:"nonnil,min=1"`
	SrtEncryption    string `json:"srtEncryption" required:"true" validate:"nonnil,min=1"`
	SrtPassPhrase    string `json:"srtPassPhrase" required:"true" validate:"nonnil,min=1"`
	UseFEC           bool   `json:"useFEC" required:"true" validate:"nonnil,min=1"`
	SrtFecCols       int    `json:"srtFecCols" required:"true" validate:"nonnil,min=1"`
}
