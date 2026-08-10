package haivision

import (
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// ErrNoDevices è restituito quando il gateway risponde alla GET /api/devices con una lista
// vuota. Fino alla v1.0.0 questo caso era un panic da index out of range: il builder faceva
// `(*deviceResponse)[0].ID` senza controllare la lunghezza.
var ErrNoDevices = errors.New("haivision: il gateway non ha restituito nessun device")

// BuildHaivision costruisce un client e apre la sessione sul gateway.
//
// Attenzione: fa DUE chiamate HTTP (POST /api/session + GET /api/devices), quindi non è un
// costruttore puro e non è utilizzabile senza rete. DeviceID e HType sono presi dal PRIMO
// device della lista: per un setup multi-device passare il deviceId esplicito ai metodi.
//
// insecure: passare un puntatore a true per NON verificare il certificato TLS del gateway.
// nil o puntatore a false ⇒ TLS verificato. Fino alla v1.0.0 il valore puntato non veniva
// mai letto, quindi anche `&false` disabilitava la verifica — l'opposto di quanto chiesto.
func BuildHaivision(url string, debug bool, username string, password string, header *HeaderConfigurator, insecure *bool) (IHaivisionClient, error) {
	// init haivision
	haivisionClient := &haivisionSdk{
		Url:        url,
		restClient: resty.New(),
	}
	// You can override all below settings and options at request level if you want to
	//--------------------------------------------------------------------------------
	// Host URL for all request. So you can use relative URL in the request
	haivisionClient.restClient.SetBaseURL(url)
	//
	if debug {
		haivisionClient.restClient.SetDebug(true)
		haivisionClient.debug = true
		// Le hook di log vanno registrate PRIMA della prima richiesta: il debug di resty
		// stampa i body, e senza redazione il POST /api/session finirebbe nei log con
		// username e password in chiaro (e la risposta con il sessionID).
		haivisionClient.restClient.OnRequestLog(func(rl *resty.RequestLog) error {
			rl.Header = redactHeader(rl.Header)
			rl.Body = redactJSON(rl.Body)
			return nil
		})
		haivisionClient.restClient.OnResponseLog(func(rl *resty.ResponseLog) error {
			rl.Header = redactHeader(rl.Header)
			rl.Body = redactJSON(rl.Body)
			return nil
		})
		haivisionClient.debugf("debug mode abilitato (credenziali e sessionID mascherati nei log)")
	}
	if insecure != nil && *insecure {
		haivisionClient.restClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		haivisionClient.debugf("verifica del certificato TLS DISABILITATA")
	}
	// Gli header custom vanno applicati PRIMA del login: servono anche alle due chiamate di
	// bootstrap qui sotto. Fino alla v1.0.0 erano impostati alla fine, quindi un Basic auth
	// richiesto da un reverse proxy davanti al gateway non partiva e il builder falliva sempre.
	if header != nil {
		// Headers for all request
		for h, v := range header.GetHeaders() {
			haivisionClient.restClient.SetHeader(h, v)
		}
	}
	respSessionId, err := haivisionClient.InitSession(username, password)
	if err != nil {
		return nil, err
	}
	// Set Cookie for all request
	haivisionClient.restClient.SetCookie(&http.Cookie{
		Name:  "sessionID",
		Value: respSessionId.Response.SessionID,
	})
	//
	deviceResponse, err := haivisionClient.GetDeviceInfo()
	if err != nil {
		return nil, err
	}
	if deviceResponse == nil || len(*deviceResponse) == 0 {
		return nil, ErrNoDevices
	}
	haivisionClient.DeviceID = (*deviceResponse)[0].ID
	haivisionClient.HType = (*deviceResponse)[0].Type
	//
	return haivisionClient, nil
}
