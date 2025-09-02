package store

import (
	"context"
	"fmt"

	"github.com/AyanokojiKiyotaka8/E-Commerce/inventory_service/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InventoryStorer interface {
	GetInventory(context.Context, map[string]interface{}) (*model.Inventory, error)
	GetInventories(context.Context, map[string]interface{}) ([]*model.Inventory, error)
	CreateInventory(context.Context, *model.Inventory) (*model.Inventory, error)
	UpdateInventory(context.Context, map[string]interface{}, map[string]interface{}) error
}

type MongoInventoryStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoInventoryStore(client *mongo.Client) *MongoInventoryStore {
	return &MongoInventoryStore{
		client: client,
		coll:   client.Database("e-commerce").Collection("inventories"),
	}
}

func (s *MongoInventoryStore) GetInventory(ctx context.Context, filter map[string]interface{}) (*model.Inventory, error) {
	dbFilter, err := parseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	var inventory model.Inventory
	if err := s.coll.FindOne(ctx, dbFilter).Decode(&inventory); err != nil {
		return nil, err
	}
	return &inventory, nil
}

func (s *MongoInventoryStore) GetInventories(ctx context.Context, filter map[string]interface{}) ([]*model.Inventory, error) {
	dbFilter, err := parseMongoFilter(filter)
	if err != nil {
		return nil, err
	}

	opts, err := getPaginationOptions(filter)
	if err != nil {
		return nil, err
	}

	cur, err := s.coll.Find(ctx, dbFilter, opts)
	if err != nil {
		return nil, err
	}

	var inventories []*model.Inventory
	if err := cur.All(ctx, &inventories); err != nil {
		return nil, err
	}
	return inventories, err
}

func (s *MongoInventoryStore) CreateInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	res, err := s.coll.InsertOne(ctx, inventory)
	if err != nil {
		return nil, err
	}

	inventory.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return inventory, nil
}

func (s *MongoInventoryStore) UpdateInventory(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) error {
	dbFilter, err := parseMongoFilter(filter)
	if err != nil {
		return err
	}

	dbUpdate := bson.M{"$set": update}
	_, err = s.coll.UpdateOne(ctx, dbFilter, dbUpdate)
	return err
}

func parseMongoFilter(filter map[string]interface{}) (bson.M, error) {
	dbFilter := bson.M{}

	for key, val := range filter {
		switch key {
		case "id":
			sid, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("invalid id type")
			}
			oid, err := primitive.ObjectIDFromHex(sid)
			if err != nil {
				return nil, fmt.Errorf("invalid ObjectID")
			}
			dbFilter["_id"] = oid

		case "page", "limit":
			// Skip pagination keys

		default:
			dbFilter[key] = val
		}
	}
	return dbFilter, nil
}

func getPaginationOptions(filter map[string]interface{}) (*options.FindOptions, error) {
	opts := &options.FindOptions{}

	limit, ok := filter["limit"].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid limit")
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}

	page, ok := filter["page"].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid page")
	}
	if page > 0 {
		opts.SetSkip((page - 1) * limit)
	}
	return opts, nil
}
