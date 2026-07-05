package usecase_test

import (
	"os"
	"path/filepath"
	"testing"

	"pp-leads-brasil/internal/usecase"
)

func TestLoadResolvesRelativePathsAndMetadata(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "use-case.json")
	content := []byte(`{
  "name": "sample-case",
  "label": "Sample case",
  "lead_table_path": "data/leads.csv",
  "output_dir": "outputs/enrichment",
  "message_field": "custom message",
  "field_map": {
    "name": "company_name",
    "cnpj": "document",
    "city_uf": "city_state",
    "contact": "email_address",
    "site": "website"
  }
}`)
	if err := os.WriteFile(configPath, content, 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := usecase.Load(configPath)
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got, want := cfg.ResolvePath(cfg.LeadTablePath), filepath.Join(root, "data", "leads.csv"); got != want {
		t.Fatalf("ResolvePath = %q, want %q", got, want)
	}
	meta := cfg.Metadata()
	if meta["name"] != "sample-case" || meta["label"] != "Sample case" {
		t.Fatalf("Metadata = %#v", meta)
	}
}

func TestLeadHelpersUseConfiguredFieldMap(t *testing.T) {
	record := map[string]string{
		"company_name":   "Acme Jr.",
		"document":       "12.345.678/0001-90",
		"city_state":     "Salvador/BA",
		"email_address":  "oi@acme.test",
		"website":        "https://acme.test",
		"custom message": "Mensagem de abordagem",
	}
	cfg := &usecase.Config{
		Name:         "sample-case",
		MessageField: "custom message",
		FieldMap: map[string]string{
			"name":    "company_name",
			"cnpj":    "document",
			"city_uf": "city_state",
			"contact": "email_address",
			"site":    "website",
		},
	}

	if got := usecase.Name(record, cfg); got != "Acme Jr." {
		t.Fatalf("Name = %q", got)
	}
	if got := usecase.CNPJ(record, cfg); got != "12.345.678/0001-90" {
		t.Fatalf("CNPJ = %q", got)
	}
	if got := usecase.Message(record, cfg); got != "Mensagem de abordagem" {
		t.Fatalf("Message = %q", got)
	}
	context := usecase.LeadContext(record)
	if context["company_name"] != "Acme Jr." {
		t.Fatalf("LeadContext = %#v", context)
	}
}
