package model

import (
	"time"

	"gorm.io/gorm"
)

type Translation struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	SourceText      string         `json:"source_text" gorm:"type:text;not null"`
	TranslatedText  string         `json:"translated_text" gorm:"type:text;not null"`
	SourceLanguage  string         `json:"source_language" gorm:"type:varchar(10);not null"`
	TargetLanguage  string         `json:"target_language" gorm:"type:varchar(10);not null"`
	Context         string         `json:"context" gorm:"type:text"`
	Category        string         `json:"category" gorm:"type:varchar(50)"`
	Votes           int           `json:"votes" gorm:"default:0"`
	CreatedBy       string         `json:"created_by" gorm:"type:varchar(255)"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Translation) TableName() string {
	return "translations"
} 