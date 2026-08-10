package haivision

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// DefaultTimeout è il timeout applicato a ogni richiesta se Config.Timeout è zero.
// Fino alla v1.x non c'era alcun timeout: una chiamata verso un gateway irraggiungibile poteva
// restare appesa a tempo indefinito.
const DefaultTimeout = 30 * time.Second

// Config raccoglie i parametri di costruzione del client.
//
// Sostituisce i sei parametri posizionali di BuildHaivision (due bool-ish e due string
// adiacenti, facili da invertire silenziosamente). In particolare Insecure è un `bool`: nella
// v1.x era `*bool` e il valore puntato non veniva letto, quindi `&false` **disabilitava** la
// verifica TLS invece di attivarla. Con un bool quel bug non è più esprimibile.
type Config struct {
	// URL base del gateway, con schema. Es. "https://gateway.example.com".
	URL string
	// Username e Password dell'utente del gateway (ruolo Administrator o Operator per le
	// operazioni di scrittura sulle route).
	Username string
	Password string

	// Insecure disabilita la verifica del certificato TLS del gateway. Da usare solo con
	// certificati self-signed in ambienti controllati.
	Insecure bool

	// Debug attiva il log delle richieste. Le credenziali, il sessionID e le passphrase SRT
	// vengono mascherate: vedi redact.go.
	Debug bool

	// Timeout per singola richiesta. Zero ⇒ DefaultTimeout.
	Timeout time.Duration

	// Headers applicati a ogni richiesta, incluse quelle di Connect. Serve per i gateway
	// dietro reverse proxy che richiedono Basic auth o header di tenant.
	Headers map[string]string

	// Logger opzionale per i messaggi di debug. Nil ⇒ il logger standard di `log`.
	Logger Logger
}

// Logger è l'interfaccia minima per i log di debug: la soddisfa `*log.Logger`.
type Logger interface {
	Printf(format string, v ...any)
}

// ErrInvalidConfig è la radice degli errori di configurazione: verificabile con errors.Is.
var ErrInvalidConfig = errors.New("haivision: configurazione non valida")

func (c Config) validate() error {
	if c.URL == "" {
		return fmt.Errorf("%w: URL mancante", ErrInvalidConfig)
	}
	u, err := url.Parse(c.URL)
	if err != nil {
		return fmt.Errorf("%w: URL `%s` non parsabile: %v", ErrInvalidConfig, c.URL, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: URL `%s` deve avere schema http o https", ErrInvalidConfig, c.URL)
	}
	if u.Host == "" {
		return fmt.Errorf("%w: URL `%s` senza host", ErrInvalidConfig, c.URL)
	}
	if c.Username == "" {
		return fmt.Errorf("%w: Username mancante", ErrInvalidConfig)
	}
	if c.Password == "" {
		return fmt.Errorf("%w: Password mancante", ErrInvalidConfig)
	}
	if c.Timeout < 0 {
		return fmt.Errorf("%w: Timeout negativo (%s)", ErrInvalidConfig, c.Timeout)
	}
	return nil
}

func (c Config) timeout() time.Duration {
	if c.Timeout == 0 {
		return DefaultTimeout
	}
	return c.Timeout
}
