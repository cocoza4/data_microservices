package services

import "github.com/cocoza4/data_microservices/models"

type ProductService interface {
	CreateProduct(*models.Product) error
	GetProduct(*string) (*models.Product, error)
	GetAll() ([]*models.Product, error)
	GetIndexes() []*string
	CreateIndex(*string) error
	GetLatest() (*models.Product, error)
}
