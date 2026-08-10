package test

import (
	"net/http"
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision"
)

// Copertura dei metodi di statistica e di lettura route contro lo stub: verifica path, query
// param e deserializzazione. Complementare a wire_stats_test.go, che copre solo le struct.

// connected costruisce un client già autenticato sullo stub.
func connected(t *testing.T, g *gatewayStub) haivision.IHaivisionClient {
	t.Helper()
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)
	c, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, nil)
	if err != nil {
		t.Fatalf("BuildHaivision: %v", err)
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

	got, err := c.GetRouteStatistics("dev-1", "route-1")
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

	got, err := c.GetSourceStatistics("dev-1", "route-1", "src-1")
	if err != nil {
		t.Fatalf("GetSourceStatistics: %v", err)
	}
	if gotQuery != "src-1" {
		t.Errorf("query sourceID = %q", gotQuery)
	}
	if got.Source.NumPackets != 7 {
		t.Errorf("NumPackets = %d", got.Source.NumPackets)
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

	if _, err := c.GetDestinationStatisticsById("dev-1", "route-1", "dst-1"); err != nil {
		t.Fatalf("ById: %v", err)
	}
	if got := <-seen; got != "id=dst-1" {
		t.Errorf("ById ha inviato %q", got)
	}
	if _, err := c.GetDestinationStatisticsByName("dev-1", "route-1", "uscita-primaria"); err != nil {
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

	_, err := c.GetRouteStatistics("dev-1", "inesistente")
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

	if _, err := c.GetRouteStatistics("dev-1", "route-1"); err == nil {
		t.Fatal("atteso errore di unmarshal su body HTML")
	}
}

func TestGetRoutesAndRouteConfiguration(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodGet, "/api/gateway/dev-1/routes", http.StatusOK,
		`{"data":[],"numPages":1,"numResults":0}`)
	g.on(http.MethodGet, "/api/gateway/dev-1/routes/route-1", http.StatusOK,
		`{"id":"route-1","name":"r1"}`)
	c := connected(t, g)

	resp, err := c.GetRoutes("dev-1")
	if err != nil {
		t.Fatalf("GetRoutes: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		t.Errorf("GetRoutes status = %d", resp.StatusCode())
	}
	cfg, err := c.GetRouteConfiguration("dev-1", "route-1")
	if err != nil {
		t.Fatalf("GetRouteConfiguration: %v", err)
	}
	if cfg.StatusCode() != http.StatusOK {
		t.Errorf("GetRouteConfiguration status = %d", cfg.StatusCode())
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

	info, err := c.GetSessionInfo()
	if err != nil {
		t.Fatalf("GetSessionInfo: %v", err)
	}
	if info.DisplayName != "Administrator" || !info.IsLicensed {
		t.Errorf("session info = %+v", info)
	}
	devs, err := c.GetDeviceInfo()
	if err != nil {
		t.Fatalf("GetDeviceInfo: %v", err)
	}
	if len(*devs) != 1 || (*devs)[0].ID != "dev-1" {
		t.Errorf("devices = %+v", *devs)
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

	c, err := haivision.BuildHaivision(g.srv.URL, true, "haiadmin", "sup3r-s3cret", nil, nil)
	if err != nil {
		t.Fatalf("BuildHaivision con debug: %v", err)
	}
	if !c.IsDebug() {
		t.Error("IsDebug() = false con debug=true")
	}
}

// Un gateway irraggiungibile deve dare errore di trasporto, non panic.
func TestBuildHaivisionUnreachableGateway(t *testing.T) {
	// porta 1 su localhost: connection refused immediato
	if _, err := haivision.BuildHaivision("http://127.0.0.1:1", false, "u", "p", nil, nil); err == nil {
		t.Fatal("atteso errore di connessione")
	}
}

// insecure=&false NON deve disabilitare la verifica TLS: contro un server HTTPS con
// certificato self-signed la connessione deve FALLIRE.
func TestInsecureFalseKeepsTLSVerification(t *testing.T) {
	g := newTLSGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	no := false
	if _, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, &no); err == nil {
		t.Fatal("insecure=&false ha accettato un certificato self-signed: " +
			"la verifica TLS è disabilitata quando non dovrebbe")
	}
	if _, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, nil); err == nil {
		t.Fatal("insecure=nil ha accettato un certificato self-signed")
	}
	// solo &true deve accettarlo
	yes := true
	if _, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, &yes); err != nil {
		t.Fatalf("insecure=&true deve accettare il self-signed, invece: %v", err)
	}
}
