package repositories

import (
	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
)

type ListRepository interface {
	//interface
	CreateList(list *models.List) error
}

type ListRepositorys struct {
	//struct

}

func NewListRepository() ListRepository {
	return &ListRepositorys{}
}

func (r *ListRepositorys) CreateList(list *models.List) error {
	return config.DB.Create(list).Error
}
