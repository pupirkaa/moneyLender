package ml

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
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
	Txs   TxsStorage
	Users UsersStorage
}

func (t *TxsController) AddTransaction(w http.ResponseWriter, req *http.Request) {
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

	lenderExist, err := t.Users.UserExist(lender)
	if err != nil {
		panic("TODO: handle error")
	}

	lendeeExist, err := t.Users.UserExist(lendee)
	if err != nil {
		panic("TODO: handle error")
	}

	if !lenderExist || !lendeeExist {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User not found")
		return
	}

	t.Txs.TransactionAdd(lender, lendee, money)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
