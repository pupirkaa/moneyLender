package ml

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"sort"
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

//go:embed distributedDebts.go.html
var htmlTemplateDebts string

func (t *TxsController) Index(w http.ResponseWriter, req *http.Request) {
	_, err := req.Cookie("user")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	io.WriteString(w, generateHTML(t.Txs.TxsGet(), t.Txs.DebtsGet(), parseTemplate(htmlTemplateMain)))
}

func (t *TxsController) DistributedDebts(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, generateHTML(t.DistributeDebts(), nil, parseTemplate(htmlTemplateDebts)))
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

func parseTemplate(s string) *template.Template {
	t, err := template.New("webpage").Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}

func (t TxsController) DistributeDebts() []Transaction {
	txs := []Transaction{}
	debts := t.Txs.DebtsGet()
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
		fmt.Println("+", posDebts, " -", negDebts)
	}

	return txs
}
