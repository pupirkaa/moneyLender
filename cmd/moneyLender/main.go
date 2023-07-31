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
	mux.HandleFunc("/api/logout", jc.Logout)
	mux.HandleFunc("/api/transaction", jc.AddTransaction)
	mux.HandleFunc("/api/txs", jc.GetTxs)
	mux.HandleFunc("/api/debts", jc.GetDebts)
	mux.HandleFunc("/api/result", jc.GetDistributedDebts)

	mux.HandleFunc("/login", c.Login)
	mux.HandleFunc("/signup", c.Signup)
	mux.HandleFunc("/logout", c.Logout)

	mux.HandleFunc("/", c.Index)
	mux.HandleFunc("/transaction", c.AddTransaction)
	mux.HandleFunc("/distributedDebts", c.DistributedDebts)

	mux.HandleFunc("/chats", c.ViewChatList)
	mux.HandleFunc("/chat1", c.UseChat)
	mux.HandleFunc("/message", c.AddMessage)
	mux.HandleFunc("/update", c.UpdateMesages)

	srv := &http.Server{Addr: "0.0.0.0:80", Handler: mux}

	go func() {
		<-exitCh

		fmt.Println("Выключаемся :(")
		c.TxsStorage.Close()
		c.Auth.Users.Close()
		c.Auth.Sessions.Close()
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
		storageFlag  = flag.String("storage", "", "kind of storage")
		sqliteFlag   = flag.String("sqlite", "", "path to file for db storage")
		usersFlag    = flag.String("users", "", "path to file for user storage")
		txsFlag      = flag.String("txs", "", "path to file for transactions storage")
		sessionsFlag = flag.String("sessions", "", "path to file for session storage")
	)
	flag.Parse()

	var (
		users    ml.UsersStorage
		txs      ml.TxsStorage
		sessions ml.SessionsStorage
	)

	//TODO:
	// Сделать дефолт или выводить, что флаг неверный

	switch *storageFlag {
	case "inmem":
		users = inmem.NewUserStorage()
		txs = inmem.NewTxsStorage()
		sessions = inmem.NewSessionsStorage()
	case "fs":
		users = fs.NewUserStorage(*usersFlag)
		txs = fs.NewTxsStorage(*txsFlag)
		sessions = fs.NewSessionsStorage(*sessionsFlag)
	case "sqlite":
		users = sqlite.NewUserStorage(*sqliteFlag)
		txs = sqlite.NewTxsStorage(*sqliteFlag)
		sessions = sqlite.NewSessionsStorage(*sqliteFlag)
	}

	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)

	hc := html.Controller{
		TxsStorage: txs,
		Auth:       ml.NewAuthServise(users, sessions, exitCh),
		Txs: ml.TxsService{
			Users: users,
			Txs:   txs,
		},
		Sessions: sessions,
		Chat:     sqlite.NewChatsStorage("fixtures/us.db"),
		//Chat: &inmem.ChatsStorage{},
	}
	//hc.Chat.AddChat(*ml.NewChat("Sima", []string{"Irina"}))

	jc := json.Controller{
		Auth: &ml.AuthService{
			Users:    users,
			Sessions: sessions,
		},
		TxsStorage: txs,
		Txs: ml.TxsService{
			Users: users,
			Txs:   txs,
		},
		Sessions: sessions,
	}

	serveHttp(exitCh, hc, jc)

}

const (
	mimetypeSqlite = "application/vnd.sqlite3"
)
