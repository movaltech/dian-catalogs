# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.0] - 2026-09-20

### Added

- `puc`: Colombia's Plan Único de Cuentas (Decreto 2650 de 1993), 2503 accounts, 4 levels
  deep -- reference data for an application to seed a COMPANY's own chart of accounts from,
  not a table meant to be duplicated as-is into every company nor pointed at directly by a
  ledger's foreign keys. `Entry` gains `ParentCode` (this catalog's own parent-relation field,
  separate from `DepartmentCode`), `Level` (1-4), `Category` (Activo, Pasivo, Patrimonio,
  Ingreso, Gasto, Costo, Costo de Producción, Cuenta de Orden Deudora/Acreedora), `IsPosting`
  (false for a grouping/header account that must never receive a posting directly) and
  `IsActive`.
- Five catalogs of Colombian tax/labor values set by law rather than by DIAN's electronic
  invoicing Anexo Técnico, ported from an ERP's own accounting/payroll seed data so every
  consumer shares one source instead of each copying it by hand:
  - `withholding_concepts`: 19 retención en la fuente / ReteIVA / ReteICA concepts (rate,
    minimum base in UVT, PUC accounts). `code` is compound,
    `"{concepto}-{JURIDICA|NATURAL|BOTH}"` (e.g. `"01-JURIDICA"`, `"04-BOTH"`), since some
    concepts have a different rate for a declarante (JURIDICA) than a no-declarante
    (NATURAL) and others share one rate for both.
  - `uvt`: Unidad de Valor Tributario by year, 2020-2025.
  - `income_tax_rates`: general corporate income tax rate by year, 2019-2026.
  - `arl_rates`: occupational-risk (ARL) contribution rate by risk class (I-V) and year,
    2024-2026.
  - `smmlv`: Salario Mínimo Mensual Legal Vigente by year, 2019-2026.

  `Entry` gains `Year`, `RateBp`, `MinBaseUVT`, `AccountPayable`, `AccountReceivable`,
  `ApplicableTo`, `RiskClass`, `ValueCents` and `Type` -- all `omitempty`, populated only by
  these five catalogs, the same "not every catalog fills every field" convention
  `DepartmentCode`/`Symbol`/`AgencyID` already established.
- Go binding (`package catalogs`, `go get github.com/movaltech/dian-catalogs`): embeds
  `index.json` and `catalogs/*.json` via `go:embed` and exposes a single generic loader
  (`Get`, `IsValid`, `IsValidMunicipality`, `All`, `CatalogIDs`) driven entirely by
  `index.json`, so it never needs a struct-per-catalog or code changes when a catalog is
  added. The JSON files remain the canonical, language-agnostic source; this is purely an
  additive Go convenience layer.

### Fixed

- `ciiu_codes`: removed a duplicate `9900` entry (the CSV source had it twice, with slightly
  different descriptions — kept the one noting "este sí pertenece a la Sección U"). `rows` in
  `index.json` corrected from 508 to 507.

## [1.0.0] - 2026-09-16

### Added

- Initial set of 14 catalogs: countries, currencies, departments, municipalities, DIAN
  document types, DIAN tax types, identification types, item standards, liability codes,
  payment methods, payment terms, tax regimes, unit measures, and CIIU codes.
- `identification_types`: code `48` "PPT" (Permiso por Protección Temporal), introduced by
  Resolución DIAN 000165 de 2023, in force with Anexo Técnico 1.9 since 2024-05-01.
- `dian_document_types`: codes `20` (Documento Equivalente Electrónico / POS), `93` and `94`
  (its debit/credit adjustment notes).
- `dian_tax_types`: codes `34` (IBUA) and `35` (ICUI), from Ley 2277 de 2022. Found via
  secondary sources, not yet verified against the official Anexo Técnico text — flagged as
  such in their `description` until confirmed.
- `index.json` manifest listing every catalog file and row count.
- `scripts/csv_to_json.php` conversion tool.

[1.2.0]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.2.0
[1.0.0]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.0.0
