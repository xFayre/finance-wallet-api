// Copyright (c) 2020, Marcelo Jorge Vieira (https://github.com/mfinancecombr)
// Licensed under the BSD 3-Clause License

package db

import (
	"context"
	"time"

	"github.com/mfinancecombr/finance-wallet-api/wallet"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	brokersCollection    = "brokers"
	portfoliosCollection = "portfolios"
	operationsCollection = "operations"
)

type mongoSession struct {
	collection Collection
}

type DB interface {
	Create(d wallet.Queryable) (*mongo.InsertOneResult, error)
	Update(id string, d wallet.Queryable) (*mongo.UpdateResult, error)
	Delete(collectionName, id string) (*mongo.DeleteResult, error)

	DeleteBrokerByID(id string) (*mongo.DeleteResult, error)
	GetAllBrokers() (*wallet.BrokersList, error)
	GetBrokerByID(id string) (*wallet.Broker, error)
	InsertBroker(d interface{}) (*mongo.InsertOneResult, error)
	UpdateBroker(id string, d interface{}) (*mongo.UpdateResult, error)

	DeletePortfolioByID(id string) (*mongo.DeleteResult, error)
	GetAllPortfolios() ([]wallet.Portfolio, error)
	GetPortfolioByID(id string) (*wallet.Portfolio, error)
	GetPortfolioItems(p *wallet.Portfolio, year int) error
	InsertPortfolio(d interface{}) (*mongo.InsertOneResult, error)
	UpdatePortfolio(id string, d interface{}) (*mongo.UpdateResult, error)

	GetPositionBySymbol(itemType, symbol string) (*wallet.PortfolioItem, error)

	GetAllOperations() (interface{}, error)
	GetAllPurchases() (interface{}, error)
	GetAllSales() (interface{}, error)

	Ping() error

	GetAllStocksOperations() (wallet.StockList, error)
	GetStockOperationByID(id string) (*wallet.Stock, error)

	GetAllFIIsOperations() (wallet.FIIList, error)
	GetFIIOperationByID(id string) (*wallet.FII, error)

	GetAllTreasuriesDirectsOperations() (wallet.TreasuryDirectList, error)
	GetTreasuryDirectOperationByID(id string) (*wallet.TreasuryDirect, error)

	GetAllCertificatesOfDepositsOperations() (wallet.CertificateOfDepositList, error)
	GetCertificateOfDepositOperationByID(id string) (*wallet.CertificateOfDeposit, error)

	GetAllStocksFundsOperations() (wallet.StockFundList, error)
	GetStockFundOperationByID(id string) (*wallet.StockFund, error)

	GetAllFICFIOperations() (wallet.FICFIList, error)
	GetFICFIOperationByID(id string) (*wallet.FICFI, error)
}

func (m *mongoSession) Create(d wallet.Queryable) (*mongo.InsertOneResult, error) {
	log.Debug("[DB] insertOperation")
	return m.collection.InsertOne(d.GetCollectionName(), d)
}

func (m *mongoSession) Update(id string, d wallet.Queryable) (*mongo.UpdateResult, error) {
	log.Debug("[DB] updateOperation")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	f := bson.D{{"_id", objectId}}
	u := bson.D{{"$set", d}}
	return m.collection.UpdateOne(d.GetCollectionName(), f, u)
}

func (m *mongoSession) Delete(id, collectionName string) (*mongo.DeleteResult, error) {
	log.Debug("[DB] Delete")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	q := bson.M{"_id": objectId}
	return m.collection.DeleteOne(collectionName, q)
}

func newDBContext() (context.Context, context.CancelFunc) {
	log.Debug("[DB] New DB context")
	timeout := viper.GetDuration("db.operation.timeout")
	return context.WithTimeout(context.Background(), timeout*time.Second)
}

func NewMongoSession() (DB, error) {
	log.Debug("[DB] New mongo session")
	dbURI := viper.GetString("mongodb.endpoint")
	dbName := viper.GetString("mongodb.name")
	ctx, _ := newDBContext()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURI))
	if err != nil {
		log.Errorf("[DB] Error on create mongo session: %s", err)
	}
	mongo := &mongoSession{
		collection: &mongoCollection{
			session: client,
			dbName:  dbName,
		},
	}
	return mongo, err
}

func (m *mongoSession) Ping() error {
	log.Debug("[DB] Ping")
	return m.collection.Ping()
}
