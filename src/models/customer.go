package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Customer struct {
	ID          string    `gorm:"type:varchar(36);primaryKey"`
	Name        string    `gorm:"not null"`
	CustomerXID string    `gorm:"unique"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoCreateTime"`
	Wallet      Wallet    // Relasi One-to-One dengan model Wallet
}

type Wallet struct {
	ID         string `gorm:"type:varchar(36);primaryKey"`
	CustomerID string `gorm:"unique"`
	Status     string
	EnabledAt  time.Time `gorm:"autoCreateTime"`
	Balance    int       `gorm:"default:0"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoCreateTime"`
}

func (c *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New().String()
	return
}

func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	w.ID = uuid.New().String()
	return
}
