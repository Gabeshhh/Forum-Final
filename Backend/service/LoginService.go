package service

import (
	"fmt"

	database "Forum/Backend/Database"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginService interface {
	LoginUser(email string, password string) bool
}

type loginInformation struct {
	db *gorm.DB
}

func StaticLoginService(db *gorm.DB) LoginService {
	return &loginInformation{
		db: db,
	}
}

func (info *loginInformation) LoginUser(email string, password string) bool {
	var user database.Users
	if err := info.db.Where("email = ?", email).First(&user).Error; err != nil {
		fmt.Println("User not found for email:", email)
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("Mot de passe incorrect")
		return false
	}

	return true
}
