package main

import (
	"log"
	"net/http"
	"pp-leads-brasil/internal/client/casadados"
	"pp-leads-brasil/internal/client/pp"
	"pp-leads-brasil/internal/factory"
	"pp-leads-brasil/internal/server"
)

func main() {
	ppClient := &pp.PPClient{}
	cdClient := &casadados.CasaDadosClient{}

	f := factory.NewFactory(ppClient, cdClient)
	s := &server.Server{Factory: f}

	mux := http.NewServeMux()

	// Register routes (Go 1.22+ routing)
	mux.HandleFunc("GET /v1/company/{cnpj}", s.GetCompanyHandler)
	mux.HandleFunc("POST /v1/company-goat/{cnpj}", s.CompanyGoatHandler)
	mux.HandleFunc("POST /v1/search", s.SearchHandler)
	mux.HandleFunc("POST /v1/enrich/{cnpj}", s.EnrichHandler)
	mux.HandleFunc("POST /v1/contact-goat", s.ContactGoatHandler)
	mux.HandleFunc("POST /v1/operations/plan", s.OperationPlanHandler)
	mux.HandleFunc("POST /v1/operations/{plan_id}/apply", s.OperationApplyHandler)

	log.Println("Backend API running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
