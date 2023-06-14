package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/julo/mini-wallet/src/models"
	resp "github.com/julo/mini-wallet/src/response"
	"github.com/julo/mini-wallet/src/services"
	"net/http"
	"os"
	"time"
)

type AuthController struct {
	AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (c *AuthController) InitializeAccount(w http.ResponseWriter, r *http.Request) {
	// Parse the customer_xid from the request body
	err := r.ParseMultipartForm(0)
	if err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}
	customerXID := r.FormValue("customer_xid")
	name := r.FormValue("user_name")

	// Initialize the account
	customer, err := c.AuthService.InitializeAccount(customerXID, name)
	if err != nil {
		http.Error(w, "Failed to initialize account", http.StatusInternalServerError)
		return
	}

	// Generate the token
	token := generateToken(customer)

	// Create the response
	response := map[string]interface{}{
		"token": token,
	}

	resp.ConstructResponse(w, http.StatusOK, response, nil)

}

func generateToken(customer *models.Customer) string {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Failed to load .env file:", err)
		return ""
	}

	// Read secret key from environment variable
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		fmt.Println("SECRET_KEY is not set in .env file")
		return ""
	}

	// Set claims
	claims := jwt.MapClaims{
		"customer_xid": customer.CustomerXID,
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
	}

	// Generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		// Handle error
		fmt.Println("Failed to generate token:", err)
		return ""
	}

	return tokenString
}
