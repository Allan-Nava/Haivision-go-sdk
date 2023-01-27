package haivision

import (
	"log"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtmp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/rtsp"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
	udprtp "github.com/Allan-Nava/Haivision-go-sdk/haivision/udp_rtp"
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

func (o *Haivision) GetRoutesSRT(deviceId string) (*route.ResponseRoutes[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt], error) {
	log.Println("GetRoutesSRT ", deviceId)
	return nil, nil
}

func (o *Haivision) GetRoutesRTMP(deviceId string) (*route.ResponseRoutes[rtmp.RequestSourceModelRTMP, rtmp.RequestDestinationModelRtmp], error) {
	log.Println("GetRoutesRTMP ", deviceId)
	return nil, nil
}

func (o *Haivision) GetRoutesRtsp(deviceId string) (*route.ResponseRoutes[rtsp.RequestSourceModelRTSP, rtsp.RequestDestinationModelRtsp], error) {
	log.Println("GetRoutesRtsp ", deviceId)
	return nil, nil
}

func (o *Haivision) GetRoutesUdpRtp(deviceId string) (*route.ResponseRoutes[udprtp.RequestSourceModelUdpRtp, udprtp.RequestDestinationModelUdpRtp], error) {
	log.Println("GetRoutesRtsp ", deviceId)
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
