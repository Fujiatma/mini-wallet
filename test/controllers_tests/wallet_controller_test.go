package controller_test

import (
	"encoding/json"
	"errors"
	"github.com/julo/mini-wallet/src/controllers"
	"github.com/julo/mini-wallet/src/models"
	"github.com/julo/mini-wallet/test/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWalletController_EnableWalletController(t *testing.T) {
	mockWalletService := &mock.MockWalletService{}
	mockCustomerService := &mock.MockAuthService{}
	walletController := controllers.NewWalletController(mockWalletService, mockCustomerService)

	request, _ := http.NewRequest("GET", "/enable_wallet", nil)
	response := httptest.NewRecorder()

	expectedCustomer := &models.Customer{
		CustomerXID: "123",
		Name:        "John",
	}
	mockCustomerService.On("GetCustomerByCustomerXID", "123").Return(expectedCustomer, nil)
	expectedWallet := &models.Wallet{
		Balance: 0,
	}
	mockWalletService.On("EnableWallet", "123").Return(expectedWallet, nil)

	walletController.EnableWalletController(response, request)

	assert.Equal(t, http.StatusUnauthorized, response.Code)

	_ = map[string]interface{}{
		"wallet": expectedWallet,
	}
	var responseBody map[string]interface{}
	json.NewDecoder(response.Body).Decode(&responseBody)
}

func TestWalletController_EnableWalletController_WalletAlreadyEnabled(t *testing.T) {
	mockWalletService := &mock.MockWalletService{}
	mockCustomerService := &mock.MockAuthService{}
	walletController := controllers.NewWalletController(mockWalletService, mockCustomerService)

	request, _ := http.NewRequest("GET", "/enable_wallet", nil)
	response := httptest.NewRecorder()

	expectedCustomer := &models.Customer{
		CustomerXID: "123",
		Name:        "John",
		Wallet:      models.Wallet{},
	}
	mockCustomerService.On("GetCustomerByCustomerXID", "123").Return(expectedCustomer, nil)

	walletController.EnableWalletController(response, request)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

func TestWalletController_EnableWalletController_GetCustomerError(t *testing.T) {
	mockWalletService := &mock.MockWalletService{}
	mockCustomerService := &mock.MockAuthService{}
	walletController := controllers.NewWalletController(mockWalletService, mockCustomerService)

	request, _ := http.NewRequest("GET", "/enable_wallet", nil)
	response := httptest.NewRecorder()

	mockCustomerService.On("GetCustomerByCustomerXID", "123").Return(nil, errors.New("failed to get customer"))

	walletController.EnableWalletController(response, request)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
}
