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

//
func main() {

	dsn := "root:root@tcp(127.0.0.1:3306)/starup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("koneksi DB berhasil *******")

	userRepository := user.NewRepository(db)

	userService := user.NewService(userRepository)

	// input := user.LoginInput{
	// 	Email:    "a@mail.com",
	// 	Password: "123456",
	// }
	// user, err := userService.Login(input)
	// if err != nil {
	// 	fmt.Println("terjadi kesalahan")
	// 	fmt.Println(err.Error())
	// }

	// fmt.Println(user)
	userHandler := handler.NewUserHandler(userService)
	// fmt.Println(">>>>>>", *userRepository)
	// userInput := user.RegisterUserInput{}

	// userInput.Name = "budi"
	// userInput.Email = "budi@mail.com"
	// userInput.Occupation = "anak sekolahan"
	// userInput.Password = "123456"
	// userService.RegisterUser(userInput)

	// userHandler := handler.NewUserHandler(userService)
	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/cek-email", userHandler.CheckEmailAvailability)
	router.Run()
}
