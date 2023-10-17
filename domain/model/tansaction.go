package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactionConfirmed string = "confirmed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []*Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"required"`
	Amount            float64  `json:"amount" valid:"required"`
	PixKeyTo          *PixKey  `valid:"-"`
	Status            string   `json:"status"`
	Description       string   `json:"description" alid:"required"`
	CancelDescription string   `json:"cancel_description" alid:"required"`
}

func (t *Transaction) isValid() error {
	if _, err := govalidator.ValidateStruct(t); err != nil {
		return err
	}

	if t.Amount <= 0 {
		return errors.New("The ammount must be higher than 0")
	}

	if t.Status != TransactionPending && t.Status != TransactionCompleted && t.Status != TransactionError {
		return errors.New("Invalid status for the transaction")
	}

	if t.PixKeyTo.Account.ID == t.AccountFrom.ID {
		return errors.New("the source and destination cannot be the same")
	}

	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	transaction := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Status:      TransactionPending,
		Description: description,
	}

	transaction.ID = uuid.NewV4().String()
	transaction.CreatedAt = time.Now()

	if err := transaction.isValid(); err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (t *Transaction) Complete() error {
	t.Status = TransactionCompleted
	t.UpdatedAt = time.Now()

	return t.isValid()
}

func (t *Transaction) Confirm() error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()

	return t.isValid()

}
func (t *Transaction) Cancel(description string) error {
	t.Status = TransactionConfirmed
	t.UpdatedAt = time.Now()
	t.CancelDescription = description

	return t.isValid()
}
