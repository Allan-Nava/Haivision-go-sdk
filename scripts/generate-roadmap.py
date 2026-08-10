#!/usr/bin/env python3
"""
generate-roadmap.py — genera docs/roadmap.md: le milestone di versione derivate dal backlog.

Legge docs/backlog.md via lib/backlog.py (fonte unica di parsing) e produce la **milestone
dinamica**: una sezione per versione target, ordinata per semver, con il bump richiesto dagli
item che contiene, i conteggi open/done e la tabella degli item. In testa calcola la
**prossima release** = la milestone di versione più bassa che ha ancora item open, con i suoi
blocker. Chiudi gli item e la prossima release avanza da sola: non c'è nessuna lista di
versioni mantenuta a mano.

La pagina è GENERATA e va committata: così il gate `--check` in CI verifica che nessuno abbia
modificato il backlog senza rigenerarla.

Uso:
  python3 scripts/generate-roadmap.py                      # scrive docs/roadmap.md
  python3 scripts/generate-roadmap.py --check              # exit 1 se l'output diverge dal file
  python3 scripts/generate-roadmap.py --release-notes v1.1.0   # sezione CHANGELOG per quella versione
Requires: solo stdlib.
"""
import argparse
import os
import sys

HERE = os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0, os.path.join(HERE, "lib"))
import backlog as B   # noqa: E402

OUT = os.path.normpath(os.path.join(HERE, "..", "docs", "roadmap.md"))
BACKLOG = os.path.normpath(os.path.join(HERE, "..", "docs", "backlog.md"))
PRIO_ORDER = {"high": 0, "medium": 1, "low": 2, "": 3}
STATUS_GLYPH = {"done": "✅ done", "open": "🟢 open"}
IMPACT_GLYPH = {"major": "💥 major", "minor": "✨ minor", "patch": "🔧 patch"}
BAR_W = 16   # larghezza fissa della barra di avanzamento nel diagramma della prossima release
# le sezioni del CHANGELOG (Keep a Changelog) in cui finisce un item, per impatto.
# `minor` → "Changed": in semver copre sia le aggiunte sia i cambi di comportamento, e
# l'impatto da solo non distingue i due casi → spostare a mano in "Added" ciò che è nuovo.
CHANGELOG_SECTION = {"major": "Changed (breaking)", "minor": "Changed", "patch": "Fixed"}


def _md_escape(s):
    return s.replace("|", "\\|")


def _sort_key(it):
    # open prima di done; poi impatto (major→patch); poi priorità (high→low); poi id
    return (it["status"].lower() == "done",
            -B.IMPACT_RANK.get(B.item_impact(it), 0),
            PRIO_ORDER.get(it["priority"].lower(), 3),
            it["id"])


def _counts(its):
    n_open = sum(1 for it in its if it["status"].lower() != "done")
    return n_open, len(its) - n_open


def _item_table(its, with_impact=True):
    L = []
    if with_impact:
        L.append("| id | Titolo | Impatto | Priorità | Status |")
        L.append("|----|--------|---------|----------|--------|")
    else:
        L.append("| id | Titolo | Impatto | Priorità |")
        L.append("|----|--------|---------|----------|")
    for it in its:
        imp = IMPACT_GLYPH.get(B.item_impact(it), B.item_impact(it))
        prio = it["priority"] or "—"
        row = f"| `{it['id']}` | {_md_escape(it['title'])} | {imp} | {prio} |"
        if with_impact:
            row += f" {STATUS_GLYPH.get(it['status'].lower(), it['status'])} |"
        L.append(row)
    return L


