package haivision

import "log"

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
func (o *Haivision) GetRoutes(deviceId string) error {
	resp, err := o.restyGet(GET_LIST_OF_ROUTES(deviceId), nil)
	if err != nil {
		return err
	}
	log.Println(resp)
	return nil
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
func (o *Haivision) CreateRoute(name string) error {
	return nil
}
