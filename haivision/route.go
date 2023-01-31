package haivision

import (
	"encoding/json"
	"errors"
	"log"

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
/*
func (o *Haivision) GetRoutesSrt(deviceId string) (*route.ResponseRoutes[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt], error) {
	log.Println("GetRoutesSrt ", deviceId)
	return nil, nil
}

func (o *Haivision) GetRoutesRtmp(deviceId string) (*route.ResponseRoutes[rtmp.RequestSourceModelRTMP, rtmp.RequestDestinationModelRtmp], error) {
	log.Println("GetRoutesRtmp ", deviceId)
	return nil, nil
}

func (o *Haivision) GetRoutesRtsp(deviceId string) (*route.ResponseRoutes[rtsp.RequestSourceModelRTSP, rtsp.RequestDestinationModelRtsp], error) {
	log.Println("GetRoutesRtsp ", deviceId)
	return nil, nil
}

func (o *Haivision) GetRoutesUdpRtp(deviceId string) (*route.ResponseRoutes[udprtp.RequestSourceModelUdpRtp, udprtp.RequestDestinationModelUdpRtp], error) {
	log.Println("GetRoutesRtsp ", deviceId)
	return nil, nil
}*/

func (o *Haivision) GetRoutes(deviceId string) (*resty.Response, error) {
	log.Println("GetRoutes ", deviceId)
	resp, err := o.restyGet(GET_LIST_OF_ROUTES(deviceId), nil)
	if err != nil {
		return nil, err
	}
	log.Println("GetRoutes ", resp)
	return nil, nil
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

func (o *Haivision) CreateRouteSrt(deviceId string, rBody *route.RouteModel[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]) (*route.ResponseCreateRoute, error) {
	log.Println("CreateRouteSrt ", deviceId)

	if errs := validator.Validate(rBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(POST_CREATE_ROUTE(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	o.debugPrint(resp)
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
func (o *Haivision) CreateRouteRtmp(deviceId string, rBody *route.RouteModel[rtmp.RequestSourceModelRTMP, rtmp.RequestDestinationModelRtmp]) (*route.ResponseCreateRoute, error) {
	log.Println("CreateRouteRtmp ", deviceId)
	if errs := validator.Validate(rBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(POST_CREATE_ROUTE(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	o.debugPrint(resp)
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
func (o *Haivision) CreateRouteRtsp(deviceId string, rBody *route.RouteModel[rtsp.RequestSourceModelRTSP, rtsp.RequestDestinationModelRtsp]) (*route.ResponseCreateRoute, error) {
	log.Println("CreateRouteRtsp ", deviceId)
	if errs := validator.Validate(rBody); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(POST_CREATE_ROUTE(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	o.debugPrint(resp)
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (o *Haivision) CreateRouteUdpRtp(deviceId string, rBody *route.RouteModel[udprtp.RequestSourceModelUdpRtp, udprtp.RequestDestinationModelUdpRtp]) (*route.ResponseCreateRoute, error) {
	log.Println("CreateRouteUdpRtp ", deviceId)
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
func (o *Haivision) StartOrStopRoute(deviceId string, routeId string, command string) (*route.ResponseStartOrRoute, error) {
	log.Println("StartOrStopRoute ", deviceId, routeId, command)
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
	//
	resp, err := o.restyPost(POST_CREATE_ROUTE(deviceId), rBody)
	if err != nil {
		return nil, err
	}
	o.debugPrint(resp)
	//
	var obj route.ResponseStartOrRoute
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
