package test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
)

// Test delle operazioni di scrittura sulle route contro lo stub: verificano il BODY che il
// gateway riceve davvero, non solo che la chiamata non dia errore.

type srtFields = route.RouteFields[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]

func sampleFields() srtFields {
	return srtFields{
		Name:       "evento-live",
		StartRoute: haivision.Bool(true),
		Source: srt.RequestSourceModelSRT{
			Name: "ingresso", Address: "0.0.0.0", Protocol: "srt", Port: 2000,
		},
		Destinations: []srt.RequestDestinationModelSrt{{
			Name: "uscita", Address: "10.0.0.9", Protocol: "srt", Port: 2001,
		}},
	}
}

// captureUpdates registra il body di ogni POST su /updates.
func captureUpdates(g *gatewayStub, ch chan map[string]any) {
	g.onFunc(http.MethodPost, "/api/devices/dev-1/updates", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		ch <- body
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
}

// CreateRoute deve produrre il body documentato: action/deviceID/elementType + wrapper fields.
// Nella v1.x veniva postato il modello di RISPOSTA, senza nessuno di quei campi.
func TestCreateRouteSendsDocumentedBody(t *testing.T) {
	g := newGatewayStub(t)
	bodies := make(chan map[string]any, 1)
	captureUpdates(g, bodies)
	c := connected(t, g)

	res, err := haivision.CreateRoute(context.Background(), c, "dev-1", sampleFields())
	if err != nil {
		t.Fatalf("CreateRoute: %v", err)
	}
	if res.Status != "ok" {
		t.Errorf("Status = %q", res.Status)
	}

	body := <-bodies
	for k, want := range map[string]any{
		"action": "create", "deviceID": "dev-1", "elementType": "route",
	} {
		if body[k] != want {
			t.Errorf("%s = %v, atteso %v", k, body[k], want)
		}
	}
	fields, ok := body["fields"].(map[string]any)
	if !ok {
		t.Fatalf("`fields` assente o non oggetto: %v", body)
	}
	if fields["name"] != "evento-live" {
		t.Errorf("fields.name = %v", fields["name"])
	}
	if fields["startRoute"] != true {
		t.Errorf("fields.startRoute = %v, atteso true", fields["startRoute"])
	}
	// nessun campo di stato della risposta deve finire nel body
	for _, k := range []string{"state", "elapsedTime", "summaryStatusCode", "pendingUpdates", "id"} {
		if _, present := fields[k]; present {
			t.Errorf("il body contiene il campo di stato %q, che appartiene alla risposta", k)
		}
	}
}

// UpdateRoute manda action=update + elementID e NON manda startRoute.
func TestUpdateRouteSendsElementIDAndDropsStartRoute(t *testing.T) {
	g := newGatewayStub(t)
	bodies := make(chan map[string]any, 1)
	captureUpdates(g, bodies)
	c := connected(t, g)

	// i fields passati hanno StartRoute impostato: UpdateRoute lo deve scartare
	if _, err := haivision.UpdateRoute(context.Background(), c, "dev-1", "route-1", sampleFields()); err != nil {
		t.Fatalf("UpdateRoute: %v", err)
	}
	body := <-bodies
	if body["action"] != "update" {
		t.Errorf("action = %v, atteso update", body["action"])
	}
	if body["elementID"] != "route-1" {
		t.Errorf("elementID = %v, atteso route-1", body["elementID"])
	}
	fields := body["fields"].(map[string]any)
	if _, present := fields["startRoute"]; present {
		t.Error("la update ha inviato startRoute: esiste solo nella create")
	}
}

func TestDeleteRoute(t *testing.T) {
	g := newGatewayStub(t)
	bodies := make(chan map[string]any, 1)
	captureUpdates(g, bodies)
	c := connected(t, g)

	if _, err := c.DeleteRoute(context.Background(), "dev-1", "route-1"); err != nil {
		t.Fatalf("DeleteRoute: %v", err)
	}
	body := <-bodies
	if body["action"] != "delete" || body["elementID"] != "route-1" {
		t.Errorf("body = %v", body)
	}
	if _, present := body["fields"]; present {
		t.Error("la delete non deve contenere fields")
	}
}

// StartRoute/StopRoute passano dall'endpoint dei comandi con il comando giusto.
func TestStartAndStopRoute(t *testing.T) {
	g := newGatewayStub(t)
	cmds := make(chan string, 2)
	g.onFunc(http.MethodPost, "/api/devices/dev-1/commands", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		cmds <- body["command"].(string)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"state":"pending","_id":"c1","result":null}]`))
	})
	c := connected(t, g)

	if _, err := c.StartRoute(context.Background(), "dev-1", "route-1"); err != nil {
		t.Fatalf("StartRoute: %v", err)
	}
	if got := <-cmds; got != "start-route" {
		t.Errorf("StartRoute ha inviato command=%q", got)
	}
	if _, err := c.StopRoute(context.Background(), "dev-1", "route-1"); err != nil {
		t.Fatalf("StopRoute: %v", err)
	}
	if got := <-cmds; got != "stop-route" {
		t.Errorf("StopRoute ha inviato command=%q", got)
	}
}

// La validazione lato client rifiuta i body incompleti prima di toccare la rete.
func TestCreateRouteValidatesBeforeSending(t *testing.T) {
	g := newGatewayStub(t)
	called := false
	g.onFunc(http.MethodPost, "/api/devices/dev-1/updates", func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	c := connected(t, g)

	cases := map[string]srtFields{
		"nome mancante":            {Source: sampleFields().Source, Destinations: sampleFields().Destinations},
		"nessuna destinazione":     {Name: "r", Source: sampleFields().Source},
		"indirizzo sorgente vuoto": {Name: "r", Source: srt.RequestSourceModelSRT{Name: "s", Protocol: "srt", Port: 1}, Destinations: sampleFields().Destinations},
	}
	for name, f := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := haivision.CreateRoute(context.Background(), c, "dev-1", f)
			if err == nil {
				t.Fatal("atteso errore di validazione")
			}
			var vErr *haivision.ValidationError
			if !errors.As(err, &vErr) {
				t.Errorf("atteso *ValidationError, ottenuto %T: %v", err, err)
			}
		})
	}
	if called {
		t.Error("una richiesta non valida è stata inviata al gateway")
	}
}

// Una route SRT senza ttl/tos/mtu è legittima: nella v1.x la validazione la rifiutava perché
// quei campi erano `string` obbligatorie.
func TestCreateRouteAcceptsRouteWithoutOptionalFields(t *testing.T) {
	g := newGatewayStub(t)
	bodies := make(chan map[string]any, 1)
	captureUpdates(g, bodies)
	c := connected(t, g)

	minimal := srtFields{
		Name:         "minimale",
		Source:       srt.RequestSourceModelSRT{Name: "in", Address: "0.0.0.0", Protocol: "srt", Port: 2000},
		Destinations: []srt.RequestDestinationModelSrt{{Name: "out", Address: "10.0.0.9", Protocol: "srt", Port: 2001}},
	}
	if _, err := haivision.CreateRoute(context.Background(), c, "dev-1", minimal); err != nil {
		t.Fatalf("una route senza campi opzionali è stata rifiutata: %v", err)
	}
	dst := (<-bodies)["fields"].(map[string]any)["destinations"].([]any)[0].(map[string]any)
	for _, k := range []string{"ttl", "tos", "mtu", "retainHeader"} {
		if _, present := dst[k]; present {
			t.Errorf("il campo opzionale %q è stato inviato comunque: %v", k, dst[k])
		}
	}
}

// StartOrStopDestination è una update con `action` sulla destinazione: senza quell'action la
// funzione rifiuta, perché una update silenziosa non avvierebbe niente.
func TestStartOrStopDestination(t *testing.T) {
	g := newGatewayStub(t)
	bodies := make(chan map[string]any, 1)
	captureUpdates(g, bodies)
	c := connected(t, g)

	f := sampleFields()
	f.Destinations[0].Action = haivision.String(route.DestinationActionStop)

	if _, err := haivision.StartOrStopDestination(context.Background(), c, "dev-1", "route-1",
		f, 0, route.DestinationActionStop); err != nil {
		t.Fatalf("StartOrStopDestination: %v", err)
	}
	body := <-bodies
	if body["action"] != "update" {
		t.Errorf("il gateway non ha un endpoint dedicato: atteso action=update, ottenuto %v", body["action"])
	}
	dst := body["fields"].(map[string]any)["destinations"].([]any)[0].(map[string]any)
	if dst["action"] != "stop" {
		t.Errorf("destinazione.action = %v, atteso stop", dst["action"])
	}

	// action non impostata sul modello → errore, non una update inutile
	f2 := sampleFields()
	if _, err := haivision.StartOrStopDestination(context.Background(), c, "dev-1", "route-1",
		f2, 0, route.DestinationActionStop); err == nil {
		t.Error("atteso errore se la destinazione non ha Action impostata")
	}
	// indice fuori range
	if _, err := haivision.StartOrStopDestination(context.Background(), c, "dev-1", "route-1",
		f, 5, route.DestinationActionStop); err == nil {
		t.Error("atteso errore con destinationIndex fuori range")
	}
	// comando non valido
	if _, err := haivision.StartOrStopDestination(context.Background(), c, "dev-1", "route-1",
		f, 0, "restart"); err == nil {
		t.Error("atteso errore con comando non valido")
	}
}

func TestLogout(t *testing.T) {
	g := newGatewayStub(t)
	hit := make(chan string, 1)
	g.onFunc(http.MethodDelete, "/api/session", func(w http.ResponseWriter, r *http.Request) {
		hit <- r.Method
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})
	c := connected(t, g)

	if err := c.Logout(context.Background()); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if got := <-hit; got != http.MethodDelete {
		t.Errorf("Logout ha usato %s, atteso DELETE", got)
	}
}

// --- context e timeout --------------------------------------------------------------------

// Un context annullato deve interrompere la chiamata: nella v1.x non era cancellabile.
func TestContextCancellationAbortsRequest(t *testing.T) {
	g := newGatewayStub(t)
	release := make(chan struct{})
	g.onFunc(http.MethodGet, "/api/gateway/dev-1/statistics", func(w http.ResponseWriter, r *http.Request) {
		<-release // non risponde finché il test non lo consente
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"collectedAt":1,"route":{}}`))
	})
	c := connected(t, g)
	t.Cleanup(func() { close(release) })

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // già annullato prima della chiamata

	if _, err := c.GetRouteStatistics(ctx, "dev-1", "route-1"); err == nil {
		t.Fatal("atteso errore con context annullato")
	} else if !errors.Is(err, context.Canceled) {
		t.Errorf("atteso context.Canceled, ottenuto %v", err)
	}
}

