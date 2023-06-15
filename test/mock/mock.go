package mock

import (
	"github.com/julo/mini-wallet/src/models"
	"github.com/stretchr/testify/mock"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) EnableWallet(customerXID string) (*models.Wallet, error) {
	args := m.Called(customerXID)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

func (m *MockWalletService) GetWalletTransactions(customerXID string) ([]models.Transaction, error) {
	args := m.Called(customerXID)
	return args.Get(0).([]models.Transaction), args.Error(1)
}

func (m *MockWalletService) AddVirtualMoney(customerXID string, amount int, referenceID string) (*models.Deposit, error) {
	args := m.Called(customerXID, amount, referenceID)
	return args.Get(0).(*models.Deposit), args.Error(1)
}

func (m *MockWalletService) UseVirtualMoney(customerXID string, amount int, referenceID string) (*models.Withdrawal, error) {
	args := m.Called(customerXID, amount, referenceID)
	return args.Get(0).(*models.Withdrawal), args.Error(1)
}

func (m *MockWalletService) DisableWallet(customerXID string, isDisabled bool) (*models.Wallet, error) {
	args := m.Called(customerXID, isDisabled)
	return args.Get(0).(*models.Wallet), args.Error(1)
}

// Auth Service
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) InitializeAccount(customerXID string, name string) (*models.Customer, *models.Wallet, error) {
	args := m.Called(customerXID, name)
	return args.Get(0).(*models.Customer), args.Get(1).(*models.Wallet), args.Error(2)
}

func (m *MockAuthService) GetCustomerByCustomerXID(customerXID string) (*models.Customer, error) {
	args := m.Called(customerXID)
	return args.Get(0).(*models.Customer), args.Error(1)
}
