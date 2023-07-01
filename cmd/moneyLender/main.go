package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	ml "main.go"
	"main.go/fs"
	"main.go/inmem"
)

func serveHttp(exitCh <-chan os.Signal, txc ml.TxsController) {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", txc.Login)
	mux.HandleFunc("/signup", txc.Signup)

	mux.HandleFunc("/", txc.Index)
	mux.HandleFunc("/transaction", txc.AddTransaction)

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
	var (
		usersFlag = flag.String("users", "", "path to file for user storage")
		txsFlag   = flag.String("txs", "", "path to file for transactions storage")
	)
	flag.Parse()

	var users ml.UsersStorage
	if *usersFlag == "" {
		users = inmem.NewUserStorage()
		fmt.Println("-users flag not provided, falling back to in memory storage")
	} else {
		users = fs.NewUserStorage(*usersFlag)
	}

	var txs ml.TxsStorage
	if *txsFlag == "" {
		txs = inmem.NewTxsStorage()
		fmt.Println("-txs flag not provided, falling back to in memory storage")
	} else {
		txs = fs.NewTxsStorage(*txsFlag)
	}

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)

	txc := ml.TxsController{
		Txs:   txs,
		Users: users,
	}

	serveHttp(exitCh, txc)
}
