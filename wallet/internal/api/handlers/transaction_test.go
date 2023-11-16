package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"wallet.com/wallet/wallet/internal/data"
)

type MockTransactionStorage struct {
	savedTransactions map[string]data.Transaction
	err               error
}

func (mts *MockTransactionStorage) Save(transaction data.Transaction) (data.Transaction, error) {
	if mts.err != nil {
		return data.Transaction{}, mts.err
	}
	if mts.savedTransactions == nil {
		mts.savedTransactions = make(map[string]data.Transaction)
	}

	mts.savedTransactions[transaction.Id] = transaction
	return transaction, nil
}
func (mts *MockTransactionStorage) GetById(transactionId string) (data.Transaction, error) {
	if mts.err != nil {
		return data.Transaction{}, mts.err
	}
	result, ok := mts.savedTransactions[transactionId]
	if !ok {
		return data.Transaction{}, errors.New("Not found")
	} else {
		return result, nil
	}
}

type MockWalletStorage struct {
	amount decimal.Decimal
}

func (mws MockWalletStorage) Get(walletId string) (data.Wallet, error) {
	return data.Wallet{walletId, mws.amount}, nil
}
func (mws MockWalletStorage) Save(wallets []data.Wallet) ([]data.Wallet, error) {
	return wallets, nil
}

func TestPostTransactionHandler(t *testing.T) {
	requestBody := TransactionReq{
		WalletId:              "some_wallet_id",
		Amount:                decimal.NewFromFloat(50.0),
		TransactionType:       data.Deposit,
		FromPaymentSystem:     "Bank",
		FromPaymentIdentifier: "bank_account_123",
		ToPaymentSystem:       "Wallet",
		ToPaymentIdentifier:   "wallet_id_456",
		Currency:              "USD",
		Description:           "Test deposit",
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal("Error encoding JSON:", err)
	}

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// Create mocks
	mockTransactionStorage := MockTransactionStorage{}
	mockWalletStorage := MockWalletStorage{decimal.NewFromFloat(100)}

	handler := TransactionHandler{
		TransactionStorage: &mockTransactionStorage,
		WalletStorage:      mockWalletStorage,
	}

	rr := httptest.NewRecorder()
	handler.PostTransactionHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, rr.Code)
	}

	var response TransactionResp
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Error decoding JSON:", err)
	}

	if response.Id == "" {
		t.Error("Expected non-empty transaction ID")
	}
	//check balance
	expectedBalance := decimal.NewFromFloat(150.0)
	if !response.Balance.Equal(expectedBalance) {
		t.Errorf("Expected balance %v, got %v", expectedBalance, response.Balance)
	}
	if response.CreatedAt.IsZero() {
		t.Error("Expected non-zero created timestamp")
	}
}

func TestPostTransactionHandlerInsufficientFunds(t *testing.T) {
	requestBody := TransactionReq{
		WalletId:              "some_wallet_id",
		UserId:                "some_user_id",
		Amount:                decimal.NewFromFloat(150.0),
		TransactionType:       data.Withdraw,
		FromPaymentSystem:     "Bank",
		FromPaymentIdentifier: "bank_account_123",
		ToPaymentSystem:       "Wallet",
		ToPaymentIdentifier:   "wallet_id_456",
		Currency:              "USD",
		Description:           "Test deposit",
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal("Error encoding JSON:", err)
	}

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// Create mocks
	mockTransactionStorage := MockTransactionStorage{}
	mockWalletStorage := MockWalletStorage{decimal.NewFromFloat(100)}

	handler := TransactionHandler{
		TransactionStorage: &mockTransactionStorage,
		WalletStorage:      mockWalletStorage,
	}

	rr := httptest.NewRecorder()
	handler.PostTransactionHandler(rr, req)

	if rr.Code != http.StatusPaymentRequired {
		t.Errorf("Expected status %v, got %v", http.StatusPaymentRequired, rr.Code)
	}

	expectedResponseBody := ErrorResponse("Insufficient funds")
	responseBody := rr.Body.String()

	if strings.TrimSpace(responseBody) != strings.TrimSpace(expectedResponseBody) {
		t.Errorf("Expected response body %v, got %v", expectedResponseBody, responseBody)
	}
}

func TestPostTransactionHandlerConflictError(t *testing.T) {
	requestBody := TransactionReq{
		WalletId:              "some_wallet_id",
		Amount:                decimal.NewFromFloat(50.0),
		TransactionType:       data.Withdraw,
		FromPaymentSystem:     "Bank",
		FromPaymentIdentifier: "bank_account_123",
		ToPaymentSystem:       "Wallet",
		ToPaymentIdentifier:   "wallet_id_456",
		Currency:              "USD",
		Description:           "Test deposit",
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal("Error encoding JSON:", err)
	}

	req, err := http.NewRequest("POST", "/transactions", bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		t.Fatal("Error creating request:", err)
	}

	// Create mocks
	mockTransactionStorage := MockTransactionStorage{err: data.BalanceWasChangedError("during transaction balance was changed")}
	mockWalletStorage := MockWalletStorage{decimal.NewFromFloat(100)}

	handler := TransactionHandler{
		TransactionStorage: &mockTransactionStorage,
		WalletStorage:      mockWalletStorage,
	}

	rr := httptest.NewRecorder()
	handler.PostTransactionHandler(rr, req)

	if rr.Code != http.StatusConflict {
		t.Errorf("Expected status %v, got %v", http.StatusConflict, rr.Code)
	}

	expectedResponseBody := ErrorResponse("Balance was changed during transaction, please retry")
	responseBody := rr.Body.String()

	if strings.TrimSpace(responseBody) != strings.TrimSpace(expectedResponseBody) {
		t.Errorf("Expected response body %v, got %v", expectedResponseBody, responseBody)
	}
}

func TestPostTransactionHandlerGetById(t *testing.T) {
	mockTransactionStorage := MockTransactionStorage{}
	transaction := data.Transaction{
		Id:                    uuid.NewString(),
		UserId:                "user_id_1",
		WalletId:              "wallet_id_1",
		Amount:                decimal.NewFromFloat(50.0),
		TransactionType:       data.Deposit,
		AmountBefore:          decimal.NewFromFloat(0.0),
		AmountAfter:           decimal.NewFromFloat(50.0),
		FromPaymentSystem:     "Bank",
		FromPaymentIdentifier: "bank_account_123",
		ToPaymentSystem:       "Wallet",
		ToPaymentIdentifier:   "wallet_id_456",
		Currency:              "USD",
		Description:           "Test deposit",
		CreatedAt:             time.Now(),
	}

	mockTransactionStorage.Save(transaction)
	handler := TransactionHandler{
		TransactionStorage: &mockTransactionStorage,
	}

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/transactions/"+transaction.Id, nil)
	vars := map[string]string{
		"transactionId": transaction.Id,
	}
	req = mux.SetURLVars(req, vars)
	handler.GetTransactionHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, rr.Code)
	}
	var response TransactionFullResp
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.Id != transaction.Id {
		t.Errorf("Expected transaction id %v, got %v", transaction.Id, response.Id)
	}
	if response.UserId != transaction.UserId {
		t.Errorf("Expected user id %v, got %v", transaction.UserId, response.UserId)
	}
}
