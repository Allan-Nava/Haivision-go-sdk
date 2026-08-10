# Changelog

Tutte le modifiche rilevanti a questo progetto sono documentate qui.
Il formato segue [Keep a Changelog](https://keepachangelog.com/it/1.1.0/) e il progetto aderisce al
[Semantic Versioning](https://semver.org/lang/it/).

> 🔗 Le sezioni di questo file **nascono dal backlog**: `make release-notes V=vX.Y.Z` genera la sezione
> di una milestone (item raggruppati per `impact`) da [`docs/backlog.md`](docs/backlog.md). Il piano
> delle versioni è in [`docs/roadmap.md`](docs/roadmap.md), generata dallo stesso backlog.
>
> Cosa significa `impact` per chi usa l'SDK: `patch` → nessun cambio di API, `minor` → aggiunte o cambi
> di comportamento retrocompatibili in compilazione, `major` → **il codice del consumer va adeguato**.

## [Unreleased]

Nessuna modifica non rilasciata. Il lavoro pianificato è in [`docs/roadmap.md`](docs/roadmap.md):

| Versione | Contenuto | Bump |
|---|---|---|
| `v1.1.0` | Correttezza client, sicurezza, wire fix — status HTTP non controllato, flag `insecure` invertito, tag JSON `fields`/`parameters`, endpoint `commands`, bump `golang.org/x/net` | minor |
| `v1.2.0` | Qualità, CI, documentazione — copertura test, gate CI, README corretto, cleanup | minor |
| `v2.0.0` | Contratto API allineato e superficie pulita — **breaking**: `CreateRoute*`, `ResponseStartOrRoute`, statistiche `float64`, `context.Context` | major |

## [1.0.0] — 2023-11-13

Prima release stabile. La storia dettagliata delle versioni `v0.1.x` → `v1.0.0` non è stata mantenuta
in questo file: è ricostruibile dai tag (`git log --oneline v0.1.29..v1.0.0`) e dalle
[release GitHub](https://github.com/Allan-Nava/Haivision-go-sdk/releases). Item di backlog che copre la
ricostruzione: `changelog-bootstrap`.

### Added

- Client per la REST API di Haivision Media Gateway / SRT Gateway: sessione (`InitSession`,
  `GetSessionInfo`), device (`GetDeviceInfo`), route (list, configuration, create per SRT/RTMP/RTSP/UDP-RTP,
  start/stop) e statistiche (route, source, destination, client SRT).
- Modelli tipizzati per protocollo con generics Go 1.18 (`route.RouteModel[TS, TD]`).
- `HeaderConfigurator` per header custom e Basic auth.
