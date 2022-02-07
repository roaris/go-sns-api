package models

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateDB() (db *gorm.DB) {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&User{}, &Post{}, &Friendship{}, &Like{})
	return db
}

func CreateTestDB() (db *gorm.DB) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln(err)
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDatabase := os.Getenv("MYSQL_TEST_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&User{}, &Post{}, &Friendship{}, &Like{})
	return db
}

func CleanUpTestDB(db *gorm.DB) {
	// 順番に注意!
	db.Migrator().DropTable(&Like{})
	db.Migrator().DropTable(&Post{})
	db.Migrator().DropTable(&Friendship{})
	db.Migrator().DropTable(&User{})
}
