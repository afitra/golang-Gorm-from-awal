package main

import (
	"fmt"
	"log"
	"net/http"
	"starup/auth"
	"starup/campaign"
	"starup/handler"
	"starup/helper"
	"starup/payment"
	"starup/transaction"
	"starup/user"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
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

	// db.Migrator().CreateTable(user.User{})
	// db.Migrator().CreateTable(campaign.Campaign{})
	// db.Migrator().CreateTable(campaign.CampaignImage{})

	// db.Migrator().CreateTable(transaction.Transaction{})
	db.Debug().AutoMigrate(
		&user.User{},
		&campaign.Campaign{},
		&campaign.CampaignImage{},
		&transaction.Transaction{},
	)

	fmt.Println("\n koneksi DB berhasil *******\n")

	campaignRepository := campaign.NewRepository(db)
	userRepository := user.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	authService := auth.NewService()
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHandler(campaignService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	router := gin.Default()
	router.Use(cors.Default())
	router.Static("/avatar", "./avatar") // kiri routenya , kanan directory folder

	api := router.Group("/api/v1")

	// ====== debug transaction

	// user, _ := userService.GetUserByID(38)
	// input := transaction.CreateTransactionInput{
	// 	CampaignID: 11,
	// 	Amount:     500000000,
	// 	User:       user,
	// }

	// transactionService.CreateTransaction(input)

	// value := os.Getenv("MIDTRANS_SERVER_KEY")
	// os package
	// value := helper.GoDotEnvVariable("MIDTRANS_SERVER_KEY")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/cek-email", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddlewere(authService, userService), userHandler.UpoadAvatar)

	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	api.POST("/campaigns", authMiddlewere(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("/campaigns/:id", authMiddlewere(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("/campaign-images", authMiddlewere(authService, userService), campaignHandler.UploadImage)

	api.GET("/campaigns/:id/transactions",
		authMiddlewere(authService, userService),
		authorizeCampaign(campaignRepository),
		transactionHandler.GetCampaignTransaction)

	api.GET("/transactions",
		authMiddlewere(authService, userService),
		transactionHandler.GetUserTransactions)

	api.POST("/transactions",
		authMiddlewere(authService, userService),
		transactionHandler.CreateTransaction)

	api.POST("transaction/notification", transactionHandler.GetNotification)

	PORT := helper.GoDotEnvVariable("PORT")
	router.Run(fmt.Sprintf(":%s", PORT))
}

func authorizeCampaign(campaignRepository campaign.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input transaction.GetCampaignTransactionsInput
		currentUser := c.MustGet("currentUser").(user.User)
		input.User = currentUser
		err := c.ShouldBindUri(&input)
		campaign, err := campaignRepository.FindByID(input.ID)

		if err != nil {

			response := helper.ApiResponse("Unauthorized", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if campaign.UserID != input.User.ID {

			response := helper.ApiResponse("Unauthorized, User not owner campaign", http.StatusUnauthorized, "error", nil)

			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return

		}

	}

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
