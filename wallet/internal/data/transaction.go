package data

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type TransactionType string

const (
	Deposit  TransactionType = "deposit"
	Withdraw TransactionType = "withdraw"
)

type Transaction struct {
	Id                    string `gorm:"primaryKey"`
	UserId                string
	WalletId              string
	Amount                decimal.Decimal
	TransactionType       TransactionType
	AmountBefore          decimal.Decimal
	AmountAfter           decimal.Decimal
	FromPaymentSystem     string
	FromPaymentIdentifier string
	ToPaymentSystem       string
	ToPaymentIdentifier   string
	Currency              string
	Description           string
	CreatedAt             time.Time
}

type BalanceWasChangedError string

func (e BalanceWasChangedError) Error() string {
	return string(e)
}

type TransactionStorage struct {
	DB            *gorm.DB
	WalletStorage WalletStorage
}

func (t *TransactionStorage) Save(transaction Transaction) (Transaction, error) {
	t.DB.Transaction(func(tx *gorm.DB) error {
		walletStorage := WalletStorage{DB: tx}
		rowsAffected, err := walletStorage.UpdateAmount(transaction.WalletId, transaction.AmountAfter, transaction.AmountBefore)
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return BalanceWasChangedError("during transaction balance was changed")
		}
		tx.Save(transaction)
		return nil
	})

	return transaction, nil
}

func (w *TransactionStorage) GetById(transactionId string) (Transaction, error) {
	var transaction Transaction
	result := w.DB.First(&transaction, "id = ?", transactionId)
	return transaction, result.Error
}
