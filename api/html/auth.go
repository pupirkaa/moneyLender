package html

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

//go:embed htmlTemplates/login.go.html
var htmlTemplateLogin string

//go:embed htmlTemplates/signup.go.html
var htmlTemplateSignup string

func (c *Controller) Login(w http.ResponseWriter, req *http.Request) {
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

		session, err := c.Auth.Login(name, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "logging in: %v\n", err)
			return
		}

		err = c.setCookie(w, &http.Cookie{
			Name:    "user",
			Value:   session,
			Expires: time.Now().Add(24 * time.Hour),
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "setting cookie: %v\n", err)
			return
		}

		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

}

func (c *Controller) Signup(w http.ResponseWriter, req *http.Request) {
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

		err := c.Auth.Signup(name, password)
		if err != nil {
			fmt.Fprintf(os.Stderr, "signing up: %v\n", err)
			return
		}

		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}
}

func (c *Controller) Logout(w http.ResponseWriter, req *http.Request) {
	session, err := req.Cookie("user")
	if err != nil {
		io.WriteString(w, "failed to get cookie")
		return
	}
	if err := c.Auth.Logout(session.Value); err != nil {
		io.WriteString(w, "failed to logout")
		return
	}

	http.Redirect(w, req, "/login", http.StatusSeeOther)
}

func (c *Controller) setCookie(w http.ResponseWriter, cookie *http.Cookie) error {
	if err := cookie.Valid(); err != nil {
		return fmt.Errorf("validating cookie: %w", err)
	}

	http.SetCookie(w, cookie)
	return nil
}
