package fs

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	ml "github.com/pupirkaa/moneyLender"
)

type TxsFileStorage struct {
	txs         []ml.Transaction
	debts       map[string]int
	sortedDebts []ml.Debt
	path        string
	exitCh      chan bool
	syncDoneCh  chan bool
}

func NewTxsStorage(path string) (tfs *TxsFileStorage) {
	tfs = &TxsFileStorage{
		txs:        parseTransactions(readFile(path)),
		path:       path,
		exitCh:     make(chan bool),
		syncDoneCh: make(chan bool),
	}
	tfs.debts = makeDebts(tfs.txs)
	tfs.sortedDebts = sortDebts(tfs.debts)

	go func() {
		t := time.NewTicker(1 * time.Minute)
		for {
			select {
			case <-tfs.exitCh:
				tfs.saveTxsToFile()
				close(tfs.syncDoneCh)
				return
			case <-t.C:
				tfs.saveTxsToFile()
			}
		}
	}()
	return tfs
}

func (tfs TxsFileStorage) DebtsGet() ([]ml.Debt, error) {
	return tfs.sortedDebts, nil
}

func (tfs TxsFileStorage) TxsGet() ([]ml.Transaction, error) {
	return tfs.txs, nil
}

func (tfs *TxsFileStorage) TransactionAdd(lender string, lendee string, money int) error {
	tfs.txs = append(tfs.txs, ml.Transaction{Lender: lender, Lendee: lendee, Money: money})
	tfs.debtAdd(lender, lendee, money)
	return nil
}

func (tfs *TxsFileStorage) debtAdd(lender string, lendee string, money int) {
	tfs.debts[lender] += money
	tfs.debts[lendee] -= money
	tfs.sortedDebts = sortDebts(tfs.debts)
}

func (tfs TxsFileStorage) Close() error {
	close(tfs.exitCh)
	<-tfs.syncDoneCh
	return nil
}

func (tfs TxsFileStorage) saveTxsToFile() {
	f, err := os.Create(tfs.path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()
	fmt.Println(tfs.txs)
	for _, t := range tfs.txs {
		_, err = f.WriteString(fmt.Sprintf("%s lent %v$ to %s\n", t.Lender, t.Money, t.Lendee))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func parseTransactions(transactionData []string) []ml.Transaction {
	var txs []ml.Transaction

	for i := 0; i < len(transactionData); i++ {
		splitedString := (strings.Split(transactionData[i], " "))
		if len(splitedString) != 5 {
			fmt.Fprintln(os.Stderr, "your data is incorrect")
			os.Exit(1)
		}

		splitedString[2] = strings.TrimSuffix(splitedString[2], "$")
		amountOfMoney, err := strconv.Atoi(splitedString[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse money amount: %v", err)
			os.Exit(1)
		}

		txs = append(txs, ml.Transaction{
			Lender: splitedString[0],
			Lendee: splitedString[4],
			Money:  amountOfMoney,
		})
	}
	return txs
}

func makeDebts(txs []ml.Transaction) (debts map[string]int) {
	debts = make(map[string]int)
	for i := range txs {
		debts[txs[i].Lender] += txs[i].Money
		debts[txs[i].Lendee] -= txs[i].Money
	}
	return
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
