package haivision

import (
	"net/http"
	"regexp"
	"strings"
)

const redactedPlaceholder = "[REDACTED]"

// sensitiveJSONRe cattura i campi JSON il cui valore non deve mai comparire in un log.
// Il gruppo 1 è `"chiave":` (con spazi), il resto è il valore stringa da sostituire; il
// pattern `(?:[^"\\]|\\.)*` regge anche i valori con virgolette o backslash escapati.
var sensitiveJSONRe = regexp.MustCompile(
	`(?i)("(?:password|sessionid|srtpassphrase|passphrase|token|secret)"\s*:\s*)"(?:[^"\\]|\\.)*"`)

// sensitiveHeaders sono gli header il cui valore va mascherato nei log. `cookie` c'è perché
// il gateway autentica con il cookie sessionID: nei log equivale a una credenziale.
var sensitiveHeaders = map[string]bool{
	"authorization":       true,
	"proxy-authorization": true,
	"cookie":              true,
	"set-cookie":          true,
}

// redactJSON maschera i valori dei campi sensibili in un body JSON.
//
// Serve perché `resty.SetDebug(true)` logga i body delle richieste: senza redazione il
// POST /api/session finisce nei log con username e password in chiaro, e la risposta con il
// sessionID. Con la redazione il debug resta utilizzabile senza esporre credenziali.
func redactJSON(s string) string {
	if s == "" {
		return s
	}
	return sensitiveJSONRe.ReplaceAllString(s, `${1}"`+redactedPlaceholder+`"`)
}

// redactHeader ritorna una COPIA degli header con i valori sensibili mascherati.
// Non muta l'input: gli header passati alle hook di log di resty sono già una copia
// (`copyHeaders` in middleware.go), ma questa funzione non ci fa affidamento.
func redactHeader(h http.Header) http.Header {
	out := make(http.Header, len(h))
	for k, vs := range h {
		if sensitiveHeaders[strings.ToLower(k)] {
			out[k] = []string{redactedPlaceholder}
			continue
		}
		cp := make([]string, len(vs))
		copy(cp, vs)
		out[k] = cp
	}
	return out
}
