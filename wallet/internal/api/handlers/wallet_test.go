package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateWalletHandler(t *testing.T) {
	requestBody := CreateWalletRequest{
		UserId: "testUser",
	}
	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Error marshaling JSON request: %v", err)
	}

	request, err := http.NewRequest("POST", "/wallet", bytes.NewBuffer(requestJSON))
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}

	responseRecorder := httptest.NewRecorder()
	CreateWalletHandler(responseRecorder, request)

	// Check the response status code
	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response CreateWalletResponse
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON response: %v", err)
	}

	// Check the response fields
	if response.UserId != requestBody.UserId {
		t.Errorf("Handler returned unexpected user ID: got %v want %v", response.UserId, requestBody.UserId)
	}

	if _, err := uuid.Parse(response.WalletId); err != nil {
		t.Errorf("Handler returned invalid Wallet ID: %v", err)
	}
}
