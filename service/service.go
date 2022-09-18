package service

import (
	common "github.com/tonyxu1/transactionhistory/common"
	storage "github.com/tonyxu1/transactionhistory/storage"
)

// Parser struct implements common.Parser interface
type Parser struct {
	Address string
	Storage *storage.Storage
}

func (p Parser) GetCurrentBlock() (int, error) {
	return p.Storage.GetCurrentBlock(p.Address)
}

// Subscribe add address into account map and retrieve all transactions from chain
func (p Parser) Subscribe() error {
	return p.Storage.CreateAccount(p.Address)
}

// GetTransactions return all transaction history record from both on chain and local storage
func (p Parser) GetTransactions() ([]common.Transaction, error) {

	return p.Storage.GetTransactions(p.Address)
}
