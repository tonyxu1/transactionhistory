package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"

	common "github.com/tonyxu1/transactionhistory/common"
	util "github.com/tonyxu1/transactionhistory/util"

	"github.com/ubiq/go-ubiq/common/hexutil"
)

type Storage struct {
	account     sync.Map
	transaction sync.Map
}

// Initiate a new storage
// Storage implementation, can be replaceed by other storage methods

func New() Storage {
	return Storage{
		account:     sync.Map{},
		transaction: sync.Map{},
	}
}

func (s *Storage) IsNewAccount(address string) bool {
	if _, ok := s.account.Load(address); ok {
		return false
	}
	return true
}

// SaveAccountInfo Save account information
func (s *Storage) CreateAccount(address string) error {
	err := util.ValidateAddress(address)
	if err != nil {
		return err
	}

	if s.IsNewAccount(address) {
		blockNum, err := getBlockNumFromChain()
		if err != nil {
			return err
		}
		s.account.Store(address, blockNum)
		s.transaction.Store(address, []common.Transaction{}) //Empty transaction for the new account
		return nil
	}
	return fmt.Errorf("account for address [%s] already subscribed", address)
}

// SaveTransactions append transactions to the account
func (s *Storage) SaveTransactions(address string, transactions []common.Transaction) error {
	err := util.ValidateAddress(address)
	if err != nil {
		return err
	}

	if s.IsNewAccount(address) {
		return fmt.Errorf("account for address [%s] does not exist", address)
	}

	data, ok := s.transaction.Load(address)

	if !ok {
		s.transaction.Store(address, transactions)
		return nil
	}
	existingTrans := data.([]common.Transaction)
	allTrans := append(existingTrans, transactions...)
	s.transaction.Store(address, allTrans)
	return nil
}

// GetCurrentBlock : get most recent block number in the storage for the given address
func (s *Storage) GetCurrentBlock(address string) (int, error) {
	err := util.ValidateAddress(address)
	if err != nil {
		return -1, err
	}

	if s.IsNewAccount(address) {
		return -1, fmt.Errorf("account for address [%s] does not exist", address)
	}

	d, _ := s.account.Load(address)
	return d.(int), nil
}

// GetTransactions : retrieve the transaction history information from the storage
func (s *Storage) GetTransactions(address string) ([]common.Transaction, error) {
	err := util.ValidateAddress(address)
	if err != nil {
		return []common.Transaction{}, err
	}
	if s.IsNewAccount(address) {
		return []common.Transaction{}, fmt.Errorf("account for address [%s] does not exist", address)
	}

	data, _ := s.transaction.Load(address)
	transHistory := data.([]common.Transaction)

	//Order by Block Number descending
	sort.Slice(transHistory, func(i, j int) bool {
		a, err := hexutil.DecodeBig(transHistory[i].BlockNumber)
		if err != nil {
			return false
		}
		b, err := hexutil.DecodeBig(transHistory[j].BlockNumber)
		if err != nil {
			return false
		}

		return a.Int64() > b.Int64()
	})

	return transHistory, nil
}

// getTransFromChain retrieve transaction data in batch manner
// for given adddress
func (s *Storage) UpdateAccountWithChainData(address string) error {

	var (
		wg         sync.WaitGroup
		mu         = &sync.Mutex{}
		blocksToGo = common.BLOCKSPERROUND
	)

	resultTrans := make([]common.Transaction, 0)
	errs := make([]error, 0)
	//search transactions from chain with  goroutines
	for i := 0; i < common.NUMOFROUTINES; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for blocksToGo > 0 {
				blockNum, _ := s.account.Load(address)
				currentBlock := blockNum.(int)

				blockNumStr := "0x" + strconv.FormatInt(int64(currentBlock), 16)
				data, err := util.GetDataFromChain(fmt.Sprintf(common.GETBLOCKBYNUMBER, blockNumStr), common.TIMEOUT)
				if err != nil {
					mu.Lock()
					errs = append(errs, err)
					blocksToGo--
					mu.Unlock()
					s.account.Store(address, currentBlock+1)
					return
				}
				var blockInfo common.Block
				err = json.Unmarshal([]byte(data), &blockInfo)
				if err != nil {
					errResp := common.ResponseError{}
					err1 := json.Unmarshal(data, &errResp)
					if err1 != nil {
						mu.Lock()
						errs = append(errs, fmt.Errorf("%s | %s", err.Error(), err1.Error()))
						blocksToGo--
						mu.Unlock()
						s.account.Store(address, currentBlock+1)
						return
					}
					mu.Lock()
					errs = append(errs, errors.New(errResp.Error.Message))
					blocksToGo--
					mu.Unlock()
					s.account.Store(address, currentBlock+1)
					return
				}
				trx := blockInfo.Result.Transactions

				for _, tr := range trx {
					if tr.From == address || tr.To == address {
						mu.Lock()
						resultTrans = append(resultTrans, tr)
						mu.Unlock()
					}
				}
				mu.Lock()
				blocksToGo--
				mu.Unlock()
				s.account.Store(address, currentBlock+1)
			}
		}()
	}
	wg.Wait()

	if len(errs) > 0 {
		errMsg := ""
		for _, v := range errs {
			errMsg += v.Error()
			if len(errMsg) > 200 {
				break
			}
		}
		log.Println("errors:", errMsg)
		return errors.New(errMsg)
	}
	if len(resultTrans) > 0 {
		return s.SaveTransactions(address, resultTrans)
	}
	return nil
}

func (s *Storage) UpdateAllAccount() {
	s.account.Range(func(key, value any) bool {
		addr := key.(string)
		log.Printf("Update account for address [%s] and block number [%d]\n", addr, value.(int))
		err := s.UpdateAccountWithChainData(addr)
		if err != nil {
			log.Println("UpdateAllAccount() err: ", err)
		}
		return true
	})
}

func getBlockNumFromChain() (int, error) {
	// get most current block number minus look
	// back blocks currently set it to 1000000
	// TODO: make it configurable
	// & transaction count for the address
	data, err := util.GetDataFromChain(common.BLOCKNUMBER, common.TIMEOUT)
	if err != nil {
		return -1, err
	}
	resp, err := util.ValidateChainData(data)
	if err != nil {
		return -1, err
	}

	blockInHex := resp.Result

	num, _ := strconv.ParseInt(blockInHex, 0, 64)
	num -= common.LOOKBACKBLOCKS
	return int(num), nil
}
