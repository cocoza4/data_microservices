package main

import (
	"context"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/cocoza4/data_microservices/controllers"
	"github.com/cocoza4/data_microservices/middleware"
	"github.com/cocoza4/data_microservices/services"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server         *gin.Engine
	productService services.ProductService
	kafkaService   services.KafkaService
	authCtrl       controllers.AuthController
	productCtrl    controllers.ProductController
	kafkaCtrl      controllers.KafkaController
	ctx            context.Context
	collection     *mongo.Collection
	client         *mongo.Client

	err error
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

func getMongoCollection() *mongo.Collection {
	mongo_uri := getEnv("MONGO_URI", "mongodb://localhost:27017")
	db := getEnv("MONGO_DBNAME", "productdb")
	log.Println("Mongo URI:", mongo_uri)
	log.Println("Mongo DB:", db)

	mongo_conn := options.Client().ApplyURI(mongo_uri)
	client, err = mongo.Connect(ctx, mongo_conn)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mongo connection established")

	collection = client.Database(db).Collection("products")
	return collection
}

func getKafkaConnAndWriter() (*kafka.Conn, *kafka.Writer) {
	uri := getEnv("KAFKA_URI", "localhost:9094")
	log.Println("Kafka URI:", uri)

	conn, err := kafka.Dial("tcp", uri)
	if err != nil {
		log.Fatal(err.Error())
	}
	controller, err := conn.Controller()
	if err != nil {
		log.Fatal(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatal(err.Error())
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(uri),
		Balancer: &kafka.LeastBytes{},
	}

	return controllerConn, writer
}

func init() {
	ctx = context.TODO()

	secret_username = getEnv("SECRET_USER", "")
	secret_password = getEnv("SECRET_PASSWORD", "")
	secret_key = getEnv("SECRET_KEY", "")

	if secret_username == "" || secret_password == "" || secret_key == "" {
		log.Fatal("`SECRET_USER`, `SECRET_PASSWORD` and `SECRET_KEY` are required")
	}

	// kafka setup
	conn, writer := getKafkaConnAndWriter()
	kafkaService = services.NewKafkaService(conn, writer)
	kafkaCtrl = controllers.NewKafkaController(kafkaService)

	// mongodb setup
	collection := getMongoCollection()
	productService = services.NewProductService(collection, ctx)
	productCtrl = controllers.ProductController{ProductService: productService}

	authCtrl = controllers.NewAuthController()

	server = gin.Default()
}

func main() {
	defer client.Disconnect(ctx)

	version := "/v1"
	basepath := server.Group(version)
	authCtrl.RegisterRoutes(basepath)

	productPath := server.Group(version, middleware.Authorize)
	productCtrl.RegisterRoutes(productPath)

	kafkaPath := server.Group(version, middleware.Authorize)
	kafkaCtrl.RegisterRoutes(kafkaPath)

	log.Fatal(server.Run(":8080"))
}