def render(items, baseline):
    chain = B.milestone_chain(items, baseline)
    planned = [it for it in items if it["milestone"]]
    tot_open, tot_done = _counts(planned)

    L = []
    L.append("---")
    L.append("layout: default")
    L.append("title: Roadmap & milestone")
    L.append("nav_order: 8")
    L.append('description: "Milestone di versione generate dal backlog"')
    L.append("permalink: /roadmap")
    L.append("---")
    L.append("")
    L.append("# Roadmap — milestone di versione dal backlog")
    L.append("")
    L.append("<!-- GENERATO da scripts/generate-roadmap.py — NON editare a mano. -->")
    L.append("> ⚙️ Pagina **generata** da [`scripts/generate-roadmap.py`](https://github.com/Allan-Nava/"
             "Haivision-go-sdk/blob/main/scripts/generate-roadmap.py) leggendo "
             "[`docs/backlog.md`](backlog.md) (unica sorgente). Rigenerala con `make roadmap`; "
             "`make roadmap-check` è il gate in CI.")
    L.append("")
    L.append(f"_Baseline rilasciata: **{B.fmt_version(baseline)}** · {len(chain)} milestone · "
             f"{len(planned)} item pianificati ({tot_open} open · {tot_done} done)._")
    L.append("")

    # --- prossima release (il pezzo "dinamico") ---
    nxt = next((ms for ms in chain if _counts(ms["items"])[0] > 0), None)
    L.append("## Prossima release")
    L.append("")
    if nxt is None:
        L.append("Nessuna milestone con item open: **il backlog pianificato è chiuso**. "
                 "La prossima release va aperta aggiungendo item con una nuova "
                 "`- **milestone**: vX.Y.Z — Titolo` in [`backlog.md`](backlog.md).")
        L.append("")
    else:
        n_open, n_done = _counts(nxt["items"])
        v = B.fmt_version(nxt["version"])
        # `actual` = bump che la versione rappresenta davvero; `required` = minimo imposto dagli
        # item. actual > required è legittimo (si può rilasciare una minor di soli patch);
        # actual < required è l'errore che blocca il linter.
        L.append(f"**{nxt['title']}** — {n_open} item da chiudere ({n_done} già fatti). "
                 f"`{nxt['actual'] or 'incremento NON valido'}` bump rispetto a "
                 f"**{B.fmt_version(nxt['prev'])}** (minimo imposto dagli item: `{nxt['required']}`).")
        L.append("")
        L.append("```")
        L.append(f"  {B.fmt_version(baseline)} (rilasciata)")
        for ms in chain:
            o, d = _counts(ms["items"])
            tot = d + o
            # barra a larghezza fissa: le milestone restano incolonnate anche con conteggi diversi
            filled = round(d / tot * BAR_W) if tot else BAR_W
            bar = "█" * filled + "░" * (BAR_W - filled)
            mark = "  ◀── PROSSIMA" if ms is nxt else ""
            L.append(f"     │   {B.fmt_version(ms['version']):<9} [{bar}] {d}/{tot:<3} "
                     f"{ms['required']:<5}{mark}".rstrip())
        L.append("     ▼")
        L.append("```")
        L.append("")
        blockers = sorted((it for it in nxt["items"]
                           if it["status"].lower() != "done"
                           and it["priority"].lower() == "high"), key=_sort_key)
        if blockers:
            L.append(f"Blocker `high` di {v}:")
            L.append("")
            L += _item_table(blockers, with_impact=False)
            L.append("")

    # --- una sezione per milestone ---
    for ms in chain:
        its = sorted(ms["items"], key=_sort_key)
        n_open, n_done = _counts(its)
        L.append(f"## {ms['title']}")
        L.append("")
        actual = ms["actual"] or "⚠️ incremento NON valido"
        state = "✅ rilasciabile" if n_open == 0 else f"🟢 {n_open} open"
        L.append(f"_{actual} bump da {B.fmt_version(ms['prev'])} · impatto richiesto dagli item: "
                 f"**{ms['required']}** · {state} · {n_done} done_")
        L.append("")
        L += _item_table(its)
        L.append("")

    # --- non pianificati: solo open (i done senza milestone sono storia) ---
    unplanned = sorted((it for it in items if not it["milestone"]
                        and it["status"].lower() != "done"), key=_sort_key)
    L.append("## Non pianificati (senza milestone)")
    L.append("")
    L.append(f"_{len(unplanned)} item open senza versione target. Assegnane una con "
             "`- **milestone**: vX.Y.Z — Titolo` in [`backlog.md`](backlog.md)._")
    L.append("")
    if unplanned:
        L += _item_table(unplanned, with_impact=False)
        L.append("")

    return "\n".join(L).rstrip() + "\n"


