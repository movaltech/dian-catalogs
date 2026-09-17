# dian-catalogs

Reference catalogs used by Colombia's DIAN (Dirección de Impuestos y Aduanas
Nacionales) electronic invoicing system — departments, municipalities, document
types, units of measure, tax types, and more — as plain JSON files, versioned with
git tags.

This repository is **data only**. There is no business logic, no framework
dependency, and no opinion about how you consume it — it's meant to be a single,
versioned source of truth that any application (in any language) can pull from,
instead of every project maintaining its own copy of the same reference data and
letting those copies drift apart over time.

---

## Contents

See [`index.json`](index.json) for the full list of available catalogs. Each one
lives at `catalogs/<name>.json` as an array of objects, most commonly shaped as
`{code, name, description}`:

```json
[
    { "code": "05", "name": "Antioquia", "description": "Departamento de Antioquia" },
    { "code": "08", "name": "Atlántico", "description": "Departamento del Atlántico" }
]
```

`municipalities.json` additionally carries `department_code`, relating each entry
back to `departments.json`.

Catalog values (codes, names, descriptions) are kept in the original Spanish used by
DIAN's own technical documentation, since that's the authoritative source — this
avoids introducing translation errors into regulatory reference data.

---

## How to use it

No installation and no server to run. Two common patterns, pick whichever fits:

### 1. On demand, via a CDN

GitHub plus a CDN such as [jsDelivr](https://www.jsdelivr.com/?docs=gh) serves any
file in a public repo as a plain HTTP endpoint, cached, with the version pinned
directly in the URL:

```
https://cdn.jsdelivr.net/gh/movaltech/dian-catalogs@v1.0.0/catalogs/departments.json
```

Fetch it with whatever HTTP client your language provides, decode the JSON, and
cache the result yourself (in memory, Redis, or wherever makes sense) — there's no
need to request it on every operation.

### 2. One-time download at deploy/seed time

Fetch the JSON once, when you deploy or seed your own database, and store it
locally from then on — the same shape as loading data from a local file, except the
source lives in one place instead of being copied by hand into every project that
needs it.

Either way, always pin an exact version (a tag, e.g. `@v1.0.0`) rather than tracking
a branch — see [Versioning](#versioning) below for why.

---

## Versioning

Every release is a git tag (`v1.0.0`, `v1.0.1`, `v1.2.0`, ...), following
[Semantic Versioning](https://semver.org/spec/v2.0.0.html):

- **Minor** version bump — a catalog gains new entries (e.g. a new tax liability
  code, a new municipality code).
- **Patch** version bump — a correction to existing, incorrect data.

Consumers pin the version they use and upgrade deliberately, never automatically —
that way a catalog update never silently changes production behavior.

See [`CHANGELOG.md`](CHANGELOG.md) for the full history.

---

## Data accuracy

These catalogs were compiled from DIAN's electronic invoicing Technical Annex and
from public standard code lists (DANE, UN/ECE Rec. 20, ISO 4217, CIIU). This is a
compilation, not an official source: **before relying on a specific code in
production, verify it against the current Technical Annex or the corresponding
DIAN/DANE source.** Some entries include a note in their `description` when a value
was added from a secondary source and hasn't yet been cross-checked against the
primary document — treat those as provisional until verified.

If you find outdated or incorrect data, please open an issue or a pull request with
a reference to the official source.

---

## Regenerating from CSV

[`scripts/csv_to_json.php`](scripts/csv_to_json.php) is the tool (plain PHP, no
dependencies) used to produce `catalogs/*.json` from source CSV files. It isn't part
of the published data itself — it's only needed when regenerating a catalog after a
source changes:

```bash
php scripts/csv_to_json.php --source=<folder with .csv files> --dest=catalogs
```

---

## License

[MIT](LICENSE) for the contents of this repository (structure, scripts). The
cataloged data itself originates from public sources of the Colombian government
(DIAN, DANE) and public international standards (UN/ECE, ISO).
