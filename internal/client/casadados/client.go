package casadados

// Client defines the interface for interacting with Casa dos Dados API
type Client interface {
	SearchCompanyByCNPJ(cnpj string) (interface{}, error)
	SearchCompanyByName(name string) (interface{}, error)
}

type CasaDadosClient struct {
	APIKey string
}

func (c *CasaDadosClient) SearchCompanyByCNPJ(cnpj string) (interface{}, error) {
	// Logic to query Casa dos Dados API
	return map[string]interface{}{
		"status": "success",
		"data": map[string]string{
			"cnpj":   cnpj,
			"name":   "Empresa Exemplo LTDA",
			"status": "ATIVA",
		},
	}, nil
}

func (c *CasaDadosClient) SearchCompanyByName(name string) (interface{}, error) {
	// Logic to query Casa dos Dados API
	return map[string]interface{}{
		"status": "success",
		"results": []map[string]interface{}{
			{
				"cnpj":   "12345678000199",
				"name":   name + " Matriz",
				"status": "ATIVA",
				"cnae":   "6204-0/00",
			},
			{
				"cnpj":   "12345678000270",
				"name":   name + " Filial",
				"status": "ATIVA",
				"cnae":   "6204-0/00",
			},
		},
	}, nil
}
