package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/shopspring/decimal"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet.com/wallet/wallet/internal/data"
)

type MockTransactionStorage struct {
	err error
}

func (mts MockTransactionStorage) Save(transaction data.Transaction) (data.Transaction, error) {
	if mts.err != nil {
		return data.Transaction{}, mts.err
	}
	return transaction, nil
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
		TransactionStorage: mockTransactionStorage,
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
		TransactionStorage: mockTransactionStorage,
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
		TransactionStorage: mockTransactionStorage,
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
