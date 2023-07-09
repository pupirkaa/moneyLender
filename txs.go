package ml

import "fmt"

type TxsService struct {
	Txs   TxsStorage
	Users UsersStorage
}

type TransactionsAndDebts struct {
	Transactions []Transaction
	Debts        []Debt
}

type Transaction struct {
	Lender string
	Lendee string
	Money  int
}

type Debt struct {
	Name  string
	Money int
}

func (s *TxsService) TxAdd(lender string, lendee string, money int) error {
	lenderExist, err := s.Users.UserExist(lender)
	if err != nil {
		return fmt.Errorf("getting lender:%v", err)
	}

	lendeeExist, err := s.Users.UserExist(lendee)
	if err != nil {
		return fmt.Errorf("getting lendee:%v", err)
	}

	if !lenderExist || !lendeeExist {
		return ErrUserNotFound
	}

	if err := s.Txs.TransactionAdd(lender, lendee, money); err != nil {
		return fmt.Errorf("adding transaction:%v", err)
	}

	return nil
}
