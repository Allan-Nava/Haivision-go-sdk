package haivision

import (
	"encoding/json"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/stats"
)

/*
GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]
cookie: sessionID: [Session ID]
*/

func (o *haivisionSdk) GetRouteStatistics(deviceId string, routeId string) (*stats.ResponseRouteStatistics, error) {
	o.debugf("GetRouteStatistics device=%s route=%s", deviceId, routeId)
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
	o.debugf("GetSourceStatistics device=%s route=%s source=%s", deviceId, routeId, sourceId)
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
	o.debugf("GetDestinationStatisticsById device=%s route=%s destination=%s", deviceId, routeId, destinationID)
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
	o.debugf("GetDestinationStatisticsByName device=%s route=%s destination=%s", deviceId, routeId, destinationName)
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
	o.debugf("GetSrtClientStatistics device=%s route=%s destination=%s client=%s:%s", deviceId, routeId, destinationID, clientAddress, clientPort)
	queryParams := map[string]string{
		"routeID":       routeId,
		"destinationID": destinationID,
		"clientAddress": clientAddress,
		"clientPort":    clientPort,
	}
	// sotto-path /statistics/client, non /statistics
	resp, err := o.restyGet(GET_ROUTES_CLIENT_STATISTICS(deviceId), queryParams)
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
