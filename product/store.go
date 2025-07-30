package main

import (
	"context"
	"fmt"

	"github.com/AyanokojiKiyotaka8/E-Commerce/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductStorer interface {
	GetProduct(context.Context, map[string]interface{}) (*types.ProductData, error)
	GetProducts(context.Context, map[string]interface{}) ([]*types.ProductData, error)
	PostProduct(context.Context, *types.ProductData) (*types.ProductData, error)
	PutProduct(context.Context, map[string]interface{}, map[string]interface{}) error
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

func (s *MongoProductStore) GetProduct(ctx context.Context, filter map[string]interface{}) (*types.ProductData, error) {
	var product types.ProductData
	dbFilter := bson.M{}
	for key, val := range filter {
		if key == "id" {
			sid, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("invalid id")
			}

			oid, err := primitive.ObjectIDFromHex(sid)
			if err != nil {
				return nil, fmt.Errorf("invalid id")
			}
			dbFilter["_id"] = oid
		} else {
			dbFilter[key] = val
		}
	}

	if err := s.coll.FindOne(ctx, dbFilter).Decode(&product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *MongoProductStore) GetProducts(ctx context.Context, filter map[string]interface{}) ([]*types.ProductData, error) {
	var products []*types.ProductData
	dbFilter := bson.M{}
	for key, val := range filter {
		dbFilter[key] = val
	}
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := curr.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (s *MongoProductStore) PostProduct(ctx context.Context, data *types.ProductData) (*types.ProductData, error) {
	res, err := s.coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}

	data.Id = res.InsertedID.(primitive.ObjectID).Hex()
	return data, nil
}

func (s *MongoProductStore) PutProduct(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) error {
	dbFilter := bson.M{}
	for key, val := range filter {
		if key == "id" {
			sid, ok := val.(string)
			if !ok {
				return fmt.Errorf("invalid id")
			}

			oid, err := primitive.ObjectIDFromHex(sid)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			dbFilter["_id"] = oid
		} else {
			dbFilter[key] = val
		}
	}

	updates := bson.M{}
	for key, val := range update {
		updates[key] = val
	}
	dbUpdate := bson.M{
		"$set": updates,
	}

	_, err := s.coll.UpdateOne(ctx, dbFilter, dbUpdate)
	return err
}

func (s *MongoProductStore) DeleteProduct(ctx context.Context, filter map[string]interface{}) error {
	dbFilter := bson.M{}
	for key, val := range filter {
		if key == "id" {
			sid, ok := val.(string)
			if !ok {
				return fmt.Errorf("invalid id")
			}

			oid, err := primitive.ObjectIDFromHex(sid)
			if err != nil {
				return fmt.Errorf("invalid id")
			}
			dbFilter["_id"] = oid
		} else {
			dbFilter[key] = val
		}
	}
	_, err := s.coll.DeleteOne(ctx, dbFilter)
	return err
}
