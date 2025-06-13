package controller_test

import (
	"encoding/json"
	"me-pague/internal/controller"
	"me-pague/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"me-pague/internal/db"
)

func setupBillingTestDB() {
	testDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{}, &models.Billing{})
	db.DB = testDB
}

func TestGetBilling_Success(t *testing.T) {
	setupBillingTestDB()

	user1, _ := controller.CreateUserHandler("Antonio")
	user2, _ := controller.CreateUserHandler("Davi")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	request_url := "/billing?payer_id=" + strconv.Itoa(int(user1.ID)) + "&receiver_id=" + strconv.Itoa(int(user2.ID))
	c.Request = httptest.NewRequest("GET", request_url, nil)
	c.Request.Header.Set("Content-Type", "application/json")
	q := c.Request.URL.Query()
	q.Add("payer_id", "1")
	q.Add("receiver_id", "2")
	c.Request.URL.RawQuery = q.Encode()

	controller.GetBilling(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var billing models.Billing
	err := json.Unmarshal(w.Body.Bytes(), &billing)
	assert.Nil(t, err)
	assert.Equal(t, int32(1), billing.PayerID)
	assert.Equal(t, int32(2), billing.ReceiverID)
	assert.Equal(t, int32(0), billing.Amount)
}

func TestGetBilling_SameIDs(t *testing.T) {
	setupBillingTestDB()

	user1, _ := controller.CreateUserHandler("Antonio")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	request_url := "/billing?payer_id=" + strconv.Itoa(int(user1.ID)) + "&receiver_id=" + strconv.Itoa(int(user1.ID))
	c.Request = httptest.NewRequest("GET", request_url, nil)
	c.Request.Header.Set("Content-Type", "application/json")
	q := c.Request.URL.Query()
	q.Add("payer_id", "1")
	q.Add("receiver_id", "2")
	c.Request.URL.RawQuery = q.Encode()

	controller.GetBilling(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "payer_id and receiver_id cannot be the same")
}

func TestGetBilling_UserNotFound(t *testing.T) {
	setupBillingTestDB()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/billing?payer_id=1&receiver_id=2", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.URL.RawQuery = "payer_id=1&receiver_id=2"

	controller.GetBilling(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Person 1 not found")
}
