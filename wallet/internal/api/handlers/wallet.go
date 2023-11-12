package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"wallet.com/wallet/wallet/internal/data"
)

type Wallet struct {
	WalletId string  `json:"walletId"`
	Amount   float64 `json:"amount"`
}

type WalletsResponse struct {
	Wallets []Wallet `json:"wallets"`
}

type WalletStorageI interface {
	Save(wallets []data.Wallet) ([]data.Wallet, error)
}

type WalletHandler struct {
	WalletStorage WalletStorageI
}

func (wh WalletHandler) PostWalletHandler(w http.ResponseWriter, r *http.Request) {
	//validate count
	var count int
	if strCount := r.FormValue("count"); strCount == "" {
		count = 1
	} else {
		var err error
		count, err = strconv.Atoi(strCount)
		if err != nil {
			http.Error(w, "invalid count", http.StatusBadRequest)
			return
		}
		if count > 10 {
			http.Error(w, "count can't be more then 10", http.StatusBadRequest)
			return
		}
	}
	var walletDatas = make([]data.Wallet, count)

	for i := 0; i < count; i++ {
		walletDatas[i] = data.Wallet{
			ID:     uuid.NewString(),
			Amount: 0.0,
		}
	}
	wh.WalletStorage.Save(walletDatas)

	var wallets = make([]Wallet, count)
	for i, w := range walletDatas {
		wallets[i] = Wallet{
			WalletId: w.ID,
			Amount:   w.Amount,
		}
	}
	response := WalletsResponse{
		Wallets: wallets,
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
