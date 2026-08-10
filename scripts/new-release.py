#!/usr/bin/env python3
"""
new-release.py — taglia una release derivando la versione DAL BACKLOG, poi committa e tagga.

La versione non si passa a mano: è la milestone pendente di versione più bassa con **0 item
open** in docs/backlog.md. Chiudi gli item e la release successiva si prepara da sola.

Flusso (in quest'ordine, si ferma al primo problema):
   1.  lint del backlog: struttura + catena semver. Errori ⇒ stop.
   2.  scelta della versione: prima milestone pendente pronta. Blocca se una milestone di
       versione inferiore ha ancora item open (non si salta una release).
   3.  controlli git: repo pulito abbastanza, branch giusto, tag non già esistente.
   4.  rigenera docs/roadmap.md e la tabella delle milestone in CHANGELOG.md.
   5.  CHANGELOG: se manca la sezione della versione, la scrive dallo scheletro generato dal
       backlog e SI FERMA (exit 3) — la prosa va rifinita a mano, poi si rilancia.
   6.  gate: gofmt, go vet, go build, go test, backlog-lint, roadmap --check.
   7.  commit di tutto + tag annotato vX.Y.Z.

NON fa MAI push: il push del tag fa scattare la release pubblica su GitHub ed è dell'utente.

Uso:
  python3 scripts/new-release.py --dry-run     # mostra cosa farebbe, non tocca nulla
  python3 scripts/new-release.py               # prepara, verifica, committa e tagga
  python3 scripts/new-release.py --version v1.2.0    # forza la versione (deve essere pronta)
  python3 scripts/new-release.py --no-commit   # si ferma dopo i gate, senza commit/tag
Requires: solo stdlib + git + go nel PATH.
"""
import argparse
import datetime
import os
import re
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.normpath(os.path.join(HERE, ".."))
sys.path.insert(0, os.path.join(HERE, "lib"))
import backlog as B   # noqa: E402

CHANGELOG = os.path.join(ROOT, "CHANGELOG.md")
BACKLOG = os.path.join(ROOT, "docs", "backlog.md")
DEFAULT_BRANCH = "main"
DATE_PLACEHOLDER = "DATA-DA-INSERIRE"
TABLE_START = "<!-- MILESTONE-TABLE:START -->"
TABLE_END = "<!-- MILESTONE-TABLE:END -->"

# file che non devono finire in un commit di release: chiavi, dump, env, artefatti
FORBIDDEN_RE = re.compile(
    r"(^|/)(\.env(\..*)?|.*\.(pem|key|p12|pfx|dump|sql|tar|tgz|zip)|id_(rsa|ed25519).*)$", re.I)
MAX_FILE_BYTES = 1_000_000


def run(cmd, **kw):
    return subprocess.run(cmd, cwd=ROOT, capture_output=True, text=True, **kw)


def git(*args):
    r = run(["git", *args])
    if r.returncode != 0:
        die(f"git {' '.join(args)} è fallito:\n{r.stderr.strip()}")
    return r.stdout.strip()


def git_status():
    """(righe_porcelain, path_modificati, path_esistenti_da_ispezionare).

    NON usa git() perché quello fa .strip() sull'output: mangerebbe lo spazio iniziale della
    prima riga di `--porcelain` (` M file` → `M file`) e il path uscirebbe troncato di un
    carattere. Gestisce anche i rename (`R  vecchio -> nuovo`) e i path fra virgolette.
    """
    r = run(["git", "status", "--porcelain"])
    if r.returncode != 0:
        die(f"git status è fallito:\n{r.stderr.strip()}")
    lines = [ln for ln in r.stdout.splitlines() if ln.strip()]
    changed, existing = [], []
    for ln in lines:
        if len(ln) < 4:
            continue
        xy, p = ln[:2], ln[3:]
        if " -> " in p:                      # rename: interessa la destinazione
            p = p.split(" -> ", 1)[1]
        p = p.strip().strip('"')
        changed.append(p)
        if "D" not in xy and os.path.exists(os.path.join(ROOT, p)):
            existing.append(p)
    return lines, changed, existing


def die(msg, code=2):
    print(f"\n❌ {msg}", file=sys.stderr)
    sys.exit(code)


def step(msg):
    print(f"▶ {msg}")


# --- 1. lint + 2. scelta della versione -------------------------------------------------

