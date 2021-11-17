package models

import (
	"avito-intern/internal/app/balance/models"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
)

func validIntFunc(value interface{}) error {
	res, ok := value.(int64)
	if !ok || res < 0 {
		return errors.New("invalid field")
	}
	return nil
}
func validFloatFunc(value interface{}) error {
	res, ok := value.(float64)
	if !ok || res < 0 {
		return errors.New("invalid field")
	}
	return nil
}

type RequestTransfer struct {
	SenderID   int64   `json:"sender_id" validate:"required"`
	ReceiverID int64   `json:"receiver_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
}

func (req *RequestTransfer) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.SenderID, validation.By(validIntFunc)),
		validation.Field(&req.ReceiverID, validation.By(validIntFunc)),
		validation.Field(&req.Amount, validation.By(validFloatFunc)))
}

type RequestUpdateBalance struct {
	Type   int64   `json:"operation"`
	Amount float64 `json:"amount"`
}

func (req *RequestUpdateBalance) Validate() error {
	validTypeFunc := func(value interface{}) error {
		res, ok := value.(int64)
		if !ok || (res != int64(models.DIFF_BALANCE) && res != int64(models.ADD_BALANCE)) {
			return errors.New("invalid field")
		}
		return nil
	}
	return validation.ValidateStruct(req,
		validation.Field(&req.Amount, validation.By(validFloatFunc)),
		validation.Field(&req.Type, validation.By(validTypeFunc)))
}
