package util

import (
	"reflect"
	"testing"
	"time"

	common "github.com/tonyxu1/transactionhistory/common"
)

func TestGetPreviousBlock(t *testing.T) {
	type args struct {
		currentBlock string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Invalid block string",
			args: args{
				currentBlock: "abd000ff",
			},
			want:    "",
			wantErr: true,
		}, {
			name: "Return previous block number hex string",
			args: args{
				currentBlock: "0xfff",
			},
			want:    "0xffe",
			wantErr: false,
		}, {
			name: " raise error for negative value",
			args: args{
				currentBlock: "0x0",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPreviousBlock(tt.args.currentBlock)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPreviousBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPreviousBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDataFromChain(t *testing.T) {
	type args struct {
		payLoad string
		timeout time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Testing post with correct payload",
			args: args{
				payLoad: `{"jsonrpc":"2.0","method":"eth_block","params":[],"id":1}`,
				timeout: time.Duration(common.TIMEOUT),
			},
			want: []byte{123, 34, 106, 115, 111, 110, 114, 112, 99, 34, 58, 34,
				50, 46, 48, 34, 44, 34, 101, 114, 114, 111, 114, 34, 58, 123, 34, 99, 111, 100, 101, 34,
				58, 45, 51, 50, 54, 48, 49, 44, 34, 109, 101, 115, 115, 97, 103, 101, 34, 58, 34, 77, 101, 116, 104, 111, 100, 32, 110, 111, 116, 32, 102, 111, 117, 110, 100, 34, 125, 44, 34, 105, 100, 34, 58, 49, 125},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDataFromChain(tt.args.payLoad, tt.args.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataFromChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDataFromChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Invalid format",
			args: args{
				address: "0x12345",
			},
			wantErr: true,
		}, {
			name: "Valid format",
			args: args{
				address: "0xE946502872DA09009Aa6dc975272AC24Ab5B4f36",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateAddress(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("ValidateAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateChainData(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    common.ResponseData
		wantErr bool
	}{
		{
			name: "Valid resposne data",
			args: args{
				data: []byte{123, 34, 106, 115, 111, 110, 114, 112, 99, 34, 58, 34, 50, 46, 48, 34, 44, 34, 114, 101, 115, 117, 108, 116, 34, 58, 34, 48, 120, 50, 51, 52, 53, 54, 34, 44, 34, 105, 100, 34, 58, 49, 125},
			},
			want: common.ResponseData{
				JsonRPC: "2.0",
				Result:  "0x23456",
				Id:      1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateChainData(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateChainData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateChainData() = %v, want %v", got, tt.want)
			}
		})
	}
}
