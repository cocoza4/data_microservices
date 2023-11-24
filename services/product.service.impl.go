package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cocoza4/data_microservices/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewProductService(collection *mongo.Collection, ctx context.Context) ProductService {
	return &ProductServiceImpl{
		collection: collection,
		ctx:        ctx,
	}
}

func (p *ProductServiceImpl) CreateProduct(product *models.Product) error {
	product.Id = primitive.NewObjectID()
	product.CreatedDate = time.Now()
	_, err := p.collection.InsertOne(p.ctx, product)
	return err
}

func (p *ProductServiceImpl) GetProduct(name *string) (*models.Product, error) {
	var product *models.Product
	query := bson.D{bson.E{Key: "name", Value: name}}
	err := p.collection.FindOne(p.ctx, query).Decode(&product)
	return product, err
}

func (p *ProductServiceImpl) GetLatest() (*models.Product, error) {
	var product *models.Product
	opts := options.FindOne().SetSort(bson.M{"$natural": -1})
	err := p.collection.FindOne(p.ctx, bson.M{}, opts).Decode(&product)
	return product, err
}

func (p *ProductServiceImpl) GetAll() ([]*models.Product, error) {
	var products []*models.Product
	cursor, err := p.collection.Find(p.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(p.ctx) {
		var product models.Product
		err := cursor.Decode(&product)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(p.ctx)
	return products, nil
}

func (p *ProductServiceImpl) GetIndexes() []*string {
	indexView := p.collection.Indexes()
	opts := options.ListIndexes()
	cursor, err := indexView.List(p.ctx, opts)

	if err != nil {
		log.Fatal(err)
	}

	var result []bson.M
	if err = cursor.All(p.ctx, &result); err != nil {
		log.Fatal(err)
	}

	var indexes []*string
	for _, v := range result {
		name := fmt.Sprintf("%v", v["name"])
		indexes = append(indexes, &name)
	}
	return indexes
}

func (p *ProductServiceImpl) CreateIndex(field *string) error {
	name, err := p.collection.Indexes().CreateOne(p.ctx, mongo.IndexModel{
		Keys: bson.M{*field: 1},
	})
	if err != nil {
		return err
	}
	log.Println("Index created:", name)
	return err
}
