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
	DeviceID   string
	HType 	   string
	debug      bool
}

type IHaivisionClient interface {
	//
	HealthCheck() error
	IsDebug() bool
	// auth stuff
	InitSession(username string, password string) (*session.BaseResponseInitSession, error)
	GetSessionInfo() (*session.ResponseSessionInfo, error)
	GetDeviceInfo() (*[]device.ResponseDeviceInfo, error)
	// Streaming
	//GetRoutes(deviceId string) (*route.ResponseRoutes, error)
	GetRoutes(deviceId string) (*resty.Response, error)
	/*GetRoutesSrt(deviceId string) (*route.ResponseRoutes[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt], error)
	GetRoutesRtmp(deviceId string) (*route.ResponseRoutes[rtmp.RequestSourceModelRTMP, rtmp.RequestDestinationModelRtmp], error)
	GetRoutesRtsp(deviceId string) (*route.ResponseRoutes[rtsp.RequestSourceModelRTSP, rtsp.RequestDestinationModelRtsp], error)
	GetRoutesUdpRtp(deviceId string) (*route.ResponseRoutes[udprtp.RequestSourceModelUdpRtp, udprtp.RequestDestinationModelUdpRtp], error)*/
	//
	GetRouteConfiguration(deviceId string, routeId string) (*resty.Response, error) // *route.RouteModel[route.RequestSource, route.RequestDestination]
	// CreateRoute() error
	CreateRouteSrt(deviceId string, rBody *route.RouteModel[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]) (*route.ResponseCreateRoute, error)
	CreateRouteRtmp(deviceId string, rBody *route.RouteModel[rtmp.RequestSourceModelRTMP, rtmp.RequestDestinationModelRtmp]) (*route.ResponseCreateRoute, error)
	CreateRouteRtsp(deviceId string, rBody *route.RouteModel[rtsp.RequestSourceModelRTSP, rtsp.RequestDestinationModelRtsp]) (*route.ResponseCreateRoute, error)
	CreateRouteUdpRtp(deviceId string, rBody *route.RouteModel[udprtp.RequestSourceModelUdpRtp, udprtp.RequestDestinationModelUdpRtp]) (*route.ResponseCreateRoute, error)
	//
	StartOrStopRoute(deviceId string, routeId string, command string) (*route.ResponseStartOrRoute, error)
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
		err = fmt.Errorf("%v", resp)
		o.debugPrint(err)
		return nil, err
	}
	return resp, nil
}
