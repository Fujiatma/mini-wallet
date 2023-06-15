package controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/julo/mini-wallet/src/controllers"
	"github.com/julo/mini-wallet/src/models"
	"github.com/julo/mini-wallet/test/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthController_InitializeAccount(t *testing.T) {
	mockAuthService := &mock.MockAuthService{}

	authController := controllers.NewAuthController(mockAuthService)

	payload := map[string]string{
		"customer_xid": "123",
		"user_name":    "John",
	}
	body, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", "/initialize", bytes.NewBuffer(body))
	response := httptest.NewRecorder()

	expectedCustomer := &models.Customer{
		CustomerXID: "123",
		Name:        "John",
	}
	mockAuthService.On("InitializeAccount", "123", "John").Return(expectedCustomer, nil, nil)

	authController.InitializeAccount(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	_ = map[string]interface{}{
		"token": "generated_token",
	}
	var responseBody map[string]interface{}
	json.NewDecoder(response.Body).Decode(&responseBody)
}

func TestGenerateToken(t *testing.T) {
	customer := &models.Customer{
		CustomerXID: "123",
	}
	_ = controllers.GenerateToken(customer)
}
