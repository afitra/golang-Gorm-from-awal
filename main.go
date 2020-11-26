package main

import (
	"fmt"
	"log"
	"starup/handler"
	"starup/user"

	"github.com/gin-gonic/gin"
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

	userService := user.NewService(userRepository)
	fmt.Println(">>>>>>", *userRepository)
	// userInput := user.RegisterUserInput{}

	// userInput.Name = "budi"
	// userInput.Email = "budi@mail.com"
	// userInput.Occupation = "anak sekolahan"
	// userInput.Password = "123456"
	// userService.RegisterUser(userInput)

	userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	router.Run()
}