def pick_version(items, baseline, forced):
    chain = B.milestone_chain(items, baseline)
    pending = B.pending_milestones(chain)
    if not pending:
        die("nessuna milestone pendente nel backlog: non c'è niente da rilasciare.\n"
            "   Aggiungi item con `- **milestone**: vX.Y.Z — Titolo` in docs/backlog.md.", 3)

    def n_open(ms):
        return sum(1 for it in ms["items"] if it["status"].lower() != "done")

    if forced:
        pm = B.parse_milestone(forced)
        if pm is None:
            die(f"--version `{forced}` non è una versione `vX.Y.Z`")
        ms = next((m for m in pending if m["version"] == pm[0]), None)
        if ms is None:
            known = ", ".join(B.fmt_version(m["version"]) for m in pending)
            die(f"nessuna milestone pendente {forced} nel backlog (pendenti: {known})")
    else:
        ms = pending[0]

    # non si salta una release: tutte le milestone pendenti più vecchie devono essere chiuse
    for older in pending:
        if older["version"] >= ms["version"]:
            break
        if n_open(older):
            die(f"{B.fmt_version(older['version'])} ha ancora {n_open(older)} item open: "
                f"va rilasciata prima di {B.fmt_version(ms['version'])}.\n"
                f"   Le release seguono la catena, non si salta una versione.", 3)

    if n_open(ms):
        remaining = sorted(it["id"] for it in ms["items"] if it["status"].lower() != "done")
        die(f"{ms['title']} ha {len(remaining)} item ancora open: non è rilasciabile.\n"
            f"   Da chiudere: {', '.join(remaining[:10])}"
            + (" …" if len(remaining) > 10 else "")
            + "\n   Metti `- **status**: done` in docs/backlog.md quando sono fatti.", 3)
    return ms, chain


# --- 3. controlli git ------------------------------------------------------------------

def check_git(version, allow_branch, dry_run):
    if not os.path.isdir(os.path.join(ROOT, ".git")):
        die("non è un repository git")
    branch = git("rev-parse", "--abbrev-ref", "HEAD")
    if branch != DEFAULT_BRANCH and not allow_branch:
        die(f"sei su `{branch}`, non su `{DEFAULT_BRANCH}`: un tag di release creato qui non "
            f"sarebbe su {DEFAULT_BRANCH}.\n   Usa --allow-branch se è voluto.")
    if git("tag", "--list", version):
        die(f"il tag {version} esiste già. Rimuovilo (`git tag -d {version}`) o correggi "
            f"la milestone nel backlog.")
    # file che non devono entrare in una release
    lines, changed, existing = git_status()
    for p in existing:
        if FORBIDDEN_RE.search(p):
            die(f"il working tree contiene `{p}`, che non deve entrare in un commit di "
                f"release (chiave/dump/env). Rimuovilo o aggiungilo a .gitignore.")
        full = os.path.join(ROOT, p)
        if os.path.isfile(full) and os.path.getsize(full) > MAX_FILE_BYTES:
            die(f"`{p}` pesa {os.path.getsize(full) // 1024} KB (> "
                f"{MAX_FILE_BYTES // 1024} KB): non committarlo in una release senza motivo.")
    if not lines and not dry_run:
        die("working tree pulito: non c'è niente da rilasciare.\n"
            "   (Se la release è già committata, crea solo il tag: "
            f"`git tag -a {version} -m \"Release {version[1:]}\"`)", 3)
    return branch, changed


# --- 4. CHANGELOG ---------------------------------------------------------------------

def milestone_table(chain, baseline):
    """Tabella delle milestone ancora da rilasciare, per la sezione [Unreleased]."""
    pending = B.pending_milestones(chain)
    if not pending:
        return ["_Nessuna milestone pendente: tutto il backlog pianificato è rilasciato._"]
    out = ["| Versione | Contenuto | Bump | Stato |", "|---|---|---|---|"]
    for ms in pending:
        n_open = sum(1 for it in ms["items"] if it["status"].lower() != "done")
        state = "✅ pronta" if n_open == 0 else f"{n_open} item open"
        label = ms["label"] or "—"
        out.append(f"| `{B.fmt_version(ms['version'])}` | {label.replace('|', chr(92) + '|')} "
                   f"| {ms['actual'] or '?'} | {state} |")
    return out


def rewrite_milestone_table(text, chain, baseline):
    """Sostituisce il contenuto fra i marker MILESTONE-TABLE. Idempotente."""
    if TABLE_START not in text or TABLE_END not in text:
        return text, False
    head, rest = text.split(TABLE_START, 1)
    _, tail = rest.split(TABLE_END, 1)
    body = "\n".join(milestone_table(chain, baseline))
    return f"{head}{TABLE_START}\n{body}\n{TABLE_END}{tail}", True


def changelog_section_state(text, version):
    """('missing'|'placeholder'|'ok', riga della sezione) per `## [X.Y.Z] — data`."""
    num = version[1:]
    m = re.search(r"^## \[" + re.escape(num) + r"\][^\n]*$", text, re.M)
    if m is None:
        return "missing", None
    return ("placeholder" if DATE_PLACEHOLDER in m.group(0) else "ok"), m.group(0)


