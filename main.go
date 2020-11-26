package main

import (
	"fmt"
	"log"
	"starup/user"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dsn := "root:root@tcp(127.0.0.1:3306)/starup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("koneksi DB berhasil *******")

	userRepository := user.NewRepository(db)

	user := user.User{
		Name: "test",
	}
	// db.AutoMigrate(&user)
	userRepository.Save(user)
}
