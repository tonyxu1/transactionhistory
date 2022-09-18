package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	common "github.com/tonyxu1/transactionhistory/common"
)

func CurrentBlockHandler(s common.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")

		blockNum, err := s.GetCurrentBlock(address)
		if err != nil {
			_, err1 := w.Write([]byte(err.Error()))
			if err1 != nil {
				log.Println("w.Write() error :", err1)
				return
			}
			return
		}
		w.Header().Add("Content-Type", "application/json")
		resp := `{"block_number":"` + strconv.Itoa(blockNum) + `"}`
		w.Write([]byte(resp))
	}
}

// SubscribeHandler : public endpoint for subscription of an account
func SubscribeHandler(s common.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")

		err := s.CreateAccount(address)
		if err != nil {
			_, err1 := w.Write([]byte(err.Error()))
			if err1 != nil {
				log.Println("w.Write() error: ", err1)
				return
			}
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"message":"subscription succeed"}`))
	}
}

// TransactionHistoryHandler : retrieve transaction for a given address from both of the
// chain and local storage.
func TransactionHistoryHandler(s common.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		trans, err := s.GetTransactions(address)
		if err != nil {
			_, err1 := w.Write([]byte(err.Error()))
			if err1 != nil {
				log.Println("w.Write() error: ", err1)
				return
			}

			return
		}
		w.Header().Add("Content-Type", "application/json")
		transBytes, err := json.Marshal(trans)
		if err != nil {
			_, err1 := w.Write([]byte(err.Error()))
			if err1 != nil {
				log.Println("w.Write() error: ", err1)
				return
			}
			return
		}
		w.Write(transBytes)

	}
}
