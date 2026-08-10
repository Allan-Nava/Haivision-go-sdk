package haivision

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision/device"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/session"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/stats"
)

// ErrNoDevices è restituito quando GET /api/devices risponde con una lista vuota.
var ErrNoDevices = errors.New("haivision: il gateway non ha restituito nessun device")

// ErrNotConnected è restituito quando si usa un client su cui Connect non è mai riuscita.
var ErrNotConnected = errors.New("haivision: client non connesso, chiama Connect prima")

// Client è il client REST del gateway Haivision.
//
// Costruirlo con New non apre alcuna connessione: la sessione si apre con Connect. Fino alla
// v1.x il costruttore faceva due chiamate HTTP, quindi non era testabile senza rete e non
// distingueva "configurazione sbagliata" da "gateway giù".
//
// Un Client è sicuro da usare da più goroutine: dopo Connect i suoi campi non vengono più
// mutati e il client resty sottostante è concorrente.
type Client struct {
	cfg        Config
	restClient *resty.Client

	deviceID  string
	hType     string
	connected bool
}

// IHaivisionClient è l'insieme delle operazioni non generiche del client, per chi ha bisogno
// di un'interfaccia da mockare nei propri test.
//
// Le operazioni sulle route tipizzate per protocollo NON stanno qui: sono funzioni generiche
// (GetRoutes, CreateRoute, UpdateRoute…), perché in Go i metodi non possono avere type
// parameter. Metterle nell'interfaccia richiederebbe quattro metodi per protocollo.
type IHaivisionClient interface {
	// connessione e stato
	Connect(ctx context.Context) error
	HealthCheck(ctx context.Context) error
	Logout(ctx context.Context) error
	IsDebug() bool
	GetDeviceID() string
	GetHType() string

	// sessione e device
	InitSession(ctx context.Context, username, password string) (*session.BaseResponseInitSession, error)
	GetSessionInfo(ctx context.Context) (*session.ResponseSessionInfo, error)
	GetDeviceInfo(ctx context.Context) ([]device.ResponseDeviceInfo, error)

	// route non tipizzate: utili per ispezionare un payload senza scegliere i tipi
	GetRoutesRaw(ctx context.Context, deviceID string) (json.RawMessage, error)
	GetRouteConfigurationRaw(ctx context.Context, deviceID, routeID string) (json.RawMessage, error)

	// route: operazioni indipendenti dal protocollo
	DeleteRoute(ctx context.Context, deviceID, routeID string) (*route.ResponseCreateRoute, error)
	StartRoute(ctx context.Context, deviceID, routeID string) (route.ResponseRouteCommand, error)
	StopRoute(ctx context.Context, deviceID, routeID string) (route.ResponseRouteCommand, error)
	StartOrStopRoute(ctx context.Context, deviceID, routeID, command string) (route.ResponseRouteCommand, error)

	// statistiche
	GetRouteStatistics(ctx context.Context, deviceID, routeID string) (*stats.ResponseRouteStatistics, error)
	GetSourceStatistics(ctx context.Context, deviceID, routeID, sourceID string) (*stats.ResponseSourceStatistics, error)
	GetDestinationStatisticsById(ctx context.Context, deviceID, routeID, destinationID string) (*stats.ResponseDestinationStatistics, error)
	GetDestinationStatisticsByName(ctx context.Context, deviceID, routeID, destinationName string) (*stats.ResponseDestinationStatistics, error)
	GetSrtClientStatistics(ctx context.Context, deviceID, routeID, destinationID, clientAddress, clientPort string) (*stats.ResponseSrtClientStatistics, error)
}

// verifica a compile time che *Client soddisfi l'interfaccia
var _ IHaivisionClient = (*Client)(nil)

// New costruisce il client senza aprire connessioni: valida la configurazione e prepara il
// client HTTP. Chiama Connect per aprire la sessione.
func New(cfg Config) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	c := &Client{cfg: cfg, restClient: resty.New()}
	c.restClient.SetBaseURL(cfg.URL)
	c.restClient.SetTimeout(cfg.timeout())

	if cfg.Debug {
		c.restClient.SetDebug(true)
		// Le hook di log vanno registrate PRIMA della prima richiesta: il debug di resty
		// stampa i body, e senza redazione il POST /api/session finirebbe nei log con
		// username e password in chiaro (e la risposta con il sessionID).
		c.restClient.OnRequestLog(func(rl *resty.RequestLog) error {
			rl.Header = redactHeader(rl.Header)
			rl.Body = redactJSON(rl.Body)
			return nil
		})
		c.restClient.OnResponseLog(func(rl *resty.ResponseLog) error {
			rl.Header = redactHeader(rl.Header)
			rl.Body = redactJSON(rl.Body)
			return nil
		})
		c.debugf("debug attivo (credenziali e sessionID mascherati nei log)")
	}
	if cfg.Insecure {
		c.restClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		c.debugf("verifica del certificato TLS DISABILITATA")
	}
	// Gli header custom valgono anche per le due chiamate di Connect: un Basic auth richiesto
	// da un reverse proxy davanti al gateway serve già al login.
	for k, v := range cfg.Headers {
		c.restClient.SetHeader(k, v)
	}
	return c, nil
}

