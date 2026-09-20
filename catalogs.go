// Package catalogs is the Go binding for this repository's reference catalogs: departments,
// municipalities, DIAN document/tax types, and the rest of catalogs/*.json. The JSON files
// remain the source of truth (consumable by any language, e.g. via jsDelivr) — this package
// only embeds them and adds typed, in-memory lookups for Go consumers.
package catalogs

import (
	"embed"
	"encoding/json"
	"fmt"
	"sort"
)

// Catalog IDs, as listed in index.json.
const (
	Countries           = "countries"
	Currencies          = "currencies"
	Departments         = "departments"
	Municipalities      = "municipalities"
	DianDocumentTypes   = "dian_document_types"
	DianTaxTypes        = "dian_tax_types"
	IdentificationTypes = "identification_types"
	ItemStandards       = "item_standards"
	LiabilityCodes      = "liability_codes"
	PaymentMethods      = "payment_methods"
	PaymentTerms        = "payment_terms"
	TaxRegimes          = "tax_regimes"
	UnitMeasures        = "unit_measures"
	CIIUCodes           = "ciiu_codes"

	// WithholdingConcepts, UVT, IncomeTaxRates, ARLRates and SMMLV are
	// Colombian tax/labor values set by law rather than by DIAN's
	// electronic-invoicing Anexo Técnico, but the same shape applies: a
	// flat, versioned lookup that changes at most once a year, with no
	// relational need beyond "look this code up". WithholdingConcepts'
	// code is compound, "{concepto}-{JURIDICA|NATURAL|BOTH}" (e.g.
	// "01-JURIDICA", "04-BOTH"), because the source table gives some
	// concepts a different rate for a declarante (JURIDICA) vs. a
	// no-declarante (NATURAL) and others a single rate for both --  a
	// caller resolves which applies first, then looks up the composite
	// code (falling back to the "-BOTH" suffix for concepts that never
	// split). UVT/IncomeTaxRates/SMMLV key by year alone (e.g. "2025");
	// ARLRates by "{year}-{risk_class}" (e.g. "2025-III").
	WithholdingConcepts = "withholding_concepts"
	UVT                 = "uvt"
	IncomeTaxRates      = "income_tax_rates"
	ARLRates            = "arl_rates"
	SMMLV               = "smmlv"
)

// Entry is a single row of a catalog. Not every catalog populates every field: DepartmentCode
// is only set for Municipalities (relating each one back to a Departments entry), Symbol only
// for Currencies, AgencyID only for ItemStandards (the DIAN table 13.3.5 agency each standard
// belongs to, e.g. "10" for UNSPSC), and the accounting/labor fields below only for
// WithholdingConcepts/UVT/IncomeTaxRates/ARLRates/SMMLV -- a field left empty for a given
// catalog simply isn't part of its shape, the same way Description is empty for Countries.
type Entry struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	DepartmentCode string `json:"department_code,omitempty"`
	Symbol         string `json:"symbol,omitempty"`
	AgencyID       string `json:"agency_id,omitempty"`

	// Year is set by UVT, IncomeTaxRates, ARLRates and SMMLV -- the same
	// value already encoded in Code, kept as its own typed field so a
	// caller doesn't need to parse it back out.
	Year int `json:"year,omitempty"`
	// RateBp is a rate in basis points (1/100 of a percent): a
	// withholding tax rate (WithholdingConcepts), the general corporate
	// income tax rate (IncomeTaxRates), or an occupational-risk (ARL)
	// contribution rate (ARLRates).
	RateBp int `json:"rate_bp,omitempty"`
	// MinBaseUVT is set only by WithholdingConcepts: the minimum
	// transaction base, in UVT, below which the withholding does not
	// apply. 0 means no minimum.
	MinBaseUVT int `json:"min_base_uvt,omitempty"`
	// AccountPayable/AccountReceivable are set only by
	// WithholdingConcepts: the PUC (Colombian chart of accounts) codes
	// this concept posts to when the company is the withholding agent
	// (AccountPayable, a liability to DIAN) versus when it is the one
	// withheld from (AccountReceivable, an advance tax credit).
	AccountPayable    string `json:"account_payable,omitempty"`
	AccountReceivable string `json:"account_receivable,omitempty"`
	// ApplicableTo is set only by WithholdingConcepts: "JURIDICA",
	// "NATURAL", or "BOTH" -- see WithholdingConcepts' own doc comment
	// for how this relates to Code.
	ApplicableTo string `json:"applicable_to,omitempty"`
	// RiskClass is set only by ARLRates: "I" through "V", lowest to
	// highest occupational risk.
	RiskClass string `json:"risk_class,omitempty"`
	// ValueCents is a monetary amount in Colombian peso cents: the UVT's
	// own value (UVT) or the minimum wage (SMMLV) for Year.
	ValueCents int64 `json:"value_cents,omitempty"`
	// Type is set only by WithholdingConcepts: "RETEFUENTE", "RETEIVA",
	// or "RETEICA".
	Type string `json:"type,omitempty"`
}

