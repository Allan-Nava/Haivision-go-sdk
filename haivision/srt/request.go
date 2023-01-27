package srt

// https://doc.haivision.com/HMG3.7.5/rest-api-integrator-s-reference/rest-api-reference/object-model-reference/source-object-model?activetab=SRT%7EUDPandRTP

type RequestSourceModelSRT struct {
	Name             string `json:"name" required:"true" validate:"nonnil,min=1"`
	ID               string `json:"id" required:"true" validate:"nonnil,min=1"`
	Address          string `json:"address" required:"true" validate:"nonnil,min=1"`
	Protocol         string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	Port             int    `json:"port" required:"true" validate:"nonnil,min=1"`
	NetworkInterface string `json:"networkInterface" required:"true" validate:"nonnil,min=1"`
}
