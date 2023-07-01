package inmem

import (
	"sort"

	ml "main.go"
)

type TxsStorage struct {
	txs         []ml.Transaction
	debts       map[string]int
	sortedDebts []ml.Debt
}

func NewTxsStorage() (tfs *TxsStorage) {
	tfs = &TxsStorage{
		txs:         []ml.Transaction{},
		debts:       make(map[string]int),
		sortedDebts: []ml.Debt{},
	}
	return tfs
}

func (tfs TxsStorage) Close() error {
	return nil
}

func (tfs *TxsStorage) TransactionAdd(lender string, lendee string, money int) {
	tfs.txs = append(tfs.txs, ml.Transaction{Lender: lender, Lendee: lendee, Money: money})
	tfs.debts[lender] += money
	tfs.debts[lendee] -= money
	tfs.sortedDebts = sortDebts(tfs.debts)
}

func (tfs *TxsStorage) DebtsGet() []ml.Debt {
	return tfs.sortedDebts
}

func (tfs *TxsStorage) TxsGet() []ml.Transaction {
	return tfs.txs
}

func sortDebts(d map[string]int) []ml.Debt {
	m := make([]string, len(d))

	i := 0
	for k := range d {
		m[i] = k
		i++
	}
	sort.Strings(m)

	var debts []ml.Debt

	for i := 0; i < len(m); i++ {
		debts = append(debts, ml.Debt{
			Name:  m[i],
			Money: d[m[i]],
		})
	}

	return debts
}
