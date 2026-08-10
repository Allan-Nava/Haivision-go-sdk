package test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
)

// Copertura dei metodi di statistica e di lettura route contro lo stub: verifica path, query
// param e deserializzazione. Complementare a wire_stats_test.go, che copre solo le struct.

// connected costruisce un client già autenticato sullo stub.
func connected(t *testing.T, g *gatewayStub) *haivision.Client {
	t.Helper()
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)
	c, err := haivision.Dial(context.Background(), haivision.Config{URL: g.srv.URL, Username: "u", Password: "p"})
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	return c
}

const routeStatsBody = `{"collectedAt":1675178018888,"route":{"name":"r1","id":"route-1",
  "state":"connected","elapsedRunningTime":"00:03:46","source":{"name":"s","id":"src-1",
  "numPackets":42},"destinations":[]}}`

func TestGetRouteStatistics(t *testing.T) {
	g := newGatewayStub(t)
	var gotQuery string
	g.onFunc(http.MethodGet, "/api/gateway/dev-1/statistics", func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(routeStatsBody))
	})
	c := connected(t, g)

	got, err := c.GetRouteStatistics(context.Background(), "dev-1", "route-1")
	if err != nil {
		t.Fatalf("GetRouteStatistics: %v", err)
	}
	if gotQuery != "routeID=route-1" {
		t.Errorf("query = %q, attesa routeID=route-1", gotQuery)
	}
	if got.Route.ID != "route-1" || got.Route.State != "connected" {
		t.Errorf("route = %+v", got.Route)
	}
	if got.CollectedAt != 1675178018888 {
		t.Errorf("CollectedAt = %d", got.CollectedAt)
	}
}

func TestGetSourceStatistics(t *testing.T) {
	g := newGatewayStub(t)
	var gotQuery string
	g.onFunc(http.MethodGet, "/api/gateway/dev-1/statistics", func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("sourceID")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"collectedAt":1,"source":{"id":"src-1","numPackets":7}}`))
	})
	c := connected(t, g)

	got, err := c.GetSourceStatistics(context.Background(), "dev-1", "route-1", "src-1")
	if err != nil {
		t.Fatalf("GetSourceStatistics: %v", err)
	}
	if gotQuery != "src-1" {
		t.Errorf("query sourceID = %q", gotQuery)
	}
	if got.Source.NumPackets != 7 {
		t.Errorf("NumPackets = %v", got.Source.NumPackets)
	}
}

func TestGetDestinationStatisticsByIdAndByName(t *testing.T) {
	g := newGatewayStub(t)
	seen := make(chan string, 2)
	g.onFunc(http.MethodGet, "/api/gateway/dev-1/statistics", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if v := q.Get("destinationID"); v != "" {
			seen <- "id=" + v
		} else {
			seen <- "name=" + q.Get("destinationName")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"collectedAt":1,"destination":[]}`))
	})
	c := connected(t, g)

	if _, err := c.GetDestinationStatisticsById(context.Background(), "dev-1", "route-1", "dst-1"); err != nil {
		t.Fatalf("ById: %v", err)
	}
	if got := <-seen; got != "id=dst-1" {
		t.Errorf("ById ha inviato %q", got)
	}
	if _, err := c.GetDestinationStatisticsByName(context.Background(), "dev-1", "route-1", "uscita-primaria"); err != nil {
		t.Fatalf("ByName: %v", err)
	}
	if got := <-seen; got != "name=uscita-primaria" {
		t.Errorf("ByName ha inviato %q", got)
	}
}

// Le statistiche devono propagare un 404 come APIError, non restituire una struct vuota.
func TestStatisticsPropagatesNotFound(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodGet, "/api/gateway/dev-1/statistics", http.StatusNotFound,
		`{"error":"route not found"}`)
	c := connected(t, g)

	_, err := c.GetRouteStatistics(context.Background(), "dev-1", "inesistente")
	if err == nil {
		t.Fatal("atteso errore su 404")
	}
	apiErr, ok := err.(*haivision.APIError)
	if !ok {
		t.Fatalf("atteso *APIError, ottenuto %T", err)
	}
	if !apiErr.IsNotFound() {
		t.Errorf("IsNotFound() = false su un 404 (status %d)", apiErr.StatusCode)
	}
}

// Un body non-JSON (es. pagina HTML di un proxy) deve dare errore di deserializzazione,
// non un oggetto a zero-value silenzioso.
func TestStatisticsRejectsNonJSONBody(t *testing.T) {
	g := newGatewayStub(t)
	g.onFunc(http.MethodGet, "/api/gateway/dev-1/statistics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte("<html><body>502 Bad Gateway</body></html>"))
	})
	c := connected(t, g)

	if _, err := c.GetRouteStatistics(context.Background(), "dev-1", "route-1"); err == nil {
		t.Fatal("atteso errore di unmarshal su body HTML")
	}
}

func TestGetRoutesTyped(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodGet, "/api/gateway/dev-1/routes", http.StatusOK,
		`{"data":[{"id":"route-1","name":"r1","state":"connected",
		  "source":{"name":"in","id":"src-1","port":2000,"protocol":"srt"},
		  "destinations":[{"name":"out","id":"dst-1","port":2001,"protocol":"srt"}]}],
		  "numPages":1,"numResults":1,"numActiveOutputConnections":1}`)
	c := connected(t, g)

	// tipi del protocollo scelti dal chiamante: nella v1.x qui tornava un *resty.Response grezzo
	got, err := haivision.GetRoutes[srt.ResponseSourceSrt, srt.ResponseDestinationSrt](
		context.Background(), c, "dev-1")
	if err != nil {
		t.Fatalf("GetRoutes: %v", err)
	}
	if got.NumResults != 1 || len(got.Data) != 1 {
		t.Fatalf("routes = %+v", got)
	}
	r := got.Data[0]
	if r.ID != "route-1" || r.State != "connected" {
		t.Errorf("route = %+v", r)
	}
	if r.Source.Port != 2000 {
		t.Errorf("source.port = %d, atteso 2000", r.Source.Port)
	}
	if len(r.Destinations) != 1 || r.Destinations[0].Port != 2001 {
		t.Errorf("destinations = %+v", r.Destinations)
	}
}

