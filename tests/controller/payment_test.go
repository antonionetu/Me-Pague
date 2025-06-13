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

func setupTestPaymentDB() {
	testDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Billing{}, &models.Payment{})
	db.DB = testDB
}

func TestCreatePayment_Success(t *testing.T) {
	setupTestPaymentDB()
	gin.SetMode(gin.TestMode)

	user1, _ := controller.CreateUserHandler("Antonio")
	user2, _ := controller.CreateUserHandler("Davi")
	billing, _ := controller.GetOrCreateBilling(request.BillingInput{PayerID: user1.ID, ReceiverID: user2.ID})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]interface{}{
		"billing_id": billing.ID,
		"amount":     50,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/payment", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreatePayment(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var payment models.Payment
	result := db.DB.First(&payment, "billing_id = ?", billing.ID)

	assert.Nil(t, result.Error)
	assert.Equal(t, int32(50), payment.Amount)
}

func TestCreatePayment_InvalidBillingID(t *testing.T) {
	setupTestPaymentDB()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]interface{}{
		"billing_id": 999,
		"amount":     50,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/payment", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreatePayment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Contains(t, responseBody["error"], "Billing not found")
}

func TestCreatePayment_InvalidPayload(t *testing.T) {
	setupTestPaymentDB()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := `{"billing_id": 1, "amount": "not_a_number"}`

	req := httptest.NewRequest("POST", "/payment", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreatePayment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreatePayment_NegativeAmount(t *testing.T) {
	setupTestPaymentDB()
	gin.SetMode(gin.TestMode)

	user1, _ := controller.CreateUserHandler("Ana")
	user2, _ := controller.CreateUserHandler("Beto")
	billing, _ := controller.GetOrCreateBilling(request.BillingInput{PayerID: user1.ID, ReceiverID: user2.ID})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]interface{}{
		"billing_id": billing.ID,
		"amount":     -30,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/payment", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreatePayment(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var responseBody map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &responseBody)
	assert.Contains(t, responseBody["error"], "amount must be greater than zero")
}

func TestCreatePayment_MultiplePayments(t *testing.T) {
	setupTestPaymentDB()
	gin.SetMode(gin.TestMode)

	user1, _ := controller.CreateUserHandler("Lara")
	user2, _ := controller.CreateUserHandler("Rafael")
	billing, _ := controller.GetOrCreateBilling(request.BillingInput{PayerID: user1.ID, ReceiverID: user2.ID})

	amounts := []int32{10, 20, 30}

	for _, amt := range amounts {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		body := map[string]interface{}{
			"billing_id": billing.ID,
			"amount":     amt,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest("POST", "/payment", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		c.Request = req

		controller.CreatePayment(c)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	var updatedBilling models.Billing
	db.DB.First(&updatedBilling, billing.ID)

	assert.Equal(t, int32(60), updatedBilling.Amount)
}

func TestCreatePayment_BalanceDifference(t *testing.T) {
	setupTestPaymentDB()
	gin.SetMode(gin.TestMode)

	user1, _ := controller.CreateUserHandler("Carlos")
	user2, _ := controller.CreateUserHandler("Fernanda")
	billing1, _ := controller.GetOrCreateBilling(request.BillingInput{PayerID: user1.ID, ReceiverID: user2.ID})
	billing2, _ := controller.GetOrCreateBilling(request.BillingInput{PayerID: user2.ID, ReceiverID: user1.ID})

	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)

	body1 := map[string]interface{}{
		"billing_id": billing1.ID,
		"amount":     100,
	}
	jsonBody1, _ := json.Marshal(body1)

	req1 := httptest.NewRequest("POST", "/payment", bytes.NewBuffer(jsonBody1))
	req1.Header.Set("Content-Type", "application/json")
	c1.Request = req1

	controller.CreatePayment(c1)
	assert.Equal(t, http.StatusOK, w1.Code)

	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)

	body2 := map[string]interface{}{
		"billing_id": billing2.ID,
		"amount":     80,
	}
	jsonBody2, _ := json.Marshal(body2)

	req2 := httptest.NewRequest("POST", "/payment", bytes.NewBuffer(jsonBody2))
	req2.Header.Set("Content-Type", "application/json")
	c2.Request = req2

	controller.CreatePayment(c2)
	assert.Equal(t, http.StatusOK, w2.Code)

	var updatedBilling1, updatedBilling2 models.Billing
	db.DB.First(&updatedBilling1, billing1.ID)
	db.DB.First(&updatedBilling2, billing2.ID)

	assert.Equal(t, int32(100), updatedBilling1.Amount)
	assert.Equal(t, int32(80), updatedBilling2.Amount)

	assert.Equal(t, updatedBilling1.Amount - updatedBilling2.Amount, int32(20))
	assert.Equal(t, updatedBilling2.Amount - updatedBilling1.Amount, int32(-20))
}
