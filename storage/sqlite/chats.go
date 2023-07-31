package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	ml "github.com/pupirkaa/moneyLender"
)

type ChatsStorage struct {
	db *sql.DB
}

func NewChatsStorage(path string) *ChatsStorage {
	dsn := "file:" + path
	d, err := sql.Open("sqlite", dsn)
	cs := &ChatsStorage{
		db: d,
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "opening db: %v\n", err)
		os.Exit(1)
	}

	_, err = cs.db.Exec("CREATE TABLE IF NOT EXISTS chats(id int, creator string, creation_date datetime, last_message_id int, FOREIGN KEY(creator) REFERENCES users(name));", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating chats table: %v\n", err)
		fmt.Println("Ошибка 1")
		os.Exit(1)
	}

	_, err = cs.db.Exec("CREATE TABLE IF NOT EXISTS messages(id int, chat_id int, creator string, creation_date datetime, message_text string, FOREIGN KEY(creator) REFERENCES users(name), FOREIGN KEY(chat_id) REFERENCES chats(id));", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating messanges table: %v\n", err)
		fmt.Println("Ошибка 2")
		os.Exit(1)
	}

	_, err = cs.db.Exec("CREATE TABLE IF NOT EXISTS chat_members(chat_id int, user_name string, FOREIGN KEY(user_name) REFERENCES users(name));", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating chat members table: %v\n", err)
		fmt.Println("Ошибка 3")
		os.Exit(1)
	}

	return cs
}

func (cs *ChatsStorage) getMessages(chatId int) ([]ml.Message, error) {
	chatMessages := []ml.Message{}
	messages, err := cs.db.Query("SELECT id, creator, creation_date, message_text FROM messages WHERE chat_id=?;", chatId)
	if err != nil {
		fmt.Println(err, "Ошибка 4")
		return []ml.Message{}, fmt.Errorf("selecting chats: %v", err)
	}
	defer messages.Close()

	for messages.Next() {
		var (
			messageId           int
			messageCreator      string
			messageCreationDate string
			messageText         string
		)

		if err := messages.Scan(&messageId, &messageCreator, &messageCreationDate, &messageText); err != nil {
			return []ml.Message{}, fmt.Errorf("getting chats: %v", err)
		}

		chatMessages = append(chatMessages, ml.Message{
			Id:      messageId,
			Message: messageText,
			Time:    messageCreationDate,
			Sender:  messageCreator,
		})
	}
	return chatMessages, nil
}

func (cs *ChatsStorage) getChatMembers(chatId int) ([]string, error) {
	chatMembers := []string{}
	members, err := cs.db.Query("SELECT user_name FROM chat_members WHERE chat_id=?;", chatId)
	if err != nil {
		fmt.Println("Ошибка 5")
		return []string{}, fmt.Errorf("selecting chats: %v", err)
	}
	defer members.Close()

	for members.Next() {
		var userName string
		if err := members.Scan(&userName); err != nil {
			return []string{}, fmt.Errorf("getting chats: %v", err)
		}
		chatMembers = append(chatMembers, userName)
	}
	return chatMembers, nil
}

func (cs *ChatsStorage) AddChat(chat ml.Chat) error {
	chatId := cs.db.QueryRow("SELECT MAX(chat_id) FROM chats;")
	var id int
	if err := chatId.Scan(&id); err != nil {
		id = 0
	}
	id += 1
	fmt.Println("id of new chat is", id)

	_, err := cs.db.Exec("INSERT INTO chats (id, creator, creation_date, last_message_id)  VALUES (?, ?, ?, ?);", id, chat.Creator, chat.CreationDate, chat.LastMessageID)
	if err != nil {
		fmt.Println("Ошибка 6")
		return fmt.Errorf("inserting into chats: %v", err)
	}
	for _, i := range chat.Members {
		_, err = cs.db.Exec("INSERT INTO chat_members (chat_id, user_name)  VALUES (?, ?);", id, i)
	}
	if err != nil {
		fmt.Println("Ошибка 7")
		return fmt.Errorf("inserting into chat members: %v", err)
	}

	return nil
}

func (cs *ChatsStorage) GetChat(chatId int) (*ml.Chat, error) {
	var (
		creator       string
		creationDate  time.Time
		lastMessageId int
		chatMembers   []string
		chatMessages  []ml.Message
	)

	chatMessages, err := cs.getMessages(chatId)
	if err != nil {
		fmt.Println("Ошибка 8")
		return &ml.Chat{}, fmt.Errorf("getting messages: %v", err)
	}

	chatMembers, err = cs.getChatMembers(chatId)
	if err != nil {
		fmt.Println("Ошибка 9")
		return &ml.Chat{}, fmt.Errorf("getting messages: %v", err)
	}

	row := cs.db.QueryRow("SELECT creator, creation_date, last_message_id FROM chats WHERE id=?;", chatId)
	if err := row.Scan(&creator, &creationDate, &lastMessageId); err != nil {
		fmt.Println(err, " Ошибка 10")
		return &ml.Chat{}, fmt.Errorf("getting chat: %v", err)
	}

	chat := ml.Chat{
		Id:            chatId,
		Creator:       creator,
		CreationDate:  creationDate,
		LastMessageID: lastMessageId,
		Members:       chatMembers,
		Messages:      chatMessages,
	}

	return &chat, nil
}

