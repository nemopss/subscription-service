package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nemopss/subscription-service/internal/db"
	"github.com/nemopss/subscription-service/internal/handlers"
)

func main() {
	r := gin.Default()
	err := db.InitDB()
	if err != nil {
		panic("db failed to init")
	}

	r.POST("/subscriptions", handlers.CreateSubscription)
	r.GET("/subscriptions/:id", handlers.GetSubscription)
	r.PUT("/subscriptions/:id", handlers.UpdateSubscription)
	r.DELETE("/subscriptions/:id", handlers.DeleteSubscription)
	r.GET("/subscriptions", handlers.ListSubscriptions)
	r.GET("/subscriptions/total", handlers.GetTotalCostByPeriod)

	r.Run()
}
