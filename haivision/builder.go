package haivision

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// Builder is used to build a new haivision client
func BuildHaivision(url string, debug bool, username string, password string, header *HeaderConfigurator, insecure *bool) (*Haivision, error) {
	// init haivision
	haivisionClient := &Haivision{
		Url:        url,
		restClient: resty.New(),
	}
	//
	if debug {
		haivisionClient.restClient.SetDebug(true)
		haivisionClient.debug = true
		log.Println("Debug mode is enabled for the haivision client ")
	}
	// You can override all below settings and options at request level if you want to
	//--------------------------------------------------------------------------------
	// Host URL for all request. So you can use relative URL in the request
	haivisionClient.restClient.SetBaseURL(url)
	if insecure != nil {
		haivisionClient.restClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		log.Println("Insecure mode is enabled for the haivision client ")
	}
	respSessionId, err := haivisionClient.InitSession(username, password)
	if err != nil {
		return nil, err
	}
	// Set Cookie for all request
	haivisionClient.restClient.SetCookie(&http.Cookie{
		Name:  "sessionID",
		Value: respSessionId.Response.SessionID,
	})
	//
	deviceResponse, err := haivisionClient.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	if deviceResponse != nil {
		haivisionClient.DeviceID = (*deviceResponse)[0].ID
		haivisionClient.HType = (*deviceResponse)[0].Type
	}
	//
	if header != nil {
		// Headers for all request
		for h, v := range header.GetHeaders() {
			haivisionClient.restClient.SetHeader(h, v)
		}
	}
	return haivisionClient, nil
}

func (o *Haivision) debugPrint(data interface{}) {
	if o.debug {
		log.Println(data)
	}
}
