#  TRANSACTION HISTORY PROJECT
## Requirement
Detailed information can be found [here](https://trustwallet.notion.site/Backend-Homework-abd431fca950427db75d73d90a0244a8) 


## Solution
Using multiple go routines to retrieve transaction history data from chain and save them into the system, 3 public enpoints for accessing data from local storage
### Run the App
To run the app is pretty staightforward, similar as the other Go applications.
- Clone the app to your local 
- Run following command from the root folder of the project: 
` $ make all`
- The servver should start at port 8485
### Endpoints
`/subscribe?address=<contract address>` : Register the given address to the system, if the address already exists, an error will be returned.

`/currentblock?address=<contract address>` : Get the current block number associated with given address that saved in current storage, error message will be returned if the given address doesn't exist in the system.

`/transaction?address=<contract address>` : Get the transaction history either from the given address or to the address.

### Timer Event
The background go routine is running under timer timer manner with configurable idle period, along with it, some other configurable items are:
```

	// Json RPC Endpoint
	RPCENDPOINT      = "https://cloudflare-eth.com"

	// Get detail information of a block ,  including transactions
	GETBLOCKBYNUMBER = `{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":2304}`

	// Get current block number 
	BLOCKNUMBER      = `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`

	// Number of blocks goes back from most recent chain block number for transaction retrieval 
	LOOKBACKBLOCKS   = 1000000 

	// Timeout conext for JSON RPC request
	TIMEOUT           = time.Duration(time.Millisecond * 1000)

	// Number of Go routines that uses for transaction retrieval from chain to loop though blocks
	NUMOFROUTINES     = 6

	// Total number of blocks that will iterate during this round, the blocks will be splitted into each go routine.
	BLOCKSPERROUND    = 100
	
	// Idle period between each round
	INTERVALINSECONDS = 10 * time.Second
    
```
## TODOs
Due to the limit of time, there some rooms to improve:
- Customized error collection
- Configuration module
- Logging framework
- Performance metrics
- Linting errors 

