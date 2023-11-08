package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

type CreateWalletRequest struct {
	UserId string `json:"userId"`
}
type CreateWalletResponse struct {
	UserId   string `json:"userId"`
	WalletId string `json:"walletId"`
}

func CreateWalletHandler(w http.ResponseWriter, r *http.Request) {

	var createWalletRequest CreateWalletRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&createWalletRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createWalletResponse := CreateWalletResponse{
		UserId:   createWalletRequest.UserId,
		WalletId: uuid.NewString(),
	}

	responseJSON, err := json.Marshal(createWalletResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
