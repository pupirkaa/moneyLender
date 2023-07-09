package ml

import (
	_ "embed"
	"fmt"
	"io"
	"sort"
)

type UsersStorage interface {
	io.Closer
	UserExist(name string) (bool, error)
	UserAdd(name string, password string) error
	UserGet(name string) (string, error)
}

type TxsStorage interface {
	io.Closer
	TransactionAdd(lender string, lendee string, money int) error
	DebtsGet() ([]Debt, error)
	TxsGet() ([]Transaction, error)
}

func DistributeDebts(debts []Debt) []Transaction {
	txs := []Transaction{}
	posDebts := []Debt{}
	negDebts := []Debt{}

	sort.Slice(debts, func(i, j int) bool { return debts[i].Money > debts[j].Money })
	for _, d := range debts {
		if d.Money > 0 {
			posDebts = append(posDebts, Debt{Name: d.Name, Money: d.Money})
		} else {
			negDebts = append(negDebts, Debt{Name: d.Name, Money: d.Money})
		}
	}

	sort.SliceStable(negDebts, func(i, j int) bool {
		return i > j
	})

	fmt.Println("+", posDebts, " -", negDebts)
	for i, k := 0, 0; i < len(posDebts) && k < len(negDebts); {
		if posDebts[i].Money > (-negDebts[k].Money) {
			txs = append(txs, Transaction{Lender: negDebts[k].Name, Lendee: posDebts[i].Name, Money: -negDebts[k].Money})
			posDebts[i].Money += negDebts[k].Money
			negDebts[k].Money = 0
			k += 1
		} else if posDebts[i].Money < (-negDebts[k].Money) {
			txs = append(txs, Transaction{Lender: negDebts[k].Name, Lendee: posDebts[i].Name, Money: posDebts[i].Money})
			negDebts[k].Money += posDebts[i].Money
			posDebts[i].Money = 0
			i += 1
		} else {
			txs = append(txs, Transaction{Lender: negDebts[k].Name, Lendee: posDebts[i].Name, Money: posDebts[i].Money})
			negDebts[k].Money = 0
			posDebts[i].Money = 0
			k += 1
			i += 1
		}
	}

	return txs
}
