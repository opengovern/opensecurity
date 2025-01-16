package model

import (
	"gorm.io/gorm"
)

type FrameworkValidation struct {
	gorm.Model
	FrameworkID    string
	FailureMessage string `gorm:"type:text"`
}