// Dial è New + Connect, per chi non ha bisogno di separare i due passi.
func Dial(ctx context.Context, cfg Config) (*Client, error) {
	c, err := New(cfg)
	if err != nil {
		return nil, err
	}
	if err := c.Connect(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

// Connect apre la sessione sul gateway e rileva il device.
//
// Fa due chiamate: POST /api/session (da cui il cookie sessionID usato da tutte le successive)
// e GET /api/devices. DeviceID e HType vengono dal PRIMO device: per un setup multi-device
// passare il deviceID esplicito ai metodi.
//
// È rieseguibile: dopo una sessione scaduta (APIError.IsUnauthorized) basta richiamarla.
func (c *Client) Connect(ctx context.Context) error {
	sess, err := c.InitSession(ctx, c.cfg.Username, c.cfg.Password)
	if err != nil {
		return err
	}
	if sess.Response.SessionID == "" {
		return errors.New("haivision: il gateway ha risposto 2xx ma senza sessionID")
	}
	c.restClient.SetCookie(&http.Cookie{Name: "sessionID", Value: sess.Response.SessionID})

	devices, err := c.GetDeviceInfo(ctx)
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return ErrNoDevices
	}
	c.deviceID = devices[0].ID
	c.hType = devices[0].Type
	c.connected = true
	return nil
}

// HealthCheck verifica che il gateway sia raggiungibile E che la sessione sia ancora valida.
// Un APIError con IsUnauthorized() == true significa sessione scaduta: richiama Connect.
func (c *Client) HealthCheck(ctx context.Context) error {
	resp, err := c.get(ctx, SESSION, nil)
	if err != nil {
		return err
	}
	c.debugResponse("HealthCheck", resp)
	return nil
}

// Logout chiude la sessione sul gateway (DELETE /api/session). Dopo la chiamata il client va
// riconnesso con Connect.
func (c *Client) Logout(ctx context.Context) error {
	resp, err := c.delete(ctx, SESSION)
	if err != nil {
		return err
	}
	c.debugResponse("Logout", resp)
	c.connected = false
	return nil
}

func (c *Client) IsDebug() bool       { return c.cfg.Debug }
func (c *Client) GetDeviceID() string { return c.deviceID }
func (c *Client) GetHType() string    { return c.hType }

// --- logging -----------------------------------------------------------------------------

func (c *Client) debugf(format string, v ...any) {
	if !c.cfg.Debug {
		return
	}
	if c.cfg.Logger != nil {
		c.cfg.Logger.Printf("[haivision] "+format, v...)
		return
	}
	log.Printf("[haivision] "+format, v...)
}

// debugResponse logga status, durata ed estratto del body, con i campi sensibili mascherati.
func (c *Client) debugResponse(label string, resp *resty.Response) {
	if !c.cfg.Debug || resp == nil {
		return
	}
	c.debugf("%s: %s (%s) %s", label, resp.Status(), resp.Time(), bodyExcerpt(resp.Body()))
}

// --- helper HTTP -------------------------------------------------------------------------
//
// Sono l'UNICO punto in cui si controlla lo status HTTP: nessun metodo dell'SDK deserializza
// mai un body d'errore. Un 4xx/5xx diventa un *APIError.

func (c *Client) post(ctx context.Context, url string, body any) (*resty.Response, error) {
	resp, err := c.restClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(url)
	return c.check(http.MethodPost, url, resp, err)
}

func (c *Client) get(ctx context.Context, url string, queryParams map[string]string) (*resty.Response, error) {
	resp, err := c.restClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		SetQueryParams(queryParams).
		Get(url)
	return c.check(http.MethodGet, url, resp, err)
}

func (c *Client) delete(ctx context.Context, url string) (*resty.Response, error) {
	resp, err := c.restClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json").
		Delete(url)
	return c.check(http.MethodDelete, url, resp, err)
}

func (c *Client) check(method, url string, resp *resty.Response, err error) (*resty.Response, error) {
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newAPIError(method, url, resp.StatusCode(), resp.Status(), resp.Body())
	}
	return resp, nil
}

// decode deserializza il body in v, arricchendo l'errore col contesto della chiamata: un
// "cannot unmarshal" nudo non dice quale endpoint ha risposto male.
func decode(resp *resty.Response, label string, v any) error {
	if err := json.Unmarshal(resp.Body(), v); err != nil {
		return &DecodeError{Label: label, Body: bodyExcerpt(resp.Body()), Err: err}
	}
	return nil
}
