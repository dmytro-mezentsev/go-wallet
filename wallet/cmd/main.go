package main

import (
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"time"
	"wallet.com/wallet/wallet/internal/api/handlers"
	"wallet.com/wallet/wallet/internal/api/middleware"
	"wallet.com/wallet/wallet/internal/api/routes"
	"wallet.com/wallet/wallet/internal/config"
	"wallet.com/wallet/wallet/internal/data"
	"wallet.com/wallet/wallet/internal/db"
)

func main() {
	log.Println("Starting server on port 8000...")

	decimal.MarshalJSONWithoutQuotes = true
	config := config.GetConfig()
	dbConnection := db.DBConnection(config.Db)
	db.MigrateSchemas(dbConnection, config.Db.DBName)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	walletStorage := data.WalletStorage{DB: dbConnection}
	transactionStorage := data.TransactionStorage{DB: dbConnection}

	walletHandler := handlers.WalletHandler{WalletStorage: &walletStorage}
	transactionHandler := handlers.TransactionHandler{TransactionStorage: &transactionStorage, WalletStorage: &walletStorage}

	routes.WalletRoute(r, &walletHandler)
	routes.TransactionRoute(r, &transactionHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//keep alive
	log.Fatal(srv.ListenAndServe())
}