// Il timeout di Config chiude una chiamata che non risponde.
func TestConfigTimeout(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)
	release := make(chan struct{})
	g.onFunc(http.MethodGet, "/api/gateway/dev-1/statistics", func(w http.ResponseWriter, r *http.Request) {
		<-release
	})
	t.Cleanup(func() { close(release) })

	c, err := haivision.Dial(context.Background(), haivision.Config{
		URL: g.srv.URL, Username: "u", Password: "p", Timeout: 150 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	start := time.Now()
	if _, err := c.GetRouteStatistics(context.Background(), "dev-1", "route-1"); err == nil {
		t.Fatal("atteso errore di timeout")
	}
	if el := time.Since(start); el > 2*time.Second {
		t.Errorf("il timeout non è stato applicato: chiamata durata %s", el)
	}
}

// --- configurazione -----------------------------------------------------------------------

func TestNewValidatesConfig(t *testing.T) {
	cases := map[string]haivision.Config{
		"URL mancante":      {Username: "u", Password: "p"},
		"schema mancante":   {URL: "gateway.example.com", Username: "u", Password: "p"},
		"schema non http":   {URL: "ftp://gateway.example.com", Username: "u", Password: "p"},
		"username mancante": {URL: "https://gw", Password: "p"},
		"password mancante": {URL: "https://gw", Username: "u"},
		"timeout negativo":  {URL: "https://gw", Username: "u", Password: "p", Timeout: -time.Second},
	}
	for name, cfg := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := haivision.New(cfg)
			if err == nil {
				t.Fatal("attesa configurazione rifiutata")
			}
			if !errors.Is(err, haivision.ErrInvalidConfig) {
				t.Errorf("atteso ErrInvalidConfig, ottenuto %v", err)
			}
		})
	}
}

