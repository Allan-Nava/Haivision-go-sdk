# AGENTS.md — Haivision-go-sdk

SDK Go per la **REST API di Haivision Media Gateway / SRT Gateway** (`github.com/Allan-Nava/Haivision-go-sdk`, modulo pubblico su GitHub, licenza MIT). Wrapper HTTP su resty v2 con modelli tipizzati per sessione, device, route (SRT / RTMP / RTSP / UDP-RTP / HLS) e statistiche. Riferimento API: https://doc.haivision.com/HMG3.7.6/rest-api-integrator-s-reference/rest-api-reference (PDF in `docs/`).

Questo file definisce le regole operative per gli agent (Copilot, Claude, altri tool AI) quando lavorano in questo repository.

## Regole di lavoro (SEMPRE)

- **Todo -> `docs/backlog.md`** (sorgente unica, item con `id` stabile, mai TODO sparsi nel codice o nei doc). Ogni item dichiara `impact` (`patch`/`minor`/`major`) e la `milestone` = **versione target** `vX.Y.Z`. Da lì si generano `docs/roadmap.md` (milestone dinamica) e le sezioni di `CHANGELOG.md`. Convenzioni e tabella di assegnazione di `impact`: in testa a `docs/backlog.md`.
- **Le release si tagliano con `make release`**, non a mano. La versione **non si sceglie**: `scripts/new-release.py` prende la prima milestone pendente con 0 item open, rigenera roadmap e tabella del CHANGELOG, esegue i gate, committa e crea il tag annotato. `make release-dry` mostra cosa farebbe senza scrivere. Se la sezione di CHANGELOG manca, la scrive dallo scheletro del backlog e **si ferma** (exit 3): i titoli degli item sono formulati come problemi, la prosa va rifinita in voce da changelog (cosa cambia per chi aggiorna, cosa e breaking), poi si rilancia lo stesso comando.
- **MAI `git push`**: lo fa sempre l'utente, e per il tag e tassativo: fa scattare `.github/workflows/tag-autorelease.yml` che crea una **release pubblica su GitHub**. MAI `Co-Authored-By` nei commit.
- **Chiudendo un item del backlog**: `status: done` (preferito, resta la traccia del perche: se la premessa era sbagliata, scrivilo) + `make roadmap` e committa la roadmap rigenerata, altrimenti il gate `roadmap-check` in CI fallisce. `make check` esegue tutti i gate in un colpo.
- **Ogni modifica sostanziale e un item del backlog**, anche il tooling: il backlog e la sorgente unica, quindi un commit che contiene lavoro non tracciato e un buco. Aggiungi l'item (anche gia `done`) prima di rilasciare.
- **E una libreria pubblica**: ogni modifica alla firma di `IHaivisionClient` o ai campi esportati dei modelli e **breaking** per i consumer (`compress-bot`, servizi HiWay). Va segnalata nel CHANGELOG e nel messaggio di commit, e giustifica un bump `minor` (o `major` su `v1`).
- **Ogni nuovo endpoint porta con se**: costante in `haivision/constants.go` + helper `fmt.Sprintf` nello stesso file, modelli request/response nel sottopacchetto di protocollo, metodo sul `haivisionSdk`, **voce in `IHaivisionClient`**, test in `test/`, riga nel README/`docs/`. Un metodo che non compare nell'interfaccia e invisibile ai consumer.
- **Documentare SEMPRE** audit, debug e verifiche sul comportamento reale del gateway: doc `.md` in `docs/` (nav Jekyll in `docs/_config.yml`), senza chiederlo. Riportare la **richiesta e la risposta reali** (redatte), non il solo riassunto della doc Haivision.
- **Prima di ogni commit**: `make check` (= `fmt-check` + `vet` + `build` + `test` + `backlog-lint` + `roadmap-check`). La CI Go **non** ha ancora gate su gofmt/vet (item `ci-quality-gates`): fino ad allora il gate e l'agent. Il gate sul backlog invece e gia in CI (`.github/workflows/backlog.yml`).
- **Niente segreti nel repo, nei log e nei test**: username/password/`sessionID` del gateway mai in doc, CHANGELOG, fixture o output di test.
- **Test deterministici e offline**: `go test ./...` non deve mai aprire una connessione verso un gateway. Per coprire i metodi HTTP usare `httptest.Server` + `BuildHaivision` puntato su quello, non l'IP di un device.
- **Un test che logga l'errore invece di fallire non e un test**: usare `t.Fatalf`/`t.Errorf`. `log.Println(err)` in un test e un bug (vedi `test/auth_test.go`).

## Pattern per aggiungere/modificare un endpoint

