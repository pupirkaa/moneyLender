package ml

import (
	_ "embed"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
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

		hashedPassword, err := t.Users.UserGet(name)
		if err != nil {
			panic("TODO: handle error")
		}

		if !comparaPasswords(hashedPassword, password) {
			io.WriteString(w, htmlTemplateLogin)
			io.WriteString(w, "incorrect password")
			return
		}

		t.setCookie(&w, name, hashedPassword)

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
		if !isNameValid(name) || !isPasswordValid(password) {
			io.WriteString(w, htmlTemplateSignup)
			io.WriteString(w, "wrong name or password")
			return
		}
		t.Users.UserAdd(name, HashAndSalt(password))
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}
}

func isNameValid(name string) bool {
	return !regexp.MustCompile(`\s`).MatchString(name)
}

func isPasswordValid(password string) bool {
	return (!regexp.MustCompile(`\s`).MatchString(password) || len(password) < 4)
}

func HashAndSalt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		fmt.Println("hashing password: ", err)
	}
	return string(hash)
}

func comparaPasswords(hashedPassword string, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		fmt.Println("comparing passwords: ", err)
		return false
	}
	return true
}

func (t *TxsController) setCookie(w *http.ResponseWriter, name string, password string) {
	cookieValue := hex.EncodeToString([]byte(name + password + time.Now().GoString()))
	cookie := http.Cookie{
		Name:  "user",
		Value: cookieValue,
	}
	t.Cookies[cookieValue] = true
	http.SetCookie(*w, &cookie)

	if err := cookie.Valid(); err != nil {
		fmt.Println("cookie is not valid ", err)
	}
}
