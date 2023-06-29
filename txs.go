package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type TransactionsAndDebts struct {
	Transactions []Transaction
	Debts        []Debt
}

type Transaction struct {
	Lender string
	Lendee string
	Money  int
}

type Debt struct {
	Name  string
	Money int
}

type TxsController struct {
	txs      []Transaction
	newTxs   []Transaction
	users    map[string]string
	newUsers map[string]string
	debts    []Debt
	t        *template.Template
}

func (t *TxsController) index(w http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("user")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	io.WriteString(w, generateHTML(t.txs, t.debts, t.t))
}

func (t *TxsController) addTransaction(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse form: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}
	form := req.Form

	var (
		lender   = form.Get("lender")
		lendee   = form.Get("lendee")
		money, _ = strconv.Atoi(form.Get("money"))
	)
	if lender == "" || lendee == "" || money <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}

	_, lenderExist := t.users[lender]
	_, lendeeExist := t.users[lendee]
	if !lenderExist || !lendeeExist {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User not found")
		return
	}

	nT := Transaction{Lender: lender, Lendee: lendee, Money: money}

	for i := range t.debts {
		if lender == t.debts[i].Name {
			t.debts[i].Money += money
		}
		if lendee == t.debts[i].Name {
			t.debts[i].Money -= money
		}
	}
	t.txs = append(t.txs, nT)
	t.newTxs = append(t.newTxs, nT)
	http.Redirect(w, req, "/", http.StatusSeeOther)

}

func parseTransactions(transactionData []string) ([]Transaction, map[string]int) {
	debts := map[string]int{}
	var transactions []Transaction

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

		tr := Transaction{
			Lender: splitedString[0],
			Lendee: splitedString[4],
			Money:  amountOfMoney,
		}

		transactions = append(transactions, tr)

		debts[tr.Lender] += tr.Money
		debts[tr.Lendee] -= tr.Money
	}
	return transactions, debts

}

func (t *TxsController) saveTxsToFile() {
	f, err := os.OpenFile(os.Args[1], os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	for i := range t.newTxs {
		_, err = f.WriteString(fmt.Sprintf("\n%s lent %v$ to %s", t.newTxs[i].Lender, t.newTxs[i].Money, t.newTxs[i].Lendee))
		if err != nil {
			fmt.Println(err)
		}
	}
	t.newTxs = []Transaction{}

	defer f.Close()
}
