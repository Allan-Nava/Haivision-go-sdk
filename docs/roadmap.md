---
layout: default
title: Roadmap & milestone
nav_order: 8
description: "Milestone di versione generate dal backlog"
permalink: /roadmap
---

# Roadmap — milestone di versione dal backlog

<!-- GENERATO da scripts/generate-roadmap.py — NON editare a mano. -->
> ⚙️ Pagina **generata** da [`scripts/generate-roadmap.py`](https://github.com/Allan-Nava/Haivision-go-sdk/blob/main/scripts/generate-roadmap.py) leggendo [`docs/backlog.md`](backlog.md) (unica sorgente). Rigenerala con `make roadmap`; `make roadmap-check` è il gate in CI.

_Baseline rilasciata: **v1.0.0** · 3 milestone · 31 item pianificati (31 open · 0 done)._

## Prossima release

**v1.1.0 — Correttezza client, sicurezza, wire fix** — 13 item da chiudere (0 già fatti). `minor` bump rispetto a **v1.0.0** (minimo imposto dagli item: `minor`).

```
  v1.0.0 (rilasciata)
     │   v1.1.0    [░░░░░░░░░░░░░░░░] 0/13  minor  ◀── PROSSIMA
     │   v1.2.0    [░░░░░░░░░░░░░░░░] 0/9   patch
     │   v2.0.0    [░░░░░░░░░░░░░░░░] 0/9   major
     ▼
```

Blocker `high` di v1.1.0:

| id | Titolo | Impatto | Priorità |
|----|--------|---------|----------|
| `http-status-check` | nessun metodo controlla lo status HTTP: 401/500 sembrano successi | ✨ minor | high |
| `insecure-flag-inverted` | `insecure: &false` DISABILITA la verifica TLS | ✨ minor | high |
| `debug-logs-credentials` | con `debug: true` username e password finiscono nei log | 🔧 patch | high |
| `deps-x-net-vuln` | `golang.org/x/net v0.7.0`: vulnerabilità raggiungibile dal codice | 🔧 patch | high |
| `device-list-empty-panic` | panic se il gateway risponde con lista device vuota | 🔧 patch | high |
| `route-json-tag-fields-parameters` | body con `"Fields"`/`"Parameters"`: il gateway non li riconosce | 🔧 patch | high |
| `srt-client-stats-path` | statistiche client SRT sul path senza `/client` | 🔧 patch | high |
| `startstop-commands-endpoint` | start/stop route postato su `/updates` invece di `/commands` | 🔧 patch | high |
| `wire-contract-fixture-tests` | test tabellari di ser/deser sui payload della doc Haivision | 🔧 patch | high |

## v1.1.0 — Correttezza client, sicurezza, wire fix

_minor bump da v1.0.0 · impatto richiesto dagli item: **minor** · 🟢 13 open · 0 done_

| id | Titolo | Impatto | Priorità | Status |
|----|--------|---------|----------|--------|
| `http-status-check` | nessun metodo controlla lo status HTTP: 401/500 sembrano successi | ✨ minor | high | 🟢 open |
| `insecure-flag-inverted` | `insecure: &false` DISABILITA la verifica TLS | ✨ minor | high | 🟢 open |
| `debug-logs-credentials` | con `debug: true` username e password finiscono nei log | 🔧 patch | high | 🟢 open |
| `deps-x-net-vuln` | `golang.org/x/net v0.7.0`: vulnerabilità raggiungibile dal codice | 🔧 patch | high | 🟢 open |
| `device-list-empty-panic` | panic se il gateway risponde con lista device vuota | 🔧 patch | high | 🟢 open |
| `route-json-tag-fields-parameters` | body con `"Fields"`/`"Parameters"`: il gateway non li riconosce | 🔧 patch | high | 🟢 open |
| `srt-client-stats-path` | statistiche client SRT sul path senza `/client` | 🔧 patch | high | 🟢 open |
| `startstop-commands-endpoint` | start/stop route postato su `/updates` invece di `/commands` | 🔧 patch | high | 🟢 open |
| `wire-contract-fixture-tests` | test tabellari di ser/deser sui payload della doc Haivision | 🔧 patch | high | 🟢 open |
| `header-configurator-order` | header custom applicati dopo il login: rotto dietro proxy autenticato | 🔧 patch | medium | 🟢 open |
| `healthcheck-always-nil` | `HealthCheck()` non può fallire | 🔧 patch | medium | 🟢 open |
| `makefile-build-and-gofmt` | `make build` è rotto e 5 file non passano gofmt | 🔧 patch | medium | 🟢 open |
| `unconditional-log-println` | la libreria scrive sul logger globale del consumer a ogni chiamata | 🔧 patch | low | 🟢 open |

## v1.2.0 — Qualità, CI, documentazione

_minor bump da v1.1.0 · impatto richiesto dagli item: **patch** · 🟢 9 open · 0 done_

| id | Titolo | Impatto | Priorità | Status |
|----|--------|---------|----------|--------|
| `httptest-client-coverage` | copertura del package `haivision`: 0% | 🔧 patch | high | 🟢 open |
| `readme-import-path-go-version` | il README documenta un import che non compila | 🔧 patch | high | 🟢 open |
| `changelog-bootstrap` | nessun CHANGELOG nonostante 30+ tag e release automatiche | 🔧 patch | medium | 🟢 open |
| `ci-go-matrix-and-actions` | matrice Go 1.18–1.21 (tutte EOL) e action obsolete | 🔧 patch | medium | 🟢 open |
| `ci-quality-gates` | la CI non ha gate su formato, lint, vulnerabilità, backlog | 🔧 patch | medium | 🟢 open |
| `deps-resty-bump` | resty v2.7.0 è del 2022 | 🔧 patch | medium | 🟢 open |
| `tag-autorelease-modernize` | release workflow su action archiviata e permessi eccessivi | 🔧 patch | medium | 🟢 open |
| `dead-code-and-stubs-cleanup` | blocchi commentati e file stub vuoti | 🔧 patch | low | 🟢 open |
| `dependabot-tests-dir` | entry dependabot su una directory che non esiste | 🔧 patch | low | 🟢 open |

## v2.0.0 — Contratto API allineato e superficie pulita

_major bump da v1.2.0 · impatto richiesto dagli item: **major** · 🟢 9 open · 0 done_

| id | Titolo | Impatto | Priorità | Status |
|----|--------|---------|----------|--------|
| `create-route-request-model` | `CreateRoute*` invia il modello di risposta, non la richiesta | 💥 major | high | 🟢 open |
| `startstop-response-slice` | `ResponseStartOrRoute` non deserializza la risposta reale | 💥 major | high | 🟢 open |
| `stats-float64` | bitrate e rate in Mbit/s tipizzati `int`: ogni valore frazionario rompe la chiamata | 💥 major | high | 🟢 open |
| `builder-options-struct` | costruttore a 6 parametri posizionali che fa I/O di rete | 💥 major | medium | 🟢 open |
| `context-and-timeout` | nessun `context.Context` e nessun timeout: chiamate non cancellabili | 💥 major | medium | 🟢 open |
| `route-update-delete` | mancano update route, delete route e gestione destinazioni | 💥 major | medium | 🟢 open |
| `typed-get-routes` | `GetRoutes` restituisce `*resty.Response`: il trasporto è nell'API pubblica | 💥 major | medium | 🟢 open |
| `validator-v10-optional-fields` | la validazione rifiuta route legittime e si rompe sui bool | 💥 major | medium | 🟢 open |
| `exported-naming-typos` | refusi in identificatori esportati | 💥 major | low | 🟢 open |

## Non pianificati (senza milestone)

_0 item open senza versione target. Assegnane una con `- **milestone**: vX.Y.Z — Titolo` in [`backlog.md`](backlog.md)._
