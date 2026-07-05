package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"pp-leads-brasil/internal/client/casadados"
	"pp-leads-brasil/internal/factory"
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
