package models

import (
	"fmt"

	"github.com/matheusm25/thoth/src/modules/db/config"
)

type Message struct {
	ID      string
	Topic   string
	Payload string
}

func (m *Message) CreateTable() {
	db, err := config.GetDB()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS messages (
		id TEXT PRIMARY KEY,
		topic TEXT NOT NULL,
		payload TEXT NOT NULL
	)`)
	if err != nil {
		panic(err)
	}
}

func GetMessageByTopic(topic string) ([]Message, error) {
	db, err := config.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	rows, err := db.Query("SELECT * FROM messages WHERE topic = ?", topic)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.Topic, &message.Payload); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}

func (m *Message) Save() error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	_, err = db.Exec("INSERT INTO messages (id, topic, payload) VALUES (?, ?, ?)", m.ID, m.Topic, m.Payload)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return nil
}

func SaveMessage(message Message) error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	_, err = db.Exec("INSERT INTO messages (id, topic, payload) VALUES (?, ?, ?)", message.ID, message.Topic, message.Payload)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return nil
}

func (m *Message) Delete() error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	_, err = db.Exec("DELETE FROM messages WHERE id = ?", m.ID)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func DeleteMessageById(id string) error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	_, err = db.Exec("DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func SaveManyMessages(messages []Message) error {
	db, err := config.GetDB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	stmt, err := tx.Prepare("INSERT INTO messages (id, topic, payload) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, message := range messages {
		if _, err := stmt.Exec(message.ID, message.Topic, message.Payload); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert message: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetAllMessages() ([]Message, error) {
	db, err := config.GetDB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	rows, err := db.Query("SELECT * FROM messages")
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.Topic, &message.Payload); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}

func ListMessageTopics() []string {
	db, err := config.GetDB()
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT DISTINCT topic FROM messages")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var topics []string
	for rows.Next() {
		var topic string
		if err := rows.Scan(&topic); err != nil {
			panic(err)
		}
		topics = append(topics, topic)
	}

	return topics
}
