package services

import (
	"github.com/julo/mini-wallet/src/models"
	"github.com/julo/mini-wallet/src/repositories"
)

type AuthService struct {
	CustomerRepository *repositories.CustomerRepository
	WalletRepository   *repositories.WalletRepository
}

func NewAuthService(customerRepository *repositories.CustomerRepository, walletRepository *repositories.WalletRepository) *AuthService {
	return &AuthService{
		CustomerRepository: customerRepository,
		WalletRepository:   walletRepository,
	}
}

func (s *AuthService) InitializeAccount(customerXID string, name string) (*models.Customer, *models.Wallet, error) {
	// Create customer
	customer := &models.Customer{
		CustomerXID: customerXID,
		Name:        name,
	}

	err := s.CustomerRepository.Create(customer)
	if err != nil {
		return nil, nil, err
	}

	// Create wallet
	wallet := &models.Wallet{
		CustomerID: customer.ID,
		Status:     models.WalletStatusDisable,
	}

	err = s.WalletRepository.Create(wallet)
	if err != nil {
		// Menghapus pelanggan yang sudah dibuat jika pembuatan wallet gagal
		s.CustomerRepository.Delete(customer)
		return nil, nil, err
	}

	// Menghubungkan wallet dengan customer
	customer.Wallet = *wallet

	return customer, wallet, nil
}

func (s *AuthService) GetCustomerByCustomerXID(customerXID string) (*models.Customer, error) {
	customer, err := s.CustomerRepository.GetByCustomerXID(customerXID)
	if err != nil {
		return nil, err
	}

	return customer, nil
}
