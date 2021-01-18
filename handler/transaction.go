package handler

import (
	"fmt"
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

func (h *transactionHandler) CreateTransaction(c *gin.Context) {

	// 1. ada input dari user
	// 2. handler angkap input trus di maping ke input struct
	// 3. panggil service untuk transaksi , manggil sistem midtrans
	// 4. panggil repository create new transaction data

	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)

	if err != nil {

		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.ApiResponse("Failed to create transaction", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return

	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)

	if err != nil {
		fmt.Println(err.Error())
		response := helper.ApiResponse("Failed to create transaction", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// response := helper.ApiResponse("User's transactions", http.StatusOK, "success", newTransaction)
	response := helper.ApiResponse("User's transactions", http.StatusOK, "success", transaction.FormatTransaction(newTransaction))

	c.JSON(http.StatusBadRequest, response)
}
