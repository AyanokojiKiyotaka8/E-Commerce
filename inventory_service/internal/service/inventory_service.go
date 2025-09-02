package service

import (
	"context"

	"github.com/AyanokojiKiyotaka8/E-Commerce/inventory_service/internal/model"
	"github.com/AyanokojiKiyotaka8/E-Commerce/inventory_service/internal/store"
)

type InventoryServicer interface {
	GetInventory(context.Context, map[string]interface{}) (*model.Inventory, error)
	GetInventories(context.Context, map[string]interface{}) ([]*model.Inventory, error)
	CreateInventory(context.Context, *model.Inventory) (*model.Inventory, error)
	UpdateInventory(context.Context, map[string]interface{}, map[string]interface{}) error
}

type InventoryService struct {
	store store.InventoryStorer
}

func NewInventoryService(store store.InventoryStorer) *InventoryService {
	return &InventoryService{
		store: store,
	}
}

func (s *InventoryService) GetInventory(ctx context.Context, filter map[string]interface{}) (*model.Inventory, error) {
	return s.store.GetInventory(ctx, filter)
}

func (s *InventoryService) GetInventories(ctx context.Context, filter map[string]interface{}) ([]*model.Inventory, error) {
	return s.store.GetInventories(ctx, filter)
}

func (s *InventoryService) CreateInventory(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	return s.store.CreateInventory(ctx, inventory)
}

func (s *InventoryService) UpdateInventory(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) error {
	return s.store.UpdateInventory(ctx, filter, update)
}
