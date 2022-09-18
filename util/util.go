package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"regexp"
	"strings"
	"time"

	common "github.com/tonyxu1/transactionhistory/common"

	"github.com/ubiq/go-ubiq/common/hexutil"
)

// GetPreviousBlock returns previous block in hex string starts with "0x"
// input parameter: currentBlock: hex string starts with "0x"
func GetPreviousBlock(currentBlock string) (string, error) {

	if !strings.HasPrefix(currentBlock, "0x") {
		return "", errors.New("Invalid current block number: " + currentBlock)
	}

	blockNum, err := hexutil.DecodeBig(currentBlock)
	if err != nil {
		return "", err
	}
	nextNum := blockNum.Add(blockNum, big.NewInt(-1))

	if nextNum.Int64() < 0 {
		return "", errors.New("Block number cannot be less than 0")
	}

	resultStr := hexutil.EncodeBig(nextNum)
	return resultStr, nil

}

// GetDataFromChain retrieve current block number and transactions from th chain
func GetDataFromChain(payLoad string, timeout time.Duration) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, common.RPCENDPOINT, bytes.NewBufferString(payLoad))
	if err != nil {
		return nil, err
	}

	// Post the request using JsonRPC endpoint
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil

}

// Validate Ethereum contract address format
func ValidateAddress(address string) error {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	if !re.MatchString(address) {
		return fmt.Errorf("input address [%s] is invalid", address)
	}

	return nil
}

// ValidateChainData validate the data format
func ValidateChainData(data []byte) (common.ResponseData, error) {
	r := common.ResponseData{}

	err := json.Unmarshal(data, &r)
	if err != nil {
		errorResp := common.ResponseError{}
		//try to unmarshal to error object.
		exception := json.Unmarshal(data, &errorResp)

		if exception != nil {
			return common.ResponseData{}, fmt.Errorf("cannot determine response data %s, %s", err.Error(), exception.Error())
		}

		return common.ResponseData{}, errors.New(errorResp.Error.Message)
	}
	return r, nil

}
