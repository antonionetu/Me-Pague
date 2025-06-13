package controller

import (
	"me-pague/internal/controller/request"
	"me-pague/internal/controller/response"
	"me-pague/internal/db"
	"me-pague/internal/models"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)


// CreatePayment godoc
// @Summary Registra um novo pagamento e atualiza o saldo
// @Tags Pagamentos
// @Accept json
// @Produce json
// @Param payment body request.PaymentInput true "Dados do pagamento"
// @Success 200 {object} request.PaymentInput
// @Failure 400 {object} response.ErrorResponse
// @Router /payment [post]
func CreatePayment(c *gin.Context) {
	var input request.PaymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: err.Error()})
		return
	}

	billing, err := getBillingByID(input.BillingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Billing not found"})
		return
	}

	payment, err := createPayment(input, billing)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{Error: "Error creating payment: " + err.Error()})
		return
	}

	billing.Amount += payment.Amount
	if err := db.DB.Save(&billing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{Error: "Error updating billing amount: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func createPayment(input request.PaymentInput, billing models.Billing) (models.Payment, error) {
	if input.Amount <= 0 {
		return models.Payment{}, fmt.Errorf("amount must be greater than zero")
	}

	var payment models.Payment
	payment.PayerID = billing.PayerID
	payment.BillingID = billing.ID
	payment.Amount = input.Amount

	if err := db.DB.Create(&payment).Error; err != nil {
		return payment, fmt.Errorf("error creating payment: %w", err)
	}

	return payment, nil
}
