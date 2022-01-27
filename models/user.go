package models

import (
	"errors"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" validate:"required,min=3"`
	Email     string    `json:"email" validate:"required,email" gorm:"unique_index"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

func Encrypt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash)
}

func CreateUser(name string, email string, password string) (err error) {
	if utf8.RuneCountInString(password) < 6 {
		err = errors.New("too short password")
		return err
	}
	user := User{}
	user.Name = name
	user.Email = email
	user.Password = Encrypt(password)
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		return err
	}
	err = db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(email string) (user User, err error) {
	err = db.First(&user, "email=?", email).Error
	return user, err
}
