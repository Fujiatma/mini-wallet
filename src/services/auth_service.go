package services

import (
	"github.com/julo/mini-wallet/src/models"
	"github.com/julo/mini-wallet/src/repositories"
)

type AuthService struct {
	CustomerRepository *repositories.CustomerRepository
}

func NewAuthService(customerRepository *repositories.CustomerRepository) *AuthService {
	return &AuthService{CustomerRepository: customerRepository}
}

func (s *AuthService) InitializeAccount(customerXID string, name string) (*models.Customer, error) {
	customer := &models.Customer{
		CustomerXID: customerXID,
		Name:        name,
	}
	err := s.CustomerRepository.Create(customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
