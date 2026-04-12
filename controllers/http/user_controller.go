package http

import (
	"auto-zen-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service services.UserService
}

func NewUserController(s services.UserService) *UserController {
	return &UserController{Service: s}
}

// POST / signup
func (ctrl *UserController) Signup(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ユーザー名とパスワードは必須です"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully!"})
}