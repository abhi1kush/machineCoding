package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthChecksHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
