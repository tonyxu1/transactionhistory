package common

import "time"

//
// Constants
//

const (
	// TODO: move these constant to config file

	// Json RPC Endpoint
	RPCENDPOINT = "https://cloudflare-eth.com"

	// Get detail information of a block ,  including transactions
	GETBLOCKBYNUMBER = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":2304}`

	// Get current block number
	BLOCKNUMBER = `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`

	// Number of blocks goes back from most recent chain block number for transaction retrieval
	LOOKBACKBLOCKS = 1000000

	// Timeout conext for JSON RPC request
	TIMEOUT = time.Duration(time.Millisecond * 1000)

	// Number of Go routines that uses for transaction retrieval from chain to loop though blocks
	NUMOFROUTINES = 6

	// Total number of blocks that will iterate during this round, the blocks will be splitted into each go routine.
	BLOCKSPERROUND = 100

	// Idle period between each round
	INTERVALINSECONDS = 10 * time.Second
)

//
// Global variables
//

var ()

//
// Interface
//

// Parser: defines common methods that expose to the caller
type Parser interface {
	// last parsed block
	GetCurrentBlock() (int, error)

	// add address to observer
	Subscribe() error

	// list of inbound or outbound transactions for an address
	GetTransactions() ([]Transaction, error)
}

// Storage defines methods that need to persist data
type Storage interface {

	//Save current block to storage
	CreateAccount(address string) error

	//Save transactions retrieved from chain to the storage
	SaveTransactions(address string, transactions []Transaction) error

	//Get most recent block number in the storage for the given address
	GetCurrentBlock(address string) (int, error)

	//Get the transaction history information from the storage
	GetTransactions(address string) ([]Transaction, error)
}

//
// Structs
//

// Error Response by Json RPC
type ResponseError struct {
	JsonRPC string `json:"jsonrpc"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Id int64 `json:"id"`
}

// Common Response by Json RPC
type ResponseData struct {
	JsonRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	Id      int64  `json:"id"`
}

// Transaction defines schema of a transaction
type Transaction struct {
	Type                 string        `json:"type"`
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	SignatureV           string        `json:"v"`
	SignatureR           string        `json:"r"`
	SignatureS           string        `json:"s"`
	GasPrice             string        `json:"gasPrice"`
	MaxFeePerGas         string        `json:"maxFeePerGas"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas"`
	ChainID              string        `json:"chainId"`
	AccessList           []interface{} `json:"accessList"`
}

// Block defines the schema of a block in Ethereum
type Block struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Number           string        `json:"number"`
		Hash             string        `json:"hash"`
		ParentHash       string        `json:"parentHash"`
		Sha3Uncles       string        `json:"sha3Uncles"`
		LogsBloom        string        `json:"logsBloom"`
		TransactionsRoot string        `json:"transactionsRoot"`
		StateRoot        string        `json:"stateRoot"`
		ReceiptsRoot     string        `json:"receiptsRoot"`
		Miner            string        `json:"miner"`
		Difficulty       string        `json:"difficulty"`
		TotalDifficulty  string        `json:"totalDifficulty"`
		ExtraData        string        `json:"extraData"`
		Size             string        `json:"size"`
		GasLimit         string        `json:"gasLimit"`
		GasUsed          string        `json:"gasUsed"`
		Timestamp        string        `json:"timestamp"`
		Transactions     []Transaction `json:"transactions"`
		Uncles           []interface{} `json:"uncles"`
		BaseFeePerGas    string        `json:"baseFeePerGas"`
		Nonce            string        `json:"nonce"`
		MixHash          string        `json:"mixHash"`
	} `json:"result"`
	ID int `json:"id"`
}
