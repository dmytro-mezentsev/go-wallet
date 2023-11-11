package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"wallet.com/wallet/wallet/internal/api/middleware"
	"wallet.com/wallet/wallet/internal/api/routes"
	"wallet.com/wallet/wallet/internal/config"
	"wallet.com/wallet/wallet/internal/db"
)

func main() {
	log.Println("Starting server on port 8000...")

	config := config.GetConfig()
	dbConnection := db.DBConnection(config.Db)
	db.MigrateSchemas(*dbConnection, config.Db.DBName)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	routes.WalletRoute(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//keep alive
	log.Fatal(srv.ListenAndServe())
}
