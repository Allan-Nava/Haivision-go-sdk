package test

import (
	"encoding/json"
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
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
	var c route.RequestCreateRoute[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]
	c.Action = "create"
	c.DeviceID = "dev-1"
	c.ElementType = "route"
	c.Fields.Name = "route-srt"
	c.Fields.StartRoute = true
	c.Fields.Source = srt.RequestSourceModelSRT{
		Name: "src", ID: "src-1", Address: "10.0.0.1", Protocol: "srt", Port: 2000,
		NetworkInterface: "eth0",
	}
	c.Fields.Destinations = []srt.RequestDestinationModelSrt{{
		Name: "dst", ID: "dst-1", Address: "10.0.0.2", Protocol: "srt", Port: 2001,
		NetworkInterface: "eth0", RetainHeader: "false", Action: "create", Ttl: "64", Tos: "0x00",
	}}

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
}

// TestKnownBug_StartStopResponseIsArray documenta un bug APERTO, non un comportamento atteso.
//
// L'API risponde con un ARRAY top-level, mentre route.ResponseStartOrRoute è una struct che
// wrappa in un campo `Response` senza tag JSON: l'unmarshal fallisce. Item di backlog:
// `startstop-response-slice` (v2.0.0, breaking perché cambia la forma di un tipo esportato).
//
// Quando quell'item sarà chiuso questo test FALLIRÀ: va riscritto per asserire i valori
// deserializzati. È voluto — è il promemoria che il bug è stato risolto.
func TestKnownBug_StartStopResponseIsArray(t *testing.T) {
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

	var obj route.ResponseStartOrRoute
	err := json.Unmarshal([]byte(documented), &obj)
	if err == nil {
		t.Fatalf("l'unmarshal ora RIESCE: il bug `startstop-response-slice` è stato corretto — " +
			"riscrivi questo test per asserire i valori deserializzati")
	}
	t.Logf("bug aperto confermato (`startstop-response-slice`): %v", err)
}

// La risposta di una create è documentata come {"status": "..."} e questa deserializza già bene.
func TestCreateRouteResponse(t *testing.T) {
	var obj route.ResponseCreateRoute
	if err := json.Unmarshal([]byte(`{"status":"Route created"}`), &obj); err != nil {
		t.Fatalf("unmarshal della risposta di create: %v", err)
	}
	if obj.Status != "Route created" {
		t.Errorf("Status = %q, atteso %q", obj.Status, "Route created")
	}
}
