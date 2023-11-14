package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"math"
	"net/http"
	"time"
	"wallet.com/wallet/wallet/internal/data"
)

type TransactionReq struct {
	WalletId              string               `json:"walletId"`
	Amount                float64              `json:"amount"`
	TransactionType       data.TransactionType `json:"transactionType"`
	FromPaymentSystem     string               `json:"fromPaymentSystem"`
	FromPaymentIdentifier string               `json:"fromPaymentIdentifier"`
	ToPaymentSystem       string               `json:"toPaymentSystem"`
	ToPaymentIdentifier   string               `json:"toPaymentIdentifier"`
	Currency              string               `json:"currency"`
	Description           string               `json:"description"`
}
type TransactionResp struct {
	Id        string    `json:"id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type TransactionStorageI interface {
	Save(transaction data.Transaction) (data.Transaction, error)
}

type TransactionHandler struct {
	TransactionStorage TransactionStorageI
	WalletStorage      WalletStorageI
}

func (th TransactionHandler) PostTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var transactionReq TransactionReq
	if err := json.NewDecoder(r.Body).Decode(&transactionReq); err != nil {
		http.Error(w, ErrorResponse("invalid json"), http.StatusBadRequest)
		return
	}
	wallet, err := th.WalletStorage.Get(transactionReq.WalletId)
	if err != nil {
		http.Error(w, ErrorResponse("wallet not found"), http.StatusNotFound)
		return
	}
	if transactionReq.TransactionType == data.Withdraw && wallet.Amount < transactionReq.Amount {
		http.Error(w, ErrorResponse("Insufficient funds"), http.StatusPaymentRequired)
		return
	}

	amountBefore := wallet.Amount
	var amountAfter float64
	if transactionReq.TransactionType == data.Withdraw {
		amountAfter = wallet.Amount - transactionReq.Amount
	} else {
		amountAfter = wallet.Amount + transactionReq.Amount
	}
	transaction := data.Transaction{
		Id:                    uuid.NewString(),
		WalletId:              transactionReq.WalletId,
		Amount:                transactionReq.Amount,
		TransactionType:       transactionReq.TransactionType,
		AmountBefore:          amountBefore,
		AmountAfter:           amountAfter,
		FromPaymentSystem:     transactionReq.FromPaymentSystem,
		FromPaymentIdentifier: transactionReq.FromPaymentIdentifier,
		ToPaymentSystem:       transactionReq.ToPaymentSystem,
		ToPaymentIdentifier:   transactionReq.ToPaymentIdentifier,
		Currency:              transactionReq.Currency,
		Description:           transactionReq.Description,
		CreatedAt:             time.Now(),
	}

	result, err := th.TransactionStorage.Save(transaction)
	if err != nil {
		log.Println("transaction error:", err)
	}
	switch err.(type) {
	case data.BalanceWasChangedError:
		http.Error(w, ErrorResponse("Balance was changed during transaction, please retry"), http.StatusConflict)
		return
	case error:
		http.Error(w, ErrorResponse("Transaction went wrong"), http.StatusConflict)
		return
	}

	response := TransactionResp{
		Id:        result.Id,
		Balance:   math.Floor(result.AmountAfter*100) / 100,
		CreatedAt: result.CreatedAt,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, InternalErrorResponse(), http.StatusInternalServerError)
		return
	}

	return
}
