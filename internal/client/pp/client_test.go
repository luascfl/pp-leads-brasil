package pp_test

import (
	"os"
	"path/filepath"
	"testing"

	"pp-leads-brasil/internal/client/pp"
)

func leadTablePath(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "lead-table-sample.csv")
	content := "lead,CNPJ,contato,site,primeira mensagem\nFacto Agência de Comunicação,11.370.755/0001-02,contato@facto.test,https://facto.example,Mensagem de prospecção\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestRunCompanyGoatUsesCLIWhenAvailable(t *testing.T) {
	companyBin := tempCLI(t, `#!/usr/bin/env bash
echo '{"snapshot":"ok","source":"company-goat"}'
`)
	client := &pp.PPClient{CompanyGoatBin: companyBin, LeadTablePath: leadTablePath(t)}

	result, err := client.RunCompanyGoat("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("RunCompanyGoat returned error: %v", err)
	}

	payload, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result type = %T, want map[string]any", result)
	}
	goat, ok := payload["company_goat"].(map[string]any)
	if !ok {
		t.Fatalf("company_goat type = %T", payload["company_goat"])
	}
	if goat["status"] != "ok" {
		t.Fatalf("company_goat.status = %v", goat["status"])
	}
}

func TestRunContactGoatUsesCLIWhenAvailable(t *testing.T) {
	contactBin := tempCLI(t, `#!/usr/bin/env bash
echo '{"coverage":"ok","source":"contact-goat"}'
`)
	client := &pp.PPClient{ContactGoatBin: contactBin, LeadTablePath: leadTablePath(t)}

	result, err := client.RunContactGoat("Facto Agência de Comunicação")
	if err != nil {
		t.Fatalf("RunContactGoat returned error: %v", err)
	}

	payload, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result type = %T, want map[string]any", result)
	}
	goat, ok := payload["contact_goat"].(map[string]any)
	if !ok {
		t.Fatalf("contact_goat type = %T", payload["contact_goat"])
	}
	if goat["status"] != "ok" {
		t.Fatalf("contact_goat.status = %v", goat["status"])
	}
}

func TestRunEnrichBuildsProspectingLinks(t *testing.T) {
	companyBin := tempCLI(t, `#!/usr/bin/env bash
echo '{"snapshot":"ok"}'
`)
	contactBin := tempCLI(t, `#!/usr/bin/env bash
echo '{"coverage":"ok","phones":["71999990000"]}'
`)
	client := &pp.PPClient{CompanyGoatBin: companyBin, ContactGoatBin: contactBin, ScrapeCreatorsBin: contactBin, LeadTablePath: leadTablePath(t)}

	result, err := client.RunEnrich("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("RunEnrich returned error: %v", err)
	}

	payload := result.(map[string]any)
	if payload["cnpj"] != "11.370.755/0001-02" {
		t.Fatalf("cnpj = %v", payload["cnpj"])
	}
	links, ok := payload["prospecting_links"].(map[string]any)
	if !ok {
		t.Fatalf("prospecting_links type = %T", payload["prospecting_links"])
	}
	if _, ok := links["mailto"].(string); !ok {
		t.Fatalf("mailto missing: %#v", links)
	}
	if payload["company_goat"].(map[string]any)["status"] != "ok" {
		t.Fatalf("company_goat status = %v", payload["company_goat"])
	}
	if payload["contact_goat"].(map[string]any)["status"] != "ok" {
		t.Fatalf("contact_goat status = %v", payload["contact_goat"])
	}
}

