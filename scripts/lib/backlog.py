#!/usr/bin/env python3
"""
backlog.py — parsing + regole di validazione condivise di docs/backlog.md.

FONTE UNICA delle regole strutturali del backlog: importata sia da `backlog-lint.py`
(che valida in CI) sia da `generate-roadmap.py` (che genera docs/roadmap.md). Niente
duplicazione di regex/convenzioni tra i due.

Un item del backlog è:

    ### `id-stabile` — Titolo
    - **status**: open|done            (default: open)
    - **priority**: low|medium|high    (opzionale)
    - **impact**: patch|minor|major    (opzionale, default patch — governa il versionamento)
    - **labels**: a, b, c              (opzionale)
    - **milestone**: vX.Y.Z — Titolo   (opzionale)
    - **owner**: <username>            (opzionale)
    - **ref**: <link>                  (opzionale)

    prosa … (descrizione dell'intervento)

Differenza rispetto al backlog di devops_hiway: qui la milestone **è una versione semver**
e ogni item dichiara il proprio `impact`. Il linter verifica che la catena delle milestone
regga il semver: una milestone che contiene un item `major` deve essere un major bump.
Questo è il "versionamento annesso" — è impossibile pianificare un breaking change dentro
una minor senza che la CI se ne accorga.

`parse_backlog()` è PURA (nessun sys.exit): ritorna la lista degli item. Ogni item porta,
oltre ai campi di contenuto, due chiavi tecniche per il linter:
  - `_line`      : numero di riga (1-based) dell'heading `### ...`
  - `_meta_seen` : lista di {key, line, value} dei bullet `- **k**: v` trovati nel *blocco
                   meta* (la sequenza di bullet subito sotto l'heading, prima della prima
                   riga di prosa) — serve a stanare chiavi sconosciute (refusi).
"""
import re
import subprocess

# --- regex strutturali (fonte unica) ---------------------------------------------------
# item: "### `id` — Titolo"  (id tra backtick; separatore — o - opzionale)
ITEM_RE = re.compile(r"^###\s+`([^`]+)`\s*[—\-]*\s*(.*)$")
META_RE = re.compile(r"^[-*]\s+\*\*(\w+)\*\*:\s*(.+)$")
# milestone: "v1.1.0" oppure "v1.1.0 — Etichetta descrittiva"
MILESTONE_RE = re.compile(r"^v(\d+)\.(\d+)\.(\d+)(?:\s*[—-]\s*(.*))?$")
TAG_RE = re.compile(r"^v(\d+)\.(\d+)\.(\d+)$")

# --- vocabolario ammesso ---------------------------------------------------------------
KNOWN_META = {"status", "priority", "impact", "labels", "milestone", "owner", "ref"}
VALID_STATUS = {"open", "done"}
VALID_PRIORITY = {"low", "medium", "high"}
VALID_IMPACT = {"patch", "minor", "major"}
IMPACT_RANK = {"patch": 0, "minor": 1, "major": 2}
ID_RE = re.compile(r"^[a-z0-9][a-z0-9._-]*$")   # kebab-case minuscolo

DEFAULT_IMPACT = "patch"


def _new_item(item_id, title, line):
    return {"id": item_id, "title": title,
            "labels": [], "status": "open", "priority": "", "impact": "", "ref": "", "owner": "",
            "milestone": "", "body": [], "_line": line, "_meta_seen": []}


def parse_backlog(path):
    """Parsa docs/backlog.md → lista di item (PURA, nessun sys.exit).

    Un bullet `- **k**: v` con `k` noto imposta il campo; qualunque altra riga (blank, prosa,
    checklist, bullet con chiave non nota) finisce nel corpo. In più registra in `_meta_seen`
    i bullet incontrati nel blocco meta iniziale (finché non parte la prosa), per permettere
    al linter di segnalare i refusi di chiave.
    """
    items, cur, in_meta = [], None, False
    with open(path, encoding="utf-8") as f:
        lines = f.read().splitlines()
    for i, ln in enumerate(lines, start=1):
        m = ITEM_RE.match(ln)
        if m:
            if cur:
                items.append(cur)
            cur = _new_item(m.group(1).strip(), m.group(2).strip(), i)
            in_meta = True
            continue
        if cur is None:
            continue
        if re.match(r"^#{1,3}\s", ln):   # nuova sezione (non-item) → chiude l'item corrente
            items.append(cur)
            cur, in_meta = None, False
            continue
        mm = META_RE.match(ln)
        if mm:
            k, v = mm.group(1).lower(), mm.group(2).strip()
            if in_meta:
                cur["_meta_seen"].append({"key": mm.group(1), "line": i, "value": v})
            if k == "labels":
                cur["labels"] = [x.strip().strip("`") for x in v.split(",") if x.strip()]
            elif k in ("status", "priority", "impact", "ref", "owner", "milestone"):
                cur[k] = v.strip().lstrip("@")
            else:
                cur["body"].append(ln)   # chiave non nota → corpo
        else:
            if ln.strip():               # prima riga di prosa non-blank → fine blocco meta
                in_meta = False
            cur["body"].append(ln)
    if cur:
        items.append(cur)
    return items


