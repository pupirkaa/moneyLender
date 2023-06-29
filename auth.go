package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func parseUsers(s []string) map[string]string {
	users := make(map[string]string)

	for i := range s {
		k := strings.Split(s[i], ",")
		users[k[0]] = k[1]
	}
	return users
}

func (t TxsController) saveUsersToFile() {
	f, err := os.OpenFile("users", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file: %v", err)
		os.Exit(1)
	}
	for i := range t.newUsers {
		_, err = f.WriteString(fmt.Sprintf("\n%v,%v", i, t.newUsers[i]))
		delete(t.newUsers, t.newUsers[i])
		if err != nil {
			fmt.Println(err)
		}
	}
	defer f.Close()
}

func (t *TxsController) login(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w, htmlTemplateLogin)

		err := req.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse form: %v", err)

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "bad form")
			return
		}

		form := req.Form
		var (
			name     = form.Get("name")
			password = form.Get("password")
		)

		upass, ok := t.users[name]
		if !ok {
			fmt.Fprintln(w, "can't find a user")
			return
		}

		if upass != password {
			fmt.Fprintln(w, "incorrect password")
			return
		}

		return
	}
	if req.Method == http.MethodPost {
		cookie := http.Cookie{
			Name:  "user",
			Value: "",
		}
		http.SetCookie(w, &cookie)
		err := cookie.Valid()
		if err != nil {
			fmt.Println("cookie is not valid ", err)
		}
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

}

func (t *TxsController) signup(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, htmlTemplateSignup)

	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse form: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}

	if req.Method == http.MethodPost {
		form := req.Form
		t.newUsers[form.Get("name")] = form.Get("password")
		t.users[form.Get("name")] = form.Get("password")
	}
}
