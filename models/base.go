package models

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

func CreateDB() (db *gorm.DB) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db.LogMode(true)                                                                                                                                                                                                                // ログの出力
	db.AutoMigrate(&User{})                                                                                                                                                                                                         // usersテーブルの作成
	db.AutoMigrate(&Post{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")                                                                                                                                            // postsテーブルの作成, 対応するuserが削除されたらpostも削除される(CASCADE), user_idの更新は認めない(RESTRICT)
	db.AutoMigrate(&Friendship{}).AddForeignKey("followee_id", "users(id)", "CASCADE", "RESTRICT").AddForeignKey("follower_id", "users(id)", "CASCADE", "RESTRICT").AddUniqueIndex("idx_friendships", "follower_id", "followee_id") // friendshipsテーブルの作成

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
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&User{})                                                                                                                                                                                                         // usersテーブルの作成
	db.AutoMigrate(&Post{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")                                                                                                                                            // postsテーブルの作成, 対応するuserが削除されたらpostも削除される(CASCADE), user_idの更新は認めない(RESTRICT)
	db.AutoMigrate(&Friendship{}).AddForeignKey("followee_id", "users(id)", "CASCADE", "RESTRICT").AddForeignKey("follower_id", "users(id)", "CASCADE", "RESTRICT").AddUniqueIndex("idx_friendships", "follower_id", "followee_id") // friendshipsテーブルの作成

	return db
}

func CleanUpTestDB(db *gorm.DB) {
	// 順番に注意!
	db.DropTable(&Post{})
	db.DropTable(&Friendship{})
	db.DropTable(&User{})
}
