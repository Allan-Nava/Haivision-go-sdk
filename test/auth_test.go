package test

import (
	"testing"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision"
)

// La copertura di sessione/device è in wire_session_test.go (deserializzazione dei payload
// documentati) e in client_http_test.go (percorsi HTTP contro httptest, inclusi 401 e lista
// device vuota). Qui resta solo la verifica che HeaderConfigurator produca l'header Basic
// atteso: è l'unico pezzo di auth che non passa dalla rete.

func TestCreateBasicAuthHeader(t *testing.T) {
	h := haivision.InitHeaderConfigurator()
	if h.HasHeaders() {
		t.Error("un configurator appena creato non deve avere header")
	}
	h.CreateBasicAuthHeader("haiadmin", "manager")
	// base64("haiadmin:manager")
	const want = "Basic aGFpYWRtaW46bWFuYWdlcg=="
	got := h.GetHeader("Authorization")
	if got == nil {
		t.Fatal("header Authorization assente")
	}
	if *got != want {
		t.Errorf("Authorization = %q, atteso %q", *got, want)
	}
	if !h.HasHeader("Authorization") {
		t.Error("HasHeader(\"Authorization\") = false")
	}
}

func TestCreateBasicAuthHeaderEncoded(t *testing.T) {
	h := haivision.InitHeaderConfigurator()
	h.CreateBasicAuthHeaderEncoded("aGFpYWRtaW46bWFuYWdlcg==")
	got := h.GetHeader("Authorization")
	if got == nil || *got != "Basic aGFpYWRtaW46bWFuYWdlcg==" {
		t.Errorf("Authorization = %v", got)
	}
}

func TestHeaderConfiguratorDelete(t *testing.T) {
	h := haivision.InitHeaderConfigurator()
	h.SetHeader("X-Tenant", "hiway")
	h.DeleteHeader("X-Tenant")
	if h.HasHeader("X-Tenant") {
		t.Error("X-Tenant presente dopo DeleteHeader")
	}
	if got := h.GetHeader("X-Tenant"); got != nil {
		t.Errorf("GetHeader su chiave assente = %v, atteso nil", *got)
	}
}
