package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision"
)

// Test dei percorsi HTTP del client contro un httptest.Server: mai un gateway reale, così
// `go test ./...` resta offline e deterministico.

const (
	sessionOKBody = `{"response":{"type":"Session","message":"ok","sessionID":"sess-1",` +
		`"lastLoginDate":1675178018888,"numLoginFailures":0}}`
	devicesOKBody = `[{"_id":"dev-1","type":"gateway","ip":"127.0.0.1","name":"GW",` +
		`"statusCode":"ok","status":"Online","serialNumber":null,"firmware":"5.5"}]`
)

// gatewayStub è un finto gateway configurabile per rotta. Registra gli header di ogni
// richiesta ricevuta, così i test possono verificare COSA è stato inviato.
type gatewayStub struct {
	mu       sync.Mutex
	sessionH http.Header // header visti sulla POST /api/session
	handlers map[string]func(w http.ResponseWriter, r *http.Request)
	srv      *httptest.Server
}

func newGatewayStub(t *testing.T) *gatewayStub {
	t.Helper()
	g := &gatewayStub{handlers: map[string]func(http.ResponseWriter, *http.Request){}}
	g.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		g.mu.Lock()
		if r.URL.Path == "/api/session" && r.Method == http.MethodPost {
			g.sessionH = r.Header.Clone()
		}
		h := g.handlers[r.Method+" "+r.URL.Path]
		g.mu.Unlock()
		if h == nil {
			http.Error(w, `{"error":"rotta non stubbata: `+r.Method+" "+r.URL.Path+`"}`, http.StatusNotFound)
			return
		}
		h(w, r)
	}))
	t.Cleanup(g.srv.Close)
	return g
}

func (g *gatewayStub) on(method, path string, status int, body string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.handlers[method+" "+path] = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}
}

// onFunc consente handler dinamici (es. risposte che cambiano durante il test).
func (g *gatewayStub) onFunc(method, path string, h func(w http.ResponseWriter, r *http.Request)) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.handlers[method+" "+path] = h
}

func (g *gatewayStub) sessionHeader(key string) string {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.sessionH == nil {
		return ""
	}
	return g.sessionH.Get(key)
}

// happy path: sessione + device, DeviceID preso dal primo device
func TestBuildHaivisionSuccess(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	c, err := haivision.BuildHaivision(g.srv.URL, false, "haiadmin", "pw", nil, nil)
	if err != nil {
		t.Fatalf("BuildHaivision: %v", err)
	}
	if got := c.GetDeviceID(); got != "dev-1" {
		t.Errorf("GetDeviceID() = %q, atteso dev-1", got)
	}
	if got := c.GetHType(); got != "gateway" {
		t.Errorf("GetHType() = %q, atteso gateway", got)
	}
}

// Credenziali sbagliate: il 401 deve diventare un errore, non un client con sessionID vuoto.
// Prima della v1.1.0 questo caso ritornava err == nil.
func TestBuildHaivisionWrongCredentialsReturnsError(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusUnauthorized,
		`{"error":"Invalid username or password"}`)

	c, err := haivision.BuildHaivision(g.srv.URL, false, "haiadmin", "sbagliata", nil, nil)
	if err == nil {
		t.Fatalf("atteso errore su 401, ottenuto client valido: %#v", c)
	}
	var apiErr *haivision.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("atteso *haivision.APIError, ottenuto %T: %v", err, err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, atteso 401", apiErr.StatusCode)
	}
	if !apiErr.IsUnauthorized() {
		t.Error("IsUnauthorized() = false su un 401")
	}
	// il body d'errore del gateway deve arrivare al chiamante: è la sola diagnostica che ha
	if apiErr.Body == "" {
		t.Error("APIError.Body vuoto: il body d'errore del gateway va riportato")
	}
}

// Il 500 su /api/devices non deve essere silenzioso.
func TestBuildHaivisionDeviceErrorReturnsError(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusInternalServerError, `{"error":"boom"}`)

	if _, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, nil); err == nil {
		t.Fatal("atteso errore su 500 da /api/devices")
	}
}

// Lista device vuota: errore esplicito, non panic da index out of range.
func TestBuildHaivisionEmptyDeviceListReturnsError(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, `[]`)

	_, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, nil)
	if err == nil {
		t.Fatal("attesa ErrNoDevices su lista device vuota")
	}
	if !errors.Is(err, haivision.ErrNoDevices) {
		t.Errorf("atteso ErrNoDevices, ottenuto %v", err)
	}
}

