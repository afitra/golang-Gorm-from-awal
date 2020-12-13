package handler

import (
	"net/http"
	"starup/helper"
	"starup/transaction"
	"starup/user"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {

	return &transactionHandler{service}

}

func (h *transactionHandler) GetCampaignTransaction(c *gin.Context) {

	var input transaction.GetCampaignTransactionsInput

	err := c.ShouldBindUri(&input)

	if err != nil {
		response := helper.ApiResponse("Failed to get  Campaign's transaction", http.StatusBadRequest, "error", nil)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	transactions, err := h.service.GetTransactionsByCampaignID(input)

	if err != nil {
		response := helper.ApiResponse("Failed to get Campaign's Transaction", http.StatusBadRequest, "error", err)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("Campaign's Transaction", http.StatusOK, "success", transaction.FormatCampaignTransactions(transactions))

	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetUserTransactions(c *gin.Context) {

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID

	transactions, err := h.service.GetTransactionsByUserID(userID)

	if err != nil {

		response := helper.ApiResponse("failed to get user's transactions", http.StatusBadRequest, "error", err)

		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.ApiResponse("User's transactions", http.StatusOK, "success", transaction.FormatUserTransactions(transactions))

	c.JSON(http.StatusBadRequest, response)

}
