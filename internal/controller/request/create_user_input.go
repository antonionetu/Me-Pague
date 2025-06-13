package request

type CreateUserInput struct {
	Name string `json:"name" binding:"required"`
}
