package routes

import (
	"github.com/gorilla/mux"
	"github.com/julo/mini-wallet/src/controllers"
	"github.com/julo/mini-wallet/src/middleware"
	"github.com/julo/mini-wallet/src/repositories"
	"github.com/julo/mini-wallet/src/services"
	"gorm.io/gorm"
)

func SetupRoutes(router *mux.Router, db *gorm.DB) {
	authRepository := repositories.NewCustomerRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	// Initialize account route
	router.HandleFunc("/api/v1/init", authController.InitializeAccount).Methods("POST")

	// Wallet routes
	walletRouter := router.PathPrefix("/api/v1/wallet").Subrouter()

	// Apply JWT middleware to protected routes
	walletRouter.Use(middleware.JWTMiddleware)

	//walletRouter.HandleFunc("", walletController.EnableWallet).Methods("POST")
}
