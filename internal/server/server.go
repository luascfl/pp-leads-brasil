package server

import (
	"encoding/json"
	"net/http"
	"pp-leads-brasil/internal/factory"
)

type Server struct {
	Factory *factory.Factory
}

func (s *Server) GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// cnpj := r.PathValue("cnpj") // Go 1.22+
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Not implemented"})
}

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Name string `json:"name"`
		CNAE string `json:"cnae"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	result, err := s.Factory.CasaDadosClient.SearchCompanyByName(reqBody.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) EnrichHandler(w http.ResponseWriter, r *http.Request) {
	// Extract CNPJ (hardcoded for now, router will provide this)
	cnpj := "12345678000199"
	
	result, err := s.Factory.PPClient.RunCompanyGoat(cnpj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
