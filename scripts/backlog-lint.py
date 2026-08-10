#!/usr/bin/env python3
"""
backlog-lint.py — valida docs/backlog.md (lint standalone per CI/PR, nessuna rete).

Verifica id univoci, vocabolario di status/priority/impact, chiavi meta note, titoli
milestone semver coerenti e — la regola che conta — che **la catena delle milestone regga
il versionamento**: una milestone che pianifica un item `impact: major` non può essere una
minor. Regole centralizzate in `lib/backlog.py` (fonte unica, usata anche dal generatore
della roadmap).

Exit code:
  0 → nessun errore (i warning sono stampati ma non bloccanti; con --strict falliscono)
  1 → almeno un errore (o warning in --strict)
  2 → errore d'uso (file mancante, ecc.)

Requires: solo stdlib.
"""
import argparse
import os
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, os.path.join(HERE, "lib"))
import backlog as B   # noqa: E402


def main():
    ap = argparse.ArgumentParser(description="Lint di docs/backlog.md (standalone, offline).")
    ap.add_argument("--backlog", default=os.path.join(HERE, "..", "docs", "backlog.md"))
    ap.add_argument("--baseline", default="",
                    help="versione già rilasciata da cui parte la catena (default: max tag git `vX.Y.Z`)")
    ap.add_argument("--strict", action="store_true", help="tratta i warning come errori (exit 1)")
    args = ap.parse_args()

    path = os.path.normpath(args.backlog)
    if not os.path.isfile(path):
        print(f"ERRORE: backlog non trovato: {path}", file=sys.stderr)
        return 2

    baseline = None
    if args.baseline:
        pm = B.parse_milestone(args.baseline)
        if pm is None:
            print(f"ERRORE: --baseline `{args.baseline}` non è una versione `vX.Y.Z`", file=sys.stderr)
            return 2
        baseline = pm[0]
    if baseline is None:
        baseline = B.released_baseline()

    items = B.parse_backlog(path)
    errors, warnings = B.lint(items, baseline=baseline)

    n_open = sum(1 for it in items if it["status"].lower() != "done")
    chain = B.milestone_chain(items, baseline)
    print(f"📋 backlog-lint: {len(items)} item ({n_open} open) · {len(chain)} milestone · "
          f"baseline rilasciata {B.fmt_version(baseline)}")
    for ms in chain:
        n_o = sum(1 for it in ms["items"] if it["status"].lower() != "done")
        actual = ms["actual"] or "NON VALIDO"
        print(f"   {B.fmt_version(ms['version'])}: {actual} bump da {B.fmt_version(ms['prev'])} "
              f"· richiesto {ms['required']} · {n_o}/{len(ms['items'])} open")

    for w in warnings:
        print(f"  ⚠️  {w}")
    for e in errors:
        print(f"  ✗ {e}")

    fail = bool(errors) or (args.strict and bool(warnings))
    if fail:
        n = len(errors) + (len(warnings) if args.strict else 0)
        print(f"\n❌ backlog-lint: {n} problema/i bloccante/i.")
        return 1
    print("\n✅ backlog-lint OK"
          + (f" ({len(warnings)} warning non bloccanti)" if warnings else ""))
    return 0


if __name__ == "__main__":
    sys.exit(main())
