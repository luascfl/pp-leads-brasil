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
	return filepath.Join("..", "..", "..", "organizejr-pp-leads", "icp", "ejs-comunicacao", "lead-table-2026-06-17-ejs-comunicacao-lucas.csv")
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
	configDir := t.TempDir()
	leadPath, err := filepath.Abs(leadTablePath())
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(leadPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "lead-table-sample.csv"), data, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PP_LEADS_ICP_DIR", configDir)

	client := &casadados.CasaDadosClient{}
	result, err := client.SearchCompanyByName("Facto")
	if err != nil {
		t.Fatalf("SearchCompanyByName returned error: %v", err)
	}
	results := result.([]map[string]any)
	if len(results) != 1 {
		t.Fatalf("len(results) = %d, want 1", len(results))
	}
	if useCase, ok := results[0]["use_case"].(map[string]any); !ok || useCase["name"] != filepath.Base(configDir) {
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
			"razao_social": "Facto Agência de Comunicação",
		})
	}))
	defer server.Close()

	client := &casadados.CasaDadosClient{APIKey: "test-key", BaseURL: server.URL}
	result, err := client.ConsultaCNPJ("11.370.755/0001-02")
	if err != nil {
		t.Fatalf("ConsultaCNPJ returned error: %v", err)
	}
	if result["razao_social"] != "Facto Agência de Comunicação" {
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
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		busca, ok := body["busca_textual"].([]any)
		if !ok || len(busca) != 1 {
			t.Fatalf("busca_textual = %#v", body["busca_textual"])
		}
		first, ok := busca[0].(map[string]any)
		if !ok {
			t.Fatalf("busca_textual[0] = %#v", busca[0])
		}
		texto, ok := first["texto"].([]any)
		if !ok || len(texto) != 1 || texto[0] != "Facto" {
			t.Fatalf("texto = %#v", first["texto"])
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"cnpjs": []map[string]any{{
				"cnpj":         "11370755000102",
				"razao_social": "Facto Agência de Comunicação",
			}},
		})
	}))
	defer server.Close()

	client := &casadados.CasaDadosClient{APIKey: "test-key", BaseURL: server.URL}
	result, err := client.PesquisaEmpresas("Facto")
	if err != nil {
		t.Fatalf("PesquisaEmpresas returned error: %v", err)
	}
	if len(result) != 1 || result[0]["razao_social"] != "Facto Agência de Comunicação" {
		t.Fatalf("result = %#v", result)
	}
}
