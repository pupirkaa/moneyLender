package ml

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
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

		exists, err := t.Users.UserExist(name)
		if err != nil {
			panic("TODO: handle error")
		}

		if !exists {
			io.WriteString(w, htmlTemplateLogin)
			io.WriteString(w, "can't find a user")
			return
		}

		pass, err := t.Users.UserGet(name)
		if err != nil {
			panic("TODO: handle error")
		}

		if pass != password {
			io.WriteString(w, htmlTemplateLogin)
			io.WriteString(w, "incorrect password")
			return
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

}

func (t *TxsController) Signup(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w, htmlTemplateSignup)
		return
	}

	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse form: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}

	if req.Method == http.MethodPost {
		form := req.Form
		name := form.Get("name")
		password := form.Get("password")
		if !isNameCorrect(name) || !isPasswordCorrect(password) {
			io.WriteString(w, htmlTemplateSignup)
			io.WriteString(w, "wrong name or password")
			return
		}
		t.Users.UserAdd(name, password)
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}
}

func isNameCorrect(name string) bool {
	return !regexp.MustCompile(`\s`).MatchString(name)
}

func isPasswordCorrect(password string) bool {
	return (!regexp.MustCompile(`\s`).MatchString(password) || len(password) < 4)
}
