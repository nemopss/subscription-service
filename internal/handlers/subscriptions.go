package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nemopss/subscription-service/internal/db"
	"github.com/nemopss/subscription-service/internal/models"
)

func CreateSubscription(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}
