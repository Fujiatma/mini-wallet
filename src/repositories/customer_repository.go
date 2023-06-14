package repositories

import (
	"github.com/julo/mini-wallet/src/models"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	Db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{Db: db}
}

func (r *CustomerRepository) Create(customer *models.Customer) error {
	err := r.Db.Create(customer).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CustomerRepository) GetByCustomerXID(customerXID string) (*models.Customer, error) {
	var customer models.Customer
	err := r.Db.Where("customer_xid = ?", customerXID).Preload("Wallet").First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Update(customer *models.Customer) error {
	err := r.Db.Save(customer).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CustomerRepository) Delete(customer *models.Customer) error {
	err := r.Db.Delete(customer).Error
	if err != nil {
		return err
	}
	return nil
}
