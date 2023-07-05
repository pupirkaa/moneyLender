package sqlite

import (
	"database/sql"
	"fmt"
	"os"

	ml "main.go"
)

type TxsStorage struct {
	db *sql.DB
}

func NewTxsStorage(path string) (tfs *TxsStorage) {
	dsn := "file:" + path
	d, err := sql.Open("sqlite", dsn)
	tfs = &TxsStorage{
		db: d,
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "opening db: %v\n", err)
		os.Exit(1)
	}

	_, err = tfs.db.Exec("CREATE TABLE IF NOT EXISTS txs(lender string, lendee string, money int,  FOREIGN KEY(lender,lendee) REFERENCES users(name,name));", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating transactions table: %v\n", err)
		os.Exit(1)
	}

	_, err = tfs.db.Exec("CREATE TABLE IF NOT EXISTS debts(name string UNIQUE, money int,  FOREIGN KEY(name) REFERENCES users(name));", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating debts table: %v\n", err)
		os.Exit(1)
	}

	return tfs
}

func (tfs TxsStorage) Close() error {
	if err := tfs.db.Close(); err != nil {
		return fmt.Errorf("closing db: %v", err)
	}
	return nil
}

func (tfs *TxsStorage) TransactionAdd(lender string, lendee string, money int) error {
	_, err := tfs.db.Exec("INSERT INTO txs (lender, lendee, money) VALUES (?, ?, ?);", lender, lendee, money)
	if err != nil {
		return fmt.Errorf("checking is transaction exist: %v", err)
	}

	_, err = tfs.db.Exec("INSERT INTO debts(name,money) VALUES(?, ?) ON CONFLICT(name) DO UPDATE SET money=money+?;", lender, money, money)
	if err != nil {
		return fmt.Errorf("checking is transaction exist: %v", err)
	}
	_, err = tfs.db.Exec("INSERT INTO debts(name,money) VALUES(?, ?) ON CONFLICT(name) DO UPDATE SET money=money+?;", lendee, -money, -money)
	if err != nil {
		return fmt.Errorf("checking is transaction exist: %v", err)
	}

	return nil
}

func (tfs *TxsStorage) DebtsGet() ([]ml.Debt, error) {
	debts := []ml.Debt{}

	rows, err := tfs.db.Query("SELECT * FROM debts ORDER BY name ASC;")
	if err != nil {
		return debts, fmt.Errorf("selecting debts: %v", err)
	}
	defer rows.Close()
	for i := 0; rows.Next(); i++ {
		debts = append(debts, ml.Debt{})
		if err := rows.Scan(&debts[i].Name, &debts[i].Money); err != nil {
			return debts, fmt.Errorf("getting debts: %v", err)
		}
	}

	return debts, nil
}

func (tfs *TxsStorage) TxsGet() ([]ml.Transaction, error) {
	txs := []ml.Transaction{}

	rows, err := tfs.db.Query("SELECT * FROM txs;")
	if err != nil {
		return txs, fmt.Errorf("selecting transactions: %v", err)
	}
	defer rows.Close()
	for i := 0; rows.Next(); i++ {
		txs = append(txs, ml.Transaction{})
		if err := rows.Scan(&txs[i].Lender, &txs[i].Lendee, &txs[i].Money); err != nil {
			return txs, fmt.Errorf("getting transactions: %v", err)
		}
	}
	return txs, nil
}
