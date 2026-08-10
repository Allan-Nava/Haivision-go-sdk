# Changelog

Tutte le modifiche rilevanti a questo progetto sono documentate qui.
Il formato segue [Keep a Changelog](https://keepachangelog.com/it/1.1.0/) e il progetto aderisce al
[Semantic Versioning](https://semver.org/lang/it/).

> 🔗 Le release **nascono dal backlog**: `make release` deriva la versione dalla prima milestone di
> [`docs/backlog.md`](docs/backlog.md) con 0 item open, scrive questa sezione, esegue i gate, committa e
> tagga (mai push). Il piano delle versioni è in [`docs/roadmap.md`](docs/roadmap.md), generata dallo
> stesso backlog.
>
> Cosa significa `impact` per chi usa l'SDK: `patch` → nessun cambio di API, `minor` → aggiunte o cambi
> di comportamento retrocompatibili in compilazione, `major` → **il codice del consumer va adeguato**.

## [Unreleased]

Il lavoro pianificato è in [`docs/roadmap.md`](docs/roadmap.md). Tabella aggiornata da `make release`:

<!-- MILESTONE-TABLE:START -->
| Versione | Contenuto | Bump | Stato |
|---|---|---|---|
| `v1.2.0` | Qualità, CI, documentazione | minor | 9 item open |
| `v2.0.0` | Contratto API allineato e superficie pulita | major | 10 item open |
<!-- MILESTONE-TABLE:END -->

## [1.1.0] — 2026-08-10

_Correttezza client, sicurezza, wire fix._ Nessun cambio di firma: il codice esistente continua a
compilare. Due cambi di **comportamento** però contano, ed entrambi possono far emergere errori che
prima erano silenziosi — vedi «Changed».

### Changed

- **I metodi ora restituiscono un errore sulle risposte HTTP di errore.** Fino alla v1.0.0 un 401/404/500
  del gateway veniva deserializzato come se fosse un successo: `err == nil` e una struct a zero-value. Con
  credenziali sbagliate `BuildHaivision` restituiva un client apparentemente valido, con cookie
  `sessionID=""`. Ora un 4xx/5xx diventa un `*APIError` (nuovo tipo esportato) con metodo, URL, status ed
  estratto del body. ⚠️ Codice che ignorava gli errori vedrà comparire errori reali: sono guasti che
  c'erano già, non nuovi.
- **`BuildHaivision(..., insecure)`: il valore puntato ora viene letto.** Prima bastava che il puntatore
  fosse non-nil per attivare `InsecureSkipVerify`, quindi passare `&false` **disabilitava** la verifica del
  certificato TLS — l'opposto dell'intento. Ora `nil` e `&false` danno entrambi TLS verificato, solo `&true`
  disabilita la verifica. ⚠️ Chi passava `&false` credendo di essere sicuro non lo era; chi passava `&false`
  contando sullo skip (es. gateway con certificato self-signed) deve passare `&true`.
- **`HealthCheck()` può fallire.** Prima ritornava `nil` in ogni caso, ramo d'errore incluso: era
  inutilizzabile come probe. Ora interroga `GET /api/session` — quindi verifica anche che la sessione non
  sia scaduta — e propaga l'errore.
- **La libreria non scrive più sul logger globale del chiamante.** I `log.Println` incondizionati di
  `route.go` e `stats.go` sono ora dietro il flag `debug`, con prefisso `[haivision]`.

### Added

- `APIError` (`Method`, `URL`, `StatusCode`, `Status`, `Body`) con `IsUnauthorized()` per 401/403 —
  tipicamente sessione scaduta, il caso in cui ricostruire il client — e `IsNotFound()` per 404.
- `ErrNoDevices`, restituito quando `GET /api/devices` risponde con lista vuota.
- Costante `ROUTES_CLIENT_STATISTICS` e helper `GET_ROUTES_CLIENT_STATISTICS`, `POST_ROUTE_COMMAND`.
- Suite di test di contratto sul wire format (payload letterali della doc Haivision) e test dei percorsi
  HTTP del client su `httptest.Server`: 32 test, tutti offline. Copertura di `./haivision/...`: 37,9%
  (era 0%).
- **Release derivata dal backlog** (`make release`, `scripts/new-release.py`): la versione non si passa a
  mano, è la prima milestone di [`docs/backlog.md`](docs/backlog.md) con 0 item open. Lo script rigenera
  roadmap e tabella delle milestone, esegue i gate, committa e crea il tag annotato — mai `git push`.
  Rifiuta di rilasciare una milestone incompleta, di saltare una versione, di taggare fuori da `main` e di
  committare file tipo chiave/dump/env.
- Target `make`: `help`, `fmt`, `fmt-check`, `vet`, `cover`, `check`, `backlog-lint`, `roadmap`,
  `roadmap-check`, `release-notes`, `release-dry`, `release`.

### Fixed

- **Corpo delle richieste con i nomi dei campi sbagliati.** Le struct annidate `Fields` (create route) e
  `Parameters` (start/stop route) non avevano tag JSON e serializzavano col nome Go: il gateway riceveva
  `"Fields"` e `"Parameters"` maiuscoli invece di `"fields"` e `"parameters"`.
- **Start/stop route inviato all'endpoint sbagliato**: `POST /api/devices/{id}/updates` invece di
  `/api/devices/{id}/commands`.
- **Statistiche del client SRT sul path sbagliato**: `/api/gateway/{id}/statistics` invece di
  `/api/gateway/{id}/statistics/client`.
- **Panic su lista device vuota**: `BuildHaivision` indicizzava il primo device senza controllare la
  lunghezza (`index out of range` invece di un errore).
- **Credenziali in chiaro nei log con `debug: true`**: `resty.SetDebug` logga i body, quindi il
  `POST /api/session` finiva nei log con username e password, e la risposta con il `sessionID`. Ora i valori
  di `password`, `sessionID`, `srtPassPhrase`, `passphrase`, `token`, `secret` e gli header
  `Authorization`/`Cookie`/`Set-Cookie` sono mascherati; il debug resta utilizzabile.
- **Header custom non inviati sulle chiamate di bootstrap**: erano applicati dopo `InitSession` e
  `GetDeviceInfo`, quindi un Basic auth richiesto da un reverse proxy davanti al gateway non partiva e la
  costruzione del client falliva sempre.
- `make build` era rotto (`go build .` sulla root, che non contiene file Go) e 5 file non passavano `gofmt`.

### Security

- Bump `golang.org/x/net` v0.7.0 → **v0.35.0**, la versione più recente compatibile con `go 1.18`.
  Chiarimento sulla vulnerabilità GO-2026-4918: **non è raggiungibile da questo SDK** — l'unico pacchetto
  `x/net` nel build è `publicsuffix`, mentre la CVE sta in `net/http/internal/http2`. L'unica istanza
  raggiungibile è quella della **stdlib** `net/http`, che si chiude compilando con go ≥ 1.25.10. Il bump a
  `x/net` v0.53.0 (prima versione col fix) richiede `go 1.25.0` e alzerebbe la direttiva `go` del modulo:
  è rinviato alla v2.0.0, che è già breaking.
- Vedi anche la correzione del flag `insecure` in «Changed» e la redazione delle credenziali in «Fixed».

### Note per chi aggiorna

`CreateRouteSrt`/`Rtmp`/`Rtsp`/`UdpRtp` e la deserializzazione della risposta di start/stop **restano
incomplete**: questa release corregge endpoint e nomi dei campi, non la forma del body (`CreateRoute*` invia
ancora il modello di risposta) né il tipo di `ResponseStartOrRoute` (l'API risponde con un array). Sono
cambiamenti breaking, pianificati in v2.0.0 — vedi [`docs/roadmap.md`](docs/roadmap.md). Nel frattempo due
test di caratterizzazione (`TestKnownBug_*`) tengono traccia dei due bug.

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
