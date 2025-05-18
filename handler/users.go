package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shorturl-service/models"
	"shorturl-service/services"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(u services.UserService) UserHandler {
	return UserHandler{
		userService: u,
	}
}

func (u *UserHandler) CreateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get ID from query parameter
		id := ctx.Query("id")
		if id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
			return
		}

		// Initialize a new UserRequest instance
		var us models.UserRequest

		// Bind the request JSON to the struct
		if err := ctx.ShouldBindJSON(&us); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid user request format: " + err.Error(),
			})
			return
		}

		if err := u.userService.SignUp(id, &us); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to sign up user: " + err.Error(),
			})
			return
		}

		// Return success response
		ctx.JSON(http.StatusOK, gin.H{
			"success": "User successfully signed up",
		})
	}
}