1. **Leggere la doc Haivision reale** (`docs/REST-API-Integrator-Reference-*.pdf` o il sito HMG 3.7.6) e incollare request+response letterali nel commento sopra il metodo: e la convenzione gia usata in `auth.go`/`route.go`/`stats.go`.
2. **Costante + helper** in `haivision/constants.go`. Path con placeholder -> `const` con `%s` + `var` funzione `fmt.Sprintf`. Non costruire path inline nei metodi.
3. **Modelli** nel sottopacchetto giusto: `session/`, `device/`, `route/`, `srt/`, `rtmp/`, `rtsp/`, `udp_rtp/`, `hls/`, `stats/`. Struct request e response **separate**: la risposta del gateway ha campi di stato (`state`, `summaryStatusCode`, `elapsedTime`) che non vanno mai inviati in POST.
4. **Verificare la serializzazione**: `json.Marshal` del body e confronto **campo per campo** con la doc. Le struct anonime annidate senza tag JSON serializzano col nome Go (`"Fields"` invece di `"fields"`): errore gia presente in `route/model.go`.
5. **Verificare la deserializzazione** con la risposta letterale della doc come fixture: array top-level != oggetto wrappato, `number` frazionari != `int`, `null` -> puntatore.
6. **Chiusura**: metodo nell'interfaccia, test con fixture, README/`docs/`, CHANGELOG, tag (senza push).

## Backlog & versionamento (milestone dinamica)

```
  docs/backlog.md ──▶ scripts/lib/backlog.py ──▶ backlog-lint.py    (gate CI: struttura + semver)
   ### `id` — Titolo    (regole condivise)   ──▶ generate-roadmap.py ─▶ docs/roadmap.md
   - impact / milestone         │                                        (generata, committata)
                                │
   git tag vX.Y.Z ──baseline────┘            new-release.py  ──▶ versione = 1a milestone pendente
   (max tag esistente)                       (make release)        con 0 item open
                                                   │                     │
                                                   │                     ├─▶ CHANGELOG.md
                                                   │                     ├─▶ gate
                                                   └─────────────────────┴─▶ commit + tag (NO push)
```

- **La milestone e dinamica**: "prossima release" = milestone pendente di versione piu bassa con almeno un item `open`. Non esiste una lista di versioni mantenuta a mano: si chiudono item e avanza da sola, roadmap e tabella del CHANGELOG comprese.
- **`impact` governa la versione**, `priority` governa l'ordine di lavoro. Sono ortogonali: un item `low`/`major` esiste (es. `exported-naming-typos`).
- **Il linter blocca un `major` pianificato dentro una minor** e verifica che la catena `baseline -> v1.2.0 -> v2.0.0` sia composta di bump semver validi (componenti inferiori azzerate). Pianificare male una versione e un errore di CI, non una scoperta post-tag.
- **Le milestone con versione <= baseline sono storia, non pianificazione**: `milestone_chain` le marca `released` e il linter non ne valida il bump. Senza questo, subito dopo il tag di `vX.Y.Z` quella milestone risulterebbe "non un incremento valido rispetto a se stessa" e il gate andrebbe rosso a ogni release. Per la stessa ragione `new-release.py` genera roadmap e tabella con **baseline = la versione in uscita** (il tag non esiste ancora, ma il commit lo portera).
- **Non si salta una release**: `new-release.py` rifiuta di taggare `v2.0.0` se `v1.2.0` ha ancora item open.
- **Accorpare i breaking**: ogni major costa un adeguamento a tutti i consumer. Se un intervento e `major`, va in `v2.0.0` insieme agli altri, non in una major propria.
- `make roadmap` dopo ogni modifica al backlog, e committa: `roadmap-check` confronta byte per byte.

## Trappole note / regole tecniche

