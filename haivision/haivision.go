package haivision

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Haivision struct {
	Url        string
	restClient *resty.Client
	debug      bool
}

type IHaivisionClient interface {
	HealthCheck() error
	IsDebug() bool
	// auth stuff
	InitSession(username string, password string) error
	// Streaming
}

func (o *Haivision) HealthCheck() error {
	resp, err := o.restClient.R().
		SetHeader("Accept", "application/json").
		Get(o.Url)
	//
	if err != nil {
		return err
	}
	//
	if !strings.Contains(resp.Status(), "200") {
		o.debugPrint(fmt.Sprintf("resp -> %v", resp))
		return errors.New("Could not connect haproxy")
	}
	return nil
}
//

func (o *Haivision) IsDebug() bool {
	return o.debug
}


// 
func (o *Haivision) IsDebug() bool {
	return o.debug
}

// Resty Methods

func (o *Haivision) restyPost(url string, body interface{}) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		return nil, err
	}
	if !strings.Contains(resp.Status(), "200") {
		err = fmt.Errorf("%v",resp)
		o.debugPrint(err)
		return nil, err
	}
	return resp, nil
}

func (o *Haivision) restyGet(url string, queryParams map[string]string) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetQueryParams(queryParams).
		Get(url)
	//
	if err != nil {
		return nil, err
	}
	if !strings.Contains(resp.Status(), "200") {
		resp.Body()
		err = fmt.Errorf("%v",resp)
		o.debugPrint(err)
		return nil, err
	}
	return resp, nil
}