func (cs *ChatsStorage) GetAllChats() ([]ml.Chat, error) {
	chats := []ml.Chat{}

	rows, err := cs.db.Query("SELECT * FROM chats")
	if err != nil {
		fmt.Println("Ошибка 11")
		return []ml.Chat{}, fmt.Errorf("selecting chats: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id            int
			creator       string
			creationDate  time.Time
			lastMessageId int
			chatMembers   []string
			chatMessages  []ml.Message
		)
		if err := rows.Scan(&id, &creator, &creationDate, &lastMessageId); err != nil {
			return []ml.Chat{}, fmt.Errorf("getting chats: %v", err)
		}

		chatMembers, err := cs.getChatMembers(id)
		if err != nil {
			fmt.Println("Ошибка 12")
			return []ml.Chat{}, fmt.Errorf("getting chat's members: %v", err)
		}

		chatMessages, err = cs.getMessages(id)
		if err != nil {
			fmt.Println("Ошибка 13")
			return []ml.Chat{}, fmt.Errorf("getting chat's messages: %v", err)
		}

		chats = append(chats, ml.Chat{
			Id:            id,
			Creator:       creator,
			CreationDate:  creationDate,
			LastMessageID: lastMessageId,
			Members:       chatMembers,
			Messages:      chatMessages,
		})
	}
	return chats, nil
}

func (cs ChatsStorage) GetLastMessages(chatId int, lastMessageId int) ([]ml.Message, error) {
	messages := []ml.Message{}
	row, err := cs.db.Query("SELECT id, creator, creation_date, message_text FROM messages WHERE chat_id=? AND id>?;", chatId, lastMessageId)
	if err != nil {
		fmt.Println("Ошибка 14")
		return []ml.Message{}, fmt.Errorf("getting last chat's messages: %v", err)
	}
	for row.Next() {
		var (
			messageId           int
			messageCreator      string
			messageCreationDate string
			messageText         string
		)
		if err = row.Scan(&messageId, &messageCreator, &messageCreationDate, &messageText); err != nil {
			return []ml.Message{}, fmt.Errorf("getting messages: %v", err)
		}
		messages = append(messages, ml.Message{
			Id:      messageId,
			Time:    messageCreationDate,
			Sender:  messageCreator,
			Message: messageText,
		})
	}

	return messages, nil
}

func (cs *ChatsStorage) AppendMessage(chatId int, message ml.Message) error {
	chat, err := cs.GetChat(1)
	if err != nil {
		fmt.Println(err, "Ошибка ")
	}
	chat.AppendMessage(message)

	var lastMessageId int
	fmt.Println("in append function")
	row := cs.db.QueryRow("SELECT last_message_id FROM chats WHERE id=?;", chatId)
	if err := row.Scan(&lastMessageId); err != nil {
		fmt.Println(err, "Ошибка 15")
		//return fmt.Errorf("checking is session exist: %v", err)
	}

	_, err = cs.db.Exec("UPDATE chats SET last_message_id=? WHERE id=?;", lastMessageId+1, chatId)
	if err != nil {
		fmt.Println(err, "Ошибка 16")
		//return fmt.Errorf("checking is session exist: %v", err)
	}

	_, err = cs.db.Exec("INSERT INTO messages (id, chat_id, creator, creation_date, message_text) VALUES (?, ?, ?, ?, ?);", lastMessageId+1, chatId, message.Sender, message.Time, message.Message)
	if err != nil {
		fmt.Println("Ошибка 17")
		return fmt.Errorf("checking is session exist: %v", err)
	}
	return nil
}

func (cs *ChatsStorage) Update(id int, updateFn func(chat *ml.Chat) error) error {
	//ОШИБКА
	chat, err := cs.GetChat(id)
	if err != nil {
		fmt.Println("Ошибка 18")
		return fmt.Errorf("gettings chat: %w", err)
	}

	if err := updateFn(chat); err != nil {
		fmt.Println("Ошибка 19")
		return fmt.Errorf("updating chat: %w", err)
	}

	return nil
}

func (cs *ChatsStorage) Close() error {
	if err := cs.db.Close(); err != nil {
		fmt.Println("Ошибка 20")
		return fmt.Errorf("closing db: %v", err)
	}
	return nil
}
