package main

import (
	"fmt"
	"log"
	"net/http"
	"starup/auth"
	"starup/campaign"
	"starup/handler"
	"starup/helper"
	"starup/user"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	// Campaign := campaign.Campaign{}
	// db.AutoMigrate(Campaign)
	fmt.Println("koneksi DB berhasil *******")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	campaignRepository := campaign.NewRepository(db)
	campaignService := campaign.NewService(campaignRepository)

	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	router := gin.Default()
	router.Static("/avatar", "./avatar") // kiri routenya , kanan directory folder

	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/cek-email", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddlewere(authService, userService), userHandler.UpoadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddlewere(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddlewere(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddlewere(authService, userService), campaignHandler.UploadImage)

	router.Run()
}

func authMiddlewere(authService auth.Service, userService user.Service) gin.HandlerFunc {
	// Midleware
	// 1. ambil nilai header authorization ->> bearer token
	// 2. dari header authorization , kita ambil token saja
	// 3. validasi tokennya
	// 4. ambil nilai user_id
	// 5. ambil user di db berdasar user_id
	// 6. set context isinya user
	return func(c *gin.Context) { // -->> gin handler adl fungsi yg punya param gin.Context

		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {

			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")

		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]

		}
		token, err := authService.ValidateToken(tokenString)

		if err != nil {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)

		if err != nil {

			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
