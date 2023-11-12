package data

import "gorm.io/gorm"

type Wallet struct {
	ID     string `gorm:"primaryKey"`
	Amount float64
}

type WalletStorage struct {
	DB *gorm.DB
}

func (w *WalletStorage) Save(wallets []Wallet) ([]Wallet, error) {
	result := w.DB.CreateInBatches(wallets, len(wallets))
	return wallets, result.Error
}

func (w *WalletStorage) Get(walletId string) (Wallet, error) {
	var wallet Wallet
	result := w.DB.First(&wallet, "id = ?", walletId)
	return wallet, result.Error
}
