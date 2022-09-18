package service

import (
	"reflect"
	"testing"

	common "github.com/tonyxu1/transactionhistory/common"
	storage "github.com/tonyxu1/transactionhistory/storage"
)

func TestParser_GetCurrentBlock(t *testing.T) {

	s := storage.New()
	s.CreateAccount("0xE946502872DA09009Aa6dc975272AC24Ab5B4f36")

	type fields struct {
		Address string
		Storage *storage.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{
			name: "Invalid address format",
			fields: fields{
				Address: "0x134856623",
				Storage: &storage.Storage{},
			},
			want:    0,
			wantErr: true,
		}, {
			name: "Should have current block number returned",
			fields: fields{
				Address: "0xE946502872DA09009Aa6dc975272AC24Ab5B4f36",
				Storage: &s,
			},
			want:    14000000,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Invalid address format" {
				p := Parser{
					Address: tt.fields.Address,
					Storage: tt.fields.Storage,
				}
				got, err := p.GetCurrentBlock()
				if (err != nil) == tt.wantErr {
					t.Logf("Parser.GetCurrentBlock() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got != tt.want {
					t.Errorf("Parser.GetCurrentBlock() = %v, want %v", got, tt.want)
				}
				return
			} else if tt.name == "Should have current block number returned" {
				p := Parser{
					Address: tt.fields.Address,
					Storage: tt.fields.Storage,
				}
				got, err := p.GetCurrentBlock()
				if (err != nil) == tt.wantErr {
					t.Logf("Parser.GetCurrentBlock() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got <= tt.want || got == -1 {
					t.Errorf("Parser.GetCurrentBlock() = %v, want %v", got, tt.want)
				}
				return
			}

		})
	}
}

func TestParser_Subscribe(t *testing.T) {
	type fields struct {
		Address string
		Storage *storage.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Subscribe succeed",
			fields: fields{
				Address: "0xE946502872DA09009Aa6dc975272AC24Ab5B4f36",
				Storage: &storage.Storage{},
			},
			wantErr: false,
		}, {
			name: "Account already subscribed",
			fields: fields{
				Address: "0xE946502872DA09009Aa6dc975272AC24Ab5B4f36",
				Storage: &storage.Storage{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Subscribe succeed" {
				p := Parser{
					Address: tt.fields.Address,
					Storage: tt.fields.Storage,
				}
				if err := p.Subscribe(); (err != nil) != tt.wantErr {
					t.Errorf("Parser.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if tt.name == "Account already subscribed" {

				s := storage.New()
				s.CreateAccount("0xE946502872DA09009Aa6dc975272AC24Ab5B4f36")
				p := Parser{
					Address: tt.fields.Address,
					Storage: &s,
				}
				if err := p.Subscribe(); (err != nil) != tt.wantErr {
					t.Errorf("Parser.Subscribe() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

		})
	}
}

func TestParser_GetTransactions(t *testing.T) {
	type fields struct {
		Address string
		Storage *storage.Storage
	}
	tests := []struct {
		name    string
		fields  fields
		want    []common.Transaction
		wantErr bool
	}{
		{
			name: "Account not subscribed",
			fields: fields{
				Address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
				Storage: &storage.Storage{},
			},
			want:    []common.Transaction{},
			wantErr: true,
		}, {
			name: "Return transaction array order by block number desc",
			fields: fields{
				Address: "0xE946502872DA09009Aa6dc975272AC24Ab5B4f36",
				Storage: &storage.Storage{},
			},
			want: []common.Transaction{
				{
					BlockNumber: "0x125",
				}, {
					BlockNumber: "0x124",
				}, {
					BlockNumber: "0x123",
				}, {
					BlockNumber: "0x121",
				},
			},
			wantErr: false,
		}, {
			name: "Return empty transaction array",
			fields: fields{
				Address: "0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b",
				Storage: &storage.Storage{},
			},
			want:    []common.Transaction{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Account not subscribed" {
				p := Parser{
					Address: tt.fields.Address,
					Storage: tt.fields.Storage,
				}
				got, err := p.GetTransactions()
				if (err != nil) != tt.wantErr {
					t.Errorf("Parser.GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parser.GetTransactions() = %v, want %v", got, tt.want)
				}
			} else if tt.name == "Return transaction array order by block number desc" {
				s := storage.New()
				s.CreateAccount("0xE946502872DA09009Aa6dc975272AC24Ab5B4f36")
				trans := []common.Transaction{
					{
						BlockNumber: "0x123",
					}, {
						BlockNumber: "0x125",
					}, {
						BlockNumber: "0x124",
					}, {
						BlockNumber: "0x121",
					},
				}
				err := s.SaveTransactions("0xE946502872DA09009Aa6dc975272AC24Ab5B4f36", trans)
				if err != nil {
					t.Errorf("s.SaveTransactions() error : %v", err)
				}
				p := Parser{
					Address: tt.fields.Address,
					Storage: &s,
				}
				got, err := p.GetTransactions()
				if (err != nil) != tt.wantErr {
					t.Errorf("Parser.GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parser.GetTransactions() = %v, want %v", got, tt.want)
				}
			} else if tt.name == "Return empty transaction array" {
				s := storage.New()
				s.CreateAccount("0x23ca95B9dE14A83CBF4A43b11C2C3825e72c7d9b")
				p := Parser{
					Address: tt.fields.Address,
					Storage: &s,
				}
				got, err := p.GetTransactions()
				if (err != nil) != tt.wantErr {
					t.Errorf("Parser.GetTransactions() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Parser.GetTransactions() = %v, want %v", got, tt.want)
				}
			}

		})
	}
}
