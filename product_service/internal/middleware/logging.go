package middleware

import (
	"context"
	"time"

	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/model"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/service"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next service.ProductServicer
}

func NewLogMiddleware(next service.ProductServicer) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) GetProduct(ctx context.Context, filter map[string]interface{}) (product *model.Product, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"product": product,
			"error":   err,
			"took":    time.Since(start),
		}).Info("Getting Product")
	}(time.Now())
	product, err = l.next.GetProduct(ctx, filter)
	return
}

func (l *LogMiddleware) GetProducts(ctx context.Context, filter map[string]interface{}) (products []*model.Product, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"count": len(products),
			"error": err,
			"took":  time.Since(start),
		}).Info("Getting Products")
	}(time.Now())
	products, err = l.next.GetProducts(ctx, filter)
	return
}

func (l *LogMiddleware) CreateProduct(ctx context.Context, product *model.Product) (p *model.Product, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"product": p,
			"error":   err,
			"took":    time.Since(start),
		}).Info("Creating Product")
	}(time.Now())
	p, err = l.next.CreateProduct(ctx, product)
	return
}

func (l *LogMiddleware) UpdateProduct(ctx context.Context, filter map[string]interface{}, update map[string]interface{}) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"took":  time.Since(start),
		}).Info("Updating Product")
	}(time.Now())
	err = l.next.UpdateProduct(ctx, filter, update)
	return
}

func (l *LogMiddleware) DeleteProduct(ctx context.Context, filter map[string]interface{}) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"error": err,
			"took":  time.Since(start),
		}).Info("Deleting Product")
	}(time.Now())
	err = l.next.DeleteProduct(ctx, filter)
	return
}