def insert_changelog_section(text, version, section):
    """Inserisce la sezione subito prima della prima `## [` che non sia [Unreleased]."""
    m = re.search(r"^## \[(?!Unreleased)", text, re.M)
    at = m.start() if m else len(text)
    return text[:at] + section.rstrip() + "\n\n" + text[at:]


# --- 6. gate --------------------------------------------------------------------------

def gates(version):
    """I gate eseguiti prima del commit. `roadmap --check` usa la versione in uscita come
    baseline, coerentemente con come le pagine sono state generate al passo 4: il tag non
    esiste ancora, quindi senza --baseline il confronto userebbe il tag precedente."""
    return [
        ("gofmt", ["gofmt", "-l", "."], lambda r: r.returncode == 0 and not r.stdout.strip()),
        ("go vet", ["go", "vet", "./..."], lambda r: r.returncode == 0),
        ("go build", ["go", "build", "./..."], lambda r: r.returncode == 0),
        ("go test", ["go", "test", "./...", "-count=1"], lambda r: r.returncode == 0),
        ("backlog-lint", [sys.executable, "scripts/backlog-lint.py", "--baseline", version],
         lambda r: r.returncode == 0),
        ("roadmap --check", [sys.executable, "scripts/generate-roadmap.py",
                             "--check", "--baseline", version], lambda r: r.returncode == 0),
    ]


def run_gates(version):
    for name, cmd, ok in gates(version):
        r = run(cmd)
        if not ok(r):
            out = (r.stdout + r.stderr).strip()
            die(f"gate `{name}` FALLITO:\n{out[:2000]}")
        print(f"   ✓ {name}")


# --- 7. commit + tag ------------------------------------------------------------------

def commit_message(ms, version):
    done = sorted((it for it in ms["items"] if it["status"].lower() == "done"),
                  key=lambda it: it["id"])
    subject = f"release: {version}"
    if ms["label"]:
        subject += f" — {ms['label']}"
    lines = [subject, "", f"{len(done)} item del backlog chiusi (docs/backlog.md):"]
    for it in done:
        title = it["title"]
        if len(title) > 64:
            title = title[:63].rstrip() + "…"
        lines.append(f"- {it['id']} — {title}")
    lines += ["", f"Impatto della release: {ms['required']}. Dettaglio in CHANGELOG.md, "
                  f"piano in docs/roadmap.md."]
    return "\n".join(lines)


