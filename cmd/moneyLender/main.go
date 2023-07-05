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

	ml "main.go"
	"main.go/fs"
	"main.go/inmem"
	"main.go/sqlite"
)

func serveHttp(exitCh <-chan os.Signal, txc ml.TxsController) {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", txc.Login)
	mux.HandleFunc("/signup", txc.Signup)

	mux.HandleFunc("/", txc.Index)
	mux.HandleFunc("/transaction", txc.AddTransaction)
	mux.HandleFunc("/distributedDebts", txc.DistributedDebts)

	srv := &http.Server{Addr: "0.0.0.0:80", Handler: mux}

	go func() {
		<-exitCh

		fmt.Println("Выключаемся :(")
		txc.Txs.Close()
		txc.Users.Close()
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

	txc := ml.TxsController{
		Txs:     txs,
		Users:   users,
		Cookies: map[string]bool{},
	}

	serveHttp(exitCh, txc)
}

const (
	mimetypeSqlite = "application/vnd.sqlite3"
)