// Gli header custom devono partire GIÀ sulla POST /api/session: è il caso del gateway dietro
// un reverse proxy con Basic auth. Prima della v1.1.0 venivano applicati dopo il login.
func TestBuildHaivisionSendsCustomHeadersOnLogin(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	h := haivision.InitHeaderConfigurator()
	h.CreateBasicAuthHeader("proxy-user", "proxy-pass")
	h.SetHeader("X-Tenant", "hiway")

	if _, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", h, nil); err != nil {
		t.Fatalf("BuildHaivision: %v", err)
	}
	if got := g.sessionHeader("Authorization"); got == "" {
		t.Error("la POST /api/session è partita senza header Authorization: " +
			"gli header custom vanno applicati prima del login")
	}
	if got := g.sessionHeader("X-Tenant"); got != "hiway" {
		t.Errorf("X-Tenant sulla POST /api/session = %q, atteso hiway", got)
	}
}

// HealthCheck deve poter fallire: interroga /api/session e propaga l'errore.
func TestHealthCheck(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	var sessionDown bool
	var mu sync.Mutex
	g.onFunc(http.MethodGet, "/api/session", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		down := sessionDown
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		if down {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"session expired"}`))
			return
		}
		_, _ = w.Write([]byte(`{"sessionID":"sess-1","displayName":"Administrator",` +
			`"roles":["Administrator"],"expireAt":1536938857529,"isLicensed":true}`))
	})

	c, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, nil)
	if err != nil {
		t.Fatalf("BuildHaivision: %v", err)
	}
	if err := c.HealthCheck(); err != nil {
		t.Errorf("HealthCheck() con sessione valida = %v, atteso nil", err)
	}

	mu.Lock()
	sessionDown = true
	mu.Unlock()
	err = c.HealthCheck()
	if err == nil {
		t.Fatal("HealthCheck() con sessione scaduta = nil, atteso errore")
	}
	var apiErr *haivision.APIError
	if !errors.As(err, &apiErr) || !apiErr.IsUnauthorized() {
		t.Errorf("atteso APIError 401, ottenuto %T: %v", err, err)
	}
}

// Le statistiche del client SRT stanno su /api/gateway/{id}/statistics/client.
func TestGetSrtClientStatisticsUsesClientSubPath(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	called := make(chan map[string]string, 1)
	g.onFunc(http.MethodGet, "/api/gateway/dev-1/statistics/client",
		func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			q := make(map[string]string, len(query))
			for k := range query {
				q[k] = query.Get(k)
			}
			called <- q
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"collectedAt":1675178018888,"clientStat":[]}`))
		})

	c, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, nil)
	if err != nil {
		t.Fatalf("BuildHaivision: %v", err)
	}
	if _, err := c.GetSrtClientStatistics("dev-1", "route-1", "dst-1", "10.0.0.9", "2000"); err != nil {
		t.Fatalf("GetSrtClientStatistics: %v (il path /statistics/client non è stato colpito?)", err)
	}
	q := <-called
	for k, want := range map[string]string{
		"routeID": "route-1", "destinationID": "dst-1",
		"clientAddress": "10.0.0.9", "clientPort": "2000",
	} {
		if q[k] != want {
			t.Errorf("query %s = %q, atteso %q", k, q[k], want)
		}
	}
}

// Start/stop route va sull'endpoint dei comandi, non su /updates.
func TestStartOrStopRouteUsesCommandsEndpoint(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	hit := make(chan string, 2)
	g.onFunc(http.MethodPost, "/api/devices/dev-1/commands", func(w http.ResponseWriter, r *http.Request) {
		hit <- "commands"
		w.Header().Set("Content-Type", "application/json")
		// la risposta reale è un array: l'unmarshal fallisce (bug `startstop-response-slice`),
		// ma questo test verifica solo QUALE endpoint viene colpito
		_, _ = w.Write([]byte(`{}`))
	})
	g.onFunc(http.MethodPost, "/api/devices/dev-1/updates", func(w http.ResponseWriter, r *http.Request) {
		hit <- "updates"
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	})

	c, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, nil)
	if err != nil {
		t.Fatalf("BuildHaivision: %v", err)
	}
	if _, err := c.StartOrStopRoute("dev-1", "route-1", "start-route"); err != nil {
		t.Fatalf("StartOrStopRoute: %v", err)
	}
	if got := <-hit; got != "commands" {
		t.Errorf("start-route ha colpito /%s, atteso /commands", got)
	}
}

// Comando non valido: errore lato client, nessuna richiesta al gateway.
func TestStartOrStopRouteRejectsUnknownCommand(t *testing.T) {
	g := newGatewayStub(t)
	g.on(http.MethodPost, "/api/session", http.StatusOK, sessionOKBody)
	g.on(http.MethodGet, "/api/devices", http.StatusOK, devicesOKBody)

	c, err := haivision.BuildHaivision(g.srv.URL, false, "u", "p", nil, nil)
	if err != nil {
		t.Fatalf("BuildHaivision: %v", err)
	}
	if _, err := c.StartOrStopRoute("dev-1", "route-1", "restart-route"); err == nil {
		t.Error("atteso errore su comando non valido")
	}
}
