package haivision

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
)

// Operazioni sulle route.
//
// Le operazioni che dipendono dal protocollo sono **funzioni generiche**, non metodi: in Go i
// metodi non possono avere type parameter, e metterle sull'interfaccia richiederebbe quattro
// metodi per ognuna (SRT, RTMP, RTSP, UDP/RTP). Si usano così:
//
//	routes, err := haivision.GetRoutes[srt.ResponseSourceSrt, srt.ResponseDestinationSrt](ctx, c, deviceID)
//	res, err := haivision.CreateRoute(ctx, c, deviceID, fields)   // tipi dedotti da `fields`
//
// Le operazioni indipendenti dal protocollo (delete, start, stop) restano metodi e stanno
// sull'interfaccia IHaivisionClient.

// --- lettura ------------------------------------------------------------------------------

// GetRoutes elenca le route di un device, deserializzate nei modelli del protocollo scelto.
//
// Sostituisce il GetRoutes della v1.x, che restituiva un `*resty.Response` grezzo: il tipo
// della libreria HTTP era esposto nell'API pubblica e il chiamante doveva deserializzare a mano.
func GetRoutes[TS route.ResponseSource, TD route.ResponseDestination](
	ctx context.Context, c *Client, deviceID string,
) (*route.ResponseRoutes[TS, TD], error) {
	if deviceID == "" {
		return nil, fmt.Errorf("%w: GetRoutes richiede un deviceID", ErrInvalidConfig)
	}
	c.debugf("GetRoutes device=%s", deviceID)
	resp, err := c.get(ctx, GET_LIST_OF_ROUTES(deviceID), nil)
	if err != nil {
		return nil, err
	}
	var obj route.ResponseRoutes[TS, TD]
	if err := decode(resp, "GetRoutes", &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// GetRouteConfiguration legge la configurazione di una singola route.
func GetRouteConfiguration[TS route.ResponseSource, TD route.ResponseDestination](
	ctx context.Context, c *Client, deviceID, routeID string,
) (*route.ResponseRouteModel[TS, TD], error) {
	if deviceID == "" || routeID == "" {
		return nil, fmt.Errorf("%w: GetRouteConfiguration richiede deviceID e routeID", ErrInvalidConfig)
	}
	c.debugf("GetRouteConfiguration device=%s route=%s", deviceID, routeID)
	resp, err := c.get(ctx, GET_ROUTE_CONFIGURATION(deviceID, routeID), nil)
	if err != nil {
		return nil, err
	}
	var obj route.ResponseRouteModel[TS, TD]
	if err := decode(resp, "GetRouteConfiguration", &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// GetRoutesRaw restituisce il payload grezzo di GET /routes: utile per ispezionare una
// risposta senza dover scegliere i tipi del protocollo.
func (c *Client) GetRoutesRaw(ctx context.Context, deviceID string) (json.RawMessage, error) {
	resp, err := c.get(ctx, GET_LIST_OF_ROUTES(deviceID), nil)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(resp.Body()), nil
}

// GetRouteConfigurationRaw è l'equivalente grezzo di GetRouteConfiguration.
func (c *Client) GetRouteConfigurationRaw(ctx context.Context, deviceID, routeID string) (json.RawMessage, error) {
	resp, err := c.get(ctx, GET_ROUTE_CONFIGURATION(deviceID, routeID), nil)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(resp.Body()), nil
}

// --- scrittura ----------------------------------------------------------------------------

/*
CreateRoute crea una route.

Costruisce il body documentato — `action`/`deviceID`/`elementType` più il wrapper `fields` — a
partire dai soli campi della route. Nella v1.x le CreateRoute* accettavano il modello di
RISPOSTA e lo postavano così com'era: il body non conteneva nessuno di quei campi, quindi la
create non poteva funzionare.

	fields := route.RouteFields[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]{
		Name:         "evento-live",
		StartRoute:   haivision.Bool(true),
		Source:       srt.RequestSourceModelSRT{Name: "in", Address: "0.0.0.0", Protocol: "srt", Port: 2000},
		Destinations: []srt.RequestDestinationModelSrt{{Name: "out", Address: "10.0.0.9", Protocol: "srt", Port: 2001}},
	}
	res, err := haivision.CreateRoute(ctx, c, deviceID, fields)
*/
func CreateRoute[TS route.RequestSource, TD route.RequestDestination](
	ctx context.Context, c *Client, deviceID string, fields route.RouteFields[TS, TD],
) (*route.ResponseCreateRoute, error) {
	body := route.RequestCreateRoute[TS, TD]{
		Action:      route.ActionCreate,
		DeviceID:    deviceID,
		ElementType: route.ElementTypeRoute,
		Fields:      fields,
	}
	if err := validateRequest("CreateRoute", &body); err != nil {
		return nil, err
	}
	c.debugf("CreateRoute device=%s name=%s", deviceID, fields.Name)
	return c.postRouteUpdate(ctx, deviceID, "CreateRoute", &body)
}

/*
UpdateRoute sostituisce la configurazione di una route esistente.

⚠️ È una sostituzione, non una modifica parziale: le destinazioni che ometti vengono RIMOSSE
dalla route. Per cambiare una sola destinazione, rileggi la route con GetRouteConfiguration e
rimanda l'elenco completo.

È anche il modo documentato per avviare o fermare una singola destinazione: vedi
StartOrStopDestination.
*/
func UpdateRoute[TS route.RequestSource, TD route.RequestDestination](
	ctx context.Context, c *Client, deviceID, routeID string, fields route.RouteFields[TS, TD],
) (*route.ResponseCreateRoute, error) {
	// `startRoute` esiste solo nella create: nella update va omesso.
	fields.StartRoute = nil
	body := route.RequestUpdateRoute[TS, TD]{
		Action:      route.ActionUpdate,
		DeviceID:    deviceID,
		ElementType: route.ElementTypeRoute,
		ElementID:   routeID,
		Fields:      fields,
	}
	if err := validateRequest("UpdateRoute", &body); err != nil {
		return nil, err
	}
	c.debugf("UpdateRoute device=%s route=%s", deviceID, routeID)
	return c.postRouteUpdate(ctx, deviceID, "UpdateRoute", &body)
}

/*
StartOrStopDestination avvia o ferma una singola destinazione di una route.

Il gateway non ha un endpoint dedicato: si manda una **update** della route con `action`
valorizzata sulla destinazione interessata e le altre destinazioni incluse senza `action`.
Questa funzione fa esattamente quello, prendendo l'elenco completo delle destinazioni.

`command` è route.DestinationActionStart o route.DestinationActionStop.
*/
func StartOrStopDestination[TS route.RequestSource, TD route.RequestDestination](
	ctx context.Context, c *Client, deviceID, routeID string,
	fields route.RouteFields[TS, TD], destinationIndex int, command string,
) (*route.ResponseCreateRoute, error) {
	if command != route.DestinationActionStart && command != route.DestinationActionStop {
		return nil, fmt.Errorf("%w: command deve essere %q o %q, ricevuto %q",
			ErrInvalidConfig, route.DestinationActionStart, route.DestinationActionStop, command)
	}
	if destinationIndex < 0 || destinationIndex >= len(fields.Destinations) {
		return nil, fmt.Errorf("%w: destinationIndex %d fuori range (%d destinazioni)",
			ErrInvalidConfig, destinationIndex, len(fields.Destinations))
	}
	// L'`action` va impostata dal chiamante sul modello concreto: i type set non permettono di
	// scrivere un campo comune. Verifichiamo che l'abbia fatto, così l'errore arriva subito e
	// non come una update silenziosa che non avvia niente.
	raw, err := json.Marshal(fields.Destinations[destinationIndex])
	if err != nil {
		return nil, err
	}
	var probe map[string]any
	if err := json.Unmarshal(raw, &probe); err != nil {
		return nil, err
	}
	if got, _ := probe["action"].(string); got != command {
		return nil, fmt.Errorf("%w: la destinazione %d ha action=%q ma il comando richiesto è %q: "+
			"imposta Action sul modello di destinazione prima di chiamare",
			ErrInvalidConfig, destinationIndex, got, command)
	}
	c.debugf("StartOrStopDestination device=%s route=%s dest=%d command=%s",
		deviceID, routeID, destinationIndex, command)
	return UpdateRoute(ctx, c, deviceID, routeID, fields)
}

// postRouteUpdate manda il body su /updates e deserializza la risposta {"status": "..."}.
func (c *Client) postRouteUpdate(ctx context.Context, deviceID, label string, body any) (*route.ResponseCreateRoute, error) {
	resp, err := c.post(ctx, POST_ROUTE_UPDATES(deviceID), body)
	if err != nil {
		return nil, err
	}
	c.debugResponse(label, resp)
	var obj route.ResponseCreateRoute
	if err := decode(resp, label, &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// DeleteRoute rimuove una route. Non dipende dal protocollo, quindi è un metodo.
func (c *Client) DeleteRoute(ctx context.Context, deviceID, routeID string) (*route.ResponseCreateRoute, error) {
	body := &route.RequestDeleteRoute{
		Action:      route.ActionDelete,
		DeviceID:    deviceID,
		ElementType: route.ElementTypeRoute,
		ElementID:   routeID,
	}
	if err := validateRequest("DeleteRoute", body); err != nil {
		return nil, err
	}
	c.debugf("DeleteRoute device=%s route=%s", deviceID, routeID)
	return c.postRouteUpdate(ctx, deviceID, "DeleteRoute", body)
}

// StartRoute avvia una route.
func (c *Client) StartRoute(ctx context.Context, deviceID, routeID string) (route.ResponseRouteCommand, error) {
	return c.StartOrStopRoute(ctx, deviceID, routeID, route.START_ROUTE)
}

// StopRoute ferma una route.
func (c *Client) StopRoute(ctx context.Context, deviceID, routeID string) (route.ResponseRouteCommand, error) {
	return c.StartOrStopRoute(ctx, deviceID, routeID, route.STOP_ROUTE)
}

/*
StartOrStopRoute manda un comando su POST /api/devices/{deviceID}/commands.

La risposta è un ARRAY di comandi accodati, e i comandi del gateway sono **asincroni**: uno
stato `pending` significa che il comando è stato accettato, non che la route sia già partita.
Usa ResponseRouteCommand.Pending() per accorgertene.
*/
func (c *Client) StartOrStopRoute(ctx context.Context, deviceID, routeID, command string) (route.ResponseRouteCommand, error) {
	body := &route.RequestStartOrStopRoutes{
		DeviceID:   deviceID,
		Command:    command,
		Parameters: route.CommandParameters{RouteID: routeID},
	}
	// il tag `oneof=start-route stop-route` rifiuta i comandi non validi
	if err := validateRequest("StartOrStopRoute", body); err != nil {
		return nil, err
	}
	c.debugf("StartOrStopRoute device=%s route=%s command=%s", deviceID, routeID, command)
	resp, err := c.post(ctx, POST_ROUTE_COMMAND(deviceID), body)
	if err != nil {
		return nil, err
	}
	c.debugResponse("StartOrStopRoute", resp)
	var obj route.ResponseRouteCommand
	if err := decode(resp, "StartOrStopRoute", &obj); err != nil {
		return nil, err
	}
	return obj, nil
}

// --- helper per i campi opzionali ---------------------------------------------------------
//
// I campi opzionali dei modelli sono puntatori, così un valore non impostato viene omesso dal
// body invece di essere inviato come zero. Questi helper evitano di dichiarare una variabile
// temporanea per ogni campo.

// Bool ritorna un puntatore a v.
func Bool(v bool) *bool { return &v }

// Int ritorna un puntatore a v.
func Int(v int) *int { return &v }

// String ritorna un puntatore a v.
func String(v string) *string { return &v }
