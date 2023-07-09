package ml

import (
	_ "embed"
	"reflect"
	"testing"
)

func TestDistributeDebts(t *testing.T) {
	type args struct {
		debts []Debt
	}
	tests := []struct {
		name string
		args args
		want []Transaction
	}{
		{
			name: "empty transactions",
			args: args{debts: []Debt{}},
			want: []Transaction{},
		},
		{
			name: "one transaction",
			args: args{debts: []Debt{{Name: "Irina", Money: 100}, {Name: "Matvei", Money: -100}}},
			want: []Transaction{{Lender: "Matvei", Lendee: "Irina", Money: 100}},
		},
		{
			name: "two transactions",
			args: args{debts: []Debt{{Name: "Irina", Money: 40}, {Name: "Matvei", Money: -20}, {Name: "Nikita", Money: -20}}},
			want: []Transaction{{Lender: "Nikita", Lendee: "Irina", Money: 20}, {Lender: "Matvei", Lendee: "Irina", Money: 20}},
		},
		{
			name: "two transactions2",
			args: args{debts: []Debt{{Name: "Irina", Money: -40}, {Name: "Matvei", Money: 20}, {Name: "Nikita", Money: 20}}},
			want: []Transaction{{Lender: "Irina", Lendee: "Matvei", Money: 20}, {Lender: "Irina", Lendee: "Nikita", Money: 20}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DistributeDebts(tt.args.debts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DistributeDebts() = %v, want %v", got, tt.want)
			}
		})
	}
}
