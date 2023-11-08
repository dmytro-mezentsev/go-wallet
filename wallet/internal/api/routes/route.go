package routes

import (
	"github.com/gorilla/mux"
	"wallet.com/wallet/wallet/internal/api/handlers"
)

func WalletRoute(r *mux.Router) {
	r.HandleFunc("/wallet", handlers.CreateWalletHandler).
		Methods("POST")
}
