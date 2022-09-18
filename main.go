package main

import (
	"log"
	"net/http"
	"time"

	common "github.com/tonyxu1/transactionhistory/common"
	handler "github.com/tonyxu1/transactionhistory/handler"
	storage "github.com/tonyxu1/transactionhistory/storage"
)

func main() {
	storage := storage.New()
	mux := http.NewServeMux()

	mux.Handle("/currentblock", handler.CurrentBlockHandler(&storage))
	mux.Handle("/subscribe", handler.SubscribeHandler(&storage))
	mux.Handle("/transaction", handler.TransactionHistoryHandler(&storage))

	//TODO: Not found handler

	go func() {
		for {
			storage.UpdateAllAccount()
			time.Sleep(common.INTERVALINSECONDS)
		}

	}()
	log.Println("Http server started at port 8485")
	log.Fatalln(http.ListenAndServe(":8485", mux))
}
