package main

import (
	"miner_game/company"
	"miner_game/server"
)

func main() {
	company := company.NewCompany()
	HTTPHandlers := server.NewHTTPHandlers(company)
	server := server.NewServer(HTTPHandlers)
	print(1)
	err := server.StartServer()
	if err != nil {
		print(err)
		return
	}
}