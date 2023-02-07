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

func (o *haivisionSdk) GetRouteStatistics(deviceId string, routeId string) (*stats.ResponseRouteStatistics, error) {
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

func (o *haivisionSdk) GetSourceStatistics(deviceId string, routeId string, sourceId string) (*stats.ResponseSourceStatistics, error) {
	log.Println("GetStats ", deviceId, routeId, sourceId)
	queryParams := map[string]string{
		"routeID":  routeId,
		"sourceID": sourceId,
	}
	resp, err := o.restyGet(GET_ROUTES_STATISTICS(deviceId), queryParams)
	if err != nil {
		return nil, err
	}
	var obj stats.ResponseSourceStatistics
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

/*

GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]&destinationID=
  [Destination ID]
cookie: sessionID: [Session ID]

or

GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]&destinationName=
  [Destination Name]
cookie: sessionID: [Session ID]

*/

func (o *haivisionSdk) GetDestinationStatisticsById(deviceId string, routeId string, destinationID string) (*stats.ResponseDestinationStatistics, error) {
	log.Println("GetStats ", deviceId, routeId, destinationID)
	queryParams := map[string]string{
		"routeID":       routeId,
		"destinationID": destinationID,
	}
	resp, err := o.restyGet(GET_ROUTES_STATISTICS(deviceId), queryParams)
	if err != nil {
		return nil, err
	}
	var obj stats.ResponseDestinationStatistics
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (o *haivisionSdk) GetDestinationStatisticsByName(deviceId string, routeId string, destinationName string) (*stats.ResponseDestinationStatistics, error) {
	log.Println("GetStats ", deviceId, routeId, destinationName)
	queryParams := map[string]string{
		"routeID":         routeId,
		"destinationName": destinationName,
	}
	resp, err := o.restyGet(GET_ROUTES_STATISTICS(deviceId), queryParams)
	if err != nil {
		return nil, err
	}
	var obj stats.ResponseDestinationStatistics
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

/*
GET /api/gateway/[Device ID]/statistics/client?routeID=[Route ID]&destinationID=
  [Destination ID]&clientAddress=[Client Address]&clientPort=[Client Port]
cookie: sessionID: [Session ID]
*/

func (o *haivisionSdk) GetSrtClientStatistics(deviceId string, routeId string, destinationID string, clientAddress string, clientPort string) (*stats.ResponseSrtClientStatistics, error) {
	log.Println("GetStats ", deviceId, routeId, destinationID)
	queryParams := map[string]string{
		"routeID":       routeId,
		"destinationID": destinationID,
		"clientAddress": clientAddress,
		"clientPort":    clientPort,
	}
	resp, err := o.restyGet(GET_ROUTES_STATISTICS(deviceId), queryParams)
	if err != nil {
		return nil, err
	}
	var obj stats.ResponseSrtClientStatistics
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

//
