package handler

import (
	"fmt"
	"net/http"
	"starup/helper"
	"starup/user"
	"time"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}
func (h *userHandler) RegisterUser(c *gin.Context) {

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {

		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := helper.ApiResponse("Register Account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.ApiResponse("Register Account failed", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "token")

	response := helper.ApiResponse("Account has been registered", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	// 1. user memasukan email dan password
	// 2. input ditangkap handler
	// 3. mapping dari input user ke input struct
	// 4. struct input di parsing ke bentuk service
	// 5. di service , akan mencari dengan bantuan repository user dengan email
	// 6. cek validasi password benar atau salah

	var input user.LoginInput

	err := c.ShouldBindJSON(&input)

	if err != nil {

		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Login Account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedInUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ApiResponse("Login Account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.FormatUser(loggedInUser, "token")

	response := helper.ApiResponse("Login successfull", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)

}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	// 1. ada input email dari user
	// 2. input imael di parsing ke struct input
	// 3. struct input di parsing ke struct service
	// 4. service akan memanggil repository dan cek email sudah ada atau belum
	// 5. repository akan membuat query ke DB

	var input user.CheckEmailInput
	err := c.ShouldBindJSON(&input)

	if err != nil {

		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Login Account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	isEmail, err := h.userService.IsEmailAvailable(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.ApiResponse("Login Account failed", http.StatusUnprocessableEntity, "error", errorMessage)

		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	formatter := user.CheckEmail(isEmail)
	var message string

	if isEmail {
		message = "Email Checked"
	} else {
		message = "Email already use in another account"
	}
	response := helper.ApiResponse(message, http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UpoadAvatar(c *gin.Context) {
	// 1. inputdari user
	// 2. simpan gambar di folder "images/"
	// 3. di servis memanggil repo
	// 4. JWT jika belum ada default pakai user ID = 1
	// 5. repo ambil data user dengan ID 1
	// 6. repo mengupdate data user dan menyimpan lokasi file (path)

	file, err := c.FormFile("avatar")

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("failed to upload avatar image", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	}
	//  next pakat jwt bukan 1
	userID := 1
	currentTime := time.Now()

	// path := "images/" + + currentTime.Format("2006#01#02") + "#" + file.Filename

	path := fmt.Sprintf("images/%d-%s-%s", userID, currentTime.Format("2006#01#02"), file.Filename)
	_, err = h.userService.SaveAvatar(userID, path)
	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("failed to upload avatar image", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	}

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("failed to upload avatar image", http.StatusOK, "success", data)
		c.JSON(http.StatusOK, response)
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

// func (h *userHandler) FetchUserr(c *gin.Context) {
// 	currentUser := c.MustGet("currentUser").(user.User)
// 	formatter := user.FormatUser(currentUser, "")
// 	response := helper.ApiResponse("successfully fetch user data", http.StatusOK, "success", formatter)
// 	c.JSON(http.StatusOK, response)

// }
