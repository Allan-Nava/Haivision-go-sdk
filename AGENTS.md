# AGENTS.md — Haivision-go-sdk

SDK Go per la **REST API di Haivision Media Gateway / SRT Gateway** (`github.com/Allan-Nava/Haivision-go-sdk`, modulo pubblico su GitHub, licenza MIT). Wrapper HTTP su resty v2 con modelli tipizzati per sessione, device, route (SRT / RTMP / RTSP / UDP-RTP / HLS) e statistiche. Riferimento API: https://doc.haivision.com/HMG3.7.6/rest-api-integrator-s-reference/rest-api-reference (PDF in `docs/`).

Questo file definisce le regole operative per gli agent (Copilot, Claude, altri tool AI) quando lavorano in questo repository.

## Regole di lavoro (SEMPRE)

- **Todo -> `docs/backlog.md`** (sorgente unica, item con `id` stabile, mai TODO sparsi nel codice o nei doc). Ogni item dichiara `impact` (`patch`/`minor`/`major`) e la `milestone` = **versione target** `vX.Y.Z`. Da lì si generano `docs/roadmap.md` (milestone dinamica) e le sezioni di `CHANGELOG.md`. Convenzioni e tabella di assegnazione di `impact`: in testa a `docs/backlog.md`.
- **Le release si tagliano con `make release`**, non a mano. La versione **non si sceglie**: `scripts/new-release.py` prende la prima milestone pendente con 0 item open, rigenera roadmap e tabella del CHANGELOG, esegue i gate, committa e crea il tag annotato. `make release-dry` mostra cosa farebbe senza scrivere. Se la sezione di CHANGELOG manca, la scrive dallo scheletro del backlog e **si ferma** (exit 3): i titoli degli item sono formulati come problemi, la prosa va rifinita in voce da changelog (cosa cambia per chi aggiorna, cosa e breaking), poi si rilancia lo stesso comando.
- **MAI `git push`**: lo fa sempre l'utente, e per il tag e tassativo: fa scattare `.github/workflows/tag-autorelease.yml` che crea una **release pubblica su GitHub**. MAI `Co-Authored-By` nei commit.
- **Chiudendo un item del backlog**: `status: done` (preferito, resta la traccia del perche: se la premessa era sbagliata, scrivilo) + `make roadmap` e committa la roadmap rigenerata, altrimenti il gate `roadmap-check` in CI fallisce. `make check` esegue tutti i gate in un colpo.
- **Ogni modifica sostanziale e un item del backlog**, anche il tooling: il backlog e la sorgente unica, quindi un commit che contiene lavoro non tracciato e un buco. Aggiungi l'item (anche gia `done`) prima di rilasciare.
- **E una libreria pubblica**: ogni modifica alla firma di `IHaivisionClient` o ai campi esportati dei modelli e **breaking** per i consumer (`compress-bot`, servizi HiWay). Va segnalata nel CHANGELOG e nel messaggio di commit, e giustifica un bump `major`.
- **Ogni nuovo endpoint porta con se**: costante in `haivision/constants.go` + helper `fmt.Sprintf` nello stesso file, modelli request/response nel sottopacchetto di protocollo, metodo sul `haivisionSdk`, **voce in `IHaivisionClient`**, test in `test/`, riga nel README/`docs/`. Un metodo che non compare nell'interfaccia e invisibile ai consumer.
- **Documentare SEMPRE** audit, debug e verifiche sul comportamento reale del gateway: doc `.md` in `docs/` (nav Jekyll in `docs/_config.yml`), senza chiederlo. Riportare la **richiesta e la risposta reali** (redatte), non il solo riassunto della doc Haivision.
- **Prima di ogni commit**: `make check` (= `fmt-check` + `vet` + `build` + `test` + `backlog-lint` + `roadmap-check`). La CI Go **non** ha ancora gate su gofmt/vet (item `ci-quality-gates`): fino ad allora il gate e l'agent. Il gate sul backlog invece e gia in CI (`.github/workflows/backlog.yml`).
- **Niente segreti nel repo, nei log e nei test**: username/password/`sessionID` del gateway mai in doc, CHANGELOG, fixture o output di test.
- **Test deterministici e offline**: `go test ./...` non deve mai aprire una connessione verso un gateway. Per coprire i metodi HTTP usare `httptest.Server` + `haivision.Dial(ctx, Config{URL: srv.URL, ...})` puntato su quello, non l'IP di un device. Lo stub riusabile è `gatewayStub` in `test/client_http_test.go`, con variante HTTPS self-signed per i test TLS.
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

