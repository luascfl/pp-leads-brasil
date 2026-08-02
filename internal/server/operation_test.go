package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"pp-leads-brasil/internal/operation"
	"pp-leads-brasil/internal/usecase"
)

func configuredOperationServer(t *testing.T) *Server {
	t.Helper()
	profileDir := t.TempDir()
	configPath := filepath.Join(profileDir, "profile.json")
	config := `{"name":"test-profile","output_dir":"operations-output"}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(usecase.EnvConfigPath, configPath)
	t.Setenv(usecase.EnvConfigDir, "")
	return &Server{}
}

func operationPlanBody() []byte {
	return []byte(`{
		"profile":"test-profile",
		"kind":"crm-update",
		"operator":"test-user",
		"targets":[{
			"key":"company:1:site",
			"system":"crm",
			"record_id":"1",
			"field":"site",
			"before":"",
			"after":"https://example.test",
			"evidence":"https://example.test/about"
		}]
	}`)
}

func TestOperationPlanRequiresConfiguredMatchingProfile(t *testing.T) {
	server := configuredOperationServer(t)
	req := httptest.NewRequest(http.MethodPost, "/v1/operations/plan", bytes.NewReader(operationPlanBody()))
	res := httptest.NewRecorder()
	server.OperationPlanHandler(res, req)
	if res.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", res.Code, res.Body.String())
	}
	var plan operation.Plan
	if err := json.NewDecoder(res.Body).Decode(&plan); err != nil {
		t.Fatal(err)
	}
	if plan.ID == "" || plan.Digest == "" {
		t.Fatalf("plan = %#v", plan)
	}

	body := bytes.Replace(operationPlanBody(), []byte("test-profile"), []byte("other-profile"), 1)
	req = httptest.NewRequest(http.MethodPost, "/v1/operations/plan", bytes.NewReader(body))
	res = httptest.NewRecorder()
	server.OperationPlanHandler(res, req)
	if res.Code != http.StatusBadRequest {
		t.Fatalf("mismatch status = %d", res.Code)
	}
}

func TestOperationApplyRequiresApprovalAndNeverUsesAnEmbeddedAdapter(t *testing.T) {
	server := configuredOperationServer(t)
	planReq := httptest.NewRequest(http.MethodPost, "/v1/operations/plan", bytes.NewReader(operationPlanBody()))
	planRes := httptest.NewRecorder()
	server.OperationPlanHandler(planRes, planReq)
	var plan operation.Plan
	if err := json.NewDecoder(planRes.Body).Decode(&plan); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/operations/"+plan.ID+"/apply", bytes.NewBufferString(`{"approved":false}`))
	req.SetPathValue("plan_id", plan.ID)
	res := httptest.NewRecorder()
	server.OperationApplyHandler(res, req)
	if res.Code != http.StatusConflict {
		t.Fatalf("unapproved status = %d", res.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/v1/operations/"+plan.ID+"/apply", bytes.NewBufferString(`{"approved":true}`))
	req.SetPathValue("plan_id", plan.ID)
	res = httptest.NewRecorder()
	server.OperationApplyHandler(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("approved status = %d, body = %s", res.Code, res.Body.String())
	}
	var result operation.Result
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result.Receipts) != 1 || result.Receipts[0].Status != "failed" {
		t.Fatalf("receipts = %#v", result.Receipts)
	}
}

func TestOperationApplyUsesExplicitProfileAdapterCommand(t *testing.T) {
	profileDir := t.TempDir()
	scriptPath := filepath.Join(profileDir, "adapter.sh")
	if err := os.WriteFile(scriptPath, []byte("#!/bin/sh\ncat >/dev/null\n"), 0o700); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(profileDir, "profile.json")
	config := `{"name":"test-profile","output_dir":"operations-output","operation_adapter_command":["adapter.sh"]}`
	if err := os.WriteFile(configPath, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv(usecase.EnvConfigPath, configPath)
	t.Setenv(usecase.EnvConfigDir, "")
	server := &Server{}

	planReq := httptest.NewRequest(http.MethodPost, "/v1/operations/plan", bytes.NewReader(operationPlanBody()))
	planRes := httptest.NewRecorder()
	server.OperationPlanHandler(planRes, planReq)
	if planRes.Code != http.StatusCreated {
		t.Fatalf("plan status = %d, body = %s", planRes.Code, planRes.Body.String())
	}
	var plan operation.Plan
	if err := json.NewDecoder(planRes.Body).Decode(&plan); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/v1/operations/"+plan.ID+"/apply", bytes.NewBufferString(`{"approved":true}`))
	req.SetPathValue("plan_id", plan.ID)
	res := httptest.NewRecorder()
	server.OperationApplyHandler(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("apply status = %d, body = %s", res.Code, res.Body.String())
	}
	var result operation.Result
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if len(result.Receipts) != 1 || result.Receipts[0].Status != "applied" {
		t.Fatalf("receipts = %#v", result.Receipts)
	}
}
