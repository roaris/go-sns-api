package models

import (
	"time"

	"github.com/roaris/go-sns-api/swagger/gen"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Name      string
	Email     string `gorm:"unique_index"`
	Password  string
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
	user.Name = name
	user.Email = email
	user.Password = Encrypt(password)
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

func UpdateUser(id int64, name string, email string, password string) (user User) {
	user, _ = GetUserById(id)
	userAfter := user
	userAfter.Name = name
	userAfter.Email = email
	userAfter.Password = Encrypt(password)
	db.Model(&user).Updates(userAfter)
	return user
}
