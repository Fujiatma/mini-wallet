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
