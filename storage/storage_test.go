package storage

import (
	"reflect"
	"sync"
	"testing"

	"github.com/tonyxu1/transactionhistory/common"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Storage
	}{
		{
			name: "New storage is created",
			want: Storage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_IsNewAccount(t *testing.T) {
	type fields struct {
		account     sync.Map
		transaction sync.Map
	}
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "New Account",
			fields: fields{
				account:     sync.Map{},
				transaction: sync.Map{},
			},
			args: args{
				address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
			},
			want: true,
		}, {
			name: "Existing Account",
			fields: fields{
				account:     sync.Map{},
				transaction: sync.Map{},
			},
			args: args{
				address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "New Account" {
				s := &Storage{
					account:     tt.fields.account,
					transaction: tt.fields.transaction,
				}
				if got := s.IsNewAccount(tt.args.address); got != tt.want {
					t.Errorf("Storage.IsNewAccount() = %v, want %v", got, tt.want)
				}
			} else if tt.name == "Existing Account" {
				s := New()
				err := s.CreateAccount("0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b")
				if err != nil {
					t.Errorf("s.CreateAccount() error : %v", err)
				}
				if got := s.IsNewAccount(tt.args.address); got != tt.want {
					t.Errorf("Storage.IsNewAccount() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

func TestStorage_CreateAccount(t *testing.T) {
	type fields struct {
		account     sync.Map
		transaction sync.Map
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "New account created",
			fields: fields{
				account:     sync.Map{},
				transaction: sync.Map{},
			},
			args: args{
				address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
			},
			wantErr: false,
		}, {
			name:   "Existing Account",
			fields: fields{},
			args: args{
				address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Existing Account" {
				s := &Storage{
					account:     tt.fields.account,
					transaction: tt.fields.transaction,
				}
				err := s.CreateAccount("0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b")
				if err != nil {
					t.Errorf("s.CreateAccount() error : %v", err)
				}
				if err := s.CreateAccount(tt.args.address); (err != nil) != tt.wantErr {
					t.Errorf("Storage.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				s := &Storage{
					account:     tt.fields.account,
					transaction: tt.fields.transaction,
				}
				if err := s.CreateAccount(tt.args.address); (err != nil) != tt.wantErr {
					t.Errorf("Storage.CreateAccount() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

		})
	}
}

func TestStorage_SaveTransactions(t *testing.T) {
	type fields struct {
		account     sync.Map
		transaction sync.Map
	}
	type args struct {
		address      string
		transactions []common.Transaction
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Append transactions to existing account",
			fields: fields{
				account:     sync.Map{},
				transaction: sync.Map{},
			},
			args: args{
				address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
				transactions: []common.Transaction{
					{BlockNumber: "0x333"}, {BlockNumber: "0x334"},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				account:     tt.fields.account,
				transaction: tt.fields.transaction,
			}
			err := s.CreateAccount(tt.args.address)
			if err != nil {
				t.Errorf("s.CreateAccount() error : %v", err)
			}
			if err := s.SaveTransactions(tt.args.address, tt.args.transactions); (err != nil) != tt.wantErr {
				t.Errorf("Storage.SaveTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
			count := 0
			s.transaction.Range(func(key, value any) bool {
				if key.(string) == "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b" {
					count = len(value.([]common.Transaction))
					return false
				}
				return true
			})

			if count != 2 {
				t.Errorf("Storage.SaveTransactions() error: got %d expected count 2", count)
				return
			}
		})
	}
}

func TestStorage_GetCurrentBlock(t *testing.T) {
	type fields struct {
		account     sync.Map
		transaction sync.Map
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "Current block in storage",
			fields: fields{
				account:     sync.Map{},
				transaction: sync.Map{},
			},
			args: args{
				address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
			},
			want:    14000000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				account:     tt.fields.account,
				transaction: tt.fields.transaction,
			}
			s.CreateAccount(tt.args.address)
			got, err := s.GetCurrentBlock(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetCurrentBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got < tt.want {
				t.Errorf("Storage.GetCurrentBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetTransactions(t *testing.T) {
	type fields struct {
		account     sync.Map
		transaction sync.Map
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []common.Transaction
		wantErr bool
	}{
		{
			name: "Return sorted transaction array by block number desc",
			fields: fields{
				account:     sync.Map{},
				transaction: sync.Map{},
			},
			args: args{
				address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
			},
			want: []common.Transaction{
				{BlockNumber: "0x1237"},
				{BlockNumber: "0x1236"},
				{BlockNumber: "0x1235"},
				{BlockNumber: "0x1234"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				account:     tt.fields.account,
				transaction: tt.fields.transaction,
			}
			err := s.CreateAccount(tt.args.address)
			if err != nil {
				t.Errorf("s.CreateAccount() error : %v", err)
			}
			err = s.SaveTransactions(tt.args.address, []common.Transaction{
				{BlockNumber: "0x1234"},
				{BlockNumber: "0x1236"},
				{BlockNumber: "0x1235"},
				{BlockNumber: "0x1237"},
			})
			if err != nil {
				t.Errorf("s.SaveTransactions() error : %v", err)
			}

			got, err := s.GetTransactions(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.GetTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: Add test cases.
func TestStorage_UpdateAccountWithChainData(t *testing.T) {
	type fields struct {
		account     sync.Map
		transaction sync.Map
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				account:     tt.fields.account,
				transaction: tt.fields.transaction,
			}
			if err := s.UpdateAccountWithChainData(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("Storage.UpdateAccountWithChainData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TODO: Add test cases.
func TestStorage_UpdateAllAccount(t *testing.T) {
	type fields struct {
		account     sync.Map
		transaction sync.Map
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				account:     tt.fields.account,
				transaction: tt.fields.transaction,
			}
			s.UpdateAllAccount()
		})
	}
}

// TODO: Add test cases.
func Test_getBlockNumFromChain(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getBlockNumFromChain()
			if (err != nil) != tt.wantErr {
				t.Errorf("getBlockNumFromChain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getBlockNumFromChain() = %v, want %v", got, tt.want)
			}
		})
	}
}
