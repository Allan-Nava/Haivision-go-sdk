# Haivision Go SDK
[![Go build](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-build.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-build.yml)
[![Go test workflow](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-test.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-test.yml)
[![Release on Tag](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/tag-autorelease.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/tag-autorelease.yml)


The Haivision Go SDK is a software development kit for interacting with Haivision's video streaming platform using the Go programming language. It allows developers to easily integrate Haivision's platform into their Go applications, providing access to a wide range of features such as live streaming, video playback, and media management.

## Installation

Serve **Go 1.25 o superiore**.

```bash
go get github.com/Allan-Nava/Haivision-go-sdk
```

### BROADCAST VIDEO ROUTING
Haivision SRT Gateway is a highly flexible and scalable broadcast solution for secure routing of live video streams across different types of IP networks. By serving as a network bridge and converting between protocols including SRT, SRT Gateway provides broadcasters with cost-effective live video streaming to one or multiple destinations for content production and distribution. SRT Gateway includes Haivision’s Path Redundancy feature ensuring uninterrupted live IP video streaming of premium content.

[Rest Api Reference](https://doc.haivision.com/HMG3.7.6/rest-api-integrator-s-reference/rest-api-reference)

## Usage

Il package importabile è **`.../haivision`**: la root del modulo non contiene file Go.

```go
package main

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/route"
	"github.com/Allan-Nava/Haivision-go-sdk/haivision/srt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// New non apre connessioni: valida la configurazione e prepara il client HTTP.
	// Dial = New + Connect, quando non serve separare i due passi.
	client, err := haivision.Dial(ctx, haivision.Config{
		URL:      "https://gateway.example.com",
		Username: "haiadmin",
		Password: "password",
		Timeout:  30 * time.Second, // zero ⇒ haivision.DefaultTimeout
		// Insecure: true,          // solo per certificati self-signed
		// Debug:    true,          // credenziali e sessionID mascherati nei log
	})
	if err != nil {
		log.Fatalf("connessione al gateway: %v", err)
	}
	defer client.Logout(ctx)

	deviceID := client.GetDeviceID() // primo device restituito da /api/devices

	// --- creare una route SRT ---
	fields := route.RouteFields[srt.RequestSourceModelSRT, srt.RequestDestinationModelSrt]{
		Name:       "evento-live",
		StartRoute: haivision.Bool(true),
		Source: srt.RequestSourceModelSRT{
			Name: "ingresso", Address: "0.0.0.0", Protocol: "srt", Port: 2000,
		},
		Destinations: []srt.RequestDestinationModelSrt{{
			Name: "uscita", Address: "10.0.0.9", Protocol: "srt", Port: 2001,
			Ttl: haivision.Int(64), // i campi opzionali non impostati vengono OMESSI dal body
		}},
	}
	if _, err := haivision.CreateRoute(ctx, client, deviceID, fields); err != nil {
		log.Fatalf("creazione route: %v", err)
	}

	// --- elencare le route, tipizzate per protocollo ---
	routes, err := haivision.GetRoutes[srt.ResponseSourceSrt, srt.ResponseDestinationSrt](
		ctx, client, deviceID)
	if err != nil {
		log.Fatalf("lettura route: %v", err)
	}
	for _, r := range routes.Data {
		log.Printf("route %s (%s): %s", r.Name, r.ID, r.State)
	}

	// --- avviare una route ---
	// i comandi del gateway sono ASINCRONI: `pending` significa accodato, non partito
	cmds, err := client.StartRoute(ctx, deviceID, routes.Data[0].ID)
	if err != nil {
		log.Fatalf("start route: %v", err)
	}
	if cmds.Pending() {
		log.Println("comando accodato, la route non è ancora partita")
	}

	// --- statistiche ---
	st, err := client.GetRouteStatistics(ctx, deviceID, routes.Data[0].ID)
	if err != nil {
		var apiErr *haivision.APIError
		if errors.As(err, &apiErr) && apiErr.IsUnauthorized() {
			// la sessione del gateway scade e l'SDK non la rinnova: riconnettiti
			if err := client.Connect(ctx); err != nil {
				log.Fatalf("riconnessione: %v", err)
			}
		} else {
			log.Fatalf("statistiche: %v", err)
		}
	} else {
		log.Printf("bitrate %.2f Mbit/s", st.Route.Source.Bitrate)
	}
}
```

### Gestione degli errori

L'SDK restituisce tre tipi di errore, tutti verificabili con `errors.As` / `errors.Is`:

| Tipo | Significato |
|---|---|
| `*haivision.APIError` | Il gateway ha risposto con uno status >= 400. `IsUnauthorized()` (401/403) di norma vuol dire **sessione scaduta**: richiama `Connect`. `IsNotFound()` per i 404. |
| `*haivision.ValidationError` | Il corpo della richiesta non è valido: **la richiesta non è stata inviata**. |
| `*haivision.DecodeError` | Il gateway ha risposto 2xx con un body che non corrisponde al modello (es. la pagina HTML di un proxy). |

Errori sentinella: `haivision.ErrInvalidConfig`, `haivision.ErrNoDevices`.

## API disponibili

| Area | Come |
|---|---|
| Connessione | `New`, `Dial`, `Connect`, `HealthCheck`, `Logout` |
| Sessione e device | `InitSession`, `GetSessionInfo`, `GetDeviceInfo`, `GetDeviceID`, `GetHType` |
| Route (lettura) | `GetRoutes[TS,TD]`, `GetRouteConfiguration[TS,TD]`, `GetRoutesRaw`, `GetRouteConfigurationRaw` |
| Route (scrittura) | `CreateRoute[TS,TD]`, `UpdateRoute[TS,TD]`, `DeleteRoute`, `StartRoute`, `StopRoute`, `StartOrStopRoute`, `StartOrStopDestination[TS,TD]` |
| Statistiche | `GetRouteStatistics`, `GetSourceStatistics`, `GetDestinationStatisticsById`/`ByName`, `GetSrtClientStatistics` |

Le operazioni che dipendono dal protocollo sono **funzioni generiche**, non metodi: in Go i metodi
non possono avere type parameter, e metterle sull'interfaccia richiederebbe quattro varianti per
ognuna (SRT, RTMP, RTSP, UDP/RTP). Il resto sono metodi su `*Client`, esposti anche
dall'interfaccia `IHaivisionClient` per chi vuole mockarli.

### Cose da sapere sul gateway

- **La sessione scade** e l'SDK non la rinnova: guarda `ExpireAt` di `GetSessionInfo`, oppure
  intercetta `APIError.IsUnauthorized()` e richiama `Connect`.
- **`UpdateRoute` sostituisce, non modifica**: le destinazioni che ometti vengono rimosse. Per
  cambiarne una sola, rileggi la route con `GetRouteConfiguration` e rimanda l'elenco completo.
- **Avviare/fermare una singola destinazione** non ha un endpoint dedicato: è una update con
  `Action` valorizzata su quella destinazione. Ci pensa `StartOrStopDestination`.
- **`DeviceID` viene dal primo device** restituito da `/api/devices`: in un setup multi-device
  passa il `deviceID` esplicito ai metodi.
- **I campi opzionali sono puntatori**: non impostati vengono **omessi** dal body, non inviati a
  zero. Helper: `haivision.Bool`, `haivision.Int`, `haivision.String`.

## Migrazione da v1.x a v2.0.0

La v2.0.0 è una release **breaking**: la v1.x non poteva funzionare per create-route e
start/stop (vedi il [CHANGELOG](CHANGELOG.md)), e sistemarlo richiedeva cambiare firme e tipi.

| v1.x | v2.0.0 |
|---|---|
| `BuildHaivision(url, debug, user, pass, header, insecure)` | `haivision.Dial(ctx, haivision.Config{...})` — oppure `New` + `Connect` per separare costruzione e I/O |
| `insecure *bool` (dove `&false` **disabilitava** il TLS) | `Config.Insecure bool` |
| `header *HeaderConfigurator` | `Config.Headers map[string]string` (`HeaderConfigurator.GetHeaders()` per riusarlo) |
| `client.GetRouteStatistics(dev, route)` | `client.GetRouteStatistics(ctx, dev, route)` — **tutti** i metodi accettano un `context.Context` come primo parametro |
| `client.CreateRouteSrt(dev, *route.RouteModel[...])` | `haivision.CreateRoute(ctx, client, dev, route.RouteFields[...])` |
| `client.GetRoutes(dev) (*resty.Response, error)` | `haivision.GetRoutes[TS,TD](ctx, client, dev)` tipizzata, o `client.GetRoutesRaw(ctx, dev)` |
| `route.ResponseStartOrRoute` (struct) | `route.ResponseRouteCommand` (slice, con `.Pending()`) |
| `route.RouteModel[TS,TD]` come corpo di richiesta | `route.RouteFields[TS,TD]`; `ResponseRouteModel[TS,TD]` è la forma di risposta |
| statistiche `int` | statistiche `float64` (i valori sono Mbit/s frazionari) |
| `ttl`/`tos`/`mtu` `string`, `shaping`/`maxBitrate` `*string` | `*int` / `*bool`, come negli esempi della doc |
| `ROUTE_COMMMAND` (tre M) | `ROUTE_COMMANDS`, con `POST_ROUTE_COMMAND` |
| campi obbligatori: tutti | solo quelli che il gateway richiede davvero; gli opzionali sono puntatori omessi se nil |
| `gopkg.in/validator.v2` | `go-playground/validator/v10` (errori come `*ValidationError`) |
| Go 1.18+ | **Go 1.25+** |

### Support
If you have any issues or need assistance using the Haivision Go SDK, please contact the developer at allan.nava@hiway.media or visit the project's issue tracker at https://github.com/Allan-Nava/Haivision-go-sdk/issues

### Contribution
We welcome contributions to the Haivision Go SDK. If you would like to contribute, please fork the repository, make your changes, and submit a pull request. When submitting a pull request, please make sure to follow the project's coding style and include tests for any new functionality.

### License
The Haivision Go SDK is released under the MIT License and can be used for both personal and commercial projects.



