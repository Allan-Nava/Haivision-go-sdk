package haivision

import (
	"log"

	"github.com/go-resty/resty/v2"
)

func BuildHaivision(url string, debug bool, header *HeaderConfigurator) (*Haivision, error) {
	haivisionClient := &Haivision{
		Url:        url,
		RestClient: resty.New(),
	}
	// You can override all below settings and options at request level if you want to
	//--------------------------------------------------------------------------------
	// Host URL for all request. So you can use relative URL in the request
	haivisionClient.restClient.SetHostURL(url)
	if header != nil {
		// Headers for all request
		for h, v := range header.GetHeaders() {
			haivisionClient.restClient.SetHeader(h, v)
		}
	}
	//
	if debug {
		haivisionClient.RestClient.SetDebug(true)
		haivisionClient.Debug = true
		log.Println("Debug mode is enabled for the haivision client ")
	}
	return haivisionClient, nil
}

func (o *Haivision) DebugPrint(data interface{}) {
	if o.Debug {
		log.Println(data)
	}
}
