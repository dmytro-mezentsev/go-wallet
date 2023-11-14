package data

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Wallet struct {
	Id     string `gorm:"primaryKey"`
	Amount decimal.Decimal
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
func (w *WalletStorage) UpdateAmount(walletId string, newAmount, amountBefore decimal.Decimal) (int64, error) {
	result := w.DB.Model(Wallet{}).Where("id = ? AND amount = ?", walletId, amountBefore).Update("amount", newAmount)
	return result.RowsAffected, result.Error
}
