package repositories

import (
	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
)

type BoardMemberRepository interface {
	//Interface
	GetMembers(boardPublicID string) ([]models.User, error)
}

type BoardMemberRepositorys struct {
	//Struct

}

func NewMemberRepository() BoardMemberRepository {
	return &BoardMemberRepositorys{}
}

func (r *BoardMemberRepositorys) GetMembers(boardPublicID string) ([]models.User, error) {
	var user []models.User
	err := config.DB.
		// Table("users").
		Joins("JOIN board_members ON board_members.user_internal_id = users.internal_id").
		Joins("JOIN boards on boards.internal_id = board_members.board_internal_id").
		Where("boards.public_id = ?", boardPublicID).
		Find(&user).Error
	return user, err
}

//board_internal_id
