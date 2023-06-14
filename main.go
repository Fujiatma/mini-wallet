package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/julo/mini-wallet/src/config"
	"github.com/julo/mini-wallet/src/routes"
	"gorm.io/gorm"
	"log"
	"net/http"
)

var (
	router *gin.Engine
	db     *gorm.DB
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	config.InitDB()

	// Initialize the router
	router := mux.NewRouter()

	// Setup API routes
	routes.SetupRoutes(router, config.DB)

	// Start the server
	log.Println("Server started on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
