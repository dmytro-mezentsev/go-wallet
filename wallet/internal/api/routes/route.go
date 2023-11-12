package routes

import (
	"github.com/gorilla/mux"
	"wallet.com/wallet/wallet/internal/api/handlers"
)

func WalletRoute(r *mux.Router, walletHandler handlers.WalletHandler) {
	r.HandleFunc("/wallet", walletHandler.PostWalletHandler).Methods("POST")
}
