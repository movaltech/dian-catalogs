package catalogs

import (
	"encoding/json"
	"testing"
)

func TestIndexMatchesEmbeddedFiles(t *testing.T) {
	raw, err := embeddedFS.ReadFile("index.json")
	if err != nil {
		t.Fatalf("reading index.json: %v", err)
	}

	var idx indexFile
	if err := json.Unmarshal(raw, &idx); err != nil {
		t.Fatalf("parsing index.json: %v", err)
	}

	if len(idx.Catalogs) == 0 {
		t.Fatal("index.json declares no catalogs")
	}

	for _, c := range idx.Catalogs {
		entries := All(c.ID)
		if entries == nil {
			t.Errorf("catalog %q: not loaded", c.ID)
			continue
		}
		if len(entries) != c.Rows {
			t.Errorf("catalog %q: index.json says %d rows, loaded %d", c.ID, c.Rows, len(entries))
		}
	}
}

func TestCatalogIDsMatchesIndex(t *testing.T) {
	ids := CatalogIDs()
	if len(ids) != 19 {
		t.Fatalf("expected 19 catalogs, got %d: %v", len(ids), ids)
	}
}

func TestGetKnownEntry(t *testing.T) {
	entry, ok := Get(Departments, "05")
	if !ok {
		t.Fatal("expected department 05 to exist")
	}
	if entry.Name != "Antioquia" {
		t.Fatalf("expected Antioquia, got %q", entry.Name)
	}
}

// TestGetCurrencyHasSymbol and TestGetItemStandardHasAgencyID guard the two catalogs whose
// shape doesn't fit plain {code, name, description}: Currencies carries Symbol instead of
// Description, and ItemStandards carries AgencyID alongside it. A field present in the JSON
// but missing from Entry fails silently (encoding/json just drops it, no compile error, no
// panic) -- these exist so that regression is caught here instead of downstream, the way
// item_standards' agency_id already caused a real DIAN rejection (FAZ12) in a consumer that
// once hardcoded this exact lookup.
func TestGetCurrencyHasSymbol(t *testing.T) {
	entry, ok := Get(Currencies, "COP")
	if !ok {
		t.Fatal("expected currency COP to exist")
	}
	if entry.Symbol != "$" {
		t.Fatalf("expected COP symbol %q, got %q", "$", entry.Symbol)
	}
}

func TestGetItemStandardHasAgencyID(t *testing.T) {
	entry, ok := Get(ItemStandards, "001")
	if !ok {
		t.Fatal("expected item standard 001 (UNSPSC) to exist")
	}
	if entry.AgencyID != "10" {
		t.Fatalf("expected UNSPSC agency_id %q, got %q", "10", entry.AgencyID)
	}
}

// TestGetWithholdingConceptHasCompositeCodeAndAccounts guards
// WithholdingConcepts' own shape: a compound code
// ("{concepto}-{JURIDICA|NATURAL|BOTH}") and the accounting-specific
// fields (RateBp, MinBaseUVT, AccountPayable, AccountReceivable,
// ApplicableTo, Type) no other catalog populates.
func TestGetWithholdingConceptHasCompositeCodeAndAccounts(t *testing.T) {
	entry, ok := Get(WithholdingConcepts, "03-NATURAL")
	if !ok {
		t.Fatal("expected withholding concept 03-NATURAL to exist")
	}
	if entry.Type != "RETEFUENTE" || entry.RateBp != 1000 || entry.ApplicableTo != "NATURAL" {
		t.Fatalf("entry = %+v, want Type=RETEFUENTE RateBp=1000 ApplicableTo=NATURAL", entry)
	}
	if entry.AccountPayable != "236515" || entry.AccountReceivable != "135505" {
		t.Fatalf("entry = %+v, want AccountPayable=236515 AccountReceivable=135505", entry)
	}
}

func TestGetUVTHasYearAndValue(t *testing.T) {
	entry, ok := Get(UVT, "2025")
	if !ok {
		t.Fatal("expected uvt 2025 to exist")
	}
	if entry.Year != 2025 || entry.ValueCents != 4979900 {
		t.Fatalf("entry = %+v, want Year=2025 ValueCents=4979900", entry)
	}
}

func TestGetARLRateHasCompositeCodeAndRiskClass(t *testing.T) {
	entry, ok := Get(ARLRates, "2025-III")
	if !ok {
		t.Fatal("expected arl_rates 2025-III to exist")
	}
	if entry.Year != 2025 || entry.RiskClass != "III" || entry.RateBp != 244 {
		t.Fatalf("entry = %+v, want Year=2025 RiskClass=III RateBp=244", entry)
	}
}

func TestGetSMMLVHasValue(t *testing.T) {
	entry, ok := Get(SMMLV, "2026")
	if !ok {
		t.Fatal("expected smmlv 2026 to exist")
	}
	if entry.ValueCents != 150000000 {
		t.Fatalf("entry.ValueCents = %d, want 150000000", entry.ValueCents)
	}
}

func TestGetIncomeTaxRateHasRate(t *testing.T) {
	entry, ok := Get(IncomeTaxRates, "2022")
	if !ok {
		t.Fatal("expected income_tax_rates 2022 to exist")
	}
	if entry.RateBp != 3500 {
		t.Fatalf("entry.RateBp = %d, want 3500", entry.RateBp)
	}
}

func TestGetUnknownCatalog(t *testing.T) {
	if _, ok := Get("does_not_exist", "05"); ok {
		t.Fatal("expected unknown catalog to return ok=false")
	}
}

func TestGetUnknownCode(t *testing.T) {
	if _, ok := Get(Departments, "99999"); ok {
		t.Fatal("expected unknown code to return ok=false")
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid(DianTaxTypes, "01") {
		t.Fatal("expected DianTaxTypes 01 (IVA) to be valid")
	}
	if IsValid(DianTaxTypes, "zz") {
		t.Fatal("expected DianTaxTypes zz to be invalid")
	}
}

func TestIsValidMunicipality(t *testing.T) {
	if !IsValidMunicipality("05001", "05") {
		t.Fatal("expected Medellín (05001) to belong to Antioquia (05)")
	}
	if IsValidMunicipality("05001", "08") {
		t.Fatal("expected Medellín (05001) to not belong to Atlántico (08)")
	}
	if IsValidMunicipality("99999", "05") {
		t.Fatal("expected unknown municipality code to be invalid")
	}
}

func TestAllSortedByCode(t *testing.T) {
	entries := All(Departments)
	for i := 1; i < len(entries); i++ {
		if entries[i-1].Code >= entries[i].Code {
			t.Fatalf("entries not sorted by code: %q >= %q", entries[i-1].Code, entries[i].Code)
		}
	}
}
