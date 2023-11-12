package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet.com/wallet/wallet/internal/data"
)

// init mock
type WalletStorageMock struct {
}

func (wst WalletStorageMock) Save(wallets []data.Wallet) ([]data.Wallet, error) {
	return wallets, nil
}

func TestCreateWalletHandlerWithOneWallet(t *testing.T) {

	walletHandler := WalletHandler{WalletStorage: WalletStorageMock{}}

	request, err := http.NewRequest("POST", "/wallet", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	responseRecorder := httptest.NewRecorder()
	walletHandler.PostWalletHandler(responseRecorder, request)

	// Check the response status code
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response WalletsResponse
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	if len(response.Wallets) != 1 {
		t.Errorf("Invalid count of wallets: %v", len(response.Wallets))
	}
	// Check walletId is a valid UUID, else panic
	uuid.MustParse(response.Wallets[0].WalletId)

	if response.Wallets[0].Amount != 0.0 {
		t.Errorf("Invalid amount: %v", response.Wallets[0].Amount)
	}
}

func TestCreateWalletHandler(t *testing.T) {

	walletHandler := WalletHandler{WalletStorage: WalletStorageMock{}}

	request, err := http.NewRequest("POST", "/wallet?count=2", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	responseRecorder := httptest.NewRecorder()
	walletHandler.PostWalletHandler(responseRecorder, request)

	// Check the response status code
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response WalletsResponse
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	if len(response.Wallets) != 2 {
		t.Errorf("Invalid count of wallets: %v", len(response.Wallets))
	}

}
