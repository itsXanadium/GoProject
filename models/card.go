package models

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	InternalID  int64      `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID    uuid.UUID  `json:"public_id" db:"public_id" gorm:"public_id"`
	ListID      int64      `json:"list_internal_id" db:"list_internal_id" gorm:"column:list_internal_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Duedate     *time.Time `json:"due_date,omitempty" gorm:"column:due_date"`
	Position    int        `json:"position" db:"position"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at" gorm:"autoCreateTime"`

	//Relation
	Assignee    []CardAssignment  `json:"assignees,omitempty" gorm:"foreignKey:CardID;references:InternalID"`
	Attatchment []CardAttatchment `json:"attachment,omitempty" gorm:"foreignKey:CardID;references:InternalID"`
	Labels      []Label           `json:"Labels,omitempty" gorm:"foreignKey:CardID;references:InternalID"`

	// Attatchment []CardAttatchment `json:"attachment,omitempty" gorm:"-"`
	// Labels      []Label           `json:"Labels,omitempty" gorm:"-"`
}
