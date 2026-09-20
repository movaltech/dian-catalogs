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
	if len(ids) != 14 {
		t.Fatalf("expected 14 catalogs, got %d: %v", len(ids), ids)
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
