package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nemopss/subscription-service/internal/db"
)

func main() {
	r := gin.Default()
	err := db.InitDB()
	if err != nil {
		panic("db failed to init")
	}

	r.Run()
}