def render_release_notes(items, baseline, version_str):
    """Sezione CHANGELOG (Keep a Changelog) per una milestone: gli item raggruppati per
    impatto. Da incollare in CHANGELOG.md al momento del tag."""
    pm = B.parse_milestone(version_str)
    if pm is None:
        return None, f"`{version_str}` non è una versione `vX.Y.Z`"
    ver = pm[0]
    chain = B.milestone_chain(items, baseline)
    ms = next((m for m in chain if m["version"] == ver), None)
    if ms is None:
        known = ", ".join(B.fmt_version(m["version"]) for m in chain) or "nessuna"
        return None, f"nessuna milestone {B.fmt_version(ver)} nel backlog (presenti: {known})"

    L = [f"## [{B.fmt_version(ver)[1:]}] — DATA-DA-INSERIRE", ""]
    if ms["label"]:
        L.append(f"_{ms['label']}._")
        L.append("")
    n_open, _ = _counts(ms["items"])
    if n_open:
        L.append(f"> ⚠️ ATTENZIONE: {n_open} item di questa milestone sono ancora `open`. "
                 f"Non taggare finché non sono chiusi.")
        L.append("")
    if ms["required"] == "major":
        L.append("> 💥 Release **breaking**: i consumer devono adeguare il codice. Dettaglio "
                 "nella sezione «Changed (breaking)».")
        L.append("")
    for imp in ("major", "minor", "patch"):
        group = sorted((it for it in ms["items"] if B.item_impact(it) == imp), key=_sort_key)
        if not group:
            continue
        L.append(f"### {CHANGELOG_SECTION[imp]}")
        L.append("")
        for it in group:
            ref = f" ({it['ref']})" if it["ref"] else ""
            L.append(f"- {it['title']} — `{it['id']}`{ref}")
        L.append("")
    return "\n".join(L).rstrip() + "\n", None


def main():
    ap = argparse.ArgumentParser(description="Genera docs/roadmap.md dalle milestone del backlog.")
    ap.add_argument("--backlog", default=BACKLOG)
    ap.add_argument("--out", default=OUT)
    ap.add_argument("--baseline", default="",
                    help="versione già rilasciata da cui parte la catena (default: max tag git)")
    ap.add_argument("--check", action="store_true",
                    help="non scrive: exit 1 se l'output diverge dal file esistente (per CI)")
    ap.add_argument("--release-notes", metavar="vX.Y.Z", default="",
                    help="stampa su stdout la sezione CHANGELOG di quella milestone e esce")
    args = ap.parse_args()

    baseline = None
    if args.baseline:
        pm = B.parse_milestone(args.baseline)
        if pm is None:
            print(f"ERRORE: --baseline `{args.baseline}` non è una versione `vX.Y.Z`", file=sys.stderr)
            return 2
        baseline = pm[0]
    if baseline is None:
        baseline = B.released_baseline()

    items = B.parse_backlog(os.path.normpath(args.backlog))

    if args.release_notes:
        notes, err = render_release_notes(items, baseline, args.release_notes)
        if err:
            print(f"ERRORE: {err}", file=sys.stderr)
            return 2
        sys.stdout.write(notes)
        return 0

    content = render(items, baseline)

    if args.check:
        existing = ""
        if os.path.isfile(args.out):
            existing = open(args.out, encoding="utf-8").read()
        if existing != content:
            print("✗ docs/roadmap.md non aggiornata: rilancia `make roadmap` e committa.")
            return 1
        print("✅ docs/roadmap.md aggiornata.")
        return 0

    open(args.out, "w", encoding="utf-8").write(content)
    chain = B.milestone_chain(items, baseline)
    nxt = next((ms for ms in chain if _counts(ms["items"])[0] > 0), None)
    print(f"✅ scritta {os.path.relpath(args.out)} · {len(chain)} milestone · "
          f"prossima: {nxt['title'] if nxt else 'nessuna (backlog chiuso)'}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
