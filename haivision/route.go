package haivision

import (
	"encoding/json"
	"errors"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtmp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtsp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
	udprtp "github.com/Allan-Nava/Haivision-go-sdk/haivision/udp_rtp"
	"github.com/go-resty/resty/v2"
	"gopkg.in/validator.v2"
)

/*
Requests
GET /api/gateway/[Device ID]/routes
cookie: sessionID: [Session ID]

Response

	{
	  "data": { [
	    <Array of Route Objects>
	  ]},
	  "numPages": 1,
	  "numResults": 3,
	  "numActiveOutputConnections": 1,
	  "pendingRouteCreates": 1
	}
*/

func (o *haivisionSdk) GetRoutes(deviceId string) (*resty.Response, error) {
	o.debugf("GetRoutes device=%s", deviceId)
	resp, err := o.restyGet(GET_LIST_OF_ROUTES(deviceId), nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (o *haivisionSdk) GetRouteConfiguration(deviceId string, routeId string) (*resty.Response, error) {
	o.debugf("GetRouteConfiguration device=%s route=%s", deviceId, routeId)
	resp, err := o.restyGet(GET_ROUTE_CONFIGURATION(deviceId, routeId), nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

/*
Use this command to create an individual route.

Authorizations: Administrator, Operator
POST /api/devices/[Device ID]/updates
cookie: sessionID: [Session ID]

	{
		 "action": "create",
		 "deviceID": "[Device ID]",
		 "elementType": "route",
		 "fields":
		   {
		    "name": "[Route name]",
		    "startRoute": [true,false],
		    "source":
		      {
		       <Source object>
		      },
		      "destinations": [
		       <Destination object list>
		      ]
		   }
	}

Response

	{
	  "status": "[success message]"
	}
*/

func (o *haivisionSdk) CreateRouteSrt(deviceId string, rBody *route.RouteModel[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]) (*route.ResponseCreateRoute, error) {
	o.debugf("CreateRouteSrt device=%s", deviceId)

	if errs := validator.Validate(rBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(POST_CREATE_ROUTE(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	o.debugResponse("CreateRoute", resp)
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
func (o *haivisionSdk) CreateRouteRtmp(deviceId string, rBody *route.RouteModel[rtmp.RequestSourceModelRTMP, rtmp.RequestDestinationModelRtmp]) (*route.ResponseCreateRoute, error) {
	o.debugf("CreateRouteRtmp device=%s", deviceId)
	if errs := validator.Validate(rBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(POST_CREATE_ROUTE(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	o.debugResponse("CreateRoute", resp)
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
func (o *haivisionSdk) CreateRouteRtsp(deviceId string, rBody *route.RouteModel[rtsp.RequestSourceModelRTSP, rtsp.RequestDestinationModelRtsp]) (*route.ResponseCreateRoute, error) {
	o.debugf("CreateRouteRtsp device=%s", deviceId)
	if errs := validator.Validate(rBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(POST_CREATE_ROUTE(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	o.debugResponse("CreateRoute", resp)
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (o *haivisionSdk) CreateRouteUdpRtp(deviceId string, rBody *route.RouteModel[udprtp.RequestSourceModelUdpRtp, udprtp.RequestDestinationModelUdpRtp]) (*route.ResponseCreateRoute, error) {
	o.debugf("CreateRouteUdpRtp device=%s", deviceId)
	if errs := validator.Validate(rBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(POST_CREATE_ROUTE(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// Start or Stop a route
func (o *haivisionSdk) StartOrStopRoute(deviceId string, routeId string, command string) (*route.ResponseStartOrRoute, error) {
	o.debugf("StartOrStopRoute device=%s route=%s command=%s", deviceId, routeId, command)
	if command != route.START_ROUTE && command != route.STOP_ROUTE {
		return nil, errors.New("command must be start-route or stop-route")
	}
	rBody := &route.RequestStartOrStopRoutes{
		Command:  command,
		DeviceID: deviceId,
		Parameters: struct {
			RouteID string "json:\"routeID\" required:\"true\" validate:\"nonnil,min=1\""
		}{
			RouteID: routeId,
		},
	}
	if errs := validator.Validate(rBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	// endpoint dei COMANDI, non /updates: sono due endpoint diversi del gateway
	resp, err := o.restyPost(POST_ROUTE_COMMAND(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	o.debugResponse("StartOrStopRoute", resp)
	//
	var obj route.ResponseStartOrRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
