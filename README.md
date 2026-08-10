# Haivision Go SDK
[![Go build](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-build.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-build.yml)
[![Go test workflow](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-test.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/go-test.yml)
[![Release on Tag](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/tag-autorelease.yml/badge.svg)](https://github.com/Allan-Nava/Haivision-go-sdk/actions/workflows/tag-autorelease.yml)


The Haivision Go SDK is a software development kit for interacting with Haivision's video streaming platform using the Go programming language. It allows developers to easily integrate Haivision's platform into their Go applications, providing access to a wide range of features such as live streaming, video playback, and media management.

## Installation

Serve **Go 1.18 o superiore**: l'SDK usa i generics (`route.RouteModel[TS, TD]`).

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
	"log"

	"github.com/Allan-Nava/Haivision-go-sdk/haivision"
)

func main() {
	// BuildHaivision apre la sessione: fa POST /api/session + GET /api/devices.
	// insecure: nil o &false ⇒ certificato TLS verificato; &true ⇒ verifica disabilitata.
	client, err := haivision.BuildHaivision(
		"https://gateway.example.com",
		false,          // debug (con true, credenziali e sessionID sono mascherati nei log)
		"haiadmin", "password",
		nil,            // *HeaderConfigurator: header custom / Basic auth di un proxy davanti al gateway
		nil,            // insecure
	)
	if err != nil {
		log.Fatalf("connessione al gateway: %v", err)
	}

	if err := client.HealthCheck(); err != nil {
		log.Fatalf("gateway non raggiungibile o sessione scaduta: %v", err)
	}

	deviceID := client.GetDeviceID() // primo device restituito da /api/devices

	stats, err := client.GetRouteStatistics(deviceID, "<route-id>")
	if err != nil {
		log.Fatalf("statistiche: %v", err)
	}
	log.Printf("route %s: stato %s, bitrate %v Mbit/s",
		stats.Route.Name, stats.Route.State, stats.Route.Source.Bitrate)
}
```

### Gestione degli errori

Ogni risposta HTTP di errore del gateway diventa un `*haivision.APIError`:

```go
var apiErr *haivision.APIError
if errors.As(err, &apiErr) {
	if apiErr.IsUnauthorized() {
		// 401/403: sessione scaduta (il gateway la fa scadere e l'SDK non la rinnova)
		// oppure credenziali/ruolo insufficienti → ricostruisci il client
	}
	log.Printf("%s %s: %d — %s", apiErr.Method, apiErr.URL, apiErr.StatusCode, apiErr.Body)
}
```

`haivision.ErrNoDevices` segnala che `GET /api/devices` ha risposto con una lista vuota.

## API disponibili

| Area | Metodi |
|---|---|
| Sessione | `InitSession`, `GetSessionInfo`, `HealthCheck` |
| Device | `GetDeviceInfo`, `GetDeviceID`, `GetHType` |
| Route | `GetRoutes`, `GetRouteConfiguration`, `CreateRouteSrt`/`Rtmp`/`Rtsp`/`UdpRtp`, `StartOrStopRoute` |
| Statistiche | `GetRouteStatistics`, `GetSourceStatistics`, `GetDestinationStatisticsById`/`ByName`, `GetSrtClientStatistics` |

### ⚠️ Limitazioni note in v1.x

- **`CreateRoute*` non funziona ancora contro il gateway**: invia il modello di risposta invece del body documentato (manca il wrapper `action`/`deviceID`/`elementType`/`fields`).
- **`StartOrStopRoute` colpisce l'endpoint giusto ma non deserializza la risposta**: l'API risponde con un array top-level, il tipo di ritorno è una struct.
- **Le statistiche con valori frazionari falliscono**: `bitrate`/`sendRate`/`usedBandwidth` sono documentati in Mbit/s ma tipizzati `int`.
- Nessun `context.Context` e nessun timeout: le chiamate non sono cancellabili.
- Mancano update/delete di una route e la gestione delle singole destinazioni.

Sono tutti cambiamenti **breaking**, pianificati in **v2.0.0**: vedi la [roadmap](docs/roadmap.md) e il [backlog](docs/backlog.md).

## Sviluppo

```bash
make help          # elenco dei target
make check         # gofmt + vet + build + test + backlog-lint + roadmap-check
make cover         # copertura per package
make release-dry   # mostra quale release verrebbe tagliata dal backlog
```

I todo stanno **solo** in [`docs/backlog.md`](docs/backlog.md) (sorgente unica): da lì si generano la
[roadmap per milestone di versione](docs/roadmap.md) e le sezioni del [CHANGELOG](CHANGELOG.md).
Le release si tagliano con `make release`, che deriva la versione dal backlog. Convenzioni per gli
agent AI: [`AGENTS.md`](AGENTS.md) / [`CLAUDE.md`](CLAUDE.md).

### Support
If you have any issues or need assistance using the Haivision Go SDK, please contact the developer at allan.nava@hiway.media or visit the project's issue tracker at https://github.com/Allan-Nava/Haivision-go-sdk/issues

### Contribution
We welcome contributions to the Haivision Go SDK. If you would like to contribute, please fork the repository, make your changes, and submit a pull request. When submitting a pull request, please make sure to follow the project's coding style and include tests for any new functionality.

### License
The Haivision Go SDK is released under the MIT License and can be used for both personal and commercial projects.



