package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet.com/wallet/wallet/internal/data"
)

type WalletStorageMock struct {
	mockedAmount float64
}

func (wst WalletStorageMock) Save(wallets []data.Wallet) ([]data.Wallet, error) {
	return wallets, nil
}
func (wst WalletStorageMock) Get(walletId string) (data.Wallet, error) {
	return data.Wallet{walletId, wst.mockedAmount}, nil
}

func TestPostWalletHandlerWithOneWallet(t *testing.T) {

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

	var response WalletsResp
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

func TestPostWalletHandlerWithTwoWallets(t *testing.T) {

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

	var response WalletsResp
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	if len(response.Wallets) != 2 {
		t.Errorf("Invalid count of wallets: %v", len(response.Wallets))
	}

}

func TestGetWalletHandler(t *testing.T) {
	amount := 100.0
	walletHandler := WalletHandler{WalletStorage: WalletStorageMock{mockedAmount: amount}}
	walletId := uuid.NewString()

	request, err := http.NewRequest("GET", "/wallet/"+walletId, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	//set path variable
	vars := map[string]string{
		"walletId": walletId,
	}
	request = mux.SetURLVars(request, vars)

	responseRecorder := httptest.NewRecorder()
	walletHandler.GetWalletHandler(responseRecorder, request)

	// Check the response status code
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response WalletResp
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}
	expectedResponse := WalletResp{walletId, amount}

	if response != expectedResponse {
		t.Errorf("Invalid response: %v", response)
	}

}
func TestGetWalletHandlerInvalidWalletId(t *testing.T) {
	amount := 100.0
	walletHandler := WalletHandler{WalletStorage: WalletStorageMock{mockedAmount: amount}}
	walletId := "invalid"

	request, err := http.NewRequest("GET", "/wallet/"+walletId, nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	//set path variable
	vars := map[string]string{
		"walletId": walletId,
	}
	request = mux.SetURLVars(request, vars)

	responseRecorder := httptest.NewRecorder()
	walletHandler.GetWalletHandler(responseRecorder, request)

	// Check the response status code
	if status := responseRecorder.Code; status != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	response := string(bytes.TrimSpace(responseRecorder.Body.Bytes()))
	expectedResponse := strings.TrimSpace(ErrorResponse("invalid walletId"))

	if response != expectedResponse {
		t.Errorf("Invalid response: %v", response)
	}
}
