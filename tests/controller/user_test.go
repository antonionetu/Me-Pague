package controller_test

import (
	"bytes"
	"encoding/json"
	"me-pague/internal/controller"
	"me-pague/internal/db"
	"me-pague/internal/models"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupUserTestDB() {
	testDB, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	testDB.AutoMigrate(&models.User{})
	db.DB = testDB
}

func TestCreateUser_Success(t *testing.T) {
	setupUserTestDB()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{"name": "Antonio"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreateUser(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var user models.User
	err := json.Unmarshal(w.Body.Bytes(), &user)
	assert.Nil(t, err)
	assert.Equal(t, "Antonio", user.Name)
}

func TestCreateUser_AlreadyExists(t *testing.T) {
	setupUserTestDB()
	gin.SetMode(gin.TestMode)

	// cria o usuário antes do teste
	_, err := controller.CreateUserHandler("Antonio")
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := map[string]string{"name": "Antonio"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "User already exists")
}

func TestCreateUser_InvalidInput(t *testing.T) {
	setupUserTestDB()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer([]byte(`{invalid json}`)))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req

	controller.CreateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid input")
}

func TestGetUser_Success(t *testing.T) {
	setupUserTestDB()
	gin.SetMode(gin.TestMode)

	// cria um usuário válido
	user, err := controller.CreateUserHandler("Davi")
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(int(user.ID))}}
	req := httptest.NewRequest("GET", "/user/"+strconv.Itoa(int(user.ID)), nil)
	c.Request = req

	controller.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var gotUser models.User
	err = json.Unmarshal(w.Body.Bytes(), &gotUser)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, gotUser.ID)
	assert.Equal(t, user.Name, gotUser.Name)
}

func TestGetUser_NotFound(t *testing.T) {
	setupUserTestDB()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "id", Value: "999"}}
	req := httptest.NewRequest("GET", "/user/999", nil)
	c.Request = req

	controller.GetUser(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "User not found")
}

func TestGetUser_InvalidID(t *testing.T) {
	setupUserTestDB()
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Params = []gin.Param{{Key: "id", Value: "abc"}}
	req := httptest.NewRequest("GET", "/user/abc", nil)
	c.Request = req

	controller.GetUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid user ID")
}
