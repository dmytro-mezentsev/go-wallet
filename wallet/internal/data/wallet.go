package data

import "gorm.io/gorm"

type Wallet struct {
	ID     string `gorm:"primaryKey"`
	Amount float64
}

type WalletStorage struct {
	DB *gorm.DB
}

func (w *WalletStorage) Save(wallet Wallet) (Wallet, error) {
	result := w.DB.Create(&wallet)
	return wallet, result.Error
}
