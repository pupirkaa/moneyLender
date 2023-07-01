package ml

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
)

//go:embed login.go.html
var htmlTemplateLogin string

//go:embed signup.go.html
var htmlTemplateSignup string

func (t *TxsController) Login(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse form: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}

	if req.Method == http.MethodGet {
		io.WriteString(w, htmlTemplateLogin)
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

		form := req.Form
		var (
			name     = form.Get("name")
			password = form.Get("password")
		)

		if !t.Users.UserExist(name) {
			io.WriteString(w, htmlTemplateLogin)
			io.WriteString(w, "can't find a user")
			return
		}

		if t.Users.UserGet(name) != password {
			io.WriteString(w, htmlTemplateLogin)
			io.WriteString(w, "incorrect password")
			return
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

}

func (t *TxsController) Signup(w http.ResponseWriter, req *http.Request) {
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

		t.Users.UserAdd(form.Get("name"), form.Get("password"))
	}
}
