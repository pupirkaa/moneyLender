package inmem

import (
	"fmt"
	"sync"
	"time"

	ml "github.com/pupirkaa/moneyLender"
)

type chat struct {
	id           int
	mu           sync.Mutex
	creationDate time.Time
	creator      string
	members      []string
	msgs         []ml.Message
	lastMsgId    int
}

type ChatsStorage struct {
	Chats []chat
	mutex sync.Mutex
}

func NewChatsStorage(creator string, members []string) ChatsStorage {
	cs := ChatsStorage{}
	return cs
}

func chatFromDomain(c *ml.Chat) *chat {
	if c == nil {
		return nil
	}

	return &chat{
		id:           c.Id,
		mu:           sync.Mutex{},
		creationDate: c.CreationDate,
		creator:      c.Creator,
		members:      c.Members,
		msgs:         c.Messages,
		lastMsgId:    c.LastMessageID,
	}
}

func (c *chat) Domain() *ml.Chat {
	if c == nil {
		return nil
	}

	msgs := make([]ml.Message, 0, len(c.msgs))
	msgs = append(msgs, c.msgs...)

	return &ml.Chat{
		Id:            c.id,
		Messages:      msgs,
		CreationDate:  c.creationDate,
		Creator:       c.creator,
		Members:       c.members,
		LastMessageID: c.lastMsgId,
	}
}

func (cs *ChatsStorage) AddChat(chat ml.Chat) error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()
	dchat := chatFromDomain(&chat)
	dchat.id = len(cs.Chats) + 1
	cs.Chats = append(cs.Chats, *dchat)
	return nil
}

func (cs *ChatsStorage) GetChat(id int) (*ml.Chat, error) {
	for _, i := range cs.Chats {
		if i.id == id {
			fmt.Println("got chat")
			return i.Domain(), nil
		}
	}
	return &ml.Chat{}, nil
}

func (cs *ChatsStorage) Update(id int, updateFn func(chat *ml.Chat) error) error {
	chat, err := cs.GetChat(id)
	if err != nil {
		return fmt.Errorf("gettings chat: %w", err)
	}

	if err := updateFn(chat); err != nil {
		return fmt.Errorf("updating chat: %w", err)
	}

	if err := cs.saveChat(chat); err != nil {
		return fmt.Errorf("saving chat: %w", err)
	}

	return nil
}

func (cs *ChatsStorage) saveChat(new *ml.Chat) error {
	for i, old := range cs.Chats {
		if old.id == new.Id {
			cs.Chats[i] = *chatFromDomain(new)
			return nil
		}
	}

	return fmt.Errorf("chat not found")
}

func (cs ChatsStorage) GetAllChats() ([]ml.Chat, error) {
	chats := []ml.Chat{}
	for _, i := range cs.Chats {
		chats = append(chats, *i.Domain())
	}
	return chats, nil
}

func (cs ChatsStorage) GetLastMessages(chatId int, lastMessageId int) ([]ml.Message, error) {
	chat, err := cs.GetChat(chatId)
	if err != nil {
		return []ml.Message{}, fmt.Errorf("can't get a chat")
	}

	messages := []ml.Message{}
	fmt.Println("last message ID", lastMessageId, "chat len", chat.LastMessageID)

	for i := lastMessageId; i < chat.LastMessageID; i++ {
		messages = append(messages, chat.Messages[i])
	}

	return messages, nil
}

func (cs *ChatsStorage) AppendMessage(chatId int, message ml.Message) error {
	chat, err := cs.GetChat(chatId)
	if err != nil {
		return fmt.Errorf("can't get a chat")
	}
	chat.AppendMessage(message)
	for i := range cs.Chats {
		if cs.Chats[i].id == chatId {
			cs.Chats[i].msgs = append(cs.Chats[i].msgs, message)
			cs.Chats[i].lastMsgId += 1
		}
	}
	return nil
}

func (cs ChatsStorage) Close() error {
	return nil
}
