# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

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

[1.0.0]: https://github.com/movaltech/dian-catalogs/releases/tag/v1.0.0
