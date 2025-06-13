package controller
import (
	"me-pague/internal/db"
	"me-pague/internal/models"
	"me-pague/internal/controller/request"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

// getUserByID grodoc
// @Summary Obtém um usuário pelo ID
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path int true "ID do usuário"
// @Success 200 {object} models.User
// @Failure 404 {object} response.ErrorResponse
// @Router /user/{id} [get]
func GetUser(c *gin.Context) {
	ID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := getUserByID(uint(ID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}


// CreateUser godoc
// @Summary Cria um novo usuário
// @Tags Usuários
// @Accept json
// @Produce json
// @Param user body request.CreateUserInput true "Dados do usuário"
// @Success 201 {object} request.CreateUserInput
// @Failure 400 {object} response.ErrorResponse
// @Router /user [post]
func CreateUser(c *gin.Context) {
	var input request.CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	
	_, err := getUserByName(input.Name)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	newUser, err := CreateUserHandler(input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, newUser) 	
}


func getUserByID(ID uint) (models.User, error) {
	var user models.User
	if err := db.DB.First(&user, ID).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func getUserByName(name string) (models.User, error) {
	var user models.User
	if err := db.DB.Where("name = ?", name).First(&user).Error; err != nil {
		return models.User{}, err
	}
	return user, nil
}

func CreateUserHandler(name string) (models.User, error) {
	newUser := models.User{Name: name}
	if err := db.DB.Create(&newUser).Error; err != nil {
		return models.User{}, err
	}
	return newUser, nil
}