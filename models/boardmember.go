package models

import "time"

type BoardMember struct {
	BoardID  int64     `json:"board_id" gorm:"column:board_internal_id;primaryKey;autoIncrement:false"`
	UserID   int64     `json:"user_id" gorm:"column:user_internal_id;primaryKey;autoIncrement:false"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at" gorm:"autoCreateTime"`
}
