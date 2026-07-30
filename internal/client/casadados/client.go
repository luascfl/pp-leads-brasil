package casadados

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"pp-leads-brasil/internal/usecase"
)

var ErrNotFound = errors.New("company not found")

// Client defines the interface for company lookup data.
type Client interface {
	SearchCompanyByCNPJ(cnpj string) (interface{}, error)
	SearchCompanyByName(name string) (interface{}, error)
}

type CasaDadosClient struct {
	APIKey        string
	LeadTablePath string
	BaseURL       string
	HTTPClient    *http.Client
}

type LeadRecord map[string]string

func (c *CasaDadosClient) SearchCompanyByCNPJ(cnpj string) (interface{}, error) {
	lead, ok, err := c.FindLeadByCNPJ(cnpj)
	if err != nil {
		return nil, err
	}
	company := map[string]any{}
	if ok {
		company = CompanyPayload(lead)
	}

	external, externalErr := c.ConsultaCNPJ(cnpj)
	if externalErr == nil && len(external) > 0 {
		company["casadosdados"] = external
		company["sources"] = []string{"local_lead_table", "casa_dos_dados"}
		return company, nil
	}
	if ok {
		company["casadosdados"] = externalStatus(externalErr)
		company["sources"] = []string{"local_lead_table"}
		return company, nil
	}
	if externalErr != nil {
		return nil, externalErr
	}
	return nil, ErrNotFound
}

func (c *CasaDadosClient) SearchCompanyByName(name string) (interface{}, error) {
	leads, err := c.LoadLeads()
	if err != nil {
		return nil, err
	}

	query := normalizeText(name)
	results := make([]map[string]any, 0, len(leads))
	for _, lead := range leads {
		if query != "" && !strings.Contains(normalizeText(usecase.SearchText(lead)), query) {
			continue
		}
		results = append(results, CompanyPayload(lead))
	}

	external, externalErr := c.PesquisaEmpresas(name)
	if externalErr == nil && len(external) > 0 {
		results = append(results, external...)
	}
	return results, nil
}

func (c *CasaDadosClient) ConsultaCNPJ(cnpj string) (map[string]any, error) {
	key := c.apiKey()
	if key == "" {
		return nil, fmt.Errorf("Casa dos Dados API key not configured; set CASA_DADOS_API_KEY or PP_LEADS_CASA_DADOS_API_KEY")
	}
	digits := OnlyDigits(cnpj)
	if len(digits) != 14 {
		return nil, fmt.Errorf("invalid cnpj %q", cnpj)
	}

	url := strings.TrimRight(c.baseURL(), "/") + "/v4/cnpj/" + digits
	var data map[string]any
	if err := c.doJSON(http.MethodGet, url, key, nil, &data); err != nil {
		return nil, err
	}
	data["source"] = "casa_dos_dados"
	data["source_endpoint"] = "/v4/cnpj/{cnpj}"
	return data, nil
}

func (c *CasaDadosClient) PesquisaEmpresas(name string) ([]map[string]any, error) {
	key := c.apiKey()
	if key == "" || strings.TrimSpace(name) == "" {
		return nil, nil
	}

	body := map[string]any{
		"busca_textual": []map[string]any{{
			"texto":         []string{name},
			"tipo_busca":    "exata",
			"razao_social":  true,
			"nome_fantasia": true,
			"nome_socio":    false,
		}},
		"situacao_cadastral": []string{"ATIVA"},
		"limite":             10,
		"pagina":             1,
	}

	var data map[string]any
	url := strings.TrimRight(c.baseURL(), "/") + "/v5/cnpj/pesquisa"
	if err := c.doJSON(http.MethodPost, url, key, body, &data); err != nil {
		return nil, err
	}
	items, _ := data["cnpjs"].([]any)
	results := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if obj, ok := item.(map[string]any); ok {
			obj["source"] = "casa_dos_dados"
			obj["source_endpoint"] = "/v5/cnpj/pesquisa"
			results = append(results, obj)
		}
	}
	return results, nil
}

func (c *CasaDadosClient) doJSON(method, url, apiKey string, body any, out any) error {
	var payload *bytes.Reader
	if body == nil {
		payload = bytes.NewReader(nil)
	} else {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		payload = bytes.NewReader(data)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return err
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := c.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("Casa dos Dados unauthorized: configure valid api-key")
	}
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("Casa dos Dados forbidden: api-key lacks access or balance")
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Casa dos Dados status %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *CasaDadosClient) apiKey() string {
	if c != nil && c.APIKey != "" {
		return c.APIKey
	}
	for _, key := range []string{"CASA_DADOS_API_KEY", "CASADOSDADOS_API_KEY", "PP_LEADS_CASA_DADOS_API_KEY"} {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value
		}
	}
	return ""
}

