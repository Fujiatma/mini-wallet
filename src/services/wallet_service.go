package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/julo/mini-wallet/src/models"
	"github.com/julo/mini-wallet/src/repositories"
	"github.com/julo/mini-wallet/src/response"
	"log"
	"time"
)

type WalletService struct {
	CustomerRepository *repositories.CustomerRepository
	WalletRepository   *repositories.WalletRepository
}

func NewWalletService(customerRepository *repositories.CustomerRepository, walletRepository *repositories.WalletRepository) *WalletService {
	return &WalletService{
		CustomerRepository: customerRepository,
		WalletRepository:   walletRepository,
	}
}

func (s *WalletService) EnableWallet(customerXID string) (*models.Wallet, error) {
	customer, err := s.CustomerRepository.GetByCustomerXID(customerXID)
	if err != nil {
		return nil, err
	}

	if customer.Wallet.Status == models.WalletStatusEnable {
		return nil, response.WalletAlreadyEnabledError{
			Code:    800,
			Message: "Wallet already enabled",
		}
	}

	if customer != nil && customer.Wallet.ID != "" {
		// Mengupdate status Wallet menjadi "enabled"
		customer.Wallet.Status = models.WalletStatusEnable
		customer.Wallet.EnabledAt = time.Now()

		err = s.WalletRepository.Update(&customer.Wallet)
		if err != nil {
			return nil, err
		}
	}

	return &customer.Wallet, nil
}

func (s *WalletService) GetWalletTransactions(customerXID string) ([]models.Transaction, error) {
	// Mendapatkan customer berdasarkan customerXID
	customer, err := s.CustomerRepository.GetByCustomerXID(customerXID)
	if err != nil {
		return nil, err
	}

	// Check status wallet
	if customer.Wallet.Status != models.WalletStatusEnable {
		return nil, response.WalletAlreadyDisabledError{
			Code:    801,
			Message: "Wallet already disabled",
		}
	}

	// Mendapatkan transactions wallet berdasarkan WalletID pada customer
	transactions, err := s.WalletRepository.GetTransactionsByWalletID(customer.Wallet.ID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (s *WalletService) AddVirtualMoney(customerXID string, amount int, referenceID string) (*models.Deposit, error) {
	// Mendapatkan data pelanggan berdasarkan CustomerXID
	customer, err := s.CustomerRepository.GetByCustomerXID(customerXID)
	if err != nil {
		return nil, err
	}

	// Check status wallet
	if customer.Wallet.Status != models.WalletStatusEnable {
		return nil, response.WalletAlreadyDisabledError{
			Code:    801,
			Message: "Wallet already disabled",
		}
	}

	// Membuat objek deposit baru
	deposit := &models.Deposit{
		ID:          uuid.New().String(),
		WalletID:    customer.Wallet.ID,
		DepositedBy: customer.ID,
		Status:      models.StatusSuccess,
		DepositedAt: time.Now(),
		Amount:      amount,
		ReferenceID: referenceID,
	}

	// Menyimpan deposit ke dalam database
	err = s.WalletRepository.CreateDeposit(deposit)
	if err != nil {
		return nil, err
	}

	// Membuat objek transaksi baru
	transaction := &models.Transaction{
		ID:              uuid.New().String(),
		TransactionType: models.TransactionTypeDeposit,
		Amount:          amount,
		ReferenceID:     referenceID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Time{},
		WalletID:        customer.Wallet.ID,
	}

	// Menyimpan transaksi ke dalam database
	err = s.WalletRepository.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Goroutine untuk menunda pembaruan saldo wallet
	go func() {
		time.Sleep(5 * time.Second) // Menunda pembaruan saldo selama 5 detik
		// Memperbarui saldo wallet dengan menambahkan jumlah deposit
		customer.Wallet.Balance += amount
		// Menyimpan perubahan saldo wallet ke dalam database
		err = s.WalletRepository.UpdateWalletBalance(customer.Wallet.ID, customer.Wallet.Balance)
		if err != nil {
			// Menangani error jika terjadi masalah saat memperbarui saldo wallet
			log.Printf("Error updating wallet balance: %v", err)
		}
	}()

	return deposit, nil
}

func (s *WalletService) UseVirtualMoney(customerXID string, amount int, referenceID string) (*models.Withdrawal, error) {
	// Mendapatkan informasi wallet dari customerXID
	customer, err := s.CustomerRepository.GetByCustomerXID(customerXID)
	if err != nil {
		return nil, err
	}

	// Check status wallet
	if customer.Wallet.Status != models.WalletStatusEnable {
		return nil, response.WalletAlreadyDisabledError{
			Code:    801,
			Message: "Wallet already disabled",
		}
	}

	// Memvalidasi saldo wallet
	if customer.Wallet.Balance < amount {
		return nil, errors.New("Insufficient balance")
	}

	// Membuat objek withdrawal
	withdrawal := &models.Withdrawal{
		ID:          uuid.New().String(),
		WalletID:    customer.Wallet.ID,
		Amount:      amount,
		ReferenceID: referenceID,
		WithdrawnBy: customer.ID,
		Status:      models.StatusSuccess,
		WithdrawnAt: time.Now(),
	}

	// Menyimpan withdrawal ke dalam database
	err = s.WalletRepository.CreateWithdrawal(withdrawal)
	if err != nil {
		return nil, err
	}

	// Membuat objek transaksi baru
	transaction := &models.Transaction{
		ID:              uuid.New().String(),
		WalletID:        customer.Wallet.ID,
		TransactionType: models.TransactionTypeWithdrawn,
		Amount:          amount,
		ReferenceID:     referenceID,
		CreatedAt:       time.Now(),
	}

	// Menyimpan transaksi ke dalam database
	err = s.WalletRepository.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	// Goroutine untuk menunda pembaruan saldo wallet
	go func() {
		time.Sleep(5 * time.Second) // Menunda pembaruan saldo selama 5 detik
		// Mengurangi saldo wallet
		err = s.WalletRepository.UpdateWalletBalance(customer.Wallet.ID, customer.Wallet.Balance-amount)
		if err != nil {
			// Menangani error jika terjadi masalah saat memperbarui saldo wallet
			log.Printf("Error updating wallet balance: %v", err)
		}
	}()
	return withdrawal, nil

}

func (s *WalletService) DisableWallet(customerXID string, isDisabled bool) (*models.Wallet, error) {
	// Mendapatkan data pelanggan berdasarkan CustomerXID
	customer, err := s.CustomerRepository.GetByCustomerXID(customerXID)
	if err != nil {
		return nil, err
	}

	// Check status wallet
	if customer.Wallet.Status == models.WalletStatusDisable {
		return nil, response.WalletAlreadyDisabledError{
			Code:    801,
			Message: "Wallet already disabled",
		}
	}

	// Mengubah status disable wallet
	if isDisabled {
		customer.Wallet.Status = models.WalletStatusDisable
	}

	// Menyimpan perubahan status wallet ke dalam database
	err = s.WalletRepository.Update(&customer.Wallet)
	if err != nil {
		return nil, err
	}

	return &customer.Wallet, nil
}
