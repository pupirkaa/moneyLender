package ml

import (
	"time"
)

//	Как работает мой чат?
// 	У меня есть вкладка чаты - там я вижу все доступные чаты и могу перейти в них
//	Каждый чат в списке назван <Name>, <Name>, <Creatin date>
//	Также я могу добавить новый чат, указав с кем он будет
//	При добавлении чат появляется в списке всех чатов и я могу перейти в него
// 	При переходе в чат я вижу всю историю сообщений в формате <Name>, <time>: <Message>
//	Также у меня есть текстовое поле с возможностью послать новое сообщение

type Message struct {
	Time    string
	Message string
	Sender  string
	Id      int
}

func NewMessage(message string, sender string) Message {
	m := Message{
		Time:    time.Now().Format(time.Kitchen),
		Message: message,
		Sender:  sender,
	}
	return m
}

type Chat struct {
	Id            int
	Messages      []Message
	CreationDate  time.Time
	Creator       string
	Members       []string
	LastMessageID int
}

func NewChat(creator string, members []string) *Chat {
	c := Chat{
		Messages:      []Message{},
		CreationDate:  time.Now(),
		Creator:       creator,
		Members:       members,
		LastMessageID: 0,
	}

	/*c.name = creator + ", "
	for i := 0; i < len(members); i++ {
		c.name = c.name + members[i] + ", "
	}
	c.name = c.name + c.creationDate.GoString()*/

	return &c

}

func (c *Chat) GetName() string {
	name := c.Creator + ", "
	for i := 0; i < len(c.Members); i++ {
		name = name + c.Members[i] + ", "
	}
	name = name + c.CreationDate.Format(time.Kitchen)
	return name
}

func (c *Chat) AppendMessage(m Message) {
	m.Id = c.LastMessageID + 1
	c.Messages = append(c.Messages, m)
	c.LastMessageID += 1
}
