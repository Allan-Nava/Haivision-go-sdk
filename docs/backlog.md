---
layout: default
title: Backlog
nav_order: 9
description: "Sorgente unica dei todo e delle milestone di versione"
permalink: /backlog
---

# Backlog — sorgente di verità delle modifiche da fare

Questo file è l'**unica sorgente di verità** dei todo del repo. Niente TODO sparsi nel codice,
nelle issue o nei doc: si aprono qui. Due script stdlib-only lo leggono:

- [`scripts/backlog-lint.py`](https://github.com/Allan-Nava/Haivision-go-sdk/blob/main/scripts/backlog-lint.py) — valida struttura e **coerenza del versionamento** (gate CI).
- [`scripts/generate-roadmap.py`](https://github.com/Allan-Nava/Haivision-go-sdk/blob/main/scripts/generate-roadmap.py) — genera [`docs/roadmap.md`](roadmap.md), la **milestone dinamica**.

## Come funziona (flusso)

```
  docs/backlog.md  ──parse(id)──▶  lib/backlog.py  ──▶  backlog-lint.py      (gate CI)
   ### `id` — Titolo                 (fonte unica     ──▶  generate-roadmap.py
   - status / impact / milestone       delle regole)         │
                                                            ▼
                          git tag vX.Y.Z ──baseline──▶  docs/roadmap.md
                          (max tag esistente)            "Prossima release: vN"
                                                              │
                                     --release-notes vX.Y.Z ──┴──▶ sezione CHANGELOG.md

  make release ──▶ new-release.py ──▶ versione = 1ª milestone pendente con 0 item open
                                  ──▶ roadmap + tabella CHANGELOG (baseline = versione in uscita)
                                  ──▶ gate (gofmt/vet/build/test/lint/roadmap-check)
                                  ──▶ commit + tag annotato        (MAI push)

  La milestone è DINAMICA: "prossima release" = milestone pendente di versione più
  bassa con almeno un item open. Chiudi gli item → avanza da sola. Nessuna lista di
  versioni mantenuta a mano.
```

## Convenzione di scrittura di un item

- Un item inizia con `### \`<id-stabile>\` — <Titolo>` (l'`id` è in backtick, kebab-case, **non cambiarlo mai**: è la chiave di continuità fra roadmap, changelog e commit).
- Metadati come bullet `- **chiave**: valore`:
  - **status**: `open` (default) | `done`.
  - **priority**: `low` | `medium` | `high` — urgenza. Non c'entra col versionamento.
  - **impact**: `patch` | `minor` | `major` (default `patch`) — **effetto sull'API pubblica**, e quindi sulla versione. È il campo che governa il semver.
  - **labels**: lista separata da virgola.
  - **milestone**: `vX.Y.Z` oppure `vX.Y.Z — Etichetta` — la **versione target**. Una versione = un titolo, identico carattere per carattere.
  - **owner**, **ref**: opzionali.
- Tutto il resto del blocco è prosa descrittiva. Non usare heading `#`/`##`/`###` nella prosa: chiudono l'item.

## Come si assegna `impact` (è una libreria pubblica)

| impact | Quando | Esempi |
|---|---|---|
| `patch` | Fix che non tocca firme né tipi esportati. Include i fix del **wire format** (tag JSON, path endpoint): il codice del consumer continua a compilare. | `route-json-tag-fields-parameters`, `device-list-empty-panic` |
| `minor` | Aggiunte retrocompatibili, o cambi di **semantica** senza cambio di firma. | `http-status-check` (metodi che prima tornavano `nil` ora tornano errore) |
| `major` | Qualsiasi cosa che rompe la compilazione di un consumer: firme, tipi/campi esportati, rinomine, **e ogni aggiunta di metodo a `IHaivisionClient`** (rompe chi la implementa per i mock). | `stats-float64`, `create-route-request-model`, `context-and-timeout` |

> ⚠️ **L'aggiunta di un metodo all'interfaccia `IHaivisionClient` è `major`, non `minor`.** In Go
> un'interfaccia esportata è un contratto a due vie: chi la implementa (tipicamente un mock nei test
> del consumer) non compila più. È la ragione per cui `route-update-delete` sta in v2.0.0 e non in
> una minor, pur essendo "solo" una feature nuova.

Il linter **blocca** una milestone il cui bump è inferiore all'`impact` massimo dei suoi item: pianificare
un `major` dentro una minor è un errore di CI, non una svista che si scopre dopo il tag.

## Chiudere un item

Metti `status: done` (preferito: resta la traccia del perché — se la premessa dell'item era sbagliata,
scrivilo, come in `deps-x-net-vuln`) oppure rimuovilo. Quando **tutti** gli item di una milestone sono
`done`, il linter stampa `milestone PRONTA` e la roadmap la marca `✅ pronta al rilascio`: da lì
**`make release`** fa tutto il resto — versione, CHANGELOG, gate, commit e tag annotato. Il push lo fa
sempre l'utente. `make release-dry` mostra cosa farebbe senza scrivere niente.

Lo script rifiuta di rilasciare una milestone con item open, di saltare una versione della catena, di
taggare da un branch diverso da `main` e di committare file tipo chiave/dump/env o più grandi di 1 MB.

> 🔍 **Prima di aprire un item nuovo, cerca il doppione per artefatto, non per parole.** Due item scritti
> in momenti diversi sullo stesso problema usano parole diverse. Cerca ciò che l'intervento *toccherebbe*:
> il file (`grep -n "builder.go" docs/backlog.md`), il simbolo (`InsecureSkipVerify`), l'endpoint
> (`/statistics`), il nome del workflow.

---

## Item attivi

Fonte di tutti gli item aperti al 2026-08-10: [audit tecnico](audit-2026-08-10.md).

### `http-status-check` — nessun metodo controlla lo status HTTP: 401/500 sembrano successi

- **status**: done
- **priority**: high
- **impact**: minor
- **labels**: correctness, client, audit-p1
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §A1](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — controllo centralizzato nei due helper resty (`haivision/haivision.go`): un 4xx/5xx diventa `*APIError` (nuovo `haivision/errors.go`) con metodo, URL, status ed estratto del body **redatto e troncato** a 512 rune. Helper `IsUnauthorized()` (401/403, tipicamente sessione scaduta) e `IsNotFound()`. Coperto da `TestBuildHaivisionWrongCredentialsReturnsError` e `TestBuildHaivisionDeviceErrorReturnsError`.

`restyGet`/`restyPost` in `haivision/haivision.go` ritornano `err == nil` per qualsiasi risposta ricevuta,
poi ogni metodo fa `json.Unmarshal` sul body. Con credenziali errate `InitSession` ritorna un oggetto a
zero-value e **nessun errore**, e `BuildHaivision` prosegue impostando il cookie `sessionID=""`.

Fix nei due soli helper resty (un punto): controllo di `resp.IsError()` e errore tipizzato con status +
body troncato. È `minor` perché metodi che prima tornavano `nil` ora tornano errore: il consumer compila
ancora, ma il comportamento cambia — va scritto nel changelog.

### `insecure-flag-inverted` — `insecure: &false` DISABILITA la verifica TLS

- **status**: done
- **priority**: high
- **impact**: minor
- **labels**: security, tls, audit-p1
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §A2](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — `if insecure != nil && *insecure`: `nil` **e** `&false` danno entrambi TLS verificato, solo `&true` disabilita la verifica. Firma invariata (non-breaking); la pulizia in struct di opzioni resta `builder-options-struct` (v2.0.0).

`builder.go` verifica solo `insecure != nil` e non legge mai il valore puntato: passare un puntatore a
`false` attiva `InsecureSkipVerify`, l'opposto di quanto chiede il chiamante. L'unico modo di avere TLS
verificato oggi è passare `nil`.

Fix minimo e non-breaking: `if insecure != nil && *insecure`. La pulizia della firma (`bool` o struct di
opzioni) è separata, in `builder-options-struct` (v2.0.0). Da annunciare nel changelog come cambio di
comportamento: chi passava `&false` credendo di disabilitare lo skip aveva TLS non verificato.

### `device-list-empty-panic` — panic se il gateway risponde con lista device vuota

- **status**: done
- **priority**: high
- **impact**: patch
- **labels**: correctness, client, audit-p1
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §A3](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — guardia su `nil`/slice vuota + nuovo errore esportato `ErrNoDevices` (verificabile con `errors.Is`). Coperto da `TestBuildHaivisionEmptyDeviceListReturnsError`.

`BuildHaivision` fa `(*deviceResponse)[0].ID` senza controllo di lunghezza: gateway senza device, o body
d'errore che deserializza in slice vuota, danno `index out of range` invece di un errore. Aggiungere la
guardia e un errore esplicito. Documentare che `DeviceID`/`HType` vengono dal **primo** device e che il
multi-device richiede di passare `deviceId` ai metodi.

### `debug-logs-credentials` — con `debug: true` username e password finiscono nei log

- **status**: done
- **priority**: high
- **impact**: patch
- **labels**: security, logging, audit-p1
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §A7](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — hook `OnRequestLog`/`OnResponseLog` registrate **prima** della prima richiesta, con redazione in `haivision/redact.go`: valori di `password`, `sessionID`, `srtPassPhrase`, `passphrase`, `token`, `secret` e header `Authorization`/`Cookie`/`Set-Cookie`. Verificato leggendo `middleware.go` di resty che le hook ricevono una **copia** degli header (`copyHeaders`), quindi la redazione non altera la richiesta reale. `debugResponse` passa da `bodyExcerpt`, che redige e tronca. Test in `haivision/redact_test.go` (in-package: le funzioni sono non esportate).

`resty.SetDebug(true)` logga i body delle richieste: `POST /api/session` finisce nei log con credenziali in
chiaro, e `debugPrint` logga la risposta con il `sessionID`. Redigere i campi sensibili (hook resty
`OnBeforeRequest`/logger custom) invece di disabilitare il debug, così il debug resta utile.

### `route-json-tag-fields-parameters` — body con `"Fields"`/`"Parameters"`: il gateway non li riconosce

- **status**: done
- **priority**: high
- **impact**: patch
- **labels**: api-contract, serialization, audit-p0
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §B1](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — aggiunti `json:"fields"` e `json:"parameters"`. Le struct anonime **non** sono state estratte in tipi con nome: i tag fanno parte del tipo di una struct anonima, quindi estrarle romperebbe chi costruisce quei literal → rinviato a v2.0.0. Coperto da `TestStartStopRouteRequestBody` e `TestCreateRouteRequestBody`, che asseriscono sui nomi dei campi JSON e falliscono se ricompare `"Fields"`/`"Parameters"`.

Le struct anonime annidate `Fields` e `Parameters` in `haivision/route/model.go` non hanno tag JSON, quindi
serializzano col nome Go. Verificato: `{"deviceID":"d1","command":"start-route","Parameters":{"routeID":"r1"}}`
mentre l'API documenta `"parameters"` minuscolo. Aggiungere `json:"fields"` / `json:"parameters"`.
Nessun tipo esportato cambia forma: `patch`.

### `startstop-commands-endpoint` — start/stop route postato su `/updates` invece di `/commands`

- **status**: done
- **priority**: high
- **impact**: patch
- **labels**: api-contract, routes, audit-p0
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §B3](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — aggiunto l'helper `POST_ROUTE_COMMAND` e usato in `StartOrStopRoute`. Coperto da `TestStartOrStopRouteUsesCommandsEndpoint`, che stubba **entrambi** gli endpoint e verifica quale viene colpito. Il refuso nel nome `ROUTE_COMMMAND` resta: rinominare una costante esportata è breaking → `exported-naming-typos` (v2.0.0).

`StartOrStopRoute` usa `POST_CREATE_ROUTE(deviceId)` = `/api/devices/%s/updates`, mentre la doc (e il
commento sopra il metodo stesso) indicano `/api/devices/{id}/commands`. La costante esiste già in
`constants.go` con un typo — `ROUTE_COMMMAND`, tre `M` — e senza helper `Sprintf`, e non è mai usata.
Aggiungere l'helper e usarlo. La rinomina della costante è breaking → sta in `exported-naming-typos`.

### `srt-client-stats-path` — statistiche client SRT sul path senza `/client`

- **status**: done
- **priority**: high
- **impact**: patch
- **labels**: api-contract, stats, audit-p0
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §B6](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — nuova costante `ROUTES_CLIENT_STATISTICS` + helper `GET_ROUTES_CLIENT_STATISTICS`. Coperto da `TestGetSrtClientStatisticsUsesClientSubPath`, che verifica path e i quattro query param.

`GetSrtClientStatistics` usa `GET_ROUTES_STATISTICS` (`/api/gateway/%s/statistics`) mentre la doc indica
`/api/gateway/{id}/statistics/client`. Serve la costante per il sotto-path.

### `deps-x-net-vuln` — `golang.org/x/net` v0.7.0 obsoleto → v0.35.0 (GO-2026-4918 non raggiungibile)

- **status**: done
- **priority**: high
- **impact**: patch
- **labels**: security, dependencies, audit-p2
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §D](audit-2026-08-10.md)

✅ **CHIUSO in v1.1.0 — ma la premessa era in parte SBAGLIATA: GO-2026-4918 non è raggiungibile da questo codice.**

La verifica: l'unico pacchetto di `golang.org/x/net` che entra nel build è **`publicsuffix`** (tirato dal
cookiejar di resty), mentre la vulnerabilità sta in `net/http/internal/http2`. `govulncheck` non produce
**nessuna trace** attraverso `x/net`: l'unica istanza raggiungibile di GO-2026-4918 è quella della
**stdlib** `net/http`, che si chiude aggiornando il toolchain (go ≥ 1.25.10), non `go.mod`.

