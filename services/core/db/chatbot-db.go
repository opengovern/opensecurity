package db

import (
	"errors"
	"github.com/opengovern/opensecurity/services/core/db/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// --- Session Functions ---

func (db *Database) CreateSession(session *models.Session) error {
	tx := db.orm.Create(session)
	return tx.Error
}

func (db *Database) GetSession(id uuid.UUID) (*models.Session, error) {
	var session models.Session
	tx := db.orm.Preload("Chats").First(&session, "id = ?", id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &session, nil
}

func (db *Database) ListSessions() ([]models.Session, error) {
	var sessions []models.Session
	tx := db.orm.Find(&sessions)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return sessions, nil
}

// --- Chat Functions ---

func (db *Database) CreateChat(chat *models.Chat) error {
	tx := db.orm.Create(chat)
	return tx.Error
}

func (db *Database) GetChat(id uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	tx := db.orm.Preload("Session").Preload("Suggestions").Preload("ClarifyingQuestions").First(&chat, "id = ?", id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &chat, nil
}

func (db *Database) ListChats(sessionID *uuid.UUID) ([]models.Chat, error) {
	var chats []models.Chat
	query := db.orm.Model(&models.Chat{})
	if sessionID != nil {
		query = query.Where("session_id = ?", *sessionID)
	}
	tx := query.Order("created_at desc").Find(&chats)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return chats, nil
}

// --- ChatSuggestion Functions ---

func (db *Database) CreateChatSuggestion(suggestion *models.ChatSuggestion) error {
	tx := db.orm.Create(suggestion)
	return tx.Error
}

func (db *Database) GetChatSuggestion(id uuid.UUID) (*models.ChatSuggestion, error) {
	var suggestion models.ChatSuggestion
	tx := db.orm.Preload("Chat").First(&suggestion, "id = ?", id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &suggestion, nil
}

func (db *Database) ListChatSuggestions(chatID *uuid.UUID) ([]models.ChatSuggestion, error) {
	var suggestions []models.ChatSuggestion
	query := db.orm.Model(&models.ChatSuggestion{})
	if chatID != nil {
		query = query.Where("chat_id = ?", *chatID)
	}
	tx := query.Find(&suggestions)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return suggestions, nil
}

// --- ChatClarification Functions ---

func (db *Database) CreateChatClarification(clarification *models.ChatClarification) error {
	tx := db.orm.Create(clarification)
	return tx.Error
}

func (db *Database) GetChatClarification(id uuid.UUID) (*models.ChatClarification, error) {
	var clarification models.ChatClarification
	tx := db.orm.Preload("Chat").First(&clarification, "id = ?", id)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &clarification, nil
}

func (db *Database) ListChatClarifications(chatID *uuid.UUID) ([]models.ChatClarification, error) {
	var clarifications []models.ChatClarification
	query := db.orm.Model(&models.ChatClarification{})
	if chatID != nil {
		query = query.Where("chat_id = ?", *chatID)
	}
	tx := query.Order("created_at asc").Find(&clarifications)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return clarifications, nil
}
