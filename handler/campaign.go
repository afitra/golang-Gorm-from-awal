package handler

import (
	"fmt"
	"net/http"
	"reflect"
	"starup/campaign"
	"starup/helper"
	"starup/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 1. tangkap parameter di handler
// 2. handler ke service
// 3. service menentukan repository(method) mana yg di panggil
// 4. repository getAll / findByUserID
// 5. db
type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {

	return &campaignHandler{service}

}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {

	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)

	if err != nil {
		response := helper.ApiResponse("Error to get Campaigns", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return

	}

	response := helper.ApiResponse("List of Campaign", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))

	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {

	// uri -->> api/v1/campaigns/2
	// 1. handler 		: maping id yang di url ke struct input trus di masukkan ke  => service dan  call formatter
	// 2. service 		: inputnya struct input => menangkap id di url , manggil repo
	// 3. repository 	: get campaign by id

	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.ApiResponse("Error to get detail Campaigns", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)

	if err != nil {
		response := helper.ApiResponse("Error to get detail Campaigns", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return

	}

	response := helper.ApiResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}

// 1. ubah parameter dari user ke struct input
// 2. ambil current user dari jwt/handler
// 3. panggil service , parameternya input struct dan buat slug
// 4. panggil repository untuk simpan campaing baru

func (h *campaignHandler) CreateCampaign(c *gin.Context) {

	var input campaign.CreateCampaignInput
	err := c.ShouldBindJSON(&input)

	if err != nil {

		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Failed to create Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return

	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newCampaign, err := h.service.CreateCampaign(input)

	if err != nil {
		response := helper.ApiResponse("Failed to create Campaign", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.ApiResponse("Succes to create  Campaign", http.StatusOK, "success", campaign.FormatCampaign(newCampaign))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) UpdateCampaign(c *gin.Context) {

	// 1. user masukkan input
	// 2. handler
	// 3. mapping dari input ke struct input (ada 2)
	// 4. input dari user dan juga input dari uri (passing ke service )
	// 5. service
	// 6. repository update data campaign

	var inputID campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)

	if err != nil {
		response := helper.ApiResponse("Failed to update Campaign", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)

	if err != nil {

		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Failed to update Campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return

	}

	currentUser := c.MustGet("currentUser").(user.User) //  data user (dari relasi db) di set agar bisa di baca di service dan dilakukan verifikasi

	inputData.User = currentUser // apakah benar yg sedang login adlah user yang punya campaign

	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)

	if err != nil {

		response := helper.ApiResponse("Failed to update Campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	response := helper.ApiResponse("Succes to update  Campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

//  handler
//  1. tangkap dan rubah struct ke struct input
//  2. save image campaign ke suatu folder
//  3. service manggil point 2  atau 1 repository
//  repository
//  1. create image/ save data image ke dalam db
//  2. cek ada gak is_primary yg true , jika ada ubah is_primary true ke false

func (h *campaignHandler) UploadImage(c *gin.Context) {

	var input campaign.CreateCampaignImageInput

	err := c.ShouldBind(&input)
	fmt.Println("========", input.IsPrimary, "xxxxx", reflect.TypeOf(input.IsPrimary), err)
	fmt.Println("<<<<<<", err)
	if err != nil {
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Failed to upload Campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return

	}
	file, err := c.FormFile("file")
	fmt.Println(">>>>>>>", file)
	if err != nil {

		data := gin.H{"is_uploaded": false}

		response := helper.ApiResponse("Failed to upload Campaign image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
		return

	}
	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)

	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("failed to upload Campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.ApiResponse("failed to upload campaign  image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.ApiResponse("Avatar successfully uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
