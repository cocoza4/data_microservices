package services

import (
	"context"

	"github.com/cocoza4/data_microservices/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductServiceImpl struct {
	productcollection *mongo.Collection
	ctx               context.Context
}

func NewProductService(productcollection *mongo.Collection, ctx context.Context) ProductService {
	return &ProductServiceImpl{
		productcollection: productcollection,
		ctx:               ctx,
	}
}

func (p *ProductServiceImpl) CreateProduct(product *models.Product) error {
	return nil
}

func (p *ProductServiceImpl) GetProduct(name *string) (*models.Product, error) {
	return nil, nil
}

func (p *ProductServiceImpl) GetAll() ([]*models.Product, error) {
	return nil, nil
}