//go:embed index.json catalogs/*.json
var embeddedFS embed.FS

type indexEntry struct {
	ID          string `json:"id"`
	File        string `json:"file"`
	Description string `json:"description"`
	Rows        int    `json:"rows"`
}

type indexFile struct {
	Version  string       `json:"version"`
	Updated  string       `json:"updated"`
	Catalogs []indexEntry `json:"catalogs"`
}

var catalogsByID map[string]map[string]Entry

func init() {
	raw, err := embeddedFS.ReadFile("index.json")
	if err != nil {
		panic(fmt.Errorf("catalogs: reading index.json: %w", err))
	}

	var idx indexFile
	if err := json.Unmarshal(raw, &idx); err != nil {
		panic(fmt.Errorf("catalogs: parsing index.json: %w", err))
	}

	catalogsByID = make(map[string]map[string]Entry, len(idx.Catalogs))
	for _, c := range idx.Catalogs {
		data, err := embeddedFS.ReadFile(c.File)
		if err != nil {
			panic(fmt.Errorf("catalogs: reading %s (catalog %q): %w", c.File, c.ID, err))
		}

		var entries []Entry
		if err := json.Unmarshal(data, &entries); err != nil {
			panic(fmt.Errorf("catalogs: parsing %s (catalog %q): %w", c.File, c.ID, err))
		}

		byCode := make(map[string]Entry, len(entries))
		for _, e := range entries {
			byCode[e.Code] = e
		}
		catalogsByID[c.ID] = byCode
	}
}

// Get returns the entry for code in the given catalog. ok is false if the catalog ID is
// unknown or the code does not exist in it.
func Get(catalogID, code string) (entry Entry, ok bool) {
	byCode, ok := catalogsByID[catalogID]
	if !ok {
		return Entry{}, false
	}
	entry, ok = byCode[code]
	return entry, ok
}

// IsValid reports whether code exists in the given catalog.
func IsValid(catalogID, code string) bool {
	_, ok := Get(catalogID, code)
	return ok
}

// IsValidMunicipality reports whether code is a known municipality belonging to department
// departmentCode.
func IsValidMunicipality(code, departmentCode string) bool {
	entry, ok := Get(Municipalities, code)
	return ok && entry.DepartmentCode == departmentCode
}

// All returns every entry in the given catalog, sorted by code. A nil slice means the catalog
// ID is unknown.
func All(catalogID string) []Entry {
	byCode, ok := catalogsByID[catalogID]
	if !ok {
		return nil
	}

	entries := make([]Entry, 0, len(byCode))
	for _, e := range byCode {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Code < entries[j].Code })
	return entries
}

// CatalogIDs returns every known catalog ID, as declared in index.json.
func CatalogIDs() []string {
	ids := make([]string, 0, len(catalogsByID))
	for id := range catalogsByID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