E il bump alla versione col fix non era comunque fattibile in una minor: **`x/net` v0.53.0 richiede
`go 1.25.0`**, quindi avrebbe alzato la direttiva `go` del modulo da 1.18 a 1.25, impedendo la compilazione
a tutti i consumer su toolchain più vecchi — un breaking change travestito da patch di sicurezza.

**Fatto**: bump a **`x/net` v0.35.0**, la più recente che resta su `go 1.18` (la v0.36.0 passa a
`go 1.23.0`). Chiude 4 vulnerabilità non raggiungibili nei moduli richiesti (26 → 22) a costo zero di
compatibilità. Resta segnalata la sola GO-2026-4918, non raggiungibile.

Il bump vero è tracciato in `x-net-http2-go-directive` (v2.0.0), dove alzare la direttiva `go` è
accettabile perché la release è già breaking.

### `wire-contract-fixture-tests` — test tabellari di ser/deser sui payload della doc Haivision

- **status**: done
- **priority**: high
- **impact**: patch
- **labels**: testing, api-contract, audit-p2
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §C](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — `test/wire_route_test.go`, `test/wire_session_test.go`, `test/wire_stats_test.go`: i payload sono quelli **letterali** della doc e le asserzioni sono sui nomi dei campi JSON (via `marshalToMap`), non sui nomi Go — è l'unico modo di vedere un `"Fields"` maiuscolo. I due bug che restano aperti sono bloccati da test di **caratterizzazione** `TestKnownBug_StartStopResponseIsArray` e `TestKnownBug_StatsFractionalBitrate`: asseriscono che l'unmarshal fallisca oggi e **falliranno** quando gli item v2.0.0 saranno chiusi, forzandone l'aggiornamento.

Tutti i bug di contratto trovati nell'audit (tag JSON, array vs oggetto, `int` vs frazionari) sono
intercettabili con `json.Marshal`/`json.Unmarshal` sui payload **letterali** della doc, offline e senza
gateway. Un caso per ogni request e ogni response documentata, con confronto campo per campo. Va fatto
**nella stessa milestone dei fix**, altrimenti i fix non hanno una rete che li tenga.

### `makefile-build-and-gofmt` — `make build` è rotto e 5 file non passano gofmt

- **status**: done
- **priority**: medium
- **impact**: patch
- **labels**: tooling, dx, audit-p2
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §C](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — `Makefile` riscritto: `build` usa `go build -v ./...`, aggiunti `help`, `fmt`, `fmt-check`, `vet`, `cover`, `check` e i target del backlog. `gofmt -w .` applicato: `gofmt -l .` ora è vuoto.

`make build` esegue `go build .` sulla root, dove non esistono file Go: `no Go files in ...`. Il target
corretto è `go build ./...` (quello che usa la CI, motivo per cui nessuno se n'è accorto). Non formattati:
`haivision/device/response.go`, `haivision/haivision.go`, `haivision/header_configurator.go`,
`haivision/rtsp/response.go`, `haivision/stats/response.go`. Aggiungere i target `fmt`, `vet`, `cover`.

### `header-configurator-order` — header custom applicati dopo il login: rotto dietro proxy autenticato

- **status**: done
- **priority**: medium
- **impact**: patch
- **labels**: correctness, client, audit-p1
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §A4](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — gli header custom sono applicati subito dopo `SetBaseURL`, prima di `InitSession`. Coperto da `TestBuildHaivisionSendsCustomHeadersOnLogin`, che verifica sugli header **realmente ricevuti** dal server sulla `POST /api/session`.

Gli header di `HeaderConfigurator` vengono impostati **dopo** `InitSession` e `GetDeviceInfo`, quindi
`CreateBasicAuthHeader(...)` e gli header richiesti da un reverse proxy davanti al gateway non partono
sulle due chiamate di bootstrap: dietro un proxy con Basic auth il costruttore fallisce sempre. Spostare
l'applicazione degli header subito dopo `SetBaseURL`.

### `healthcheck-always-nil` — `HealthCheck()` non può fallire

- **status**: done
- **priority**: medium
- **impact**: patch
- **labels**: correctness, client, audit-p1
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §A5](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — interroga `GET /api/session` (valida raggiungibilità **e** sessione) e propaga l'errore; risolto anche il path concatenato, dato che `o.Url` era già la BaseURL. Coperto da `TestHealthCheck`, che esercita sessione valida e sessione scaduta (401 → `APIError.IsUnauthorized()`).

Il ramo d'errore fa `return nil`: la funzione ritorna `nil` in ogni caso, inutilizzabile come probe. Inoltre
fa `GET` su `o.Url` che, essendo già la BaseURL di resty, produce un path concatenato non intenzionale.
Correggere entrambe le cose e decidere l'endpoint di liveness (`GET /api/session` è il candidato: verifica
anche che la sessione sia ancora valida).

### `unconditional-log-println` — la libreria scrive sul logger globale del consumer a ogni chiamata

- **status**: done
- **priority**: low
- **impact**: patch
- **labels**: logging, dx, audit-p1
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [audit §A7](audit-2026-08-10.md)

✅ **FATTO in v1.1.0** — tutti i `log.Println` di `route.go`/`stats.go` sostituiti con `o.debugf` (attivo solo con `debug=true`, prefisso `[haivision]`); `debugPrint` sostituito da `debugResponse`, che logga status, durata ed estratto redatto del body.

`route.go` e `stats.go` fanno `log.Println` **incondizionato**, ignorando il flag `debug`: una libreria non
deve inquinare lo stdout del chiamante. Sostituire con `o.debugPrint`.

### `release-tooling-dynamic` — release derivata dal backlog: versione, CHANGELOG, gate, commit e tag

- **status**: done
- **priority**: medium
- **impact**: patch
- **labels**: tooling, release, dx
- **milestone**: v1.1.0 — Correttezza client, sicurezza, wire fix
- **ref**: [roadmap](roadmap.md)

✅ **FATTO in v1.1.0** — `scripts/new-release.py` (+ `make release` / `make release-dry`). La versione **non
si passa a mano**: è la prima milestone pendente con 0 item open. Lo script fa lint del backlog, sceglie la
versione, controlla git (branch, tag libero, nessun file tipo chiave/dump/env, nessun file > 1 MB),
rigenera `docs/roadmap.md` e la tabella delle milestone in `CHANGELOG.md`, esegue i gate
(gofmt/vet/build/test/backlog-lint/roadmap-check), committa e crea il tag annotato. **Mai push.**

Se la sezione di CHANGELOG della versione manca, la scrive dallo scheletro generato dal backlog e **si
ferma** (exit 3): i titoli degli item sono formulati come problemi, la prosa va rifinita in voce da
changelog prima di taggare.

Due difetti trovati durante la messa in opera, entrambi corretti:

- le milestone con versione <= baseline vanno trattate come **storia**, non come pianificazione: subito
  dopo il tag di `vX.Y.Z` quella milestone diventava "non un incremento valido rispetto a sé stessa" e il
  gate `backlog-lint` sarebbe andato rosso a ogni release (`released` in `milestone_chain`);
- roadmap e tabella vanno generate con **baseline = la versione in uscita**, non col max tag attuale: il tag
  non esiste ancora ma il commit lo porterà, e generandole con la baseline vecchia il gate `roadmap-check`
  diventava rosso sul commit di release stesso.

### `httptest-client-coverage` — copertura del package `haivision`: 0%

- **status**: open
- **priority**: high
- **impact**: patch
- **labels**: testing, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §C](audit-2026-08-10.md)

