// Copyright (c) 2020, Marcelo Jorge Vieira (https://github.com/mfinancecombr)
// Licensed under the BSD 3-Clause License

package db

import (
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *mongoSession) insertOperation(d interface{}) (*mongo.InsertOneResult, error) {
	log.Debug("[DB] insertOperation")
	return m.collection.InsertOne(operationsCollection, d)
}

func (m *mongoSession) DeleteOperationByID(id string) (*mongo.DeleteResult, error) {
	log.Debug("[DB] DeleteOperationByID")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	q := bson.M{"_id": objectId}
	return m.collection.DeleteOne(operationsCollection, q)
}

func (m *mongoSession) getOperationsSymbols(filter bson.M) ([]interface{}, error) {
	log.Debug("[DB] getOperationSymbols")
	return m.collection.Distinct(operationsCollection, "symbol", filter)
}

func (m *mongoSession) getAllOperationsBySymbol(symbol string, year int) ([]bson.M, error) {
	log.Debug("[DB] getAllOperationsBySymbol")
	date := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)
	query := bson.M{"symbol": symbol, "date": bson.M{"$lte": date}}
	opts := options.Find().SetSort(bson.D{{"date", 1}})
	return m.collection.FindAll(operationsCollection, query, opts)
}

func (m *mongoSession) GetAllOperations() (interface{}, error) {
	log.Debug("[DB] GetAllOperations")
	return m.collection.FindAll(operationsCollection, bson.M{})
}

func (m *mongoSession) GetAllPurchases() (interface{}, error) {
	log.Debug("[DB] GetAllOperations")
	query := bson.M{"type": "purchase"}
	opts := options.Find().SetSort(bson.D{{"date", -1}})
	return m.collection.FindAll(operationsCollection, query, opts)
}

func (m *mongoSession) GetAllSales() (interface{}, error) {
	log.Debug("[DB] GetAllOperations")
	query := bson.M{"type": "sale"}
	opts := options.Find().SetSort(bson.D{{"date", -1}})
	return m.collection.FindAll(operationsCollection, query, opts)
}
