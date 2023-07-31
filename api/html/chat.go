package html

import (
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	ml "github.com/pupirkaa/moneyLender"
)

//go:embed htmlTemplates/chatlist.go.html
var htmlChatList string

//go:embed htmlTemplates/chat.go.html
var htmlChat string

var tChat = ParseTemplate(htmlChat)
var tChatList = ParseTemplate(htmlChatList)

func (c *Controller) UseChat(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("user")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if _, ok := c.Sessions.SessionExist(cookie.Value); !ok {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	if req.Method == http.MethodGet {
		chat, _ := c.Chat.GetChat(1)

		strB := strings.Builder{}
		err := tChat.ExecuteTemplate(&strB, "Full", chat.Messages)

		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate html: %v", err)
			os.Exit(1)
		}

		io.WriteString(w, strB.String())

	}
}

func (c *Controller) AddMessage(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("user")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if _, ok := c.Sessions.SessionExist(cookie.Value); !ok {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	sessions, _ := c.Sessions.GetSessions()
	name := sessions[cookie.Value].Name

	err = req.ParseForm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse form: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "bad form")
		return
	}

	form := req.Form
	var message = form.Get("message")
	strB := strings.Builder{}
	c.Chat.Update(1, func(chat *ml.Chat) error {
		c.Chat.AppendMessage(1, ml.NewMessage(message, name))
		fmt.Println("Appended message to chat")
		ch, _ := c.Chat.GetChat(1)
		m, _ := c.Chat.GetLastMessages(1, ch.LastMessageID-1)
		err = tChat.ExecuteTemplate(&strB, "Part", m)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate html: %v", err)
			os.Exit(1)
		}
		return nil
	})

	io.WriteString(w, strB.String())
}

// Если в чате более чем один пользователь, то лишь только первый пользователь сделавший запрос на новые сообщения
// и поменявший, соответственно флаг "newMessages", их получит. Другим же пользователям
// будет казаться, будто новых сообщений нет.
// 1. Пользователь отправляет все сообщения, которые есть у него, а потом мы сравниваем каких у него нет
// 2. Пользователь присылает только последнее сообщение, и мы возвращаем ему все сообщения, дата отправки которых позже
// 3. Пользователь отправляет дату последнего полученного сообщения, а мы возвращаем ему все сообщения
// 4. Пользователь отправляет идентификатор последнего отправленного сообщения, а мы ему возвращаем те, которые поновее - V
func (c *Controller) UpdateMesages(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		id, _ := strconv.Atoi(req.URL.Query().Get("lastMessageId"))
		fmt.Println("id of the last message in the chat is ", id)
		chat, _ := c.Chat.GetChat(1)
		for i := 0; i < 10; i++ {
			if chat.LastMessageID > id {
				fmt.Println("Нашли новое сообщение")
				strB := strings.Builder{}

				m, err := c.Chat.GetLastMessages(chat.Id, id)
				fmt.Println("get messages from ", id, "to", chat.LastMessageID)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to get last messages: %v", err)
				}

				err = tChat.ExecuteTemplate(&strB, "Part", m)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to generate html: %v", err)
					os.Exit(1)
				}
				io.WriteString(w, strB.String())
				return
			}
			fmt.Println("Спим")
			time.Sleep(time.Second)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func (c *Controller) ViewChatList(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("user")
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	if _, ok := c.Sessions.SessionExist(cookie.Value); !ok {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	allChats, _ := c.Chat.GetAllChats()
	strB := strings.Builder{}

	err = tChatList.Execute(&strB, allChats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate html: %v", err)
		os.Exit(1)
	}
	io.WriteString(w, strB.String())
}
