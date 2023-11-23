package main

import (
	"context"
	"log"
	"os"

	"github.com/cocoza4/data_microservices/controllers"
	"github.com/cocoza4/data_microservices/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server     *gin.Engine
	service    services.ProductService
	ctr        controllers.ProductController
	ctx        context.Context
	collection *mongo.Collection
	client     *mongo.Client
	err        error
)

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		value = fallback
	}
	return value
}

func init() {
	ctx = context.TODO()

	mongo_uri := getEnv("MONGO_URI", "mongodb://localhost:27017")
	db := getEnv("MONGO_DBNAME", "productdb")
	log.Println("Mongo URI:", mongo_uri)
	log.Println("Mongo DB:", db)
	conn := options.Client().ApplyURI(mongo_uri)
	client, err = mongo.Connect(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mongo connection established")

	collection = client.Database(db).Collection("products")
	service = services.NewProductService(collection, ctx)
	ctr = controllers.ProductController{ProductService: service}
	server = gin.Default()
}

func main() {
	defer client.Disconnect(ctx)

	basepath := server.Group("/v1")
	ctr.RegisterRoutes(basepath)
	log.Fatal(server.Run(":8080"))
}
