package rtmp

type ResponseSourceRtmp struct {
	Name                 string `json:"name"`
	ID                   string `json:"id"`
	Address              string `json:"address"`
	Protocol             string `json:"protocol"`
	Port                 int    `json:"port"`
	NetworkInterface     string `json:"networkInterface"`
	StreamName           string `json:"streamName"`
	RtmpMode             string `json:"rtmpMode"`
	State                string `json:"state"`
	SummaryStatusCode    string `json:"summaryStatusCode"`
	SummaryStatusDetails string `json:"summaryStatusDetails"`
}