// New non deve aprire connessioni: costruire un client verso un host inesistente riesce,
// è Connect che fallisce. Nella v1.x il costruttore faceva due chiamate HTTP.
func TestNewDoesNotConnect(t *testing.T) {
	c, err := haivision.New(haivision.Config{
		URL: "http://127.0.0.1:1", Username: "u", Password: "p",
	})
	if err != nil {
		t.Fatalf("New non deve fare I/O di rete: %v", err)
	}
	if c.GetDeviceID() != "" {
		t.Errorf("GetDeviceID() = %q prima di Connect, atteso vuoto", c.GetDeviceID())
	}
	if err := c.Connect(context.Background()); err == nil {
		t.Error("Connect verso una porta chiusa deve fallire")
	}
}

// Connect è rieseguibile: dopo una sessione scaduta basta richiamarla.
func TestConnectIsRepeatable(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	c, err := haivision.New(haivision.Config{URL: g.srv.URL, Username: "u", Password: "p"})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for i := 0; i < 2; i++ {
		if err := c.Connect(context.Background()); err != nil {
			t.Fatalf("Connect #%d: %v", i+1, err)
		}
	}
	if c.GetDeviceID() != "dev-1" {
		t.Errorf("GetDeviceID = %q", c.GetDeviceID())
	}
}

// Un 2xx senza sessionID non deve produrre un client apparentemente valido.
func TestConnectRejectsEmptySessionID(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, `{"response":{"type":"Session"}}`)

	_, err := haivision.Dial(context.Background(), haivision.Config{
		URL: g.srv.URL, Username: "u", Password: "p",
	})
	if err == nil {
		t.Fatal("atteso errore quando il gateway risponde 2xx senza sessionID")
	}
}
