package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"

	"pp-leads-brasil/internal/client/casadados"
	"pp-leads-brasil/internal/factory"
	"pp-leads-brasil/internal/operation"
	"pp-leads-brasil/internal/usecase"
)

type Server struct {
	Factory *factory.Factory
}

func (s *Server) GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.Factory.CasaDadosClient.SearchCompanyByCNPJ(r.PathValue("cnpj"))
	if err != nil {
		writeClientError(w, err)
		return
	}
	writeJSON(w, result)
}

func (s *Server) CompanyGoatHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.Factory.PPClient.RunCompanyGoat(r.PathValue("cnpj"))
	if err != nil {
		writeClientError(w, err)
		return
	}
	writeJSON(w, result)
}

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name string `json:"name"`
		CNAE string `json:"cnae"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	query := reqBody.Name
	if query == "" {
		query = reqBody.CNAE
	}
	result, err := s.Factory.CasaDadosClient.SearchCompanyByName(query)
	if err != nil {
		writeClientError(w, err)
		return
	}
	writeJSON(w, result)
}

func (s *Server) EnrichHandler(w http.ResponseWriter, r *http.Request) {
	result, err := s.Factory.PPClient.RunEnrich(r.PathValue("cnpj"))
	if err != nil {
		writeClientError(w, err)
		return
	}
	writeJSON(w, result)
}

func (s *Server) ContactGoatHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		URL     string `json:"url"`
		Email   string `json:"email"`
		Contact string `json:"contact"`
		Company string `json:"company"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	contact := reqBody.URL
	if contact == "" {
		contact = reqBody.Email
	}
	if contact == "" {
		contact = reqBody.Contact
	}
	if contact == "" {
		contact = reqBody.Company
	}
	if contact == "" {
		http.Error(w, "url, email, contact, or company is required", http.StatusBadRequest)
		return
	}

	result, err := s.Factory.PPClient.RunContactGoat(contact)
	if err != nil {
		writeClientError(w, err)
		return
	}
	writeJSON(w, result)
}

func (s *Server) OperationPlanHandler(w http.ResponseWriter, r *http.Request) {
	var input operation.PlanInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	service, err := profileOperationService(input.Profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	plan, err := service.Plan(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, plan)
}

func (s *Server) OperationApplyHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Approved bool `json:"approved"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	service, err := profileOperationService("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := service.Apply(r.Context(), r.PathValue("plan_id"), request.Approved)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, operation.ErrApprovalRequired) || errors.Is(err, operation.ErrPlanExpired) {
			status = http.StatusConflict
		}
		http.Error(w, err.Error(), status)
		return
	}
	writeJSON(w, result)
}

func profileOperationService(requestProfile string) (operation.Service, error) {
	cfg, err := usecase.LoadFromEnv()
	if err != nil {
		return operation.Service{}, err
	}
	if cfg == nil || cfg.Name == "" {
		return operation.Service{}, fmt.Errorf("an explicit private profile is required for external operations")
	}
	if requestProfile != "" && requestProfile != cfg.Name {
		return operation.Service{}, fmt.Errorf("requested profile does not match configured profile")
	}
	outputDir := cfg.OutputDir
	if outputDir == "" {
		outputDir = "outputs"
	}
	store, err := operation.NewFileStore(filepath.Join(cfg.ResolvePath(outputDir), "operations"))
	if err != nil {
		return operation.Service{}, err
	}
	return operation.Service{Store: store}, nil
}

func writeClientError(w http.ResponseWriter, err error) {
	if errors.Is(err, casadados.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}
