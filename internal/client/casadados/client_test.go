package casadados_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"pp-leads-brasil/internal/client/casadados"
)

func leadTablePath() string {
	return filepath.Join("..", "..", "..", ".context", "organizejr", "lead-table-2026-06-17-ejs-comunicacao-lucas.csv")
}

func TestSearchCompanyByNameUsesLocalLeadTable(t *testing.T) {
	client := &casadados.CasaDadosClient{LeadTablePath: leadTablePath()}

	result, err := client.SearchCompanyByName("Facto")
	if err != nil {
		t.Fatalf("SearchCompanyByName returned error: %v", err)
	}

	results, ok := result.([]map[string]any)
	if !ok {
		t.Fatalf("result type = %T, want []map[string]any", result)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if results[0]["name"] != "Facto Agência de Comunicação" {
		t.Fatalf("name = %v", results[0]["name"])
	}
}

func TestSearchCompanyByCNPJUsesLocalLeadTable(t *testing.T) {
	client := &casadados.CasaDadosClient{LeadTablePath: leadTablePath()}

	result, err := client.SearchCompanyByCNPJ("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("SearchCompanyByCNPJ returned error: %v", err)
	}

	company, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("result type = %T, want map[string]any", result)
	}
	if company["name"] != "Facto Agência de Comunicação" {
		t.Fatalf("name = %v", company["name"])
	}
	if company["cnpj"] != "11.370.755/0001-02" {
		t.Fatalf("cnpj = %v", company["cnpj"])
	}
}

func TestSearchCompanyByNameUsesUseCaseConfigLeadTable(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "use-case.json")
	leadPath, err := filepath.Abs(leadTablePath())
	if err != nil {
		t.Fatal(err)
	}
	configContent := []byte("{\n  \"name\": \"sample-case\",\n  \"lead_table_path\": \"" + leadPath + "\"\n}")
	if err := os.WriteFile(configPath, configContent, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PP_LEADS_USE_CASE_CONFIG", configPath)

	client := &casadados.CasaDadosClient{}
	result, err := client.SearchCompanyByName("Facto")
	if err != nil {
		t.Fatalf("SearchCompanyByName returned error: %v", err)
	}
	results := result.([]map[string]any)
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if useCase, ok := results[0]["use_case"].(map[string]any); !ok || useCase["name"] != "sample-case" {
		t.Fatalf("use_case = %#v", results[0]["use_case"])
	}
}

func TestConsultaCNPJUsesCasaDosDadosAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v4/cnpj/11370755000102" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Header.Get("api-key") != "test-key" {
			t.Fatalf("api-key header = %q", r.Header.Get("api-key"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"cnpj":         "11370755000102",
			"razao_social": "FACTO AGENCIA DE COMUNICACAO",
			"situacao_cadastral": map[string]any{
				"situacao_cadastral": "ATIVA",
			},
		})
	}))
	defer server.Close()

	client := &casadados.CasaDadosClient{APIKey: "test-key", BaseURL: server.URL}
	result, err := client.ConsultaCNPJ("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("ConsultaCNPJ returned error: %v", err)
	}
	if result["source"] != "casa_dos_dados" {
		t.Fatalf("source = %v", result["source"])
	}
	if result["razao_social"] != "FACTO AGENCIA DE COMUNICACAO" {
		t.Fatalf("razao_social = %v", result["razao_social"])
	}
}

func TestPesquisaEmpresasUsesCasaDosDadosAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v5/cnpj/pesquisa" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s", r.Method)
		}
		if r.Header.Get("api-key") != "test-key" {
			t.Fatalf("api-key header = %q", r.Header.Get("api-key"))
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"total": 1,
			"cnpjs": []map[string]any{{
				"cnpj":         "11370755000102",
				"razao_social": "FACTO AGENCIA DE COMUNICACAO",
			}},
		})
	}))
	defer server.Close()

	client := &casadados.CasaDadosClient{APIKey: "test-key", BaseURL: server.URL}
	results, err := client.PesquisaEmpresas("Facto")
	if err != nil {
		t.Fatalf("PesquisaEmpresas returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("len(results) = %d", len(results))
	}
	if results[0]["source"] != "casa_dos_dados" {
		t.Fatalf("source = %v", results[0]["source"])
	}
}
