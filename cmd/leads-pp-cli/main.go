package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("Leads PP CLI - Ferramenta de controle da fábrica de scrapers")
		fmt.Println("Uso: pp-leads [comando]")
		return
	}
	fmt.Println("Hello World - Leads PP CLI pronto para operação")
}
