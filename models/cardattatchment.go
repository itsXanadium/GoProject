package models

import (
	"time"

	"github.com/google/uuid"
)

type CardAttatchment struct {
	InternalID int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID `json:"public_id" gorm:"type:uuid;not null"`
	File       string    `json:"file" db:"file"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`

	FileURL string `json:"file_url" gorm:"-"`
}

func (CardAttatchment) TableName() string {
	return "card_attachment"
}
