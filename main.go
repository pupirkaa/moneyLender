package main

import (
	"bufio"
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"time"

	"main.go/fs"
	"main.go/inmem"
)

//go:embed index.go.html
var htmlTemplateMain string

//go:embed login.go.html
var htmlTemplateLogin string

//go:embed signup.go.html
var htmlTemplateSignup string

func syncDataPeriodically(t TxsController) {
	time.Sleep(5 * time.Minute)
	t.saveTxsToFile()
	go syncDataPeriodically(t)
}

func serveHttp(exitCh <-chan os.Signal, txc TxsController) {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", txc.login)
	mux.HandleFunc("/signup", txc.signup)

	mux.HandleFunc("/", txc.index)
	mux.HandleFunc("/transaction", txc.addTransaction)

	srv := &http.Server{Addr: "0.0.0.0:80", Handler: mux}

	go syncDataPeriodically(txc)

	go func() {
		<-exitCh

		fmt.Println("Выключаемся :(")
		txc.saveTxsToFile()
		txc.users.Close()
		srv.Shutdown(context.TODO())
	}()

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "failed to listen and serve: %v\n", err)
			os.Exit(1)
		}
	}
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

func readFile(s string) []string {
	var fileData []string
	f, err := os.Open(s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fileData = append(fileData, scanner.Text())
	}
	return fileData
}

func parseTemplate() *template.Template {
	t, err := template.New("webpage").Parse(htmlTemplateMain)
	if err != nil {
		panic(err)
	}
	return t
}

func getResult(m map[string]int) []Debt {
	mk := make([]string, len(m))
	i := 0
	for k := range m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)

	var debts []Debt

	for i := 0; i < len(mk); i++ {
		d := Debt{
			Name:  mk[i],
			Money: m[mk[i]],
		}

		debts = append(debts, d)
	}
	return debts
}

type usersStorage interface {
	io.Closer
	UserExist(name string) bool
	UserAdd(name string, password string)
	UserGet(name string) string
}

func main() {
	var (
		usersFlag = flag.String("users", "", "path to file for user storage")
		//txsFlag   = flag.String("txs", "", "path to file for transactions storage")
	)
	flag.Parse()

	var users usersStorage
	if *usersFlag == "" {
		users = inmem.NewUserStorage()
		fmt.Println("-users flag not provided, falling back to in memory storage")
	} else {
		users = fs.NewUserStorage(*usersFlag)
	}

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)

	transactions, debts := parseTransactions(readFile(flag.Arg(0)))

	txc := TxsController{
		txs:   transactions,
		debts: getResult(debts),
		t:     parseTemplate(),
		users: users,
	}

	serveHttp(exitCh, txc)
}