func (c *CasaDadosClient) baseURL() string {
	if c != nil && c.BaseURL != "" {
		return c.BaseURL
	}
	if value := strings.TrimSpace(os.Getenv("CASA_DADOS_BASE_URL")); value != "" {
		return value
	}
	if value := strings.TrimSpace(os.Getenv("PP_LEADS_CASA_DADOS_BASE_URL")); value != "" {
		return value
	}
	return "https://api.casadosdados.com.br"
}

func externalStatus(err error) map[string]any {
	if err == nil {
		return map[string]any{"status": "not_requested"}
	}
	return map[string]any{"status": "unavailable", "error": err.Error()}
}

func (c *CasaDadosClient) FindLeadByCNPJ(cnpj string) (LeadRecord, bool, error) {
	want := OnlyDigits(cnpj)
	leads, err := c.LoadLeads()
	if err != nil {
		return nil, false, err
	}
	for _, lead := range leads {
		if OnlyDigits(lead["CNPJ"]) == want {
			return lead, true, nil
		}
	}
	return nil, false, nil
}

func (c *CasaDadosClient) LoadLeads() ([]LeadRecord, error) {
	path, err := c.resolveLeadTablePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening local lead table %s: %w", path, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading local lead table %s: %w", path, err)
	}
	if len(records) == 0 {
		return nil, nil
	}

	headers := records[0]
	leads := make([]LeadRecord, 0, len(records)-1)
	for _, record := range records[1:] {
		lead := make(LeadRecord, len(headers))
		for i, header := range headers {
			if i < len(record) {
				lead[header] = record[i]
			}
		}
		leads = append(leads, lead)
	}
	return leads, nil
}

func (c *CasaDadosClient) resolveLeadTablePath() (string, error) {
	if c != nil && c.LeadTablePath != "" {
		return c.LeadTablePath, nil
	}
	if path := os.Getenv("PP_LEADS_LOCAL_TABLE"); path != "" {
		return path, nil
	}
	cfg, err := usecase.LoadFromEnv()
	if err != nil {
		return "", err
	}
	if cfg != nil && strings.TrimSpace(cfg.LeadTablePath) != "" {
		return cfg.ResolvePath(cfg.LeadTablePath), nil
	}
	return "", fmt.Errorf("local lead table not configured; set PP_LEADS_LOCAL_TABLE, %s, or %s", usecase.EnvConfigPath, usecase.EnvConfigDir)
}

func CompanyPayload(lead LeadRecord) map[string]any {
	cfg, _ := usecase.LoadFromEnv()
	city, state := SplitCityUF(usecase.CityUF(lead, cfg))
	payload := map[string]any{
		"name":   usecase.Name(lead, cfg),
		"cnpj":   usecase.CNPJ(lead, cfg),
		"status": "lead-local",
		"address": map[string]string{
			"city":  city,
			"state": state,
		},
		"contacts": map[string]string{
			"email": usecase.Contact(lead, cfg),
			"site":  usecase.Site(lead, cfg),
		},
		"lead_context": usecase.LeadContext(lead),
	}
	if meta := cfg.Metadata(); meta != nil {
		payload["use_case"] = meta
	}
	return payload
}

func SplitCityUF(value string) (string, string) {
	parts := strings.Split(value, "/")
	if len(parts) < 2 {
		return value, ""
	}
	return parts[0], parts[1]
}

func OnlyDigits(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func normalizeText(value string) string {
	return foldLatin(strings.ToLower(strings.TrimSpace(value)))
}

func foldLatin(value string) string {
	replacer := strings.NewReplacer(
		"á", "a", "à", "a", "ã", "a", "â", "a", "ä", "a",
		"é", "e", "è", "e", "ê", "e", "ë", "e",
		"í", "i", "ì", "i", "î", "i", "ï", "i",
		"ó", "o", "ò", "o", "õ", "o", "ô", "o", "ö", "o",
		"ú", "u", "ù", "u", "û", "u", "ü", "u",
		"ç", "c",
	)
	return replacer.Replace(value)
}
