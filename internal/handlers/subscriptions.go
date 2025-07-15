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

func GetSubscription(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription
	if err := db.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

func UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var sub models.Subscription
	if err := db.DB.First(&sub, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Save(&sub)

	c.JSON(http.StatusOK, sub)

}

func DeleteSubscription(c *gin.Context) {
	id := c.Param("id")
	if rows := db.DB.Delete(models.Subscription{}, id).RowsAffected; rows == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
