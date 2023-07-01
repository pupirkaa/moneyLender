package ml

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

type UsersStorage interface {
	io.Closer
	UserExist(name string) bool
	UserAdd(name string, password string)
	UserGet(name string) string
}

type TxsStorage interface {
	io.Closer
	TransactionAdd(lender string, lendee string, money int)
	DebtsGet() []Debt
	TxsGet() []Transaction
}

//go:embed index.go.html
var htmlTemplateMain string

func (t *TxsController) Index(w http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("user")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	io.WriteString(w, generateHTML(t.Txs.TxsGet(), t.Txs.DebtsGet(), parseTemplate()))
}

func generateHTML(transactions []Transaction, debts []Debt, t *template.Template) string {
	tad := TransactionsAndDebts{
		Transactions: transactions,
		Debts:        debts,
	}
	strB := strings.Builder{}

	err := t.Execute(&strB, tad)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate html: %v", err)
		os.Exit(1)
	}

	return strB.String()
}

func parseTemplate() *template.Template {
	t, err := template.New("webpage").Parse(htmlTemplateMain)
	if err != nil {
		panic(err)
	}
	return t
}
