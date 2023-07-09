package html

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	ml "github.com/pupirkaa/moneyLender"
)

//go:embed htmlTemplates/index.go.html
var htmlTemplateMain string

//go:embed htmlTemplates/distributedDebts.go.html
var htmlTemplateDebts string

type Controller struct {
	Txs        ml.TxsService
	TxsStorage ml.TxsStorage
	Auth       ml.AuthService
	Sessions   map[string]bool
}

func (t *Controller) Index(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("user")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if _, ok := t.Sessions[cookie.Value]; !ok {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	txs, err := t.TxsStorage.TxsGet()
	if err != nil {
		fmt.Println("getting transactions: ", err)
	}

	debts, err := t.TxsStorage.DebtsGet()
	if err != nil {
		fmt.Println("getting debts: ", err)
	}

	io.WriteString(w, GenerateHTML(txs, debts, ParseTemplate(htmlTemplateMain)))
}

func (t *Controller) DistributedDebts(w http.ResponseWriter, req *http.Request) {
	debts, err := t.TxsStorage.DebtsGet()
	if err != nil {
		fmt.Println("getting debts: ", err)
	}
	io.WriteString(w, GenerateHTML(ml.DistributeDebts(debts), nil, ParseTemplate(htmlTemplateDebts)))
}

func GenerateHTML(transactions []ml.Transaction, debts []ml.Debt, t *template.Template) string {
	tad := ml.TransactionsAndDebts{
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

func ParseTemplate(s string) *template.Template {
	t, err := template.New("webpage").Parse(s)
	if err != nil {
		panic(err)
	}
	return t
}

func RunCookieStorage() map[string]bool {
	return map[string]bool{}
}
