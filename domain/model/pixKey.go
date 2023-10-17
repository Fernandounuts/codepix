package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type PixKeyRepositoryInterface interface {
	Register(pixKey *PixKey) (*PixKey, error)
	FindKeyByKind(key string, kind string) (*PixKey, error)
	AddBank(bank *Bank) error
	AddAccount(account *Account) error
	FindAccount(id string) (*Account, error)
}

type PixKey struct {
	Base      `valid:"required" valid:"notnull"`
	Kind      string `json:"kind" valid:"notnull"`
	Key       string `json:"key" valid:"notnull"`
	AccountId string `json:"account_id" valid:"notnull"`
	Account   *Account
	Status    string `json:"status" valid:"notnull"`
}

func (pixKey *PixKey) isValid() error {

	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid type of key")
	}

	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("invalid status")
	}

	if _, err := govalidator.ValidateStruct(pixKey); err != nil {
		return err
	}

	return nil
}

func NewPixKey(kind string, account *Account, key string) (*PixKey, error) {
	pixKey := PixKey{
		Kind:    kind,
		Key:     key,
		Account: account,
		Status:  "active",
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now()

	if err := pixKey.isValid(); err != nil {
		return nil, err
	}

	return &pixKey, nil
}
