package factory

import (
	"pp-leads-brasil/internal/client/pp"
	"pp-leads-brasil/internal/client/casadados"
)

type Factory struct {
	PPClient         pp.Client
	CasaDadosClient  casadados.Client
}

func NewFactory(ppClient pp.Client, casaDadosClient casadados.Client) *Factory {
	return &Factory{
		PPClient:        ppClient,
		CasaDadosClient: casaDadosClient,
	}
}
