package models

import (
	"time"

	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeSwap     TransactionType = "swap"
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeBorrow   TransactionType = "borrow"
	TransactionTypeRepay    TransactionType = "repay"
)

type Transaction struct {
	gorm.Model
	UserID          uint
	Type            TransactionType
	TokenIn         string
	TokenOut        string
	AmountIn        string
	AmountOut       string
	Price           string
	Status          string
	BlockNumber     uint64
	TransactionHash string
	Timestamp       time.Time
}

type TransactionService struct {
	db *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{db: db}
}

func (s *TransactionService) CreateTransaction(tx *Transaction) error {
	return s.db.Create(tx).Error
}

func (s *TransactionService) GetUserTransactions(userID uint, limit int) ([]Transaction, error) {
	var transactions []Transaction
	err := s.db.Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}

func (s *TransactionService) GetTransactionByHash(hash string) (*Transaction, error) {
	var transaction Transaction
	err := s.db.Where("transaction_hash = ?", hash).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (s *TransactionService) UpdateTransactionStatus(hash string, status string) error {
	return s.db.Model(&Transaction{}).
		Where("transaction_hash = ?", hash).
		Update("status", status).Error
}

func (s *TransactionService) GetRecentTransactions(limit int) ([]Transaction, error) {
	var transactions []Transaction
	err := s.db.Order("created_at desc").
		Limit(limit).
		Find(&transactions).Error
	return transactions, err
}
