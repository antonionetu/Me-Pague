package controller
import (
	"me-pague/internal/db"
	"me-pague/internal/models"
	"me-pague/internal/controller/request"
	"net/http"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
	"fmt"
)

// GetBilling godoc
// @Summary Obtém ou cria uma cobrança entre dois usuários
// @Tags Cobranças
// @Accept json
// @Produce json
// @Param payer_id query string true "ID do pagador"
// @Param receiver_id query string true "ID do recebedor"
// @Success 200 {object} models.Billing
// @Failure 400 {object} response.ErrorResponse
// @Router /billing [get]
func GetBilling(c *gin.Context) {
	var billingInput request.BillingInput
	payerID, _ := strconv.Atoi(c.Query("payer_id"))
	receiverID, _ := strconv.Atoi(c.Query("receiver_id"))
	billingInput.PayerID = int32(payerID)
	billingInput.ReceiverID = int32(receiverID)

	_, err := validationError(billingInput)	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	billing, err := GetOrCreateBilling(billingInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, billing)
}


func GetOrCreateBilling(billingInput request.BillingInput) (models.Billing, error) {
    var billing models.Billing
	db.DB.Where("payer_id = ? AND receiver_id = ?", billingInput.PayerID, billingInput.ReceiverID).First(&billing)
	
	if billing.ID == 0 {
		billing = models.Billing{
			PayerID: billingInput.PayerID,
			ReceiverID: billingInput.ReceiverID,
			Amount:    0,
			CreatedAt: time.Now(),
		}
		db.DB.Create(&billing)
	}
	return billing, nil
}

func getBillingByID(id int32) (models.Billing, error) {
	var billing models.Billing
	if err := db.DB.Where("id = ?", id).First(&billing).Error; err != nil {
		return billing, fmt.Errorf("Billing not found")
	}
	return billing, nil
}

func validationError(billingInput request.BillingInput) (request.BillingInput, error) {
	if billingInput.PayerID == billingInput.ReceiverID {
		return billingInput, fmt.Errorf("payer_id and receiver_id cannot be the same")
	}

	if db.DB.Where("id = ?", billingInput.PayerID).First(&models.User{}).Error != nil {
        return billingInput, fmt.Errorf("Person 1 not found")
    }
    if db.DB.Where("id = ?", billingInput.ReceiverID).First(&models.User{}).Error != nil {
        return billingInput, fmt.Errorf("Person 2 not found")
    }

	return billingInput, nil
}
