package repositories

import (
	"errors"
	"github.com/julo/mini-wallet/src/models"
	"gorm.io/gorm"
)

type WalletRepository struct {
	Db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{Db: db}
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	err := r.Db.Create(wallet).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WalletRepository) GetByCustomerID(customerID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.Db.Where("customer_id = ?", customerID).First(&wallet).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Update(wallet *models.Wallet) error {
	err := r.Db.Save(wallet).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WalletRepository) GetTransactionsByWalletID(walletID string) ([]models.Transaction, error) {
	transactions := []models.Transaction{}
	result := r.Db.Where("wallet_id = ?", walletID).Find(&transactions)
	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}

func (r *WalletRepository) CreateDeposit(deposit *models.Deposit) error {
	// Memulai transaksi
	tx := r.Db.Begin()

	// Menyimpan deposit ke dalam database dengan transaksi
	if err := tx.Create(deposit).Error; err != nil {
		// Rollback transaksi jika terjadi error
		tx.Rollback()
		return err
	}

	// Menjalankan komit transaksi
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *WalletRepository) CreateWithdrawal(withdrawal *models.Withdrawal) error {
	tx := r.Db.Begin()

	if err := tx.Create(withdrawal).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Mengupdate saldo wallet setelah withdrawal
	wallet := &models.Wallet{}
	if err := tx.Where("id = ?", withdrawal.WalletID).First(wallet).Error; err != nil {
		tx.Rollback()
		return err
	}

	updatedBalance := wallet.Balance - withdrawal.Amount
	if updatedBalance < 0 {
		tx.Rollback()
		return errors.New("insufficient balance")
	}

	if err := tx.Model(wallet).Update("balance", updatedBalance).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *WalletRepository) UpdateWalletBalance(walletID string, newBalance int) error {
	tx := r.Db.Begin()

	// Mengambil wallet berdasarkan ID
	wallet := &models.Wallet{}
	if err := tx.Where("id = ?", walletID).First(wallet).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Mengupdate saldo wallet
	wallet.Balance = newBalance
	if err := tx.Save(wallet).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *WalletRepository) CreateTransaction(transaction *models.Transaction) error {
	// Membuat instance database transaction
	tx := r.Db.Begin()

	// Menyimpan transaksi ke dalam database
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback() // Rollback transaction jika terjadi error
		return err
	}

	// Commit transaction jika tidak ada error
	return tx.Commit().Error
}

func (r *WalletRepository) Delete(wallet *models.Wallet) error {
	err := r.Db.Delete(wallet).Error
	if err != nil {
		return err
	}
	return nil
}
