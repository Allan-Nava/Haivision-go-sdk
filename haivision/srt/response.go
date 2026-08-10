package srt

// need to be finished
type ResponseSourceSrt struct {
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
}

type ResponseDestinationSrt struct {
	Name             string `json:"name"`
	ID               string `json:"id"`
	Address          string `json:"address"`
	Protocol         string `json:"protocol"`
	Port             int    `json:"port"`
	NetworkInterface string `json:"networkInterface"`
	RetainHeader     string `json:"retainHeader"`
	Mtu              int    `json:"mtu"`
	Ttl              string `json:"ttl"`
	Tos              string `json:"tos"`
	SrtEncryption    string `json:"srtEncryption"`
	SrtPassPhrase    string `json:"srtPassPhrase"`
	UseFEC           bool   `json:"useFEC"`
	SrtFecCols       int    `json:"srtFecCols"`
}
