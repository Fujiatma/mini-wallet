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
	walletRepository := repositories.NewWalletRepository(db)

	authService := services.NewAuthService(authRepository, walletRepository)
	authController := controllers.NewAuthController(authService)

	// Initialize account route
	router.HandleFunc("/api/v1/init", authController.InitializeAccount).Methods("POST")

	walletService := services.NewWalletService(authRepository, walletRepository)
	walletController := controllers.NewWalletController(walletService, authService)

	// Wallet routes
	walletRouter := router.PathPrefix("/api/v1/wallet").Subrouter()

	// Apply JWT middleware to protected routes
	walletRouter.Use(middleware.JWTMiddleware)
	walletRouter.HandleFunc("", walletController.EnableWalletController).Methods("POST")
	walletRouter.HandleFunc("", walletController.GetWalletBalanceController).Methods("GET")
	walletRouter.HandleFunc("/transactions", walletController.GetWalletTransactionController).Methods("GET")
	walletRouter.HandleFunc("/deposits", walletController.AddVirtualMoneyController).Methods("POST")
	walletRouter.HandleFunc("/withdrawals", walletController.UseVirtualMoneyController).Methods("POST")
	walletRouter.HandleFunc("", walletController.DisableWalletController).Methods("PATCH")

}
