package main

import (
	"flag"
	"log"

	"github.com/sreway/gophermart/config"
	"github.com/sreway/gophermart/internal/app"
)

func main() {
	flag.StringVar(
		&config.DefaultHTTPAddress, "a", config.DefaultHTTPAddress, "Server address: host:port")
	flag.StringVar(
		&config.DefaultPostgresDSN, "d", config.DefaultPostgresDSN, "PostgreSQL data source name")
	flag.StringVar(
		&config.DefaultAccrualAddress, "r", config.DefaultAccrualAddress, "Accrual system address")
	flag.Parse()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	app.Migrate(cfg)
	app.Run(cfg)
}
