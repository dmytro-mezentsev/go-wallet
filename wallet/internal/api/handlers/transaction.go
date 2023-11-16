package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"time"
	"wallet.com/wallet/wallet/internal/data"
)

type TransactionReq struct {
	WalletId              string               `json:"walletId"`
	UserId                string               `json:"userId"`
	Amount                decimal.Decimal      `json:"amount"`
	TransactionType       data.TransactionType `json:"transactionType"`
	FromPaymentSystem     string               `json:"fromPaymentSystem"`
	FromPaymentIdentifier string               `json:"fromPaymentIdentifier"`
	ToPaymentSystem       string               `json:"toPaymentSystem"`
	ToPaymentIdentifier   string               `json:"toPaymentIdentifier"`
	Currency              string               `json:"currency"`
	Description           string               `json:"description"`
}
type TransactionResp struct {
	Id        string          `json:"id"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"createdAt"`
}
type TransactionFullResp struct {
	Id                    string               `json:"id"`
	UserId                string               `json:"userId"`
	WalletId              string               `json:"walletId"`
	Amount                decimal.Decimal      `json:"amount"`
	TransactionType       data.TransactionType `json:"transactionType"`
	AmountBefore          decimal.Decimal      `json:"amountBefore"`
	AmountAfter           decimal.Decimal      `json:"amountAfter"`
	FromPaymentSystem     string               `json:"fromPaymentSystem"`
	FromPaymentIdentifier string               `json:"fromPaymentIdentifier"`
	ToPaymentSystem       string               `json:"toPaymentSystem"`
	ToPaymentIdentifier   string               `json:"toPaymentIdentifier"`
	Currency              string               `json:"currency"`
	Description           string               `json:"description"`
	CreatedAt             time.Time            `json:"createdAt"`
}

type TransactionStorageI interface {
	Save(transaction data.Transaction) (data.Transaction, error)
	GetById(transactionId string) (data.Transaction, error)
}

type TransactionHandler struct {
	TransactionStorage TransactionStorageI
	WalletStorage      WalletStorageI
}

func (th *TransactionHandler) PostTransactionHandler(w http.ResponseWriter, r *http.Request) {
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
	if transactionReq.TransactionType == data.Withdraw && wallet.Amount.LessThan(transactionReq.Amount) {
		http.Error(w, ErrorResponse("Insufficient funds"), http.StatusPaymentRequired)
		return
	}

	amountBefore := wallet.Amount.RoundFloor(4)
	var amountAfter decimal.Decimal
	if transactionReq.TransactionType == data.Withdraw {
		amountAfter = wallet.Amount.Sub(transactionReq.Amount)
	} else {
		amountAfter = wallet.Amount.Add(transactionReq.Amount)
	}
	transaction := data.Transaction{
		Id:                    uuid.NewString(),
		UserId:                transactionReq.UserId,
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
		Balance:   result.AmountAfter.RoundFloor(2),
		CreatedAt: result.CreatedAt,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, InternalErrorResponse(), http.StatusInternalServerError)
		return
	}

	return
}

func (th *TransactionHandler) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	transactionId := vars["transactionId"]

	if transactionId == "" {
		http.Error(w, ErrorResponse("transactionId is required"), http.StatusBadRequest)
		return
	}

	transaction, err := th.TransactionStorage.GetById(transactionId)
	if err != nil {
		http.Error(w, ErrorResponse("transaction not found"), http.StatusNotFound)
		return
	}
	response := TransactionFullResp{
		Id:                    transaction.Id,
		UserId:                transaction.UserId,
		WalletId:              transaction.WalletId,
		Amount:                transaction.Amount,
		TransactionType:       transaction.TransactionType,
		AmountBefore:          transaction.AmountBefore,
		AmountAfter:           transaction.AmountAfter,
		FromPaymentSystem:     transaction.FromPaymentSystem,
		FromPaymentIdentifier: transaction.FromPaymentIdentifier,
		ToPaymentSystem:       transaction.ToPaymentSystem,
		ToPaymentIdentifier:   transaction.ToPaymentIdentifier,
		Currency:              transaction.Currency,
		Description:           transaction.Description,
		CreatedAt:             transaction.CreatedAt,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, InternalErrorResponse(), http.StatusInternalServerError)
		return
	}
	return
}
