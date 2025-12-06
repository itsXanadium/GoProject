package models

type Label struct {
	InternalID int64  `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID   string `json:"public_id" db:"public_id" `
	CardID     int64  `json:"card_internal_id" db:"card_internal_id" gorm:"column:card_internal_id"`
	UserID     int64  `json:"user_internal_id" db:"user_internal_id" gorm:"column:user_internal_id"`
	LabelName  string `json:"label_name" db:"label_name"`
	Color      string `json:"color" db:"color"`
}
