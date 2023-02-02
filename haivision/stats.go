package haivision

import (
	"encoding/json"
	"log"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/stats"
)

/*
GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]
cookie: sessionID: [Session ID]
*/

func (o *Haivision) GetRouteStatistics(deviceId string, routeId string) (*stats.ResponseRouteStatistics, error) {
	log.Println("GetStats ", deviceId, routeId)
	queryParams := map[string]string{
		"routeID": routeId,
	}
	resp, err := o.restyGet(GET_ROUTES_STATISTICS(deviceId), queryParams)
	if err != nil {
		return nil, err
	}
	var obj stats.ResponseRouteStatistics
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

/*
Requests
GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]&sourceID=[Source ID]
cookie: sessionID: [Session ID]
*/
