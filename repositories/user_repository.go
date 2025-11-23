package repositories

import (
	"strings"

	"github.com/ADMex1/GoProject/config"
	"github.com/ADMex1/GoProject/models"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindById(id uint) (*models.User, error)
	FindByPublicID(publicID string) (*models.User, error)
	FetchAllWPagination(filter, sort string, limit, offset int) ([]models.User, int64, error)
}

type userRepository struct {
}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Create(user *models.User) error {
	return config.DB.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("email = ?", email).First(&user)
	return &user, err.Error
}

func (r *userRepository) FindById(id uint) (*models.User, error) {
	var user models.User
	err := config.DB.First(&user, id).Error
	return &user, err
}

func (r *userRepository) FindByPublicID(publicID string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("public_id= ?", publicID).First(&user).Error
	return &user, err
}

func (r *userRepository) FetchAllWPagination(filter, sort string, limit, offset int) ([]models.User, int64, error) {
	var user []models.User
	var total int64

	db := config.DB.Model(&models.User{})

	//Filters
	if filter != "" {
		filterPattern := "%" + filter + "%"
		db = db.Where("name Ilike ? OR email Ilike ?", filterPattern, filterPattern)
	}
	//Data Count
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//Sorting sort=name (Ascending) while sort =- Descending
	if sort != "" {
		if sort == "-id" {
			sort = "-internal_id "
		} else if sort == "id" {
			sort = "internal_id "
		}

		//Trimming
		if strings.HasPrefix(sort, "") {
			sort = strings.TrimPrefix(sort, "-") + " DESC"
		} else {
			sort += " ASC"
		}
		db = db.Order(sort)
	}
	err := db.Limit(limit).Offset(offset).Find(&user).Error
	return user, total, err

}