Coprire i metodi HTTP con `httptest.Server` + `BuildHaivision` puntato su quello — mai un gateway reale.

⚠️ **Parzialmente anticipato in v1.1.0**: `test/client_http_test.go` esiste già (stub `gatewayStub`
configurabile per rotta, che registra gli header ricevuti) e copre bootstrap, 401, 500, lista device vuota,
ordine degli header, `HealthCheck` e i due endpoint corretti. I 2 test triviali originali sono stati
sostituiti (`TestDeviceInfo`, che faceva `log.Println(err)` invece di `t.Fatalf`, non c'è più). Copertura di
`./haivision/...` dai test: **37,9%**.

**Resta da coprire**: le quattro `CreateRoute*` (bloccate da `create-route-request-model`, v2.0.0), le
statistiche destination/source, `GetRoutes`/`GetRouteConfiguration`, i body non-JSON e i timeout di rete.
Obiettivo: ≥70% su `./haivision/...` con un gate di coverage in CI.

### `readme-import-path-go-version` — il README documenta un import che non compila

- **status**: open
- **priority**: high
- **impact**: patch
- **labels**: docs, dx, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §D](audit-2026-08-10.md)

`import "github.com/Allan-Nava/Haivision-go-sdk"` non compila: la root del modulo non ha file Go, il
package è `.../haivision`. Il README dichiara anche "Go 1.13 or later" mentre il codice usa generics
(serve 1.18+), e promette funzionalità inesistenti ("stop stream, play stream"). Stesso errore in
`docs/index.md`. Aggiungere un esempio d'uso reale e compilabile.