func TestRunEnrichIncludesScrapeCreators(t *testing.T) {
	companyBin := tempCLI(t, `#!/usr/bin/env bash
echo '{"snapshot":"ok"}'
`)
	contactBin := tempCLI(t, `#!/usr/bin/env bash
echo '{"coverage":"ok","phones":["71999990000"]}'
`)
	scrapeBin := tempCLI(t, `#!/usr/bin/env bash
if [ "$1" = "google" ]; then
  echo '{"meta":{"source":"live"},"results":[{"url":"https://www.linkedin.com/company/facto-agencia"},{"url":"https://www.instagram.com/factoagencia/"}]}'
elif [ "$1" = "linkedin" ]; then
  echo '{"results":{"url":"https://www.linkedin.com/company/facto-agencia","name":"Facto"}}'
elif [ "$1" = "instagram" ]; then
  echo '{"results":{"handle":"factoagencia","url":"https://www.instagram.com/factoagencia/"}}'
else
  echo '{"results":{}}'
fi
`)
	client := &pp.PPClient{CompanyGoatBin: companyBin, ContactGoatBin: contactBin, ScrapeCreatorsBin: scrapeBin, LeadTablePath: leadTablePath(t)}

	result, err := client.RunEnrich("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("RunEnrich returned error: %v", err)
	}
	payload := result.(map[string]any)
	scrape, ok := payload["scrape_creators"].(map[string]any)
	if !ok {
		t.Fatalf("scrape_creators type = %T", payload["scrape_creators"])
	}
	if scrape["status"] != "ok" {
		t.Fatalf("scrape_creators status = %v", scrape["status"])
	}
}

func TestRunEnrichFallsBackToApifyWhenScrapeCreatorsFails(t *testing.T) {
	companyBin := tempCLI(t, "#!/usr/bin/env bash\necho '{\"snapshot\":\"ok\"}'\n")
	contactBin := tempCLI(t, "#!/usr/bin/env bash\necho '{\"coverage\":\"ok\"}'\n")
	scrapeBin := tempCLI(t, "#!/usr/bin/env bash\necho 'billing unavailable' >&2\nexit 1\n")
	apifyBin := tempCLI(t, "#!/usr/bin/env bash\necho '{\"platform\":\"instagram\",\"profile\":\"factoagencia\"}'\n")
	client := &pp.PPClient{
		CompanyGoatBin:     companyBin,
		ContactGoatBin:     contactBin,
		ScrapeCreatorsBin:  scrapeBin,
		ApifyBin:           apifyBin,
		ApifySocialActorID: "example/social-actor",
		LeadTablePath:      leadTablePath(t),
	}

	result, err := client.RunEnrich("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("RunEnrich returned error: %v", err)
	}
	payload := result.(map[string]any)
	if payload["social_provider"] != "apify" {
		t.Fatalf("social_provider = %v, want apify", payload["social_provider"])
	}
	apify, ok := payload["apify"].(map[string]any)
	if !ok || apify["status"] != "ok" {
		t.Fatalf("apify = %#v", payload["apify"])
	}
}

func TestRunEnrichAddsUseCaseMetadataWhenConfigured(t *testing.T) {
	companyBin := tempCLI(t, `#!/usr/bin/env bash
echo '{"snapshot":"ok"}'
`)
	contactBin := tempCLI(t, `#!/usr/bin/env bash
echo '{"coverage":"ok","phones":["71999990000"]}'
`)
	configDir := t.TempDir()
	data, err := os.ReadFile(leadTablePath(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "lead-table-sample.csv"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PP_LEADS_ICP_DIR", configDir)

	client := &pp.PPClient{CompanyGoatBin: companyBin, ContactGoatBin: contactBin, ScrapeCreatorsBin: contactBin, LeadTablePath: leadTablePath(t)}
	result, err := client.RunEnrich("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("RunEnrich returned error: %v", err)
	}

	payload := result.(map[string]any)
	useCase, ok := payload["use_case"].(map[string]any)
	if !ok {
		t.Fatalf("use_case type = %T", payload["use_case"])
	}
	if useCase["name"] != filepath.Base(configDir) {
		t.Fatalf("use_case.name = %v", useCase["name"])
	}
	leadContext, ok := payload["lead_context"].(map[string]any)
	if !ok {
		t.Fatalf("lead_context type = %T", payload["lead_context"])
	}
	if leadContext["primeira_mensagem"] == "" {
		t.Fatalf("lead_context missing primeira_mensagem: %#v", leadContext)
	}
}

func tempCLI(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fake-cli.sh")
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}
