package routes

import (
	"github.com/gorilla/mux"
	"wallet.com/wallet/wallet/internal/api/handlers"
)

func HelloWorldRoute(r *mux.Router) {
	r.HandleFunc("/hello", handlers.HelloWorldHandler)
}
