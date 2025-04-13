package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID      uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4();"`
	AgentID string    `gorm:"type:varchar(100);not null"`
	Chats   []Chat    `gorm:"foreignKey:SessionID"`
}

type Chat struct {
	gorm.Model
	ID                  uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4();"`
	Question            string
	Query               *string
	QueryError          *string
	NeedClarification   bool `gorm:"default:false"`
	AssistantText       *string
	TimeTaken           *float64
	AgentID             *string
	Result              pgtype.JSONB
	SessionID           uuid.UUID
	Session             Session             `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE"`
	Suggestions         []ChatSuggestion    `gorm:"foreignKey:ChatID"`
	ClarifyingQuestions []ChatClarification `gorm:"foreignKey:ChatID"`
}

type ChatSuggestion struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Suggestion string    `gorm:"type:text;not null"`
	ChatID     uuid.UUID `gorm:"type:uuid;not null"`
	Chat       Chat      `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}

type ChatClarification struct {
	gorm.Model
	ID        uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4();"`
	Questions string
	Answer    *string
	ChatID    uuid.UUID `gorm:"type:uuid;not null"`
	Chat      Chat      `gorm:"foreignKey:ChatID;constraint:OnDelete:CASCADE"`
}

type ChatbotSecret struct {
	Key    string `gorm:"primaryKey"`
	Secret string
}

type CachedQuery struct {
	gorm.Model
	ID    uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4();"`
	Query string
}
