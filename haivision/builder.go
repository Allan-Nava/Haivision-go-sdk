package haivision

import (
	"log"

	"github.com/go-resty/resty/v2"
)

func BuildHaivision(url string, debug bool) (*Haivision, error) {
	ovenClient := &Haivision{
		Url:        url,
		RestClient: resty.New(),
	}
	//
	if debug {
		ovenClient.RestClient.SetDebug(true)
		ovenClient.Debug = true
		log.Println("Debug mode is enabled for the Haproxy client ")
	}
	return ovenClient, nil
}

func (o *Haivision) DebugPrint(data interface{}) {
	if o.Debug {
		log.Println(data)
	}
}
