package test

import (
	"encoding/json"
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
	udprtp "github.com/Allan-Nava/Haivision-go-sdk/haivision/udp_rtp"
)

// Test di contratto sul WIRE FORMAT: verificano che i body prodotti e le risposte accettate
// corrispondano ai payload LETTERALI della documentazione Haivision
// (https://doc.haivision.com/HMG3.7.6/rest-api-integrator-s-reference/rest-api-reference).
// Sono offline e non toccano nessun gateway.
//
// Le asserzioni sono sui nomi dei campi JSON, non su quelli Go: è l'unico modo di accorgersi
// che una struct anonima annidata senza tag serializza come "Fields" invece di "fields".

// marshalToMap serializza v e la rilegge come mappa generica, così le asserzioni sono su ciò
// che il gateway vede davvero.
func marshalToMap(t *testing.T, v interface{}) map[string]interface{} {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("rilettura del body serializzato: %v (body: %s)", err, b)
	}
	return m
}

func nestedMap(t *testing.T, m map[string]interface{}, key string) map[string]interface{} {
	t.Helper()
	v, ok := m[key]
	if !ok {
		t.Fatalf("campo %q assente nel body: %v", key, m)
	}
	nested, ok := v.(map[string]interface{})
	if !ok {
		t.Fatalf("campo %q non è un oggetto JSON: %T", key, v)
	}
	return nested
}

// La doc: POST /api/devices/[Device ID]/commands
//
//	{"deviceID": "...", "command": "...", "parameters": {"routeID": "..."}}
func TestStartStopRouteRequestBody(t *testing.T) {
	r := route.RequestStartOrStopRoutes{DeviceID: "dev-1", Command: route.START_ROUTE}
	r.Parameters.RouteID = "route-1"

	m := marshalToMap(t, r)
	if _, wrong := m["Parameters"]; wrong {
		t.Errorf("il body usa \"Parameters\" (nome del campo Go) invece di \"parameters\": %v", m)
	}
	if got := m["deviceID"]; got != "dev-1" {
		t.Errorf("deviceID = %v, atteso dev-1", got)
	}
	if got := m["command"]; got != "start-route" {
		t.Errorf("command = %v, atteso start-route", got)
	}
	if got := nestedMap(t, m, "parameters")["routeID"]; got != "route-1" {
		t.Errorf("parameters.routeID = %v, atteso route-1", got)
	}
}

// La doc: POST /api/devices/[Device ID]/updates
//
//	{"action":"create","deviceID":"...","elementType":"route","fields":{"name":...,"startRoute":...,"source":{...},"destinations":[...]}}
func TestCreateRouteRequestBody(t *testing.T) {
	c := route.RequestCreateRoute[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]{
		Action:      route.ActionCreate,
		DeviceID:    "dev-1",
		ElementType: route.ElementTypeRoute,
		Fields: route.RouteFields[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]{
			Name:       "route-srt",
			StartRoute: haivision.Bool(true),
			Source: srt.RequestSourceModelSRT{
				Name: "src", Address: "10.0.0.1", Protocol: "srt", Port: 2000,
			},
			Destinations: []srt.RequestDestinationModelSrt{{
				Name: "dst", Address: "10.0.0.2", Protocol: "srt", Port: 2001,
			}},
		},
	}

	m := marshalToMap(t, c)
	if _, wrong := m["Fields"]; wrong {
		t.Errorf("il body usa \"Fields\" (nome del campo Go) invece di \"fields\": %v", m)
	}
	for k, want := range map[string]interface{}{
		"action": "create", "deviceID": "dev-1", "elementType": "route",
	} {
		if got := m[k]; got != want {
			t.Errorf("%s = %v, atteso %v", k, got, want)
		}
	}
	fields := nestedMap(t, m, "fields")
	if got := fields["name"]; got != "route-srt" {
		t.Errorf("fields.name = %v, atteso route-srt", got)
	}
	if got := fields["startRoute"]; got != true {
		t.Errorf("fields.startRoute = %v, atteso true", got)
	}
	if got := nestedMap(t, fields, "source")["port"]; got != float64(2000) {
		t.Errorf("fields.source.port = %v, atteso 2000", got)
	}
	dsts, ok := fields["destinations"].([]interface{})
	if !ok || len(dsts) != 1 {
		t.Fatalf("fields.destinations non è un array di 1 elemento: %v", fields["destinations"])
	}
	// i campi opzionali non impostati devono essere OMESSI, non inviati a zero: il gateway
	// interpreterebbe ttl=0 o retainHeader=false come valori voluti
	dst := dsts[0].(map[string]interface{})
	for _, k := range []string{"ttl", "tos", "mtu", "retainHeader", "action", "srtLatency"} {
		if _, present := dst[k]; present {
			t.Errorf("il campo opzionale %q non impostato è stato inviato comunque: %v", k, dst[k])
		}
	}
}