- **Aggiungere un metodo a `IHaivisionClient` e `major`, non `minor`.** In Go un'interfaccia esportata e un contratto a due vie: chi la implementa (i mock nei test dei consumer) non compila piu. Per questo le operazioni dipendenti dal protocollo sono **funzioni generiche** (`GetRoutes[TS,TD]`, `CreateRoute[TS,TD]`, ...) e non metodi: i metodi in Go non possono avere type parameter, e metterle nell'interfaccia richiederebbe quattro varianti per ognuna.
- **Autenticazione = cookie di sessione, non token.** `Connect` fa `POST /api/session`, prende `response.sessionID` e lo mette come cookie per tutte le richieste. La sessione **scade** (`ExpireAt` di `GetSessionInfo`) e l'SDK non la rinnova: intercettare `APIError.IsUnauthorized()` e richiamare `Connect`, che e rieseguibile.
- **`New` non fa I/O, `Connect` si.** Non rimettere chiamate di rete nel costruttore: era il problema della v1.x, dove `BuildHaivision` non era testabile senza rete e non distingueva "config sbagliata" da "gateway giu".
- **Lo status HTTP si controlla in `client.go`, in un solo punto** (`check`): nessun metodo deve deserializzare un body d'errore. Ogni nuovo metodo passa da `c.get`/`c.post`/`c.delete` e da `decode`, che avvolge gli errori in `*DecodeError` con l'etichetta della chiamata.
- **`UpdateRoute` sostituisce, non modifica**: le destinazioni omesse vengono rimosse dalla route. Avviare/fermare una singola destinazione non ha un endpoint dedicato: e una update con `Action` sulla destinazione (`StartOrStopDestination`).
- **I tre endpoint di scrittura sono lo stesso path.** `create`, `update` e `delete` sono tutti `POST /api/devices/{id}/updates` e si distinguono per il campo `action`; `start-route`/`stop-route` stanno su `/commands`. Non cercare path REST separati: non esistono.
- **I comandi del gateway sono asincroni**: `state: pending` significa accodato, non eseguito. `ResponseRouteCommand.Pending()` lo dice.
- **I campi opzionali dei modelli di richiesta sono puntatori con `omitempty`**: non impostati vengono **omessi**, non inviati a zero (il gateway interpreterebbe `ttl: 0` come voluto). Helper `Bool`/`Int`/`String`. I tipi seguono gli esempi **letterali** della doc: `ttl`/`tos`/`mtu` numerici, `shaping` booleano : non stringhe.
- **Le statistiche sono tutte `float64`** (tranne `port`/`localPort`): la doc le da come `number` in Mbit/s e con `int` un solo 4.5 faceva fallire l'intera chiamata. Non riconvertirle a interi "perche sono contatori".
- **`validator/v10` valida solo le richieste.** I modelli di risposta non hanno tag `validate`: v10 parsa i tag di tutto l'albero e va in **panic** su una sintassi che non conosce (e successo migrando da `validator.v2`: `nonnil` -> `Undefined validation function`). Se aggiungi un tag, testa `validate.Struct` sulla struct.
- **`debug: true` e sicuro ma verificalo**: `resty.SetDebug` logga i body, e la redazione in `redact.go` e cio che tiene fuori password, `sessionID` e passphrase SRT. Le hook vanno registrate **prima** della prima richiesta. Se aggiungi un campo sensibile ai modelli, aggiungilo a `sensitiveJSONRe`.
- **Estrarre i payload dal PDF della doc funziona**: `docs/REST-API-Integrator-Reference-*.pdf` si legge decomprimendo gli stream con zlib ed estraendo i literal fra parentesi. E cosi che sono stati verificati `elementID`, i tipi di `ttl`/`mtu`/`tos` e il meccanismo start/stop destinazione, invece di indovinarli.
- **Requisito Go = 1.25+** (`x/net` con il fix HTTP/2 e `validator/v10` lo impongono). La matrice CI testa solo il floor: aggiungere `1.26.x` quando esce.

## Puntatori

- **Backlog operativo**: `docs/backlog.md` (sorgente unica) - **Roadmap per milestone di versione**: `docs/roadmap.md` (generata) - `CHANGELOG.md`
- **Tooling**: `scripts/lib/backlog.py` (parser + regole semver, fonte unica), `scripts/backlog-lint.py`, `scripts/generate-roadmap.py` (`--check`, `--release-notes`, `--baseline`), `scripts/new-release.py`. Target: `make check`, `make backlog-lint`, `make roadmap`, `make release-dry`, `make release`. Solo stdlib Python 3.
- Codice client: `haivision/client.go` (`Client`, `New`/`Dial`/`Connect`, interfaccia, helper HTTP), `haivision/config.go` (`Config` + validazione), `haivision/routes.go` (funzioni generiche + comandi), `haivision/auth.go`, `haivision/stats.go`, `haivision/constants.go` (tutti i path API), `haivision/errors.go` (`APIError`, `DecodeError`), `haivision/validate.go`, `haivision/redact.go`, `haivision/header_configurator.go`.
- Modelli per protocollo: `haivision/{srt,rtmp,rtsp,udp_rtp,hls}/`; comuni: `haivision/{session,device,route,stats}/`.
- Test: `test/` (package `test`, esterno alla libreria): oggi 2 test triviali, copertura di `haivision/` **zero**.
- Doc: `docs/` (Jekyll/GitHub Pages, `_config.yml`) + PDF di riferimento API + `docs/audit-2026-08-10.md` (audit tecnico corrente).
- CI: `.github/workflows/go-build.yml`, `go-test.yml`, `tag-autorelease.yml`; dipendenze `.github/dependabot.yml` + `renovate.json`.
- Repo affini HiWay: `devops_hiway` (infra/doc), `compress-bot` (consumer tipico di SDK Go interni).
