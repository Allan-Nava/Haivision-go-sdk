package haivision

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/device"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtmp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtsp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/session"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
	udprtp "github.com/Allan-Nava/Haivision-go-sdk/haivision/udp_rtp"
)

type Haivision struct {
	Url        string
	restClient *resty.Client
	debug      bool
}

type IHaivisionClient interface {
	//
	HealthCheck() error
	IsDebug() bool
	// auth stuff
	InitSession(username string, password string) (*session.BaseResponseInitSession, error)
	GetSessionInfo() (*session.ResponseSessionInfo, error)
	GetDeviceInfo() (*device.BaseResponseDeviceInfo, error)
	// Streaming
	// GetRoutes(deviceId string) (route.ResponseRoutes[TS, TD], error)
	GetRoutesSRT(deviceId string) (*route.ResponseRoutes[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt], error)
	GetRoutesRTMP(deviceId string) (*route.ResponseRoutes[rtmp.RequestSourceModelRTMP, rtmp.RequestDestinationModelRtmp], error)
	GetRoutesRtsp(deviceId string) (*route.ResponseRoutes[rtsp.RequestSourceModelRTSP, rtsp.RequestDestinationModelRtsp], error)
	GetRoutesUdpRtp(deviceId string) (*route.ResponseRoutes[udprtp.RequestSourceModelUdpRtp, udprtp.RequestDestinationModelUdpRtp], error)
	// CreateRoute() error
	//
}

func (o *Haivision) HealthCheck() error {
	_, err := o.restyGet(o.Url, nil)
	if err != nil {
		return nil
	}
	return nil
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
		err = fmt.Errorf("%v", resp)
		o.debugPrint(err)
		return nil, err
	}
	return resp, nil
}

// get request
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
		err = fmt.Errorf("%v", resp)
		o.debugPrint(err)
		return nil, err
	}
	return resp, nil
}
