package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const (
	WalletStatusEnable  = "Enabled"
	WalletStatusDisable = "Disabled"

	StatusSuccess = "Success"

	TransactionTypeDeposit   = "deposit"
	TransactionTypeWithdrawn = "withdrawal"
)

type Customer struct {
	ID          string    `gorm:"type:varchar(36);primaryKey"`
	Name        string    `gorm:"not null"`
	CustomerXID string    `gorm:"column:customer_xid;type:varchar(75);unique"`
	Wallet      Wallet    `gorm:"foreignKey:CustomerID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoCreateTime"`
}

type Wallet struct {
	ID           string `gorm:"type:varchar(36);primaryKey"`
	CustomerID   string `gorm:"unique;column:customer_id" json:"owned_by"`
	Status       string
	EnabledAt    time.Time     `gorm:"autoCreateTime"`
	Balance      int           `gorm:"default:0"`
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"autoCreateTime"`
	Deposits     []Deposit     `gorm:"-"`
	Withdrawals  []Withdrawal  `gorm:"-"`
	Transactions []Transaction `gorm:"-"`
}

type Transaction struct {
	ID              string    `gorm:"type:varchar(36);primaryKey"`
	WalletID        string    `gorm:"type:varchar(36);column:wallet_id"`
	TransactionType string    `gorm:"type:enum('deposit','withdrawal')"`
	Amount          int       `gorm:"not null"`
	ReferenceID     string    `gorm:"not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoCreateTime"`
}

type Deposit struct {
	ID          string `gorm:"type:varchar(36);primaryKey"`
	WalletID    string `gorm:"column:wallet_id"`
	DepositedBy string `gorm:"column:deposited_by"`
	Status      string
	DepositedAt time.Time `gorm:"autoCreateTime"`
	Amount      int       `gorm:"not null"`
	ReferenceID string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoCreateTime"`
}

type Withdrawal struct {
	ID          string `gorm:"type:varchar(36);primaryKey"`
	WalletID    string `gorm:"column:wallet_id"`
	WithdrawnBy string `gorm:"column:withdrawn_by"`
	Status      string
	WithdrawnAt time.Time `gorm:"autoCreateTime"`
	Amount      int       `gorm:"not null"`
	ReferenceID string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoCreateTime"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New().String()
	return
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New().String()
	return
}

func (d *Deposit) BeforeCreate(tx *gorm.DB) (err error) {
	d.ID = uuid.New().String()
	return
}

func (w *Withdrawal) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New().String()
	return
}
