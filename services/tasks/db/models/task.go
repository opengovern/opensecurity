package models

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ID          string `gorm:"primarykey"`
	Name        string `gorm:"unique;not null"` // Enforces uniqueness and non-null constraint
	ResultType  string
	Description string
	ImageUrl    string
	Interval    uint64
	Timeout     uint64
	NatsConfig  pgtype.JSONB
	ScaleConfig pgtype.JSONB
}
