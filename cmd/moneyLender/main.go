package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"

	ml "github.com/pupirkaa/moneyLender"
	"github.com/pupirkaa/moneyLender/api/html"
	"github.com/pupirkaa/moneyLender/api/json"
	"github.com/pupirkaa/moneyLender/storage/fs"
	"github.com/pupirkaa/moneyLender/storage/inmem"
	"github.com/pupirkaa/moneyLender/storage/sqlite"
)

func serveHttp(exitCh <-chan os.Signal, c html.Controller, jc json.Controller) {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", jc.Login)
	mux.HandleFunc("/api/signup", jc.Signup)
	mux.HandleFunc("/api/transaction", jc.AddTransaction)
	mux.HandleFunc("/api/txs", jc.GetTxs)
	mux.HandleFunc("/api/debts", jc.GetDebts)
	mux.HandleFunc("/api/result", jc.GetDistributedDebts)

	mux.HandleFunc("/login", c.Login)
	mux.HandleFunc("/signup", c.Signup)

	mux.HandleFunc("/", c.Index)
	mux.HandleFunc("/transaction", c.AddTransaction)
	mux.HandleFunc("/distributedDebts", c.DistributedDebts)

	srv := &http.Server{Addr: "0.0.0.0:80", Handler: mux}

	go func() {
		<-exitCh

		fmt.Println("Выключаемся :(")
		c.TxsStorage.Close()
		c.Auth.Users.Close()
		srv.Shutdown(context.TODO())
	}()

	if err := srv.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "failed to listen and serve: %v\n", err)
			os.Exit(1)
		}
	}
}

func main() {
	mime.AddExtensionType(".db", mimetypeSqlite)
	_ = mime.TypeByExtension(filepath.Ext(""))

	var (
		storageFlag = flag.String("storage", "", "kind of storage")
		sqliteFlag  = flag.String("sqlite", "", "path to file for db storage")
		usersFlag   = flag.String("users", "", "path to file for user storage")
		txsFlag     = flag.String("txs", "", "path to file for transactions storage")
	)
	flag.Parse()

	var (
		users ml.UsersStorage
		txs   ml.TxsStorage
	)

	switch *storageFlag {
	case "inmem":
		users = inmem.NewUserStorage()
		txs = inmem.NewTxsStorage()
	case "fs":
		users = fs.NewUserStorage(*usersFlag)
		txs = fs.NewTxsStorage(*txsFlag)
	case "sqlite":
		users = sqlite.NewUserStorage(*sqliteFlag)
		txs = sqlite.NewTxsStorage(*sqliteFlag)

	}

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)

	txc := html.Controller{
		TxsStorage: txs,
		Auth: ml.AuthService{
			Users: users,
		},
		Sessions: map[string]bool{},
	}

	jc := json.Controller{
		Auth: &ml.AuthService{
			Users: users,
		},
		TxsStorage: txs,
		Txs: ml.TxsService{
			Users: users,
			Txs:   txs,
		},
		Sessions: map[string]bool{},
	}

	serveHttp(exitCh, txc, jc)
}

const (
	mimetypeSqlite = "application/vnd.sqlite3"
)