### `ci-quality-gates` — la CI non ha gate su formato, lint, vulnerabilità, backlog

- **status**: open
- **priority**: medium
- **impact**: patch
- **labels**: ci, tooling, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §D](audit-2026-08-10.md)

Oggi la CI fa solo `go build` + `go test`. Aggiungere: `gofmt -l` (fallisce se non vuoto), `go vet`,
`staticcheck`, `govulncheck`, `backlog-lint.py` e `generate-roadmap.py --check` — quest'ultimo perché la
roadmap è generata e committata, quindi va verificato che sia in pari col backlog. Il gate `backlog` è già
in `.github/workflows/backlog.yml`; restano quelli Go.

### `ci-go-matrix-and-actions` — matrice Go 1.18–1.21 (tutte EOL) e action obsolete

- **status**: open
- **priority**: medium
- **impact**: patch
- **labels**: ci, dependencies, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §D](audit-2026-08-10.md)

Aggiornare la matrice a 1.22–1.25, `actions/checkout@v3`→v4, `actions/setup-go@v4`→v5. Rimuovere
`cache-dependency-path: subdir/go.sum`, residuo di template che punta a un path inesistente: è inerte solo
perché `cache:` è commentato, quindi oggi **non c'è cache dei moduli**. Decidere la versione minima
supportata e allinearla a `go.mod`.

