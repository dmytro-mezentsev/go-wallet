package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"net/http"
	"strconv"
	"wallet.com/wallet/wallet/internal/data"
)

type WalletResp struct {
	WalletId string          `json:"walletId"`
	Amount   decimal.Decimal `json:"amount"`
}

type WalletsResp struct {
	Wallets []WalletResp `json:"wallets"`
}

type WalletStorageI interface {
	Save(wallets []data.Wallet) ([]data.Wallet, error)
	Get(walletId string) (data.Wallet, error)
}

type WalletHandler struct {
	WalletStorage WalletStorageI
}

func (wh WalletHandler) PostWalletHandler(w http.ResponseWriter, r *http.Request) {
	var count int
	if strCount := r.FormValue("count"); strCount == "" {
		count = 1
	} else {
		//validate count
		var err error
		count, err = strconv.Atoi(strCount)
		if err != nil {
			http.Error(w, ErrorResponse("invalid count"), http.StatusBadRequest)
			return
		}
		if count > 10 {
			http.Error(w, ErrorResponse("count can't be more then 10"), http.StatusBadRequest)
			return
		}
	}
	var walletDatas = make([]data.Wallet, count)

	for i := 0; i < count; i++ {
		walletDatas[i] = data.Wallet{
			Id:     uuid.NewString(),
			Amount: decimal.Zero,
		}
	}
	wh.WalletStorage.Save(walletDatas)

	var wallets = make([]WalletResp, count)
	for i, w := range walletDatas {
		wallets[i] = WalletResp{
			WalletId: w.Id,
			Amount:   w.Amount,
		}
	}
	response := WalletsResp{
		Wallets: wallets,
	}
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, InternalErrorResponse(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

func (wh WalletHandler) GetWalletHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	walletId := vars["walletId"]

	if walletId == "" {
		http.Error(w, ErrorResponse("walletId is required"), http.StatusBadRequest)
		return
	}
	//validate walletId
	_, err := uuid.Parse(walletId)
	if err != nil {
		http.Error(w, ErrorResponse("invalid walletId"), http.StatusBadRequest)
		return
	}

	walletData, err := wh.WalletStorage.Get(walletId)
	if err != nil {
		http.Error(w, ErrorResponse("wallet not found"), http.StatusNotFound)
		return
	}
	response := WalletResp{
		WalletId: walletData.Id,
		Amount:   walletData.Amount,
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, InternalErrorResponse(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, ErrorResponse("Internal error"), http.StatusInternalServerError)
		return
	}
}