func TestGetRouteConfigurationTyped(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodGet, "/api/gateway/dev-1/routes/route-1", http.StatusOK,
		`{"id":"route-1","name":"r1","state":"connected",
		  "source":{"name":"in","id":"src-1","port":2000},"destinations":[]}`)
	c := connected(t, g)

	got, err := haivision.GetRouteConfiguration[srt.ResponseSourceSrt, srt.ResponseDestinationSrt](
		context.Background(), c, "dev-1", "route-1")
	if err != nil {
		t.Fatalf("GetRouteConfiguration: %v", err)
	}
	if got.ID != "route-1" || got.Source.Name != "in" {
		t.Errorf("route = %+v", got)
	}
}

// Gli accessi grezzi restano disponibili per ispezionare un payload senza scegliere i tipi.
func TestGetRoutesRaw(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodGet, "/api/gateway/dev-1/routes", http.StatusOK, `{"data":[],"numResults":0}`)
	c := connected(t, g)

	raw, err := c.GetRoutesRaw(context.Background(), "dev-1")
	if err != nil {
		t.Fatalf("GetRoutesRaw: %v", err)
	}
	if !json.Valid(raw) {
		t.Errorf("payload non JSON valido: %s", raw)
	}
}

// deviceID vuoto: errore lato client, nessuna richiesta al gateway.
func TestGetRoutesRejectsEmptyDeviceID(t *testing.T) {
	g := newGatewayStub(t)
	c := connected(t, g)
	if _, err := haivision.GetRoutes[srt.ResponseSourceSrt, srt.ResponseDestinationSrt](
		context.Background(), c, ""); err == nil {
		t.Error("atteso errore con deviceID vuoto")
	}
}

func TestGetSessionInfoAndDeviceInfo(t *testing.T) {
	g := newGatewayStub(t)
	g.onFunc(http.MethodGet, "/api/session", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sessionID":"sess-1","displayName":"Administrator",
		  "roles":["Administrator"],"expireAt":1536938857529,"isLicensed":true}`))
	})
	c := connected(t, g)

	info, err := c.GetSessionInfo(context.Background())
	if err != nil {
		t.Fatalf("GetSessionInfo: %v", err)
	}
	if info.DisplayName != "Administrator" || !info.IsLicensed {
		t.Errorf("session info = %+v", info)
	}
	devs, err := c.GetDeviceInfo(context.Background())
	if err != nil {
		t.Fatalf("GetDeviceInfo: %v", err)
	}
	if len(devs) != 1 || devs[0].ID != "dev-1" {
		t.Errorf("devices = %+v", devs)
	}
}

// Il client espone i getter usati dai consumer per non ripetere il deviceID.
func TestClientGetters(t *testing.T) {
	g := newGatewayStub(t)
	c := connected(t, g)
	if c.GetDeviceID() != "dev-1" {
		t.Errorf("GetDeviceID = %q", c.GetDeviceID())
	}
	if c.GetHType() != "gateway" {
		t.Errorf("GetHType = %q", c.GetHType())
	}
	if c.IsDebug() {
		t.Error("IsDebug() = true su un client costruito con debug=false")
	}
}

// Con debug=true il client funziona e le hook di redazione sono registrate: il test verifica
// che il percorso di debug non rompa nulla (la redazione ha test dedicati in-package).
func TestDebugModeDoesNotBreakRequests(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	c, err := haivision.Dial(context.Background(), haivision.Config{URL: g.srv.URL, Username: "haiadmin", Password: "sup3r-s3cret", Debug: true})
	if err != nil {
		t.Fatalf("Dial con debug: %v", err)
	}
	if !c.IsDebug() {
		t.Error("IsDebug() = false con debug=true")
	}
}

// Un gateway irraggiungibile deve dare errore di trasporto, non panic.
func TestDialUnreachableGateway(t *testing.T) {
	// porta 1 su localhost: connection refused immediato
	if _, err := haivision.Dial(context.Background(), haivision.Config{URL: "http://127.0.0.1:1", Username: "u", Password: "p"}); err == nil {
		t.Fatal("atteso errore di connessione")
	}
}

// Insecure:false NON deve disabilitare la verifica TLS: contro un server HTTPS con
// certificato self-signed la connessione deve FALLIRE.
func TestInsecureFalseKeepsTLSVerification(t *testing.T) {
	g := newTLSGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	if _, err := haivision.Dial(context.Background(), haivision.Config{URL: g.srv.URL, Username: "u", Password: "p", Insecure: false}); err == nil {
		t.Fatal("Insecure:false ha accettato un certificato self-signed: " +
			"la verifica TLS è disabilitata quando non dovrebbe")
	}
	if _, err := haivision.Dial(context.Background(), haivision.Config{URL: g.srv.URL, Username: "u", Password: "p"}); err == nil {
		t.Fatal("Insecure non impostato ha accettato un certificato self-signed")
	}
	// solo Insecure:true deve accettarlo
	if _, err := haivision.Dial(context.Background(), haivision.Config{URL: g.srv.URL, Username: "u", Password: "p", Insecure: true}); err != nil {
		t.Fatalf("Insecure:true deve accettare il self-signed, invece: %v", err)
	}
}
