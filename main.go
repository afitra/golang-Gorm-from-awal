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

	// userService.SaveAvatar(1, "yesss")

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/cek-email", userHandler.CheckEmailAvailability)
	api.POST("/avatars", userHandler.UpoadAvatar)
	router.Run()
}
