package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	ml "github.com/pupirkaa/moneyLender"
)

type Controller struct {
	Auth       *ml.AuthService
	TxsStorage ml.TxsStorage
	Txs        ml.TxsService
	Sessions   map[string]bool
}

func (c *Controller) Login(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse body: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}

	var body struct {
		Name     string
		Password string
	}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Invalid request: %v\n", err)
	}

	session, err := c.Auth.Login(body.Name, body.Password)
	switch {
	case err == nil:
		//continue
	case errors.Is(err, ml.ErrInvalidPassword):
		w.WriteHeader(http.StatusUnauthorized)
		io.WriteString(w, "wrong name or password")
	case errors.Is(err, ml.ErrUserNotFound):
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "wrong name or password")
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Login failed: %v\n", err)
	}

	c.Sessions[session] = true
	resp, err := json.Marshal(map[string]string{"session": session, "name": body.Name})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to marshal: %v\n", err)
	}

	_, err = w.Write(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to write: %v\n", err)
	}
}

func (c *Controller) Signup(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse body: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}
	var body struct {
		Name     string
		Password string
	}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Invalid request: %v\n", err)
	}

	err = c.Auth.Signup(body.Name, body.Password)
	switch {
	case err == nil:
		//continue
	case errors.Is(err, ml.ErrInvalidSignup):
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "invalid name or password")
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Signup failed: %v\n", err)
	}

	w.WriteHeader(http.StatusCreated)
}
