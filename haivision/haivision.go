package haivision

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Haivision struct {
	Url        string
	RestClient *resty.Client
	Debug      bool
}

type IHaivisionClient interface {
	HealthCheck() error
}

func (o *Haivision) HealthCheck() error {
	resp, err := o.RestClient.R().
		SetHeader("Accept", "application/json").
		Get(o.Url)
	//
	if err != nil {
		return err
	}
	//
	if !strings.Contains(resp.Status(), "200") {
		o.DebugPrint(fmt.Sprintf("resp -> %v", resp))
		return errors.New("Could not connect haproxy")
	}
	return nil
}