// La update ha elementID e NON ha startRoute.
func TestUpdateRouteRequestBody(t *testing.T) {
	u := route.RequestUpdateRoute[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]{
		Action:      route.ActionUpdate,
		DeviceID:    "dev-1",
		ElementType: route.ElementTypeRoute,
		ElementID:   "route-1",
		Fields: route.RouteFields[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]{
			Name:         "route-srt",
			Source:       srt.RequestSourceModelSRT{Name: "src", Address: "10.0.0.1", Protocol: "srt", Port: 2000},
			Destinations: []srt.RequestDestinationModelSrt{{Name: "dst", Address: "10.0.0.2", Protocol: "srt", Port: 2001}},
		},
	}
	m := marshalToMap(t, u)
	if got := m["action"]; got != "update" {
		t.Errorf("action = %v", got)
	}
	if got := m["elementID"]; got != "route-1" {
		t.Errorf("elementID = %v, atteso route-1", got)
	}
	if _, present := nestedMap(t, m, "fields")["startRoute"]; present {
		t.Error("la update non deve contenere startRoute: esiste solo nella create")
	}
}

// La delete non ha fields.
func TestDeleteRouteRequestBody(t *testing.T) {
	d := route.RequestDeleteRoute{
		Action:      route.ActionDelete,
		DeviceID:    "dev-1",
		ElementType: route.ElementTypeRoute,
		ElementID:   "route-1",
	}
	m := marshalToMap(t, d)
	if got := m["action"]; got != "delete" {
		t.Errorf("action = %v", got)
	}
	if _, present := m["fields"]; present {
		t.Error("la delete non deve contenere fields")
	}
}

// I tipi dei campi devono corrispondere agli esempi LETTERALI della doc: ttl/tos/mtu numerici,
// shaping booleano, maxBitrate numerico. Fino alla v1.x erano stringhe.
func TestDestinationFieldTypesMatchDoc(t *testing.T) {
	d := udprtp.RequestDestinationModelUdpRtp{
		Name: "Destination1Name", Address: "10.0.65.10", Protocol: "udp", Port: 1111,
		Action: haivision.String("stop"),
		Ttl:    haivision.Int(64), Mtu: haivision.Int(1496), Tos: haivision.Int(136),
		Fec: haivision.String("none"), Encryption: haivision.String("none"),
		Shaping: haivision.Bool(false), MaxBitrate: haivision.Int(10000),
	}
	m := marshalToMap(t, d)
	for _, k := range []string{"ttl", "mtu", "tos", "maxBitrate", "port"} {
		if _, ok := m[k].(float64); !ok {
			t.Errorf("%s = %T (%v), atteso un numero JSON", k, m[k], m[k])
		}
	}
	if _, ok := m["shaping"].(bool); !ok {
		t.Errorf("shaping = %T (%v), atteso un booleano JSON", m["shaping"], m["shaping"])
	}
	if got := m["action"]; got != "stop" {
		t.Errorf("action = %v", got)
	}
}

// La risposta di start/stop è un ARRAY top-level: nella v1.x il tipo era una struct e
// l'unmarshal della risposta reale falliva sempre.
func TestStartStopResponseIsArray(t *testing.T) {
	documented := `[{
	  "action": "command",
	  "command": "start-route",
	  "parameters": { "routeID": "route-1" },
	  "deviceID": "dev-1",
	  "createdAt": 1675178018888,
	  "completedAt": 0,
	  "result": null,
	  "state": "pending",
	  "_id": "a5x4-7KEApdS0UuAUUCSog"
	}]`

	var obj route.ResponseRouteCommand
	if err := json.Unmarshal([]byte(documented), &obj); err != nil {
		t.Fatalf("unmarshal della risposta documentata: %v", err)
	}
	if len(obj) != 1 {
		t.Fatalf("attesi 1 comando, ottenuti %d", len(obj))
	}
	cmd := obj[0]
	if cmd.Command != "start-route" || cmd.Parameters.RouteID != "route-1" {
		t.Errorf("comando = %+v", cmd)
	}
	if cmd.State != "pending" || !obj.Pending() {
		t.Errorf("Pending() = %v con state=%q", obj.Pending(), cmd.State)
	}
	if cmd.ID != "a5x4-7KEApdS0UuAUUCSog" {
		t.Errorf("_id = %q", cmd.ID)
	}
}

// `result` può essere null o un oggetto: json.RawMessage regge entrambi.
func TestRouteCommandResultAcceptsObject(t *testing.T) {
	var obj route.ResponseRouteCommand
	body := `[{"command":"stop-route","state":"completed","result":{"ok":true}}]`
	if err := json.Unmarshal([]byte(body), &obj); err != nil {
		t.Fatalf("unmarshal con result oggetto: %v", err)
	}
	if obj.Pending() {
		t.Error("Pending() = true con state=completed")
	}
	if string(obj[0].Result) != `{"ok":true}` {
		t.Errorf("Result = %s", obj[0].Result)
	}
}

// La risposta di una create è documentata come {"status": "..."}.
func TestCreateRouteResponse(t *testing.T) {
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal([]byte(`{"status":"Route created"}`), &obj); err != nil {
		t.Fatalf("unmarshal della risposta di create: %v", err)
	}
	if obj.Status != "Route created" {
		t.Errorf("Status = %q, atteso %q", obj.Status, "Route created")
	}
}
