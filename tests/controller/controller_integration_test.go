package controller_test

import (
	"bytes"
	"encoding/json"
	"me-pague/internal/controller"
	"me-pague/internal/controller/request"
	"me-pague/internal/db"
	"me-pague/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIntegrationDB() {
	testDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Billing{}, &models.Payment{})
	db.DB = testDB
}

func TestIntegration_FullPaymentFlow(t *testing.T) {
	setupIntegrationDB()
	gin.SetMode(gin.TestMode)

	userA, err := controller.CreateUserHandler("Ana")
	assert.Nil(t, err)
	userB, err := controller.CreateUserHandler("Bruno")
	assert.Nil(t, err)

	billing, err := controller.GetOrCreateBilling(request.BillingInput{
		PayerID:    userA.ID,
		ReceiverID: userB.ID,
	})
	assert.Nil(t, err)
	assert.Equal(t, int32(0), billing.Amount)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	paymentBody := map[string]interface{}{
		"billing_id": billing.ID,
		"amount":     150,
	}
	jsonBody, _ := json.Marshal(paymentBody)

	req := httptest.NewRequest("POST", "/payment", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreatePayment(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedBilling models.Billing
	err = db.DB.First(&updatedBilling, billing.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, int32(150), updatedBilling.Amount)

	var payment models.Payment
	err = db.DB.First(&payment, "billing_id = ?", billing.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, userA.ID, payment.PayerID)
	assert.Equal(t, int32(150), payment.Amount)
}
