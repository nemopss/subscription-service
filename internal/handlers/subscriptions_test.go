package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/nemopss/subscription-service/internal/config"
	"github.com/nemopss/subscription-service/internal/db"
	"github.com/nemopss/subscription-service/internal/handlers"
	"github.com/nemopss/subscription-service/internal/models"
)

// setupTestDB подключается к тестовой базе данных и выполняет миграции
func setupTestDB() {
	config.LoadConfig("../../.env")
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME_TEST"),
		os.Getenv("DB_PASSWORD"),
	)

	var err error
	db.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	db.DB.AutoMigrate(&models.Subscription{})
}

// setupRouter настраивает маршрутизатор Gin для тестов
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/subscriptions", handlers.CreateSubscription)
	r.GET("/subscriptions/:id", handlers.GetSubscription)
	r.PUT("/subscriptions/:id", handlers.UpdateSubscription)
	r.DELETE("/subscriptions/:id", handlers.DeleteSubscription)
	r.GET("/subscriptions", handlers.ListSubscriptions)
	r.GET("/subscriptions/total", handlers.GetTotalCostByPeriod)
	return r
}

// withTransaction выполняет тестовую функцию в рамках транзакции и откатывает её после завершения
func withTransaction(t *testing.T, testFunc func(tx *gorm.DB)) {
	tx := db.DB.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	// Временно подменяем глобальный db.DB на транзакционный
	originalDB := db.DB
	db.DB = tx
	defer func() {
		db.DB = originalDB
		tx.Rollback()
	}()

	testFunc(tx)
}

// Тест для CreateSubscription
func TestCreateSubscription(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		userID := uuid.New()
		sub := models.CreateSubscription{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      userID,
			StartDate:   "07-2025",
		}
		jsonData, _ := json.Marshal(sub)
		req, _ := http.NewRequest(
			"POST",
			"/subscriptions",
			bytes.NewBuffer(jsonData),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
		var createdSub models.Subscription
		json.Unmarshal(w.Body.Bytes(), &createdSub)
		assert.Equal(t, sub.ServiceName, createdSub.ServiceName)
		assert.Equal(t, sub.Price, createdSub.Price)
		assert.Equal(t, sub.UserID, createdSub.UserID)
		assert.Equal(t, "2025-07-01", createdSub.StartDate)
	})
}

// Тест для CreateSubscription с некорректной датой
func TestCreateSubscriptionInvalidDate(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		sub := models.CreateSubscription{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      uuid.New(),
			StartDate:   "invalid-date",
		}
		jsonData, _ := json.Marshal(sub)
		req, _ := http.NewRequest(
			"POST",
			"/subscriptions",
			bytes.NewBuffer(jsonData),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

// Тест для GetSubscription
func TestGetSubscription(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		userID := uuid.New()
		sub := models.Subscription{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      userID,
			StartDate:   "2025-07-01",
		}
		tx.Create(&sub)
		req, _ := http.NewRequest(
			"GET",
			"/subscriptions/"+strconv.Itoa(sub.ID),
			nil,
		)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedSub models.Subscription
		json.Unmarshal(w.Body.Bytes(), &fetchedSub)
		assert.Equal(t, sub.ID, fetchedSub.ID)
		assert.Equal(t, sub.ServiceName, fetchedSub.ServiceName)
	})
}

// Тест для GetSubscription с несуществующим ID
func TestGetSubscriptionNotFound(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		req, _ := http.NewRequest("GET", "/subscriptions/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Тест для UpdateSubscription
func TestUpdateSubscription(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		userID := uuid.New()
		sub := models.Subscription{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      userID,
			StartDate:   "2025-07-01",
		}
		tx.Create(&sub)
		updatedSub := models.CreateSubscription{
			ServiceName: "Netflix",
			Price:       600,
			UserID:      userID,
			StartDate:   "08-2025",
		}
		jsonData, _ := json.Marshal(updatedSub)
		req, _ := http.NewRequest(
			"PUT",
			"/subscriptions/"+strconv.Itoa(sub.ID),
			bytes.NewBuffer(jsonData),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var fetchedSub models.Subscription
		tx.First(&fetchedSub, sub.ID)
		assert.Equal(t, updatedSub.ServiceName, fetchedSub.ServiceName)
		assert.Equal(t, updatedSub.Price, fetchedSub.Price)
		assert.Equal(t, "2025-08-01", fetchedSub.StartDate)
	})
}

// Тест для UpdateSubscription с несуществующим ID
func TestUpdateSubscriptionNotFound(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		updatedSub := models.CreateSubscription{
			ServiceName: "Netflix",
			Price:       600,
			UserID:      uuid.New(),
			StartDate:   "08-2025",
		}
		jsonData, _ := json.Marshal(updatedSub)
		req, _ := http.NewRequest(
			"PUT",
			"/subscriptions/999",
			bytes.NewBuffer(jsonData),
		)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Тест для DeleteSubscription
func TestDeleteSubscription(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		userID := uuid.New()
		sub := models.Subscription{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      userID,
			StartDate:   "2025-07-01",
		}
		tx.Create(&sub)
		req, _ := http.NewRequest(
			"DELETE",
			"/subscriptions/"+strconv.Itoa(sub.ID),
			nil,
		)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		var deletedSub models.Subscription
		result := tx.First(&deletedSub, sub.ID)
		assert.Error(t, result.Error)
	})
}

// Тест для DeleteSubscription с несуществующим ID
func TestDeleteSubscriptionNotFound(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		req, _ := http.NewRequest("DELETE", "/subscriptions/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// Тест для ListSubscriptions
func TestListSubscriptions(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		userID := uuid.New()
		tx.Create(
			&models.Subscription{
				ServiceName: "Yandex Plus",
				Price:       400,
				UserID:      userID,
				StartDate:   "2025-07-01",
			},
		)
		tx.Create(
			&models.Subscription{
				ServiceName: "Netflix",
				Price:       600,
				UserID:      userID,
				StartDate:   "2025-08-01",
			},
		)
		req, _ := http.NewRequest("GET", "/subscriptions", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var subs []models.Subscription
		json.Unmarshal(w.Body.Bytes(), &subs)
		assert.Len(t, subs, 2)
	})
}

// Тест для GetTotalCostByPeriod
func TestGetTotalCostByPeriod(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		userID := uuid.New()
		tx.Create(&models.Subscription{
			ServiceName: "Yandex Plus",
			Price:       400,
			UserID:      userID,
			StartDate:   "2025-06-01",
		})
		tx.Create(&models.Subscription{
			ServiceName: "Netflix",
			Price:       600,
			UserID:      userID,
			StartDate:   "2025-07-01",
		})
		req, _ := http.NewRequest(
			"GET",
			"/subscriptions/total?user_id="+userID.String()+"&start_date=06-2025&end_date=07-2025",
			nil,
		)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		var response models.TotalCostResponse
		json.Unmarshal(w.Body.Bytes(), &response)
		// Проверяем: Yandex Plus (2 месяца * 400 = 800) + Netflix (1 месяц * 600 = 600) = 1400
		assert.Equal(t, 1400, response.Total)
	})
}

// Тест для GetTotalCostByPeriod с некорректной датой
func TestGetTotalCostByPeriodInvalidDate(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	withTransaction(t, func(tx *gorm.DB) {
		req, _ := http.NewRequest(
			"GET",
			"/subscriptions/total?user_id="+uuid.New().
				String()+
				"&start_date=invalid-date",
			nil,
		)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
