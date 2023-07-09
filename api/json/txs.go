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

func (t *Controller) AddTransaction(w http.ResponseWriter, req *http.Request) {
	rawBody, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse body: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}

	var body struct {
		Lender string
		Lendee string
		Money  int
	}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Invalid request: %v\n", err)
	}

	switch err := t.Txs.TxAdd(body.Lender, body.Lendee, body.Money); {
	case err == nil:
		//continue
	case errors.Is(err, ml.ErrUserNotFound):
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(os.Stderr, "Failed find user: %v", err)
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Failed to add transaction: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (c *Controller) GetTxs(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
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
		Session string
	}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Invalid request: %v\n", err)
	}

	if _, ok := c.Sessions.SessionExist(body.Session); !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	txs, err := c.TxsStorage.TxsGet()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Getting transactions failed: %v\n", err)
	}

	resp, err := json.Marshal(txs)
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

func (c *Controller) GetDebts(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
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
		Session string
	}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Invalid request: %v\n", err)
	}

	if _, ok := c.Sessions.SessionExist(body.Session); !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	debts, err := c.TxsStorage.DebtsGet()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Getting debts failed: %v\n", err)
	}

	resp, err := json.Marshal(debts)
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

func (c *Controller) GetDistributedDebts(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
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
		Session string
	}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(os.Stderr, "Invalid request: %v\n", err)
	}

	if _, ok := c.Sessions.SessionExist(body.Session); !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	debts, err := c.TxsStorage.DebtsGet()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Getting debts failed: %v\n", err)
	}

	resp, err := json.Marshal(ml.DistributeDebts(debts))
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