# --- semver ----------------------------------------------------------------------------

def parse_milestone(title):
    """'v1.1.0 — Etichetta' → ((1,1,0), 'Etichetta'). None se il titolo non è una versione."""
    m = MILESTONE_RE.match(title.strip())
    if not m:
        return None
    return (int(m.group(1)), int(m.group(2)), int(m.group(3))), (m.group(4) or "").strip()


def fmt_version(v):
    return "v%d.%d.%d" % v


def bump_kind(v_from, v_to):
    """Tipo di incremento semver da v_from a v_to: 'major'|'minor'|'patch', o None se v_to
    non è un incremento valido (uguale, minore, o salto incoerente tipo 1.0.0 → 2.1.0)."""
    if v_to <= v_from:
        return None
    if v_to[0] > v_from[0]:
        # major: le componenti inferiori devono azzerarsi (1.4.2 → 2.0.0, non 2.3.1)
        return "major" if v_to[1] == 0 and v_to[2] == 0 else None
    if v_to[1] > v_from[1]:
        return "minor" if v_to[2] == 0 else None
    return "patch"


def item_impact(item):
    return (item.get("impact") or DEFAULT_IMPACT).lower()


def required_bump(items):
    """Impatto massimo di un gruppo di item = bump minimo che la milestone deve rappresentare."""
    ranks = [IMPACT_RANK.get(item_impact(it), 0) for it in items]
    top = max(ranks) if ranks else 0
    return next(k for k, v in IMPACT_RANK.items() if v == top)


def released_baseline(default=(0, 0, 0)):
    """Ultima versione già rilasciata, dai tag git `vX.Y.Z`. Fallback `default` se git non
    è disponibile o non ci sono tag (la generazione della roadmap non deve mai dipendere da git)."""
    try:
        out = subprocess.run(["git", "tag", "--list", "v*"], capture_output=True, text=True,
                             timeout=10, check=True).stdout
    except Exception:
        return default
    vs = []
    for ln in out.splitlines():
        m = TAG_RE.match(ln.strip())
        if m:
            vs.append((int(m.group(1)), int(m.group(2)), int(m.group(3))))
    return max(vs) if vs else default


def milestone_chain(items, baseline):
    """Milestone ordinate per versione, con il bump richiesto e quello effettivo.

    Ritorna una lista di dict: version, label, title, items, required, actual, prev.
    `actual` è calcolato rispetto alla milestone precedente della catena (o alla baseline
    per la prima): è così che una v1.1.0 e una v1.2.0 restano entrambe minor bump legittimi.
    """
    planned = {}
    for it in items:
        if it["milestone"]:
            planned.setdefault(it["milestone"], []).append(it)

    parsed = []
    for title, its in planned.items():
        pm = parse_milestone(title)
        if pm is None:
            continue
        parsed.append({"version": pm[0], "label": pm[1], "title": title, "items": its})
    parsed.sort(key=lambda m: m["version"])

    prev = baseline
    for ms in parsed:
        ms["prev"] = prev
        ms["required"] = required_bump(ms["items"])
        ms["actual"] = bump_kind(prev, ms["version"])
        prev = ms["version"]
    return parsed


def _norm_milestone(title):
    """Normalizza un titolo milestone per il match 'stessa milestone' (case + spazi collassati)."""
    return re.sub(r"\s+", " ", title).strip().casefold()


