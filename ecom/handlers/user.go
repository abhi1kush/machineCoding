package handlers

import (
	"net/http"
	"strconv"

	"ecom.com/models"
	"ecom.com/services"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	// Bind JSON request body to struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call service layer to create user
	if err := h.Service.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Get the "id" parameter from URL path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Fetch user from service layer
	user, err := h.Service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return user as JSON response
	c.JSON(http.StatusOK, user)
}