### `tag-autorelease-modernize` — release workflow su action archiviata e permessi eccessivi

- **status**: open
- **priority**: medium
- **impact**: patch
- **labels**: ci, security, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §D](audit-2026-08-10.md)

`tag-autorelease.yml` usa `actions/create-release@v1` (archiviata dal 2021) con `permissions: write-all`
(basta `contents: write`) e installa **ffmpeg** senza che nulla nel repo lo usi. Sostituire con
`gh release create` o `softprops/action-gh-release`, e popolare il corpo della release dalla sezione di
CHANGELOG generata da `make release-notes`.

### `changelog-bootstrap` — nessun CHANGELOG nonostante 30+ tag e release automatiche

- **status**: open
- **priority**: medium
- **impact**: patch
- **labels**: docs, release, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §C](audit-2026-08-10.md)

Lo scheletro Keep a Changelog è in `CHANGELOG.md`, ma la storia da `v0.1.0` a `v1.0.0` è da ricostruire
dai tag (`git log --oneline v0.1.29..v1.0.0`) almeno a grana grossa. Da lì in avanti ogni release nasce da
`make release-notes V=vX.Y.Z`.

### `deps-resty-bump` — resty v2.7.0 è del 2022

- **status**: open
- **priority**: medium
- **impact**: patch
- **labels**: dependencies, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §D](audit-2026-08-10.md)

