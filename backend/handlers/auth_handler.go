package handlers

import (
	"brute-force-login/models"
	"brute-force-login/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get IP address from request
	ipAddress := c.ClientIP()
	fmt.Println("Client IP from c.ClientIP():", ipAddress)
	if ipAddress == "" {
		ipAddress = c.GetHeader("X-Forwarded-For")
		if ipAddress == "" {
			ipAddress = c.GetHeader("X-Real-IP")
		}
	}

	response, err := h.authService.Login(req.Email, req.Password, ipAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, response)
	} else {
		c.JSON(http.StatusUnauthorized, response)
	}
}
