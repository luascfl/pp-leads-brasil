package pp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRunEnrichFallsBackToWebsitePhone(t *testing.T) {
	site := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			_, _ = w.Write([]byte(`<a href="/contato">Contato</a>`))
			return
		}
		if r.URL.Path == "/contato" || r.URL.Path == "/contato/" {
			_, _ = w.Write([]byte(`<a href="tel:+55 (21) 97611-0940">WhatsApp</a>`))
			return
		}
		http.NotFound(w, r)
	}))
	defer site.Close()

	root := t.TempDir()
	leadTablePath := filepath.Join(root, "leads.csv")
	configPath := filepath.Join(root, "use-case.json")
	csvContent := "lead,CNPJ,contato,site,primeira mensagem\nTeste Comunicação,11.370.755/0001-02,contato@teste.com," + site.URL + ",Olá pelo WhatsApp\n"
	if err := os.WriteFile(leadTablePath, []byte(csvContent), 0o644); err != nil {
		t.Fatal(err)
	}
	configContent := []byte("{\n  \"name\": \"sample\",\n  \"lead_table_path\": \"" + leadTablePath + "\",\n  \"message_field\": \"primeira mensagem\",\n  \"field_map\": {\n    \"name\": \"lead\",\n    \"cnpj\": \"CNPJ\",\n    \"contact\": \"contato\",\n    \"site\": \"site\"\n  }\n}\n")
	if err := os.WriteFile(configPath, configContent, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PP_LEADS_USE_CASE_CONFIG", configPath)

	companyBin := filepath.Join(root, "company.sh")
	if err := os.WriteFile(companyBin, []byte("#!/usr/bin/env bash\necho '{\"snapshot\":\"ok\"}'\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	contactBin := filepath.Join(root, "contact.sh")
	if err := os.WriteFile(contactBin, []byte("#!/usr/bin/env bash\necho '{\"coverage\":\"ok\",\"results\":[]}'\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	client := &PPClient{CompanyGoatBin: companyBin, ContactGoatBin: contactBin, ScrapeCreatorsBin: contactBin, LeadTablePath: leadTablePath}
	result, err := client.RunEnrich("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("RunEnrich returned error: %v", err)
	}
	payload := result.(map[string]any)
	links, ok := payload["prospecting_links"].(map[string]any)
	if !ok {
		t.Fatalf("prospecting_links type = %T", payload["prospecting_links"])
	}
	whatsapp, ok := links["whatsapp"].([]map[string]string)
	if !ok {
		items, okAny := links["whatsapp"].([]any)
		if !okAny || len(items) == 0 {
			t.Fatalf("whatsapp type = %T value=%#v", links["whatsapp"], links["whatsapp"])
		}
		first, okMap := items[0].(map[string]any)
		if !okMap {
			t.Fatalf("whatsapp item type = %T", items[0])
		}
		if first["phone"] != "5521976110940" {
			t.Fatalf("whatsapp phone = %v", first["phone"])
		}
		if first["wa_me"] == "" {
			t.Fatalf("wa_me missing: %#v", first)
		}
		return
	}
	if len(whatsapp) == 0 {
		t.Fatalf("whatsapp empty: %#v", links)
	}
	if whatsapp[0]["phone"] != "5521976110940" {
		t.Fatalf("whatsapp phone = %v", whatsapp[0]["phone"])
	}
	if whatsapp[0]["wa_me"] == "" {
		t.Fatalf("wa_me missing: %#v", whatsapp[0])
	}
}