Bump all'ultima v2 (API-compatibile nel major, quindi `patch` per noi). Da fare **dopo**
`httptest-client-coverage`: senza test sui metodi HTTP un cambio di comportamento del client resty passa
inosservato. Valutare in quella sede se `net/http` puro basti, riducendo la superficie di dipendenze di
una libreria.

### `dead-code-and-stubs-cleanup` — blocchi commentati e file stub vuoti

- **status**: open
- **priority**: low
- **impact**: patch
- **labels**: cleanup, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §C](audit-2026-08-10.md)

Da rimuovere o completare: i `GetRoutes*` per protocollo commentati, le interface `Response`/`Route` in
`route/response.go`, `RequestUdpRtpCreateRoute`, `BaseSource`. File con solo `package`:
`haivision/stats/request.go`, `haivision/device/request.go`; `haivision/rtsp/response.go` è vuoto (0 byte).
`RequestCreateRoute` è dead code oggi ma **serve** a `create-route-request-model`: non rimuoverla.

### `dependabot-tests-dir` — entry dependabot su una directory che non esiste

- **status**: open
- **priority**: low
- **impact**: patch
- **labels**: ci, cleanup, audit-p2
- **milestone**: v1.2.0 — Qualità, CI, documentazione
- **ref**: [audit §D](audit-2026-08-10.md)

