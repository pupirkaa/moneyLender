package ml

import (
	"reflect"
	"testing"
)

type mock struct {
	debts []Debt
}

func (m mock) Close() (e error)                                       { return nil }
func (m mock) TransactionAdd(lender string, lendee string, money int) {}
func (m mock) DebtsGet() (d []Debt)                                   { return m.debts }
func (m mock) TxsGet() (t []Transaction)                              { return nil }

func TestTxsController_DistributeDebts(t *testing.T) {
	type fields struct {
		Txs   TxsStorage
		Users UsersStorage
	}
	tests := []struct {
		name   string
		fields fields
		want   []Transaction
	}{
		{
			name: "empty transactions",
			fields: fields{
				Txs:   mock{debts: []Debt{}},
				Users: nil,
			},
			want: []Transaction{},
		},
		{
			name: "one transaction",
			fields: fields{
				Txs:   mock{debts: []Debt{{Name: "Irina", Money: 100}, {Name: "Matvei", Money: -100}}},
				Users: nil,
			},
			want: []Transaction{{Lender: "Matvei", Lendee: "Irina", Money: 100}},
		},
		{
			name: "two transactions",
			fields: fields{
				Txs:   mock{debts: []Debt{{Name: "Irina", Money: 40}, {Name: "Matvei", Money: -20}, {Name: "Nikita", Money: -20}}},
				Users: nil,
			},
			want: []Transaction{{Lender: "Nikita", Lendee: "Irina", Money: 20}, {Lender: "Matvei", Lendee: "Irina", Money: 20}},
		},
		{
			name: "two transactions",
			fields: fields{
				Txs:   mock{debts: []Debt{{Name: "Irina", Money: 20}, {Name: "Matvei", Money: -20}, {Name: "Nikita", Money: -20}, {Name: "Kirill", Money: 20}}},
				Users: nil,
			},
			want: []Transaction{{Lender: "Nikita", Lendee: "Irina", Money: 20}, {Lender: "Matvei", Lendee: "Kirill", Money: 20}},
		},
		{
			name: "two transactions2",
			fields: fields{
				Txs:   mock{debts: []Debt{{Name: "Irina", Money: -40}, {Name: "Matvei", Money: 20}, {Name: "Nikita", Money: 20}}},
				Users: nil,
			},
			want: []Transaction{{Lender: "Irina", Lendee: "Matvei", Money: 20}, {Lender: "Irina", Lendee: "Nikita", Money: 20}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := TxsController{
				Txs:   tt.fields.Txs,
				Users: tt.fields.Users,
			}
			if got := tr.DistributeDebts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TxsController.DistributeDebts() = %v, want %v", got, tt.want)
			}
		})
	}
}
