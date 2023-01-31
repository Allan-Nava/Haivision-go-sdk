package rtmp

type ResponseSourceRtmp struct {
	Name                 string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID                   string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address              string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol             string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port                 int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface     string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
	StreamName           string `json:"streamName" required:"true" validate:"nonnil,min=1"`
	RtmpMode             string `json:"rtmpMode" required:"true" validate:"nonnil,min=1"`
	State                string `json:"state" required:"true" validate:"nonnil,min=1"`
	SummaryStatusCode    string `json:"summaryStatusCode" required:"true" validate:"nonnil,min=1"`
	SummaryStatusDetails string `json:"summaryStatusDetails" required:"true" validate:"nonnil,min=1"`
}
