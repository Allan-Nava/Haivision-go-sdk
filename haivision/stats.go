package haivision

import (
	"context"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/stats"
)

/*
GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]
cookie: sessionID: [Session ID]
*/

func (c *Client) GetRouteStatistics(ctx context.Context, deviceID, routeID string) (*stats.ResponseRouteStatistics, error) {
	c.debugf("GetRouteStatistics device=%s route=%s", deviceID, routeID)
	queryParams := map[string]string{
		"routeID": routeID,
	}
	resp, err := c.get(ctx, GET_ROUTES_STATISTICS(deviceID), queryParams)
	if err != nil {
		return nil, err
	}
	c.debugResponse("statistics", resp)
	var obj stats.ResponseRouteStatistics
	if err := decode(resp, "GetRouteStatistics", &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

/*
Requests
GET /api/gateway/[Device ID]/statistics?routeID=[Route ID]&sourceID=[Source ID]
cookie: sessionID: [Session ID]
*/

func (c *Client) GetSourceStatistics(ctx context.Context, deviceID, routeID, sourceID string) (*stats.ResponseSourceStatistics, error) {
	c.debugf("GetSourceStatistics device=%s route=%s source=%s", deviceID, routeID, sourceID)
	queryParams := map[string]string{
		"routeID":  routeID,
		"sourceID": sourceID,
	}
	resp, err := c.get(ctx, GET_ROUTES_STATISTICS(deviceID), queryParams)
	if err != nil {
		return nil, err
	}
	c.debugResponse("statistics", resp)
	var obj stats.ResponseSourceStatistics
	if err := decode(resp, "GetSourceStatistics", &obj); err != nil {
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

func (c *Client) GetDestinationStatisticsById(ctx context.Context, deviceID, routeID, destinationID string) (*stats.ResponseDestinationStatistics, error) {
	c.debugf("GetDestinationStatisticsById device=%s route=%s destination=%s", deviceID, routeID, destinationID)
	queryParams := map[string]string{
		"routeID":       routeID,
		"destinationID": destinationID,
	}
	resp, err := c.get(ctx, GET_ROUTES_STATISTICS(deviceID), queryParams)
	if err != nil {
		return nil, err
	}
	c.debugResponse("statistics", resp)
	var obj stats.ResponseDestinationStatistics
	if err := decode(resp, "GetDestinationStatistics", &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (c *Client) GetDestinationStatisticsByName(ctx context.Context, deviceID, routeID, destinationName string) (*stats.ResponseDestinationStatistics, error) {
	c.debugf("GetDestinationStatisticsByName device=%s route=%s destination=%s", deviceID, routeID, destinationName)
	queryParams := map[string]string{
		"routeID":         routeID,
		"destinationName": destinationName,
	}
	resp, err := c.get(ctx, GET_ROUTES_STATISTICS(deviceID), queryParams)
	if err != nil {
		return nil, err
	}
	c.debugResponse("statistics", resp)
	var obj stats.ResponseDestinationStatistics
	if err := decode(resp, "GetDestinationStatistics", &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

/*
GET /api/gateway/[Device ID]/statistics/client?routeID=[Route ID]&destinationID=
  [Destination ID]&clientAddress=[Client Address]&clientPort=[Client Port]
cookie: sessionID: [Session ID]
*/

func (c *Client) GetSrtClientStatistics(ctx context.Context, deviceID, routeID, destinationID, clientAddress, clientPort string) (*stats.ResponseSrtClientStatistics, error) {
	c.debugf("GetSrtClientStatistics device=%s route=%s destination=%s client=%s:%s", deviceID, routeID, destinationID, clientAddress, clientPort)
	queryParams := map[string]string{
		"routeID":       routeID,
		"destinationID": destinationID,
		"clientAddress": clientAddress,
		"clientPort":    clientPort,
	}
	// sotto-path /statistics/client, non /statistics
	resp, err := c.get(ctx, GET_ROUTES_CLIENT_STATISTICS(deviceID), queryParams)
	if err != nil {
		return nil, err
	}
	c.debugResponse("statistics", resp)
	var obj stats.ResponseSrtClientStatistics
	if err := decode(resp, "GetSrtClientStatistics", &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

//
