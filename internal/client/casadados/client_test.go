package casadados_test

import (
	"pp-leads-brasil/internal/client/casadados"
	"testing"
)

func TestSearchCompanyByName(t *testing.T) {
	client := &casadados.CasaDadosClient{}
	query := "Empresa Teste"

	result, err := client.SearchCompanyByName(query)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected result to be map[string]interface{}, got %T", result)
	}

	if resultMap["status"] != "success" {
		t.Errorf("Expected status 'success', got %v", resultMap["status"])
	}

	results, ok := resultMap["results"].([]map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'results' to be []map[string]interface{}, got %T", resultMap["results"])
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	if results[0]["name"] != "Empresa Teste Matriz" {
		t.Errorf("Expected first result name to be 'Empresa Teste Matriz', got %v", results[0]["name"])
	}
}