`.github/dependabot.yml` monitora `/tests`: la directory è `test/` e non ha un `go.mod` proprio, quindi
quell'entry è morta. Inoltre Dependabot e Renovate sono entrambi attivi sullo stesso `gomod` della root:
scegliere uno dei due e rimuovere l'altro.

### `create-route-request-model` — `CreateRoute*` invia il modello di risposta, non la richiesta

- **status**: open
- **priority**: high
- **impact**: major
- **labels**: api-contract, routes, breaking, audit-p0
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §B2](audit-2026-08-10.md)

Le quattro `CreateRoute*` accettano `*route.RouteModel[TS,TD]` — la forma **di risposta**, con `id`,
`state`, `elapsedTime`, `summaryStatusCode` — e la postano così com'è: il body non contiene `action`,
`deviceID`, `elementType` né il wrapper `fields` richiesti da `POST /api/devices/{id}/updates`. La struct
corretta `RequestCreateRoute` esiste in `route/model.go` ma non è referenziata da nessuna parte.

Il cambio di tipo del parametro è breaking. Attenzione: appena `RequestCreateRoute` entra in uso, il suo tag
`validate:"nonnil,min=1"` su `Fields.StartRoute bool` fa fallire **sempre** la validazione
(`unsupported type`) — va risolto insieme, vedi `validator-v10-optional-fields`.

### `startstop-response-slice` — `ResponseStartOrRoute` non deserializza la risposta reale

- **status**: open
- **priority**: high
- **impact**: major
- **labels**: api-contract, routes, breaking, audit-p0
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §B4](audit-2026-08-10.md)

L'API risponde con un **array top-level**; la struct wrappa in un campo `Response []struct{...}` senza tag
JSON. Verificato: `json: cannot unmarshal array into Go value of type route.ResponseStartOrRoute`. Il tipo
deve diventare uno slice — cambio di forma di un tipo esportato, quindi breaking. Estrarre anche la struct
anonima interna in un tipo con nome, così i consumer possono dichiararla.

### `stats-float64` — bitrate e rate in Mbit/s tipizzati `int`: ogni valore frazionario rompe la chiamata

- **status**: open
- **priority**: high
- **impact**: major
- **labels**: api-contract, stats, breaking, audit-p0
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §B5](audit-2026-08-10.md)

`bitrate`, `sendRate`, `usedBandwidth` sono documentati come `number` in **Mbit/s** (quindi frazionari) ma
sono `int`. Verificato: `json: cannot unmarshal number 4.5 into Go struct field
SourceStatisticsModel.bitrate of type int` — una sola route con bitrate non intero fa fallire **l'intera**
chiamata `Get*Statistics`. Passare a `float64` su tutti i modelli di `haivision/stats/`, rileggendo la doc
campo per campo per non convertire a caso i contatori (che restano interi).

### `validator-v10-optional-fields` — la validazione rifiuta route legittime e si rompe sui bool

- **status**: open
- **priority**: medium
- **impact**: major
- **labels**: validation, dependencies, breaking, audit-p0
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §B7](audit-2026-08-10.md)

`validate:"nonnil,min=1"` è su **tutti** i campi, inclusi gli opzionali (`ttl`, `tos`, `retainHeader`, tutti
`string`): verificato che una route SRT senza `ttl`/`tos` viene rifiutata lato client prima di partire. E
`min=1` su un `bool` dà `unsupported type`, cioè validazione sempre fallita. `gopkg.in/validator.v2` è di
fatto non manutenuto: migrare a `go-playground/validator/v10`, tenere obbligatori solo i campi che l'API
richiede davvero e usare puntatori per gli opzionali. Breaking sui tag e sul tipo degli errori restituiti.

### `typed-get-routes` — `GetRoutes` restituisce `*resty.Response`: il trasporto è nell'API pubblica

- **status**: open
- **priority**: medium
- **impact**: major
- **labels**: api-surface, routes, breaking, audit-p0
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §B8](audit-2026-08-10.md)

