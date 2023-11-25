package services

import "github.com/cocoza4/data_microservices/models"

type KafkaService interface {
	GetTopics() ([]string, error)
	CreateTopic(*string) error
	PublishMessage(*string, *models.Product) error
	GetLatest(*string) (*models.Product, error)
}
