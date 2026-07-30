package pp

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestPublicWebsitePhonesAndProspectingLinks(t *testing.T) {
	site := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(`<!doctype html><html><body><a href="tel:+55 71 3333-4444">Ligue</a><a href="https://wa.me/5571999988888">WhatsApp</a></body></html>`))
			return
		}
		http.NotFound(w, r)
	}))
	defer site.Close()

	root := t.TempDir()
	leadTablePath := filepath.Join(root, "lead-table-sample.csv")
	csvContent := "lead,CNPJ,contato,site,primeira mensagem\nTeste Comunicação,11.370.755/0001-02,contato@teste.com," + site.URL + ",Olá pelo WhatsApp\n"
	if err := os.WriteFile(leadTablePath, []byte(csvContent), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PP_LEADS_ICP_DIR", root)

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
	payload, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result type = %T", result)
	}
	links, ok := payload["prospecting_links"].(map[string]any)
	if !ok {
		t.Fatalf("prospecting_links type = %T", payload["prospecting_links"])
	}
	email, ok := links["email"].(string)
	if !ok || email == "" {
		t.Fatalf("email = %#v", links["email"])
	}
	mailto, ok := links["mailto"].(string)
	if !ok || mailto == "" {
		t.Fatalf("mailto = %#v", links["mailto"])
	}
	wa, ok := links["whatsapp"].([]map[string]string)
	if !ok {
		t.Fatalf("whatsapp type = %T", links["whatsapp"])
	}
	if len(wa) == 0 {
		t.Fatalf("whatsapp links empty: %#v", links)
	}
}
