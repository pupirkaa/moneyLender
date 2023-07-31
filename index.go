package ml

import (
	_ "embed"
	"fmt"
	"io"
	"sort"
	"time"
)

//	Когда я вошла в систему, то на главной странице у меня список групп, в которых я состою
// 	Я могу добавить новую группу (для этого в неё нужно добавить как минимум одного участника)
//	Я ввожу уникальный код участника, у него появляется приглашение от меня на главной и
//	он может подтвердить/отказаться. До этого момента он не является участником группы
//	Для группы я должна выбрать валюту, чтобы транзакции внутри группы были в одной валюте
//	Я могу удалить группу, тогда она пропадёт у меня, но не у остальных участников
//	(и транзакции со мной будет нельзя добавлять)
//	Я могу перейти в группу и тогда увижу список всех транзаций между участниками + задолжности каждого
//	Я могу добавить транзакцию - между двумя людьми или между человеком и группой людей
//	(где долг распределяется равномерно)
//	Я могу нажать на "Распределить долги" и перейду на страницу, где мне будет показано,
//	кто сколько кому должен отдать
// 	На этом же моменте я хочу интегрировать страйп и сделать возможность каждому перевести долг
//	С этой вкладки я могу вернуться назад
//	Внутри группы есть заранее созданный чат, где каждый участник может присылать текстовые сообщения
//	Также с главной страницы я могу перейти в свои чаты, где есть отдельные чаты с разными людьми

type ChatStorage interface {
	io.Closer
	AddChat(chat Chat) error
	GetChat(id int) (*Chat, error)
	GetAllChats() ([]Chat, error)
	Update(id int, updateFn func(chat *Chat) error) error
	GetLastMessages(chatId int, lastMessageId int) ([]Message, error)
	AppendMessage(chatId int, message Message) error
}

type UsersStorage interface {
	io.Closer
	UserExist(name string) (bool, error)
	UserAdd(name string, password string) error
	UserGet(name string) (string, error)
}

type TxsStorage interface {
	io.Closer
	TransactionAdd(lender string, lendee string, money int) error
	DebtsGet() ([]Debt, error)
	TxsGet() ([]Transaction, error)
}

type SessionsStorage interface {
	io.Closer
	SessionExist(session string) (error, bool)
	AddSession(session string, name string, creationDate time.Time) error
	DeleteSession(session string) error
	GetSessions() (map[string]SessionArgs, error)
}

func DistributeDebts(debts []Debt) []Transaction {
	txs := []Transaction{}
	posDebts := []Debt{}
	negDebts := []Debt{}

	sort.Slice(debts, func(i, j int) bool { return debts[i].Money > debts[j].Money })
	for _, d := range debts {
		if d.Money > 0 {
			posDebts = append(posDebts, Debt{Name: d.Name, Money: d.Money})
		} else {
			negDebts = append(negDebts, Debt{Name: d.Name, Money: d.Money})
		}
	}

	sort.SliceStable(negDebts, func(i, j int) bool {
		return i > j
	})

	fmt.Println("+", posDebts, " -", negDebts)
	for i, k := 0, 0; i < len(posDebts) && k < len(negDebts); {
		if posDebts[i].Money > (-negDebts[k].Money) {
			txs = append(txs, Transaction{Lender: negDebts[k].Name, Lendee: posDebts[i].Name, Money: -negDebts[k].Money})
			posDebts[i].Money += negDebts[k].Money
			negDebts[k].Money = 0
			k += 1
		} else if posDebts[i].Money < (-negDebts[k].Money) {
			txs = append(txs, Transaction{Lender: negDebts[k].Name, Lendee: posDebts[i].Name, Money: posDebts[i].Money})
			negDebts[k].Money += posDebts[i].Money
			posDebts[i].Money = 0
			i += 1
		} else {
			txs = append(txs, Transaction{Lender: negDebts[k].Name, Lendee: posDebts[i].Name, Money: posDebts[i].Money})
			negDebts[k].Money = 0
			posDebts[i].Money = 0
			k += 1
			i += 1
		}
	}

	return txs
}
