package main

import (
	"context"
	"log"
	"os"

	"github.com/cocoza4/data_microservices/controllers"
	"github.com/cocoza4/data_microservices/middleware"
	"github.com/cocoza4/data_microservices/services"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	service     services.ProductService
	authCtrl    controllers.AuthController
	productCtrl controllers.ProductController
	ctx         context.Context
	collection  *mongo.Collection
	client      *mongo.Client
	err         error
)

// mock basic auth user
var (
	secret_username string
	secret_password string
	secret_key      string
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

	secret_username = getEnv("SECRET_USER", "")
	secret_password = getEnv("SECRET_PASSWORD", "")
	secret_key = getEnv("SECRET_KEY", "")

	if secret_username == "" || secret_password == "" || secret_key == "" {
		log.Fatal("`SECRET_USER`, `SECRET_PASSWORD` and `SECRET_KEY` are required")
	}

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
	productCtrl = controllers.ProductController{ProductService: service}
	authCtrl = controllers.AuthController{}
	server = gin.Default()
}

func main() {
	defer client.Disconnect(ctx)

	version := "/v1"
	basepath := server.Group(version)
	authCtrl.RegisterRoutes(basepath)

	productPath := server.Group(version, middleware.Authorize)
	productCtrl.RegisterRoutes(productPath)
	log.Fatal(server.Run(":8080"))
}
