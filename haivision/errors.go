package haivision

import (
	"fmt"
	"net/http"
	"strings"
)

// bodyExcerptLimit limita quanto del body finisce in un messaggio d'errore o in un log:
// abbastanza per diagnosticare, non tanto da riversare una pagina HTML nei log del consumer.
const bodyExcerptLimit = 512

// APIError è l'errore restituito quando il gateway risponde con uno status HTTP >= 400.
//
// Prima della v1.1.0 questi casi passavano inosservati: gli helper resty ritornavano
// err == nil per qualsiasi risposta ricevuta e il body d'errore veniva deserializzato in una
// struct a zero-value. Con credenziali sbagliate InitSession restituiva quindi un oggetto
// vuoto e nessun errore, e il client partiva con un cookie sessionID="".
type APIError struct {
	Method     string
	URL        string
	StatusCode int
	Status     string
	// Body è un estratto del body di risposta, troncato e con i campi sensibili mascherati.
	Body string
}

func (e *APIError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("haivision: %s %s: %s", e.Method, e.URL, e.Status)
	}
	return fmt.Sprintf("haivision: %s %s: %s: %s", e.Method, e.URL, e.Status, e.Body)
}

// IsUnauthorized segnala 401/403: sessione scaduta, credenziali errate o ruolo insufficiente.
// La sessione del gateway scade (vedi ResponseSessionInfo.ExpireAt) e l'SDK non la rinnova:
// è questo il caso da intercettare per ricostruire il client.
func (e *APIError) IsUnauthorized() bool {
	return e.StatusCode == http.StatusUnauthorized || e.StatusCode == http.StatusForbidden
}

// IsNotFound segnala 404: device, route o destinazione inesistenti.
func (e *APIError) IsNotFound() bool {
	return e.StatusCode == http.StatusNotFound
}

func newAPIError(method, url string, statusCode int, status string, body []byte) *APIError {
	return &APIError{
		Method:     method,
		URL:        url,
		StatusCode: statusCode,
		Status:     status,
		Body:       bodyExcerpt(body),
	}
}

// bodyExcerpt maschera i campi sensibili e tronca il body a bodyExcerptLimit caratteri.
// Il taglio è per rune, non per byte: un body UTF-8 non viene spezzato a metà carattere.
func bodyExcerpt(b []byte) string {
	s := strings.TrimSpace(redactJSON(string(b)))
	r := []rune(s)
	if len(r) > bodyExcerptLimit {
		return string(r[:bodyExcerptLimit]) + "… (troncato)"
	}
	return s
}
