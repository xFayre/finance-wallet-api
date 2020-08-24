// Copyright (c) 2020, Marcelo Jorge Vieira
// Licensed under the BSD 3-Clause License

package db

import (
	"github.com/mfinancecombr/finance-wallet-api/wallet"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *mongoSession) GetAllFIIsOperations() (wallet.FIIList, error) {
	log.Debug("[DB] GetAllFIIsOperations")
	query := bson.M{"itemType": "fiis"}
	results, err := m.collection.FindAll(operationsCollection, query)
	if err != nil {
		return nil, err
	}
	operationsList := wallet.FIIList{}
	for _, result := range results {
		bsonBytes, _ := bson.Marshal(result)
		buy := wallet.FII{}
		bson.Unmarshal(bsonBytes, &buy)
		operationsList = append(operationsList, buy)
	}
	return operationsList, nil
}

func (m *mongoSession) GetFIIOperationByID(id string) (*wallet.FII, error) {
	log.Debug("[DB] GetFIIOperationByID")
	fii := &wallet.FII{}
	if err := m.getOperationByID(operationsCollection, id, fii); err != nil {
		return nil, err
	}
	if fii.Symbol == "" {
		return nil, nil
	}
	return fii, nil
}
