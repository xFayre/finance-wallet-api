// Copyright (c) 2020, Marcelo Jorge Vieira
// Licensed under the BSD 3-Clause License

package wallet

import (
	"time"
)

type TreasuryDirect struct {
	BrokerID   string     `json:"brokerId" bson:"brokerId" validate:"required"`
	Commission float64    `json:"commission" bson:"commission"`
	Date       *time.Time `json:"date" bson:"date" validate:"required"`
	//DueDate           *time.Time `json:"dueDate" bson:"dueDate" validate:"required"`
	FixedInterestRate float64 `json:"fixedInterestRate" bson:"fixedInterestRate" validate:"required"`
	ID                string  `json:"id,omitempty" bson:"_id,omitempty"`
	ItemType          string  `json:"itemType" bson:"itemType" validate:"required"`
	PortfolioID       string  `json:"portfolioId" bson:"portfolioId" validate:"required"`
	Price             float64 `json:"price" bson:"price" validate:"required"`
	Shares            float64 `json:"shares" bson:"shares" validate:"required"`
	Symbol            string  `json:"symbol" bson:"symbol" validate:"required"`
	Type              string  `json:"type" bson:"type" validate:"required"`
}

type TreasuryDirectList []TreasuryDirect

func NewTreasuryDirect() *TreasuryDirect {
	return &TreasuryDirect{ItemType: "treasury-direct"}
}

func (s TreasuryDirect) GetPrice() float64 {
	return s.Price
}

func (s TreasuryDirect) GetShares() float64 {
	return s.Shares
}

func (s TreasuryDirect) GetComission() float64 {
	return s.Commission
}

func (s TreasuryDirect) GetType() string {
	return s.Type
}

func (s TreasuryDirect) GetBrokerID() string {
	return s.BrokerID
}
