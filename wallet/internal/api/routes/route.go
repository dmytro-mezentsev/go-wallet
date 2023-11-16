package routes

import (
	"github.com/gorilla/mux"
	"wallet.com/wallet/wallet/internal/api/handlers"
)

func WalletRoute(r *mux.Router, walletHandler *handlers.WalletHandler) {
	r.HandleFunc("/wallet", walletHandler.PostWalletHandler).
		Methods("POST").
		Queries("count", "{[0-9]*?}")

	r.HandleFunc("/wallet/{walletId}", walletHandler.GetWalletHandler).
		Methods("GET")
}

func TransactionRoute(r *mux.Router, transactionHandler *handlers.TransactionHandler) {
	r.HandleFunc("/transactions", transactionHandler.PostTransactionHandler).
		Methods("POST")
	r.HandleFunc("/transactions/{transactionId}", transactionHandler.GetTransactionHandler).
		Methods("GET")
}
