package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"wallet.com/wallet/wallet/internal/data"
)

type CreateWalletResponse struct {
	WalletId string  `json:"walletId"`
	Amount   float64 `json:"amount"`
}

type WalletStorageI interface {
	Save(wallet data.Wallet) (data.Wallet, error)
}

type WalletHandler struct {
	WalletStorage WalletStorageI
}

func (wh WalletHandler) PostWalletHandler(w http.ResponseWriter, r *http.Request) {

	wallet, err := wh.WalletStorage.Save(data.Wallet{
		ID:     uuid.NewString(),
		Amount: 0.0,
	})
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	response := CreateWalletResponse{
		WalletId: wallet.ID,
		Amount:   wallet.Amount,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