- **Aggiungere un metodo a `IHaivisionClient` e `major`, non `minor`.** In Go un'interfaccia esportata e un contratto a due vie: chi la implementa (i mock nei test dei consumer) non compila piu. Vale anche per le feature "solo aggiuntive": e la ragione per cui `route-update-delete` sta in `v2.0.0`.
- **Autenticazione = cookie di sessione, non token.** `BuildHaivision` fa `POST /api/session`, prende `response.sessionID` e lo mette come cookie `sessionID` sul client resty per tutte le richieste. Conseguenze: il costruttore fa **2 chiamate HTTP** (session + `/api/devices`), quindi non e un costruttore puro; e la sessione **scade** (`expireAt` in `GetSessionInfo`) senza che l'SDK la rinnovi. Client long-running vanno ricostruiti o serve un refresh.
- **Nessun metodo controlla lo status HTTP.** `restyGet`/`restyPost` ritornano `err == nil` anche su 401/404/500, poi si fa `json.Unmarshal` sul body d'errore. Effetto pratico: credenziali sbagliate -> `InitSession` ritorna un oggetto con `SessionID` vuoto e **nessun errore**, e il client parte con un cookie vuoto. Qualsiasi nuovo metodo deve controllare `resp.IsError()`/`resp.StatusCode()` prima di deserializzare.
- **`insecure *bool` e a tre stati ma ne usa due**: `builder.go` verifica solo `insecure != nil`, quindi passare un puntatore a `false` **abilita** `InsecureSkipVerify`. Per disabilitare il TLS-verify skip si deve passare `nil`. Non replicare il pattern: se lo si tocca, passare a `bool` o controllare `*insecure`.
- **`HealthCheck()` ritorna sempre `nil`** (`return nil` anche nel ramo d'errore): non usarlo come liveness probe finche non e corretto.
- **`BuildHaivision` indicizza `(*deviceResponse)[0]`**: se il gateway risponde `[]` (nessun device, o body d'errore) e panic da index out of range, non errore. `DeviceID`/`HType` vengono dal **primo** device: setup multi-device non e supportato, passare il `deviceId` esplicito ai metodi.
- **`debug: true` stampa le credenziali.** `resty.SetDebug(true)` logga i body delle richieste, incluso `POST /api/session` con username/password in chiaro, e `debugPrint` logga la risposta con il `sessionID`. Mai `debug: true` in produzione ne incollare quei log in doc/issue.
- **`log.Println` incondizionato** in `route.go`/`stats.go`: la libreria scrive sul logger globale del consumer a ogni chiamata, ignorando il flag `debug`. Non aggiungerne altri: usare `o.debugPrint`.
- **Niente `context.Context` e niente timeout** sul client resty: una chiamata puo bloccarsi indefinitamente e non e cancellabile. Se si estende l'API, aggiungere varianti `...WithContext` invece di cambiare le firme esistenti (breaking).
- **Generics `route.RouteModel[TS, TD]`**: i type-set `RequestSource`/`RequestDestination`/`ResponseSource`/`ResponseDestination` in `route/model.go` sono union di struct concrete. Aggiungere un protocollo = estendere il type-set **e** aggiungere il metodo `CreateRoute<Proto>` (un metodo per protocollo, non esiste dispatch dinamico). `hls` e solo destinazione, `rtsp` non ha response model.
- **`validator.v2` e severo e semi-abbandonato**: `validate:"nonnil,min=1"` su **tutti** i campi source/destination significa che una `CreateRoute*` rifiuta qualsiasi route con un campo opzionale vuoto (`ttl`, `tos`, `retainHeader` sono `string` obbligatorie). `min=1` su un `bool` produce `unsupported type` -> la validazione fallisce sempre. Prima di aggiungere tag, testare `validator.Validate` sulla struct.
- **`go build .` non funziona** (nessun file Go nella root del modulo, `make build` fallisce): il target giusto e `go build ./...`. Il package importabile e `.../haivision`, non la root: il README sbaglia.
- **Requisito Go reale = 1.18+** (generics), non 1.13 come dice il README. Matrice CI 1.18-1.21: tutte EOL.
- **`golang.org/x/net v0.7.0` e vulnerabile** (GO-2026-4918, fix in v0.53.0) e `resty v2.7.0` e del 2022. Dependabot/Renovate sono configurati ma non hanno prodotto bump: verificare che siano attivi prima di dire "aggiornato". `dependabot.yml` punta a `/tests` che non esiste (la dir e `test/`).
- **CI: `cache-dependency-path: subdir/go.sum`** e residuo di template (path inesistente); inerte solo perche `cache:` e commentato. `tag-autorelease.yml` usa `actions/create-release@v1` (archiviato) con `permissions: write-all` e installa ffmpeg senza motivo.

## Puntatori

- **Backlog operativo**: `docs/backlog.md` (sorgente unica) - **Roadmap per milestone di versione**: `docs/roadmap.md` (generata) - `CHANGELOG.md`
- **Tooling**: `scripts/lib/backlog.py` (parser + regole semver, fonte unica), `scripts/backlog-lint.py`, `scripts/generate-roadmap.py` (`--check`, `--release-notes`, `--baseline`), `scripts/new-release.py`. Target: `make check`, `make backlog-lint`, `make roadmap`, `make release-dry`, `make release`. Solo stdlib Python 3.
- Codice client: `haivision/haivision.go` (struct + interfaccia `IHaivisionClient` + helper resty), `haivision/builder.go` (costruttore/login), `haivision/auth.go`, `haivision/route.go`, `haivision/stats.go`, `haivision/constants.go` (tutti i path API), `haivision/header_configurator.go` (header custom + Basic auth).
- Modelli per protocollo: `haivision/{srt,rtmp,rtsp,udp_rtp,hls}/`; comuni: `haivision/{session,device,route,stats}/`.
- Test: `test/` (package `test`, esterno alla libreria): oggi 2 test triviali, copertura di `haivision/` **zero**.
- Doc: `docs/` (Jekyll/GitHub Pages, `_config.yml`) + PDF di riferimento API + `docs/audit-2026-08-10.md` (audit tecnico corrente).
- CI: `.github/workflows/go-build.yml`, `go-test.yml`, `tag-autorelease.yml`; dipendenze `.github/dependabot.yml` + `renovate.json`.
- Repo affini HiWay: `devops_hiway` (infra/doc), `compress-bot` (consumer tipico di SDK Go interni).
