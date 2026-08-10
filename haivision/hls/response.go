package hls

type ResponseSourceHls struct {
	Name             string `json:"name"`
	ID               string `json:"id"`
	Address          string `json:"address"`
	Protocol         string `json:"protocol"`
	Port             int    `json:"port"`
	NetworkInterface string `json:"networkInterface"`
	SrtPassPhrase    string `json:"srtPassPhrase"`
	SrtLatency       int    `json:"srtLatency"`
	SrtRcvBuf        int    `json:"srtRcvBuf"`
	SrtStreamId      string `json:"srtStreamId"`
	UseFec           bool   `json:"useFec"`
	SrtFecCols       int    `json:"srtFecCols"`
	SrtFecRows       int    `json:"srtFecRows"`
	SrtFecArq        string `json:"srtFecArq"`
}

type ResponseDestinationHls struct {
	Name            string `json:"name"`
	ID              string `json:"id"`
	Address         string `json:"address"`
	Protocol        string `json:"protocol"`
	SegmentDuration string `json:"segmentDuration"`
	UseEncryption   bool   `json:"useEncryption"`
	SegmentsPerKey  int    `json:"segmentsPerKey"`
}
