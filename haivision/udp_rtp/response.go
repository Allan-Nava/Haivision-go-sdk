package udprtp

type ResponseSourceUdpRtp struct {
	Name                 string `json:"name"`
	ID                   string `json:"id"`
	Address              string `json:"address"`
	Protocol             string `json:"protocol"`
	Port                 int    `json:"port"`
	NetworkInterface     string `json:"networkInterface"`
	RetainHeader         string `json:"retainHeader"`
	SourceAddress        string `json:"sourceAddress"`
	State                string `json:"state"`
	SummaryStatusCode    string `json:"summaryStatusCode"`
	SummaryStatusDetails string `json:"summaryStatusDetails"`
}

type ResponseDestinationUdpRtp struct {
	Name                     string `json:"name"`
	ID                       string `json:"id"`
	Address                  string `json:"address"`
	Protocol                 string `json:"protocol"`
	Port                     int    `json:"port"`
	NetworkInterface         string `json:"networkInterface"`
	RetainHeader             string `json:"retainHeader"`
	Action                   string `json:"action"`
	Mtu                      int    `json:"mtu"`
	Ttl                      string `json:"ttl"`
	Tos                      string `json:"tos"`
	Fec                      string `json:"fec"`
	PrompegFecLevel          string `json:"prompegFecLevel"`
	PrompegFeclsBlockAligned bool   `json:"prompegFeclsBlockAligned"`
	PrompegFecColumns        int    `json:"prompegFecColumns"`
	PrompegFecRows           int    `json:"prompegFecRows"`
	Shaping                  bool   `json:"shaping"`
	MaxBitrate               int    `json:"maxBitrate"`
}