`GetRoutes` e `GetRouteConfiguration` ritornano la risposta resty grezza: il consumer deserializza a mano e
il tipo della libreria HTTP è esposto nell'interfaccia (impossibile cambiarla senza rompere tutti). I
modelli tipizzati `ResponseRoutes[TS,TD]` esistono già ma sono commentati. Completarli e restituire quelli,
un metodo per protocollo come per `CreateRoute*`.

### `context-and-timeout` — nessun `context.Context` e nessun timeout: chiamate non cancellabili

- **status**: open
- **priority**: medium
- **impact**: major
- **labels**: api-surface, reliability, breaking, audit-p1
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §A6](audit-2026-08-10.md)

Nessun metodo accetta un context e il client resty non ha `SetTimeout`: verso un gateway irraggiungibile una
chiamata può bloccarsi a lungo, e chi usa l'SDK dentro un handler HTTP non può propagare la cancellazione.
Aggiungere il `ctx` come primo parametro (rifacimento dell'interfaccia) e un timeout di default
configurabile. È `major` perché tocca tutte le firme **e** perché ogni aggiunta a `IHaivisionClient` rompe
chi la implementa nei mock.

### `builder-options-struct` — costruttore a 6 parametri posizionali che fa I/O di rete

- **status**: open
- **priority**: medium
- **impact**: major
- **labels**: api-surface, dx, breaking, audit-p1
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §A2](audit-2026-08-10.md)

`BuildHaivision(url, debug, username, password, header, insecure)` mescola tipi ambigui (due `bool`-ish, due
`string` adiacenti) e fa **2 chiamate HTTP** dentro il costruttore, quindi non è testabile senza rete e non
distingue "config errata" da "gateway giù". Passare a una struct di opzioni (`Config`) con `Insecure bool`
— chiudendo per costruzione il bug `insecure-flag-inverted` — e separare costruzione da `Connect(ctx)`.
Coordinare con `context-and-timeout`: stesso ciclo di refactoring, un solo breaking per i consumer.

### `route-update-delete` — mancano update route, delete route e gestione destinazioni

- **status**: open
- **priority**: medium
- **impact**: major
- **labels**: feature, routes, breaking, audit-p0
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §B8](audit-2026-08-10.md)

L'SDK copre create/list/start/stop e le statistiche; mancano update di una route, delete, start/stop della
singola destinazione e il logout di sessione. Sono aggiunte funzionali, ma passando per
`IHaivisionClient` rompono chi implementa l'interfaccia → `major`. Da fare nello stesso ciclo di
`context-and-timeout` e `builder-options-struct` per non spendere due major.

### `x-net-http2-go-directive` — bump `x/net` alla versione col fix HTTP/2: alza la direttiva `go` a 1.25

- **status**: open
- **priority**: low
- **impact**: major
- **labels**: security, dependencies, breaking
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §D](audit-2026-08-10.md)

Scorporato da `deps-x-net-vuln` (v1.1.0), dove il bump non era fattibile: **`golang.org/x/net` v0.53.0 —
la prima con il fix di GO-2026-4918 — richiede `go 1.25.0`**, quindi alzerebbe la direttiva `go` del modulo
da 1.18 a 1.25 e impedirebbe la compilazione a ogni consumer su un toolchain più vecchio. In v2.0.0 la
release è già breaking, quindi il costo è accettabile.

Priorità `low` perché la vulnerabilità **non è raggiungibile** da questo codice: l'unico pacchetto `x/net`
nel build è `publicsuffix`, mentre la CVE sta in `net/http/internal/http2`. Da fare insieme alla scelta
della versione minima di Go supportata (`ci-go-matrix-and-actions`): è la stessa decisione.

### `exported-naming-typos` — refusi in identificatori esportati

- **status**: open
- **priority**: low
- **impact**: major
- **labels**: cleanup, naming, breaking, audit-p2
- **milestone**: v2.0.0 — Contratto API allineato e superficie pulita
- **ref**: [audit §C](audit-2026-08-10.md)

`ROUTE_COMMMAND` (tre `M`), `ResponseStartOrRoute` (manca "Stop"), e in `udp_rtp/response.go`
`PrompegFeclsBlockAligned` — `ls` invece di `Is`, che cambia anche il campo JSON atteso rispetto a
`prompegFecIsBlockAligned` usato nella request, quindi è **anche** un bug di deserializzazione. Rinominare
identificatori esportati è breaking: da accorpare all'unico major insieme agli altri item di v2.0.0.
