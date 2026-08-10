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

_Baseline rilasciata: **v2.0.0** · 3 milestone · 34 item pianificati (0 open · 34 done)._

## Prossima release

Nessuna milestone pendente: **il backlog pianificato è tutto rilasciato**. La prossima release va aperta aggiungendo item con una nuova `- **milestone**: vX.Y.Z — Titolo` in [`backlog.md`](backlog.md).

## v1.1.0 — Correttezza client, sicurezza, wire fix

_🏷️ **rilasciata** (tag v1.1.0) · impatto: **minor** · 14 item_

| id | Titolo | Impatto | Priorità | Status |
|----|--------|---------|----------|--------|
| `http-status-check` | nessun metodo controlla lo status HTTP: 401/500 sembrano successi | ✨ minor | high | ✅ done |
| `insecure-flag-inverted` | `insecure: &false` DISABILITA la verifica TLS | ✨ minor | high | ✅ done |
| `debug-logs-credentials` | con `debug: true` username e password finiscono nei log | 🔧 patch | high | ✅ done |
| `deps-x-net-vuln` | `golang.org/x/net` v0.7.0 obsoleto → v0.35.0 (GO-2026-4918 non raggiungibile) | 🔧 patch | high | ✅ done |
| `device-list-empty-panic` | panic se il gateway risponde con lista device vuota | 🔧 patch | high | ✅ done |
| `route-json-tag-fields-parameters` | body con `"Fields"`/`"Parameters"`: il gateway non li riconosce | 🔧 patch | high | ✅ done |
| `srt-client-stats-path` | statistiche client SRT sul path senza `/client` | 🔧 patch | high | ✅ done |
| `startstop-commands-endpoint` | start/stop route postato su `/updates` invece di `/commands` | 🔧 patch | high | ✅ done |
| `wire-contract-fixture-tests` | test tabellari di ser/deser sui payload della doc Haivision | 🔧 patch | high | ✅ done |
| `header-configurator-order` | header custom applicati dopo il login: rotto dietro proxy autenticato | 🔧 patch | medium | ✅ done |
| `healthcheck-always-nil` | `HealthCheck()` non può fallire | 🔧 patch | medium | ✅ done |
| `makefile-build-and-gofmt` | `make build` è rotto e 5 file non passano gofmt | 🔧 patch | medium | ✅ done |
| `release-tooling-dynamic` | release derivata dal backlog: versione, CHANGELOG, gate, commit e tag | 🔧 patch | medium | ✅ done |
| `unconditional-log-println` | la libreria scrive sul logger globale del consumer a ogni chiamata | 🔧 patch | low | ✅ done |

## v1.2.0 — Qualità, CI, documentazione

_🏷️ **rilasciata** (tag v1.2.0) · impatto: **patch** · 9 item_

| id | Titolo | Impatto | Priorità | Status |
|----|--------|---------|----------|--------|
| `httptest-client-coverage` | copertura del package `haivision`: 0% | 🔧 patch | high | ✅ done |
| `readme-import-path-go-version` | il README documenta un import che non compila | 🔧 patch | high | ✅ done |
| `changelog-bootstrap` | nessun CHANGELOG nonostante 30+ tag e release automatiche | 🔧 patch | medium | ✅ done |
| `ci-go-matrix-and-actions` | matrice Go 1.18–1.21 (tutte EOL) e action obsolete | 🔧 patch | medium | ✅ done |
| `ci-quality-gates` | la CI non ha gate su formato, lint, vulnerabilità, backlog | 🔧 patch | medium | ✅ done |
| `deps-resty-bump` | resty v2.7.0 è del 2022 | 🔧 patch | medium | ✅ done |
| `tag-autorelease-modernize` | release workflow su action archiviata e permessi eccessivi | 🔧 patch | medium | ✅ done |
| `dead-code-and-stubs-cleanup` | blocchi commentati e file stub vuoti | 🔧 patch | low | ✅ done |
| `dependabot-tests-dir` | entry dependabot su una directory che non esiste | 🔧 patch | low | ✅ done |

## v2.0.0 — Contratto API allineato e superficie pulita

_🏷️ **rilasciata** (tag v2.0.0) · impatto: **major** · 11 item_

| id | Titolo | Impatto | Priorità | Status |
|----|--------|---------|----------|--------|
| `create-route-request-model` | `CreateRoute*` invia il modello di risposta, non la richiesta | 💥 major | high | ✅ done |
| `request-field-types-mismatch` | tipi dei campi di richiesta diversi da quelli della doc | 💥 major | high | ✅ done |
| `startstop-response-slice` | `ResponseStartOrRoute` non deserializza la risposta reale | 💥 major | high | ✅ done |
| `stats-float64` | bitrate e rate in Mbit/s tipizzati `int`: ogni valore frazionario rompe la chiamata | 💥 major | high | ✅ done |
| `builder-options-struct` | costruttore a 6 parametri posizionali che fa I/O di rete | 💥 major | medium | ✅ done |
| `context-and-timeout` | nessun `context.Context` e nessun timeout: chiamate non cancellabili | 💥 major | medium | ✅ done |
| `route-update-delete` | mancano update route, delete route e gestione destinazioni | 💥 major | medium | ✅ done |
| `typed-get-routes` | `GetRoutes` restituisce `*resty.Response`: il trasporto è nell'API pubblica | 💥 major | medium | ✅ done |
| `validator-v10-optional-fields` | la validazione rifiuta route legittime e si rompe sui bool | 💥 major | medium | ✅ done |
| `exported-naming-typos` | refusi in identificatori esportati | 💥 major | low | ✅ done |
| `x-net-http2-go-directive` | bump `x/net` alla versione col fix HTTP/2: alza la direttiva `go` a 1.25 | 💥 major | low | ✅ done |

## Non pianificati (senza milestone)

_0 item open senza versione target. Assegnane una con `- **milestone**: vX.Y.Z — Titolo` in [`backlog.md`](backlog.md)._
