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

func (wst WalletStorageMock) Save(wallet data.Wallet) (data.Wallet, error) {
	return wallet, nil
}

func TestCreateWalletHandler(t *testing.T) {

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

	var response CreateWalletResponse
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	// Check walletId is a valid UUID, else panic
	uuid.MustParse(response.WalletId)

	if response.Amount != 0.0 {
		t.Errorf("Invalid amount: %v", response.Amount)
	}
}
