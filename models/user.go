package models

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/roaris/go-sns-api/swagger/gen"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	ID        int64
	Name      string `validate:"required,min=3"`
	Email     string `validate:"required,email" gorm:"unique_index"`
	Password  string `validate:"required"`
	Posts     []Post
	CreatedAt time.Time
}

func (u *User) SwaggerModel() *gen.User {
	return &gen.User{
		ID:   u.ID,
		Name: u.Name,
	}
}

func (u *User) SwaggerModelWithEmail() *gen.User {
	return &gen.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func Encrypt(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash)
}

func CreateUser(name string, email string, password string) (user User, err error) {
	if utf8.RuneCountInString(password) < 6 {
		err = errors.New("too short password")
		return user, err
	}
	user.Name = name
	user.Email = email
	user.Password = Encrypt(password)
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		return user, err
	}
	err = db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserById(id int64) (user User, err error) {
	err = db.First(&user, "id=?", id).Error
	return user, err
}

func GetUserByEmail(email string) (user User, err error) {
	err = db.First(&user, "email=?", email).Error
	return user, err
}

func UpdateUser(id int64, name string, email string, password string) (user User){
	user, _ = GetUserById(id)
	userAfter := user
	userAfter.Name = name
	userAfter.Email = email
	userAfter.Password = Encrypt(password)
	db.Model(&user).Updates(userAfter)
	return user
}