def lint(items, baseline=None):
    """Valida gli item parsati. Ritorna (errors, warnings): liste di stringhe già formattate
    (con numero di riga). `errors` non vuoto ⇒ il linter esce con status ≠0."""
    errors, warnings = [], []
    if baseline is None:
        baseline = released_baseline()

    # 1) id duplicati
    by_id = {}
    for it in items:
        by_id.setdefault(it["id"], []).append(it["_line"])
    for iid, ls in by_id.items():
        if len(ls) > 1:
            errors.append(f"id duplicato `{iid}` (righe {', '.join(map(str, ls))})")

    for it in items:
        loc = f"riga {it['_line']} [{it['id']}]"

        # 2) formato id
        if not ID_RE.match(it["id"]):
            errors.append(f"{loc}: id non valido (atteso kebab-case minuscolo `[a-z0-9._-]`)")

        # 3) titolo non vuoto
        if not it["title"].strip():
            errors.append(f"{loc}: titolo vuoto")

        # 4) status / priority / impact nel vocabolario
        if it["status"] and it["status"].lower() not in VALID_STATUS:
            errors.append(f"{loc}: status `{it['status']}` non valido (atteso {sorted(VALID_STATUS)})")
        if it["priority"] and it["priority"].lower() not in VALID_PRIORITY:
            errors.append(f"{loc}: priority `{it['priority']}` non valida (attesa {sorted(VALID_PRIORITY)})")
        if it["impact"] and it["impact"].lower() not in VALID_IMPACT:
            errors.append(f"{loc}: impact `{it['impact']}` non valido (atteso {sorted(VALID_IMPACT)})")

        # 5) chiavi meta sconosciute nel blocco iniziale (refusi tipo `- **lables**:`)
        for meta in it["_meta_seen"]:
            if meta["key"].lower() not in KNOWN_META:
                errors.append(f"riga {meta['line']} [{it['id']}]: chiave meta sconosciuta "
                              f"`{meta['key']}` (note: {sorted(KNOWN_META)}) — refuso?")

        # 6) titolo milestone = versione semver
        if it["milestone"] and parse_milestone(it["milestone"]) is None:
            errors.append(f"{loc}: milestone `{it['milestone']}` non è una versione "
                          f"(atteso `vX.Y.Z` o `vX.Y.Z — Etichetta`)")

    # 7) coerenza titoli milestone: stessi caratteri ovunque
    variants = {}
    for it in items:
        if it["milestone"]:
            variants.setdefault(_norm_milestone(it["milestone"]), set()).add(it["milestone"])
    for raws in variants.values():
        if len(raws) > 1:
            warnings.append("milestone scritta in modi diversi (stessa milestone?): "
                            + " · ".join(f"«{r}»" for r in sorted(raws))
                            + " → uniformare il titolo (match esatto per carattere)")

    # 8) etichetta della milestone coerente fra gli item della stessa versione
    by_version = {}
    for it in items:
        pm = parse_milestone(it["milestone"]) if it["milestone"] else None
        if pm:
            by_version.setdefault(pm[0], set()).add(it["milestone"])
    for ver, titles in by_version.items():
        if len(titles) > 1:
            errors.append(f"{fmt_version(ver)}: la stessa versione ha titoli diversi "
                          + " · ".join(f"«{t}»" for t in sorted(titles))
                          + " → una versione = un titolo")

    # 9) VERSIONAMENTO: la catena delle milestone deve reggere il semver
    chain = milestone_chain(items, baseline)
    for ms in chain:
        v, prev = fmt_version(ms["version"]), fmt_version(ms["prev"])
        if ms["actual"] is None:
            errors.append(f"{v}: non è un incremento semver valido rispetto a {prev} "
                          f"(atteso {prev}→patch/minor/major con le componenti inferiori azzerate)")
            continue
        if IMPACT_RANK[ms["actual"]] < IMPACT_RANK[ms["required"]]:
            culprits = sorted(it["id"] for it in ms["items"]
                              if IMPACT_RANK[item_impact(it)] > IMPACT_RANK[ms["actual"]])
            errors.append(
                f"{v} è un {ms['actual']} bump rispetto a {prev} ma contiene item "
                f"`{ms['required']}` → serve un {ms['required']} bump. "
                f"Item incompatibili: {', '.join('`%s`' % c for c in culprits)}")

    # 10) segnali dinamici (non bloccanti)
    for ms in chain:
        n_open = sum(1 for it in ms["items"] if it["status"].lower() != "done")
        if n_open == 0:
            warnings.append(f"{ms['title']}: 0 item open → milestone PRONTA, taggare "
                            f"{fmt_version(ms['version'])} (il push del tag lo fa l'utente)")
    unplanned = [it for it in items if not it["milestone"] and it["status"].lower() != "done"]
    if unplanned:
        warnings.append(f"{len(unplanned)} item open senza milestone (non entrano in nessuna "
                        f"release): {', '.join('`%s`' % it['id'] for it in unplanned[:8])}"
                        + (" …" if len(unplanned) > 8 else ""))

    return errors, warnings
