package udprtp

type ResponseSourceUdpRtp struct {
	Name                 string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID                   string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address              string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol             string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port                 int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface     string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	RetainHeader         string `json:"retainHeader" required:"true" validate:"nonnil,min=1"`
	SourceAddress        string `json:"sourceAddress" required:"true" validate:"nonnil,min=1"`
	State                string `json:"state" required:"true" validate:"nonnil,min=1"`
	SummaryStatusCode    string `json:"summaryStatusCode" required:"true" validate:"nonnil,min=1"`
	SummaryStatusDetails string `json:"summaryStatusDetails" required:"true" validate:"nonnil,min=1"`
}

type ResponseDestinationUdpRtp struct {
	Name                     string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID                       string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address                  string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol                 string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port                     int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface         string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	RetainHeader             string `json:"retainHeader" required:"true" validate:"nonnil,min=1"`
	Action                   string `json:"action" required:"true" validate:"nonnil,min=1"`
	Mtu                      int    `json:"mtu" required:"true" validate:"nonnil,min=1"`
	Ttl                      string `json:"ttl" required:"true" validate:"nonnil,min=1"`
	Tos                      string `json:"tos" required:"true" validate:"nonnil,min=1"`
	Fec                      string `json:"fec" required:"true" validate:"nonnil,min=1"`
	PrompegFecLevel          string `json:"prompegFecLevel"`
	PrompegFeclsBlockAligned bool   `json:"prompegFeclsBlockAligned"`
	PrompegFecColumns        int    `json:"prompegFecColumns"`
	PrompegFecRows           int    `json:"prompegFecRows"`
	Shaping                  bool   `json:"shaping"`
	MaxBitrate               int    `json:"maxBitrate"`
}
