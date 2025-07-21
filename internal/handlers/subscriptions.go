package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/nemopss/subscription-service/internal/db"
	"github.com/nemopss/subscription-service/internal/models"
	"github.com/nemopss/subscription-service/pkg/logger"
)

// @Summary      Create a new subscription
// @Description  Create a new subscription with the provided details
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription  body      models.CreateSubscription  true  "Subscription data"
// @Success      201           {object}  models.Subscription "CreatedSubscription"
// @Failure      400           {object}  models.ErrorResponse "Invalid request"
// @Failure      500           {object}  models.ErrorResponse "Internal error"
// @Router       /subscriptions [post]
func CreateSubscription(c *gin.Context) {
	logger.Log.Info("Creating new subscription")
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		logger.Log.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		logger.Log.WithError(err).Error("Invalid start_date format")
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid start_date format, expected MM-YYYY"},
		)
		return
	}
	sub.StartDate = startDate.Format("2006-01-02")

	if sub.EndDate != nil {
		endDate, err := time.Parse("01-2006", *sub.EndDate)
		if err != nil {
			logger.Log.WithError(err).Error("Invalid end_date format")
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid end_date format, expected MM-YYYY"},
			)
			return
		}
		endDateStr := endDate.Format("2006-01-02")
		sub.EndDate = &endDateStr
	}

	if err := db.DB.Create(&sub).Error; err != nil {
		logger.Log.WithError(err).Error("Failed to create subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithField("id", sub.ID).Info("Subscription created")
	c.JSON(http.StatusCreated, sub)
}

// @Summary      Get a subscription by ID
// @Description  Retrieve a subscription by its ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  models.Subscription
// @Failure      400  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /subscriptions/{id} [get]
func GetSubscription(c *gin.Context) {
	logger.Log.Info("Getting subscription")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.WithError(err).Error("Invalid subscription id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var sub models.Subscription
	if err := db.DB.First(&sub, id).Error; err != nil {
		logger.Log.WithError(err).Error("Subscription not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}
	logger.Log.WithField("id", sub.ID).Info("Subscription found")
	c.JSON(http.StatusOK, sub)
}

// @Summary      Update a subscription
// @Description  Updates an existing subscription by ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id           path      int                      true  "Subscription ID"
// @Param        subscription body      models.CreateSubscription  true  "Updated subscription data"
// @Success      200          {object}  models.Subscription      "Updated subscription"
// @Failure      400          {object}  models.ErrorResponse       "Invalid subscription ID or request body"
// @Failure      404          {object}  models.ErrorResponse       "Subscription not found"
// @Failure      500          {object}  models.ErrorResponse       "Internal server error"
// @Router       /subscriptions/{id} [put]
func UpdateSubscription(c *gin.Context) {
	logger.Log.Info("Updating subscription")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.WithError(err).Error("Invalid subscription id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var sub models.Subscription
	if err := db.DB.First(&sub, id).Error; err != nil {
		logger.Log.WithError(err).Error("Subscription not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	var updatedSub models.Subscription
	if err := c.ShouldBindJSON(&updatedSub); err != nil {
		logger.Log.WithError(err).Error("Invalid request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Парсим start_date, если обновлено
	if updatedSub.StartDate != "" {
		startDate, err := time.Parse("01-2006", updatedSub.StartDate)
		if err != nil {
			logger.Log.WithError(err).Error("Invalid start_date format")
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid start_date format, expected MM-YYYY"},
			)
			return
		}
		sub.StartDate = startDate.Format("2006-01-02")
	}

	// Парсим end_date, если обновлено
	if updatedSub.EndDate != nil {
		endDate, err := time.Parse("01-2006", *updatedSub.EndDate)
		if err != nil {
			logger.Log.WithError(err).Error("Invalid end_date format")
			c.JSON(
				http.StatusBadRequest,
				gin.H{"error": "invalid end_date format, expected MM-YYYY"},
			)
			return
		}
		endDateStr := endDate.Format("2006-01-02")
		sub.EndDate = &endDateStr
	} else {
		sub.EndDate = nil
	}

	// Обновляем остальные поля
	sub.ServiceName = updatedSub.ServiceName
	sub.Price = updatedSub.Price
	sub.UserID = updatedSub.UserID

	// Сохраняем изменения
	if err := db.DB.Save(&sub).Error; err != nil {
		logger.Log.WithError(err).Error("Failed to update subscription")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithField("id", sub.ID).Info("Subscription updated")
	c.JSON(http.StatusOK, sub)
}

// @Summary      Delete a subscription by ID
// @Description  Delete a subscription by its ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "Subscription ID"
// @Success      204
// @Failure      401  {object}  models.ErrorResponse
// @Failure      404  {object}  models.ErrorResponse
// @Router       /subscriptions/{id} [delete]
func DeleteSubscription(c *gin.Context) {
	logger.Log.Info("Deleting subscription")
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Log.WithError(err).Error("Invalid subscription id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}
	if rows := db.DB.Delete(models.Subscription{}, id).RowsAffected; rows == 0 {
		logger.Log.WithError(err).Error("Subscription not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	logger.Log.WithField("id", id).Info("Subscription deleted")
	c.Status(http.StatusNoContent)
}

// @Summary      Get all subscriptions
// @Description  Retrieve all of the subscriptions
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}  models.Subscription
// @Failure      400  {object}  models.ErrorResponse
// @Router       /subscriptions [get]
func ListSubscriptions(c *gin.Context) {
	logger.Log.Info("Listing subscriptions")
	var subs []models.Subscription

	if err := db.DB.Find(&subs).Error; err != nil {
		logger.Log.WithError(err).Error("Failed to list subscriptions")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Log.WithField("count", len(subs)).Info("Subscriptions listed")
	c.JSON(http.StatusOK, subs)
}

// GetTotalCostByPeriod вычисляет общую стоимость подписок за период
// @Summary      Get total cost of subscriptions
// @Description  Calculates the total cost of subscriptions for a user over a specified period, optionally filtered by service name and date range
// @Tags         subscriptions
// @Produce      json
// @Param        user_id       query     string  true   "User ID (UUID)"
// @Param        service_name  query     string  false  "Service name filter"
// @Param        start_date    query     string  false  "Start date (MM-YYYY)"
// @Param        end_date      query     string  false  "End date (MM-YYYY)"
// @Success      200           {object}  models.TotalCostResponse  "Total cost"
// @Failure      400           {object}  models.ErrorResponse      "Invalid user ID or date format"
// @Failure      500           {object}  models.ErrorResponse      "Internal server error"
// @Router       /subscriptions/total [get]
func GetTotalCostByPeriod(c *gin.Context) {
	logger.Log.Info("Getting total cost by period")
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		logger.Log.WithError(err).Error("Invalid user id")
		c.JSON(
			http.StatusBadRequest,
			models.ErrorResponse{Error: "invalid user id"},
		)
		return
	}

	defaultStart := time.Date(
		1970,
		1,
		1,
		0,
		0,
		0,
		0,
		time.UTC,
	)
	defaultEnd := time.Now()

	var periodStart, periodEnd time.Time
	if startDateStr != "" {
		periodStart, err = time.Parse("01-2006", startDateStr)
		if err != nil {
			logger.Log.WithError(err).Error("Invalid start_date format")
			c.JSON(
				http.StatusBadRequest,
				models.ErrorResponse{
					Error: "invalid start_date format, expected MM-YYYY",
				},
			)
			return
		}
	} else {
		periodStart = defaultStart
	}

	if endDateStr != "" {
		periodEnd, err = time.Parse("01-2006", endDateStr)
		if err != nil {
			logger.Log.WithError(err).Error("Invalid end_date format")
			c.JSON(
				http.StatusBadRequest,
				models.ErrorResponse{
					Error: "invalid end_date format, expected MM-YYYY",
				},
			)
			return
		}
		// Устанавливаем конец месяца
		periodEnd = periodEnd.AddDate(0, 1, -1)
	} else {
		periodEnd = defaultEnd
	}

	// Получаем все подписки, соответствующие фильтрам
	var subscriptions []models.Subscription
	dbQuery := db.DB.Model(&models.Subscription{}).
		Where("user_id = ?", userID)

	if serviceName != "" {
		dbQuery = dbQuery.Where("service_name = ?", serviceName)
	}

	if err := dbQuery.Find(&subscriptions).Error; err != nil {
		logger.Log.WithError(err).Error("Failed to fetch subscriptions")
		c.JSON(
			http.StatusInternalServerError,
			models.ErrorResponse{Error: "database error"},
		)
		return
	}

	// Рассчитываем общую стоимость
	totalCost := 0
	for _, sub := range subscriptions {
		// Парсим даты подписки
		subStart, err := time.Parse("2006-01-02", sub.StartDate)
		if err != nil {
			logger.Log.WithError(err).
				WithField("id", sub.ID).
				Error("Invalid subscription start_date")
			continue
		}

		var subEnd time.Time
		if sub.EndDate != nil {
			subEnd, err = time.Parse("2006-01-02", *sub.EndDate)
			if err != nil {
				logger.Log.WithError(err).
					WithField("id", sub.ID).
					Error("Invalid subscription end_date")
				continue
			}
		} else {
			subEnd = defaultEnd
		}

		// Определяем пересечение периода подписки с запрошенным периодом
		effectiveStart := maxTime(subStart, periodStart)
		effectiveEnd := minTime(subEnd, periodEnd)

		// Проверяем, есть ли пересечение
		if effectiveStart.Before(effectiveEnd) ||
			effectiveStart.Equal(effectiveEnd) {
			// Рассчитываем количество месяцев
			months := (effectiveEnd.Year()-effectiveStart.Year())*12 +
				int(effectiveEnd.Month()-effectiveStart.Month()) + 1
			if months < 1 {
				months = 1 // Минимально учитываем один месяц
			}
			totalCost += sub.Price * months
			logger.Log.WithFields(logrus.Fields{
				"id":     sub.ID,
				"months": months,
				"cost":   sub.Price * months,
			}).Info("Calculated cost for subscription")
		}
	}

	logger.Log.WithField("count", len(subscriptions)).
		Info("Fetched subscriptions")
	logger.Log.WithField("total", totalCost).Info("Total cost calculated")
	c.JSON(http.StatusOK, models.TotalCostResponse{Total: totalCost})
}

// Вспомогательные функции для выбора максимальной и минимальной даты
func maxTime(t1, t2 time.Time) time.Time {
	if t1.After(t2) {
		return t1
	}
	return t2
}

func minTime(t1, t2 time.Time) time.Time {
	if t1.Before(t2) {
		return t1
	}
	return t2
}
