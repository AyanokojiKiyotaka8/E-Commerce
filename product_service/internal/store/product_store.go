package store

import (
	"context"
	"fmt"

	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductStorer interface {
	GetProduct(context.Context, map[string]interface{}) (*model.Product, error)
	GetProducts(context.Context, map[string]interface{}) ([]*model.Product, error)
	CreateProduct(context.Context, *model.Product) (*model.Product, error)
	UpdateProduct(context.Context, map[string]interface{}, map[string]interface{}) error
	DeleteProduct(context.Context, map[string]interface{}) error
}

type MongoProductStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoProductStore(client *mongo.Client) *MongoProductStore {
	return &MongoProductStore{
		client: client,
		coll:   client.Database("e-commerce").Collection("products"),
	}
}

func (s *MongoProductStore) GetProduct(ctx context.Context, filter map[string]interface{}) (*model.Product, error) {
	dbFilter, err := parseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	var product model.Product
	if err := s.coll.FindOne(ctx, dbFilter).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *MongoProductStore) GetProducts(ctx context.Context, filter map[string]interface{}) ([]*model.Product, error) {
	dbFilter, err := parseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	var products []*model.Product
	curr, err := s.coll.Find(ctx, dbFilter)
	if err != nil {
		return nil, err
	}

	if err := curr.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (s *MongoProductStore) CreateProduct(ctx context.Context, product *model.Product) (*model.Product, error) {
	res, err := s.coll.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}

	product.Id = res.InsertedID.(primitive.ObjectID).Hex()
	return product, nil
}

func (s *MongoProductStore) UpdateProduct(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) error {
	dbFilter, err := parseMongoFilter(filter)
	if err != nil {
		return err
	}

	dbUpdate := bson.M{"$set": update}
	_, err = s.coll.UpdateOne(ctx, dbFilter, dbUpdate)
	return err
}

func (s *MongoProductStore) DeleteProduct(ctx context.Context, filter map[string]interface{}) error {
	dbFilter, err := parseMongoFilter(filter)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, dbFilter)
	return err
}

func parseMongoFilter(filter map[string]interface{}) (bson.M, error) {
	dbFilter := bson.M{}
	for key, val := range filter {
		if key == "id" {
			sid, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("invalid id type")
			}
			oid, err := primitive.ObjectIDFromHex(sid)
			if err != nil {
				return nil, fmt.Errorf("invalid ObjectID")
			}
			dbFilter["_id"] = oid
		} else {
			dbFilter[key] = val
		}
	}
	return dbFilter, nil
}
