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
_Nessuna milestone pendente: tutto il backlog pianificato è rilasciato._
<!-- MILESTONE-TABLE:END -->

## [2.0.0] — 2026-08-10

_Contratto API allineato e superficie pulita._

💥 **Release breaking: il codice dei consumer va adeguato.** La guida di migrazione tabellare è nel
[README](README.md#migrazione-da-v1x-a-v200). Il motivo per cui è necessaria: nella v1.x
**create-route non poteva funzionare** (veniva inviato il modello di risposta invece del body
documentato) e **start/stop non deserializzava la risposta** (l'API risponde con un array, il tipo
era una struct). Sistemarlo richiedeva cambiare firme e tipi esportati, quindi tutti i breaking
sono stati raccolti in questa singola major invece di essere spalmati su più release.

⚠️ **Richiede Go 1.25+** (era 1.18).

### Changed (breaking)

- **Costruttore**: `BuildHaivision(url, debug, user, pass, header, insecure)` → `Config` +
  **`New`** (nessun I/O di rete) e **`Connect`** (le due chiamate HTTP), più **`Dial`** = New+Connect.
  Ora è possibile testare la costruzione senza rete e distinguere "configurazione sbagliata" da
  "gateway giù". `Config.Insecure` è un `bool`: il bug del `*bool` non letto della v1.0.0 **non è
  più esprimibile**. `Config.validate()` rifiuta URL senza schema/host, credenziali mancanti e
  timeout negativi con `ErrInvalidConfig`.
- **`context.Context` come primo parametro di tutti i metodi e le funzioni**, propagato a resty.
  Aggiunto `Config.Timeout` (default `DefaultTimeout` = 30s): nella v1.x una chiamata verso un
  gateway irraggiungibile poteva restare appesa indefinitamente e non era cancellabile.
- **`CreateRoute` è una funzione generica** e costruisce il body documentato dai soli campi della
  route: `haivision.CreateRoute(ctx, client, deviceID, fields)`. Nuovo tipo `route.RouteFields[TS,TD]`
  per il corpo delle richieste, distinto da `route.ResponseRouteModel[TS,TD]` che porta i campi di
  stato. Le quattro `CreateRouteSrt`/`Rtmp`/`Rtsp`/`UdpRtp` non esistono più.
- **`GetRoutes` e `GetRouteConfiguration` sono generiche e tipizzate**: `*resty.Response` non è più
  esposto nell'API pubblica. Gli accessi grezzi restano come `GetRoutesRaw` /
  `GetRouteConfigurationRaw`, che tornano `json.RawMessage`.
- **`route.ResponseStartOrRoute` → `route.ResponseRouteCommand`**, che è uno **slice**: l'API
  risponde con un array top-level. Tipi interni estratti in `RouteCommand` e `CommandParameters`.
  `Result` è `json.RawMessage` perché la doc lo mostra a `null` ma può essere un oggetto.
- **Statistiche: 58 campi da `int` a `float64`.** La doc li dà come `number` in Mbit/s, quindi
  frazionari: con `int` un solo valore come 4.5 faceva fallire l'**intera** chiamata. Anche i
  contatori, perché JSON non distingue `42` da `42.0`. Restano `int` solo `port` e `localPort`.
- **Tipi dei campi di richiesta allineati alla doc**: `ttl`/`tos`/`mtu` sono `*int`, `shaping` è
  `*bool`, `maxBitrate` è `*int` (erano `string`/`*string`). I campi opzionali sono ora **puntatori
  con `omitempty`**: non impostati vengono **omessi** dal body invece di essere inviati a zero, che
  il gateway interpreterebbe come valore voluto. Helper `haivision.Bool`/`Int`/`String`.
- **Validazione**: da `gopkg.in/validator.v2` (non manutenuto) a `go-playground/validator/v10`.
  Obbligatori solo i campi che il gateway richiede davvero: nella v1.x `ttl`/`tos`/`retainHeader`
  erano `string` obbligatorie, quindi una route SRT che non le specificava veniva **rifiutata dal
  client** prima di partire. Gli errori sono `*ValidationError`.
- **Rinomine**: `ROUTE_COMMMAND` (tre M) → `ROUTE_COMMANDS`, più `ROUTE_UPDATES` e
  `POST_ROUTE_UPDATES` che rendono espliciti i due endpoint distinti;
  `PrompegFeclsBlockAligned` → `PrompegFecIsBlockAligned`.
- `GetDeviceInfo` restituisce `[]device.ResponseDeviceInfo` invece di `*[]...`.
- Direttiva `go` 1.18 → **1.25.0**; matrice CI ridotta a `1.25.x`.

### Added

- **`UpdateRoute[TS,TD]`, `DeleteRoute`, `StartRoute`, `StopRoute`,
  `StartOrStopDestination[TS,TD]`, `Logout`.** I payload sono stati estratti dal PDF della doc in
  `docs/`, non indovinati: la update ha `elementID` e non ha `startRoute`, la delete non ha `fields`.
- **`ResponseRouteCommand.Pending()`**: i comandi del gateway sono **asincroni**, `state: pending`
  significa accodato e non eseguito. Prima non c'era modo di accorgersene.
- **`*DecodeError`**: il gateway ha risposto 2xx con un body che non corrisponde al modello (es. la
  pagina HTML di un proxy). Porta l'etichetta della chiamata, perché un "cannot unmarshal" nudo non
  dice quale endpoint ha risposto male.
- `Config.Logger` per dirottare i log di debug, `ErrInvalidConfig`, `ErrNotConnected`,
  `route.ActionCreate`/`ActionUpdate`/`ActionDelete`, `route.DestinationActionStart`/`Stop`.
- `Connect` è **rieseguibile**: dopo una sessione scaduta (`APIError.IsUnauthorized()`) basta
  richiamarla. Rifiuta un 2xx senza `sessionID`.

### Fixed

- **`CreateRoute*` non poteva funzionare**: il body non conteneva `action`, `deviceID`,
  `elementType` né il wrapper `fields`. Ora è verificato ispezionando il body che lo stub riceve
  davvero, e che nessun campo di stato della risposta finisca nella richiesta.
- **La risposta di start/stop route non era deserializzabile**: `cannot unmarshal array into Go
  value`. I due test di caratterizzazione `TestKnownBug_*` della v1.1.0 sono stati **invertiti** in
  test che asseriscono il comportamento corretto, come previsto quando furono scritti.
- **`rtmpMode` era serializzato come `RtmpMode`** (maiuscolo) nella source RTMP: il gateway non lo
  riconosceva.

### Security

- `golang.org/x/net` v0.35.0 → **v0.57.0**: `govulncheck` non segnala più **nessuna** vulnerabilità
  in `x/net` (GO-2026-4918 inclusa), e le vulnerabilità nei moduli richiesti scendono da 22 a 9.
  Era il bump rinviato dalla v1.1.0 perché richiedeva `go 1.25`. Resty v2.14.0 → **v2.17.2**.
- Le 22 vulnerabilità stdlib che `govulncheck` continua a segnalare dipendono dal **toolchain** con
  cui si compila (serve go ≥ 1.25.12), non da `go.mod`: non c'è niente da correggere nel repo.

### Test

- Copertura di `./haivision/...` da 67,2% a **81,3%** (74 test, tutti offline e deterministici).
- Nuovi: body reale di create/update/delete, `StartOrStopDestination`, validazione che rifiuta prima
  della rete, route minimale senza campi opzionali, `Logout`, cancellazione via context, timeout,
  validazione di `Config`, `New` che non apre connessioni, `Connect` rieseguibile, tipi dei campi
  contro l'esempio della doc.

### Note per chi aggiorna

Se hai bisogno di restare su Go < 1.25, la **v1.2.0** resta disponibile e riceve i fix di
correttezza della v1.1.0 — ma non quelli di create-route e start/stop, che sono breaking per
costruzione.

## [1.2.0] — 2026-08-10

_Qualità, CI, documentazione._ Nessun cambio di API: solo aggiornamenti di dipendenze compatibili,
test, CI e documentazione. Chi usa la v1.1.0 può aggiornare senza toccare il proprio codice.

### Added

- **Gate di qualità in CI** (`quality.yml`): `gofmt` (blocca), `go vet`, `staticcheck` e
  `govulncheck`. Fino alla v1.1.0 la CI faceva solo `go build` + `go test`, per questo cinque file
  non formattati sono rimasti in `main` per mesi senza che nessuno lo segnalasse. `govulncheck` è
  informativo di proposito: segnala anche le vulnerabilità della stdlib del toolchain, che si
  chiudono aggiornando Go e non toccando `go.mod`.
- **Soglia di copertura al 60%** su `./haivision/...` come job dedicato, con l'artefatto `cover.out`.
- **Verifica che `go.mod`/`go.sum` siano in pari** in CI. Prima `go mod tidy` girava come step di
  setup, quindi un disallineamento veniva silenziosamente corretto invece di essere segnalato.
- `scripts/changelog-extract.py`: estrae la sezione di CHANGELOG di una versione. È quello che
  finisce nel corpo della GitHub Release.
- README: **esempio d'uso completo e compilabile** (verificato compilando un modulo separato),
  sezione sulla gestione degli errori, tabella delle API disponibili, sezione «Limitazioni note».

### Changed

- **Matrice CI**: `1.18.x` (il floor dichiarato in `go.mod`, che va testato, altrimenti è una
  promessa non verificata) più `1.23`/`1.24`/`1.25`. Prima era `1.18`–`1.21`, tutte EOL.
  `actions/checkout@v4`, `actions/setup-go@v5`, test con `-race`.
- **Cache dei moduli finalmente attiva**: `cache-dependency-path` puntava a `subdir/go.sum`, un
  path inesistente residuo di template, ed era inerte solo perché `cache:` era commentato.
- **Workflow di release**: `actions/create-release@v1` (archiviata dal 2021) →
  `softprops/action-gh-release@v2`; `permissions: write-all` → `contents: write`; rimossa
  l'installazione di **ffmpeg**, che nessuna parte del repo usa. Il corpo della release è ora la
  sezione rifinita a mano del CHANGELOG. Un tag con suffisso (`v2.0.0-rc.1`) diventa una prerelease.
- **Dependabot è l'unico bot** per le dipendenze: `renovate.json` è disabilitato con la motivazione
  scritta dentro (due bot sullo stesso `go.mod` aprono PR concorrenti). Aggiunti raggruppamento
  minor/patch in una sola PR e un **ignore su `golang.org/x/net`** per minor/major, perché dalla
  v0.36 richiede `go >= 1.23` e alzerebbe la direttiva del modulo: non è una decisione da bot.
- Rimossa l'entry Dependabot su `/tests`, directory che non esiste (è `test/`, senza `go.mod` proprio).

### Fixed

- **README e `docs/index.md` documentavano un import che non compila**
  (`github.com/Allan-Nava/Haivision-go-sdk` invece di `.../haivision`) e dichiaravano «Go 1.13 or
  later» quando il codice usa i generics e richiede 1.18+. Rimosse anche funzionalità mai esistite
  («stop stream», «play stream»).
- **Data errata di `[1.0.0]` nel CHANGELOG**: era 2023-11-13 (la data dei file), il tag è del
  2023-07-27.

### Removed

- Dead code: i `GetRoutes*` per protocollo commentati, le interface `Response`/`Route`, `BaseSource`
  e le due `RequestUdpRtpCreateRoute` mai usate. `RequestCreateRoute` è stata **tenuta**: serve al fix
  di `CreateRoute*` in v2.0.0.
- I tre file che contenevano solo `package`: `haivision/stats/request.go`,
  `haivision/device/request.go`, `haivision/rtsp/response.go`.

### Security

- Bump resty `v2.7.0` → **`v2.14.0`** (del 2022 → 2024). Non oltre: dalla v2.15.0 resty richiede
  `go 1.20`, che alzerebbe la direttiva del modulo e romperebbe i consumer su toolchain più vecchi.
  Il resto del bump è agganciato alla v2.0.0, che alza `go` a 1.25. Verificato che le hook
  `OnRequestLog`/`OnResponseLog`, su cui poggia la redazione delle credenziali introdotta in v1.1.0,
  esistano ancora nella v2.14.0.

### Test

- Copertura di `./haivision/...` da **37,9% a 67,2%** (43 test, tutti offline e deterministici).
- Lo stub di gateway ha ora una variante **HTTPS con certificato self-signed**, che verifica
  *comportamentalmente* il fix del flag `insecure` della v1.1.0: con `&false` e `nil` la connessione
  viene rifiutata, con `&true` accettata. Prima era verificato solo leggendo il codice.
- Coperti: tutte le statistiche (path + query param), `GetRoutes`, `GetRouteConfiguration`,
  `GetSessionInfo`, `GetDeviceInfo`, 404 → `IsNotFound()`, body HTML non-JSON, gateway irraggiungibile.

### Note per chi aggiorna

Le limitazioni della v1.1.0 **restano tutte**: `CreateRoute*` invia ancora il modello di risposta,
`StartOrStopRoute` non deserializza l'array di risposta, le statistiche frazionarie falliscono, non
c'è `context.Context`. Sono cambiamenti breaking, tutti raccolti nella **v2.0.0**. Questa release
sistema il contorno — CI, test, documentazione — non il contratto API.

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

## [1.0.0] — 2023-07-27

Prima release stabile: l'SDK copre sessione, device, route e statistiche del gateway.

### Added

- Client per la REST API di Haivision Media Gateway / SRT Gateway: sessione (`InitSession`,
  `GetSessionInfo`), device (`GetDeviceInfo`), route (list, configuration, create per SRT/RTMP/RTSP/UDP-RTP,
  start/stop) e statistiche (route, source, destination, client SRT).
- Modelli tipizzati per protocollo con generics Go 1.18 (`route.RouteModel[TS, TD]`).
- `HeaderConfigurator` per header custom e Basic auth.
- `Makefile` con i target `build` e `test`.

### Changed

- CI: `actions/setup-go` alla v4, workflow di build/test rivisti, cache dei moduli tentata
  (`cache-dependency-path` è rimasto rotto fino alla v1.2.0).

### Security

- Bump `golang.org/x/net` da `0.0.0-20211029224645` a `v0.7.0`.

## [0.1.29] — 2023-02-07

_Linea `0.1.x`: 24 tag fra il 2023-01-27 e il 2023-02-07, lo sviluppo iniziale dell'SDK._ Le voci qui
sotto sono una **ricostruzione a grana grossa** dai messaggi di commit (`git log v0.1.0..v0.1.29`): a
quel tempo il CHANGELOG non veniva mantenuto. Per il dettaglio per-tag vedi i
[tag su GitHub](https://github.com/Allan-Nava/Haivision-go-sdk/tags).

### Added

- Autenticazione: `InitSession` (`POST /api/session`) con cookie di sessione, `GetDeviceInfo`
  (`GET /api/devices`), `GetSessionInfo`.
- Route: list, configuration e `CreateRoute*` per SRT / RTMP / RTSP / UDP-RTP; `StartOrStopRoute`.
- Statistiche: `GetRouteStatistics`, `GetSourceStatistics`, `GetDestinationStatisticsById` /
  `ByName`, `GetSrtClientStatistics`.
- `GetDeviceID()` / `GetHType()` per non ripetere il deviceID a ogni chiamata.
- Opzione TLS `insecure` e modalità debug sul client resty.
- Documentazione GitHub Pages (`docs/`, tema just-the-docs): login/device info, route, statistiche.
- Primi test in `test/`.

### Fixed

- `collectedAt` delle statistiche tipizzato `int64` (prima andava in overflow/errore).
- Modello di risposta di `GetDeviceInfo` allineato al payload reale.
- Vari fix a `builder.go` e alla costruzione del body di `InitSession`.

### Note

Tre tag di questa linea **non sono semver validi** — `v0.1.01`, `v0.1.02`, `v0.1.03`, con lo zero
iniziale — e per Go valgono come `v0.1.1`/`v0.1.2`/`v0.1.3`. Sono lasciati come sono per non rompere
eventuali `go.sum` esistenti: da `v0.1.10` in poi la numerazione è regolare. Da v1.1.0 il formato del
tag è imposto da `scripts/new-release.py`, che accetta solo `vX.Y.Z`.

## [0.1.0] — 2023-01-27

### Added

- Primo scheletro del progetto: modulo Go, licenza MIT, wrapper resty e primi modelli di richiesta.