def main():
    ap = argparse.ArgumentParser(
        description="Taglia una release derivando la versione dal backlog, poi committa e tagga.")
    ap.add_argument("--version", default="", metavar="vX.Y.Z",
                    help="forza la versione invece di derivarla (deve essere una milestone pronta)")
    ap.add_argument("--date", default="", help="data della release (default: oggi)")
    ap.add_argument("--dry-run", action="store_true", help="mostra cosa farebbe, non scrive nulla")
    ap.add_argument("--no-commit", action="store_true", help="si ferma dopo i gate")
    ap.add_argument("--no-verify", action="store_true", help="salta i gate (sconsigliato)")
    ap.add_argument("--allow-branch", action="store_true",
                    help=f"consenti il tag da un branch diverso da {DEFAULT_BRANCH}")
    args = ap.parse_args()

    date = args.date or datetime.date.today().isoformat()
    if not re.match(r"^\d{4}-\d{2}-\d{2}$", date):
        die(f"--date `{date}` non è nel formato YYYY-MM-DD")

    # 1. lint
    step("lint del backlog")
    items = B.parse_backlog(BACKLOG)
    baseline = B.released_baseline()
    errors, warnings = B.lint(items, baseline=baseline)
    for w in warnings:
        print(f"   ⚠️  {w}")
    if errors:
        for e in errors:
            print(f"   ✗ {e}", file=sys.stderr)
        die("il backlog non è valido: correggilo prima di rilasciare.")
    print(f"   ✓ {len(items)} item · baseline rilasciata {B.fmt_version(baseline)}")

    # 2. versione
    step("scelta della versione dal backlog")
    ms, chain = pick_version(items, baseline, args.version)
    version = B.fmt_version(ms["version"])
    n_done = sum(1 for it in ms["items"] if it["status"].lower() == "done")
    print(f"   ✓ {version} — {ms['label'] or 'senza etichetta'} "
          f"({n_done} item, impatto {ms['required']}, {ms['actual']} bump da "
          f"{B.fmt_version(ms['prev'])})")

    # 3. git
    step("controlli git")
    branch, paths = check_git(version, args.allow_branch, args.dry_run)
    print(f"   ✓ branch {branch}, tag {version} libero, {len(paths)} path da committare")

    # 4. rigenerazione roadmap + tabella CHANGELOG
    #
    # Generate con baseline = la versione IN USCITA, non il max tag attuale: il tag non esiste
    # ancora, ma il commit che stiamo per creare lo porta. Senza questo, roadmap.md e la tabella
    # descriverebbero lo stato PRE-tag (v1.1.0 ancora "prossima release") e il gate
    # `roadmap-check` diventerebbe rosso subito dopo ogni release.
    step("rigenerazione delle pagine generate (baseline = " + version + ")")
    gen = [sys.executable, "scripts/generate-roadmap.py", "--baseline", version]
    if args.dry_run:
        print("   (dry-run) docs/roadmap.md non riscritta")
    else:
        r = run(gen)
        if r.returncode != 0:
            die(f"generate-roadmap.py è fallito:\n{r.stdout}{r.stderr}")
        print(f"   ✓ {r.stdout.strip()}")

    released_chain = B.milestone_chain(items, ms["version"])
    text = open(CHANGELOG, encoding="utf-8").read()
    new_text, had_markers = rewrite_milestone_table(text, released_chain, ms["version"])
    if not had_markers:
        print(f"   ⚠️  marker {TABLE_START} assenti in CHANGELOG.md: tabella non aggiornata")
    elif new_text != text:
        print("   ✓ tabella delle milestone in CHANGELOG.md"
              + (" da aggiornare (dry-run)" if args.dry_run else " aggiornata"))
    text = new_text

    # 5. sezione CHANGELOG della versione
    step(f"sezione CHANGELOG di {version}")
    state, line = changelog_section_state(text, version)
    if state == "ok":
        print(f"   ✓ già presente: {line}")
        if not args.dry_run:
            open(CHANGELOG, "w", encoding="utf-8").write(text)
    else:
        r = run([sys.executable, "scripts/generate-roadmap.py", "--release-notes", version])
        if r.returncode != 0:
            die(f"generazione delle release notes fallita:\n{r.stderr}")
        section = r.stdout.replace(DATE_PLACEHOLDER, date)
        if state == "placeholder":
            # sostituisce solo la data placeholder, conservando la prosa già scritta
            text = text.replace(line, line.replace(DATE_PLACEHOLDER, date))
            if not args.dry_run:
                open(CHANGELOG, "w", encoding="utf-8").write(text)
            print(f"   ✓ data della sezione impostata a {date}")
        else:
            if args.dry_run:
                print(f"   (dry-run) inserirei la sezione [{version[1:]}] — {date} "
                      f"({len(section.splitlines())} righe) e mi fermerei per la rifinitura")
                return 0
            text = insert_changelog_section(text, version, section)
            open(CHANGELOG, "w", encoding="utf-8").write(text)
            print(f"   ✓ inserita la sezione [{version[1:]}] — {date} dallo scheletro del backlog")
            print(f"\n⏸  FERMO QUI: la sezione è generata dai TITOLI degli item, che sono "
                  f"formulati come problemi.\n"
                  f"   Rifinisci la prosa di CHANGELOG.md in voce da changelog (cosa cambia per "
                  f"chi aggiorna,\n   cosa è breaking), poi rilancia lo stesso comando per "
                  f"committare e taggare.")
            return 3

    # 6. gate
    if args.no_verify:
        print("▶ gate SALTATI (--no-verify)")
    else:
        step("gate")
        if args.dry_run:
            print("   (dry-run) gate non eseguiti")
        else:
            run_gates(version)

    msg = commit_message(ms, version)

    # 7. commit + tag
    if args.no_commit or args.dry_run:
        step("commit + tag" + (" (dry-run)" if args.dry_run else " SALTATI (--no-commit)"))
        print("   file che verrebbero committati:")
        for p in paths:
            print(f"     {p}")
        print("\n   messaggio di commit:\n")
        print("   " + msg.replace("\n", "\n   "))
        print(f"\n   poi: git tag -a {version} -m \"Release {version[1:]}\"")
        return 0

    step("commit + tag")
    git("add", "-A")
    r = run(["git", "commit", "-m", msg])
    if r.returncode != 0:
        die(f"il commit è fallito:\n{r.stdout}{r.stderr}")
    sha = git("rev-parse", "--short", "HEAD")
    print(f"   ✓ commit {sha}")
    git("tag", "-a", version, "-m", f"Release {version[1:]}")
    print(f"   ✓ tag annotato {version} su {sha}")

    print(f"\n✅ {version} rilasciata localmente.")
    print(f"   Il push lo fai TU (fa scattare la release pubblica su GitHub):")
    print(f"     git push origin {branch} && git push origin {version}")
    nxt = next((m for m in B.pending_milestones(B.milestone_chain(items, ms["version"]))), None)
    if nxt:
        n_open = sum(1 for it in nxt["items"] if it["status"].lower() != "done")
        print(f"   Prossima milestone: {nxt['title']} ({n_open} item open).")
    return 0


if __name__ == "__main__":
    sys.exit(main())
