#!/usr/bin/env python3
"""
changelog-extract.py — stampa la sezione di CHANGELOG.md di una versione.

Serve al workflow di release: il corpo della GitHub Release è la sezione **rifinita a mano** del
CHANGELOG, non lo scheletro generato dal backlog (`generate-roadmap.py --release-notes`, che si usa
per *creare* la sezione, non per pubblicarla).

Uso:
  python3 scripts/changelog-extract.py v1.2.0        # sezione di quella versione
  python3 scripts/changelog-extract.py --latest      # prima sezione [X.Y.Z] del file
Exit code: 0 trovata · 1 non trovata · 2 errore d'uso.
Requires: solo stdlib.
"""
import argparse
import os
import re
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
CHANGELOG = os.path.normpath(os.path.join(HERE, "..", "CHANGELOG.md"))
VERSION_HEADING = re.compile(r"^## \[(\d+\.\d+\.\d+)\]")


def extract(text, version=None):
    """(versione, corpo) della sezione richiesta, o (None, None). `version` senza la `v`."""
    lines = text.splitlines()
    starts = [(i, m.group(1)) for i, ln in enumerate(lines)
              if (m := VERSION_HEADING.match(ln))]
    if not starts:
        return None, None
    if version is None:
        idx, found = starts[0]
    else:
        hit = next(((i, v) for i, v in starts if v == version), None)
        if hit is None:
            return None, None
        idx, found = hit
    # il corpo va fino al prossimo heading di livello 2 (versione o [Unreleased])
    end = next((i for i in range(idx + 1, len(lines)) if lines[i].startswith("## ")), len(lines))
    body = "\n".join(lines[idx + 1:end]).strip()
    return found, body


def main():
    ap = argparse.ArgumentParser(description="Estrae una sezione da CHANGELOG.md.")
    ap.add_argument("version", nargs="?", default="", metavar="vX.Y.Z")
    ap.add_argument("--latest", action="store_true", help="la prima sezione di versione del file")
    ap.add_argument("--file", default=CHANGELOG)
    args = ap.parse_args()

    if not args.version and not args.latest:
        ap.error("indica una versione vX.Y.Z oppure --latest")
    if not os.path.isfile(args.file):
        print(f"ERRORE: {args.file} non trovato", file=sys.stderr)
        return 2

    want = None
    if args.version:
        m = re.match(r"^v?(\d+\.\d+\.\d+)$", args.version)
        if not m:
            print(f"ERRORE: `{args.version}` non è una versione vX.Y.Z", file=sys.stderr)
            return 2
        want = m.group(1)

    text = open(args.file, encoding="utf-8").read()
    found, body = extract(text, want)
    if found is None:
        target = f"[{want}]" if want else "nessuna sezione di versione"
        print(f"ERRORE: {target} non trovata in {os.path.relpath(args.file)}", file=sys.stderr)
        return 1
    sys.stdout.write(body + "\n")
    return 0


if __name__ == "__main__":
    sys.exit(main())
