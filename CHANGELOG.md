# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.6.0] - 2026-11-22

### Added

- `adjustment_note_concepts` (5 codes): the correction concepts of the adjustment note to the
  support document (nota de ajuste al documento soporte), table 16.2.4 of the support document
  annex V1.1 (Resolución 000167 de 2021), `cac:DiscrepancyResponse/cbc:ResponseCode`. Its
  `cbc:Description` is free text of 20 to 5000 characters (rule NSBF04), not the code's name.

## [1.5.1] - 2026-11-20

### Fixed

- `liability_codes`: `R-99-PN` is named "No aplica – Otros", the text of DIAN table 13.2.6.1,
  instead of "No aplica (persona natural)": the table's note applies it to legal persons,
  natural persons and final consumers alike. `R-99-PJ` and `O-49` now say they are not in
  that table; DIAN habilitacion answers notice FAK26 for a receiver with `R-99-PJ`.

## [1.5.0] - 2026-11-19

### Added

- `tax_rates`: the percent rates DIAN table 13.3.11 allows in `cbc:Percent`, copied from the
  Excel file of the Caja de Herramientas FE V1.9 (2026 edition), for IVA (`01`: 0, 5, 16, 19),
  INC (`04`: 2, 4, 8, 16), ReteIVA (`05`: 15, 100) and ReteRenta (`06`, 14 rates). The code
  is `{tax}-{rate}` with the rate at 2 decimals, e.g. `01-19.00`; the description lists the
  table's concepts for each rate. Taxes with nominal values (INPP, IBUA, ICL) are left out.

## [1.4.0] - 2026-11-19

### Added

- `allowance_charge_codes`: the discount and charge codes of DIAN table 13.3.8
  (`cac:AllowanceCharge/cbc:AllowanceChargeReasonCode`), copied from the Excel file of the Caja
  de Herramientas FE V1.9 (2026 edition): `00` and `01` discounts, `02` and `03` charges. The
  footnote markers of the official names were dropped.

## [1.3.1] - 2026-11-19

### Fixed

- `credit_note_concepts`: code `3` is "Rebaja  o descuento parcial o total", with the double
  space of DIAN's official table. Release 1.3.0 reduced it to one space, and DIAN habilitacion
  accepts such a note with notice CBF04 ("La Descripción informada no se encuentra en la
  lista"); with the double space it accepts it without notices. Every other concept of both
  catalogs was sent to DIAN habilitacion and accepted without notices.

## [1.3.0] - 2026-11-19

### Added

- `credit_note_concepts` (6 codes) and `debit_note_concepts` (4 codes): the correction
  concepts of DIAN tables 13.2.4 and 13.2.5 (`cac:DiscrepancyResponse/cbc:ResponseCode`),
  copied from the Excel files of the Caja de Herramientas FE V1.9 (2026 edition). DIAN rejects
  a note whose `cbc:Description` does not match its code (rule CBF04), accents included, so
  consumers should send `name` as is. Trailing spaces of the official file were trimmed and
  the double space in credit note concept 3 ("Rebaja  o descuento") was reduced to one.

### Fixed

- `tax_regimes`: code `49` is named "No responsable de IVA" (the pair of `48`) and its
  description covers its use in sales invoices and their notes, with R-99-PN, besides the
  support document. It used to say only "no obligado a facturar".
- `index.json` reported version 1.2.0 in release 1.2.1.
- A `.gitattributes` keeps LF line endings on every platform; CRLF checkouts on Windows broke
  `gofmt` and the golden comparisons.

## [1.2.1] - 2026-11-19

### Fixed

- `tax_regimes`: add code `48` (Responsable de IVA), the @listName of `cbc:TaxLevelCode` that
  DIAN's own sample invoices in the 2026 Caja de Herramientas FE V1.9 use together with `O-13`,
  paired with `49`. The catalog had `49` but not `48`, so a consumer validating against it
  rejected the most common value. Entries are now sorted by code (`49` came before `02`).

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

[1.5.1]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.5.1
[1.5.0]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.5.0
[1.4.0]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.4.0
[1.3.1]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.3.1
[1.3.0]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.3.0
[1.2.1]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.2.1
[1.2.0]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.2.0
[1.0.0]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.0.0
