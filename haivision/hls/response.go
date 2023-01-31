package hls

type ResponseSourceHls struct {
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
	SrtFecCols       int    `json:"srtFecCols" required:"true" validate:"nonnil,min=1"`
	SrtFecRows       int    `json:"srtFecRows" required:"true" validate:"nonnil,min=1"`
	SrtFecArq        string `json:"srtFecArq" required:"true" validate:"nonnil,min=1"`
}

type ResponseDestinationHls struct {
	Name            string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID              string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address         string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol        string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	SegmentDuration string `json:"segmentDuration" required:"true" validate:"nonnil,min=1"`
	UseEncryption   bool   `json:"useEncryption" required:"true" validate:"nonnil,min=1"`
	SegmentsPerKey  int    `json:"segmentsPerKey" required:"true" validate:"nonnil,min=1"`
}
