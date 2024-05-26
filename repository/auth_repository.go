package repository

import (
	"GO-SOCMED/entity"

	"gorm.io/gorm"
)

type AuthRepository interface {
	EmailExist(email string) bool
	Register(req *entity.User) error
}

/*
	auth menggunakan huruf kecil berarti

hanya bisa di panggil di dalam file
auth_repository
*/
type authRepository struct {
	db *gorm.DB
}

/* fungsi constractor yang diawali dgn New */
func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Register(user *entity.User) error {
	err := r.db.Create(&user).Error

	return err
}

func (r *authRepository) EmailExist(email string) bool {
	var user entity.User
	err := r.db.First(&user, "email = ?", email).Error

	return err == nil
}
