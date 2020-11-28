package handler

import (
	"net/http"
	"starup/campaign"
	"starup/helper"
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
