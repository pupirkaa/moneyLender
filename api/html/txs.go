package html

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	ml "github.com/pupirkaa/moneyLender"
)

func (t *Controller) AddTransaction(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse form: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}
	form := req.Form

	var (
		lender   = form.Get("lender")
		lendee   = form.Get("lendee")
		money, _ = strconv.Atoi(form.Get("money"))
	)

	switch err := t.Txs.TxAdd(lender, lendee, money); {
	case err == nil:
		//continue
	case errors.Is(err, ml.ErrUserNotFound):
		fmt.Fprintf(os.Stderr, "Failed find user: %v", err)
		return
	default:
		fmt.Fprintf(os.Stderr, "Failed to add transaction: %v", err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
