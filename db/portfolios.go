// Copyright (c) 2020, Marcelo Jorge Vieira (https://github.com/mfinancecombr)
// Licensed under the BSD 3-Clause License

package db

import (
	"reflect"

	"github.com/mfinancecombr/finance-wallet-api/financeapi"
	"github.com/mfinancecombr/finance-wallet-api/wallet"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *mongoSession) InsertPortfolio(d interface{}) (*mongo.InsertOneResult, error) {
	log.Debug("[DB] InsertPortfolio")
	return m.collection.InsertOne(portfoliosCollection, d)
}

func (m *mongoSession) DeletePortfolioByID(id string) (*mongo.DeleteResult, error) {
	log.Debug("[DB] DeletePortfolioByID")
	q := bson.M{"_id": id}
	return m.collection.DeleteOne(portfoliosCollection, q)
}

func (m *mongoSession) UpdatePortfolio(id string, d interface{}) (*mongo.UpdateResult, error) {
	log.Debug("[DB] UpdatePortfolio")
	f := bson.D{{"_id", id}}
	u := bson.D{{"$set", d}}
	return m.collection.UpdateOne(portfoliosCollection, f, u)
}

func (m *mongoSession) GetAllPortfolios() ([]wallet.Portfolio, error) {
	log.Debug("[DB] GetAllPortfolios")
	results, err := m.collection.FindAll(portfoliosCollection, bson.M{})
	if err != nil {
		return nil, err
	}
	portfolioTemp := []wallet.Portfolio{}
	for _, result := range results {
		bsonBytes, _ := bson.Marshal(result)
		portfolio := wallet.Portfolio{}
		bson.Unmarshal(bsonBytes, &portfolio)
		portfolioTemp = append(portfolioTemp, portfolio)
	}
	return portfolioTemp, nil
}

func (m *mongoSession) GetPortfolioByID(id string) (*wallet.Portfolio, error) {
	log.Debug("[DB] GetPortfolioByID")
	h := &wallet.Portfolio{}
	query := bson.M{"_id": id}
	if err := m.collection.FindOne(portfoliosCollection, query, h); err != nil {
		return nil, err
	}
	if h.ID == "" {
		return nil, nil
	}
	return h, nil
}

// FIXME: maybe unnecessary
func contains(s []interface{}, e interface{}) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (m *mongoSession) getPortfolioItem(itemType string, item wallet.Tradable, year int) (map[string]wallet.PortfolioItem, error) {
	log.Debugf("[DB] Getting portfolio item %s", itemType)
	operationsSymbols, err := m.getOperationsSymbols(bson.M{"itemType": itemType})
	if err != nil {
		return nil, err
	}

	items := map[string]wallet.PortfolioItem{}
	for _, s := range operationsSymbols {
		symbol := s.(string)
		portfolioItem := &wallet.PortfolioItem{}
		// FIXME: one request
		if err := financeapi.GetJSON("/"+itemType+"/"+symbol, portfolioItem); err != nil {
			log.Errorf("Error on get stock item: %s", err)
		}

		operations, err := m.getAllOperationsBySymbol(symbol, year)
		if err != nil {
			return nil, err
		}

		typeOfitem := reflect.TypeOf(item)
		operationsList := wallet.OperationsList{}
		for _, result := range operations {
			newItem := reflect.New(typeOfitem).Interface().(wallet.Tradable)
			bsonBytes, _ := bson.Marshal(result)
			bson.Unmarshal(bsonBytes, newItem)
			operationsList = append(operationsList, newItem)
		}

		// FIXME
		broker := ""
		if len(operationsList) > 0 {
			operation := operationsList[0]
			if operation != nil {
				broker = operation.GetBrokerID()
			}
		}

		portfolioItem.BrokerID = broker
		portfolioItem.ItemType = itemType
		portfolioItem.Operations = operationsList
		portfolioItem.Recalculate()
		items[symbol] = *portfolioItem
	}

	return items, nil
}

// FIXME
func (m *mongoSession) GetPortfolioItems(portfolio *wallet.Portfolio, year int) error {
	log.Debug("[DB] GetPortfolioItems")

	// FIXME
	slugs := map[string]wallet.Tradable{
		"certificates-of-deposit": wallet.CertificateOfDeposit{},
		"ficfi":                   wallet.FICFI{},
		"fiis":                    wallet.FII{},
		"stocks":                  wallet.Stock{},
		"stocks-funds":            wallet.StockFund{},
		"treasuries-direct":       wallet.TreasuryDirect{},
	}

	portfolio.Items = map[string]wallet.PortfolioItem{}
	for slug, item := range slugs {
		stocks, err := m.getPortfolioItem(slug, item, year)
		if err != nil {
			log.Errorf("[DB] Error on get portfolio items: %v", err)
			continue
		}
		for symbol, portfolioItem := range stocks {
			portfolio.Items[symbol] = portfolioItem
		}
	}

	portfolio.Recalculate()

	return nil
}
