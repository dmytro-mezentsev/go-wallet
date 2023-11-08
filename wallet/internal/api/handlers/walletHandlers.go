package walletHandlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type CreateWalletRequest struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	BirthDate JsonDate `json:"birthDate"`
	Address   string   `json:"address"`
}

func (t *JsonDate) UnmarshalJSON(d []byte) (err error) {
	if d[0] == '"' && d[len(d)-1] == '"' {
		d = d[1 : len(d)-1]
	}
	t.Time, err = time.Parse(time.DateOnly, string(d))
	return
}

type JsonDate struct {
	time.Time
}

func CreateWalletHandler(w http.ResponseWriter, r *http.Request) {

	var wallet CreateWalletRequest

	// Create a new JSON decoder for the request body
	decoder := json.NewDecoder(r.Body)

	// Use the decoder to unmarshal the JSON data into the 'wallet' struct
	if err := decoder.Decode(&wallet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	responseJSON, err := json.Marshal(wallet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(responseJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
