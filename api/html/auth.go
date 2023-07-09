package html

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
)

//go:embed htmlTemplates/login.go.html
var htmlTemplateLogin string

//go:embed htmlTemplates/signup.go.html
var htmlTemplateSignup string

func (t *Controller) Login(w http.ResponseWriter, req *http.Request) {
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
		form := req.Form
		var (
			name     = form.Get("name")
			password = form.Get("password")
		)

		session, err := t.Auth.Login(name, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logging in: %v\n", err)
			return
		}

		t.Sessions[session] = true
		err = t.setCookie(w, &http.Cookie{
			Name:  "user",
			Value: session,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "setting cookie: %v\n", err)
			return
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

}

func (t *Controller) Signup(w http.ResponseWriter, req *http.Request) {
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

		err := t.Auth.Signup(name, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "signing up: %v\n", err)
			return
		}

		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}
}

func (t *Controller) setCookie(w http.ResponseWriter, cookie *http.Cookie) error {
	if err := cookie.Valid(); err != nil {
		return fmt.Errorf("validating cookie: %w", err)
	}

	http.SetCookie(w, cookie)
	return nil
}
