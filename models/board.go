package models

import (
	"time"

	"github.com/google/uuid"
)

type Board struct {
	InternalID    int64      `json:"internal_id" db:"board_id" gorm:"primaryKey"`
	PublicID      uuid.UUID  `json:"public_id" db:"public_id"`
	Title         string     `json:"title" db:"title"`
	Description   string     `json:"description" db:"description"`
	OwnerID       int64      `json:"owner_internal_id" gorm:"column:owner_internal_id"`
	OwnerPublicID uuid.UUID  `json:"owner_public_id" db:"owner_public_id"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	Duedate       *time.Time `json:"due_date,omitempty" gorm:"column:due_date"`
}
