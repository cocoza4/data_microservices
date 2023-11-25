package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cocoza4/data_microservices/models"
	"github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KafkaServiceImpl struct {
	conn   *kafka.Conn
	writer *kafka.Writer
}

func NewKafkaService(conn *kafka.Conn, writer *kafka.Writer) KafkaService {
	return &KafkaServiceImpl{
		conn:   conn,
		writer: writer,
	}
}

func (k *KafkaServiceImpl) GetTopics() ([]string, error) {
	partitions, err := k.conn.ReadPartitions()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var topics []string
	for _, p := range partitions {
		topics = append(topics, p.Topic)
	}
	return topics, nil
}

func (k *KafkaServiceImpl) CreateTopic(name *string) error {
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             *name,
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	err := k.conn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Printf("Topic `%s` created\n", *name)
	return nil
}

func (k *KafkaServiceImpl) PublishMessage(topic *string, product *models.Product) error {
	product.Id = primitive.NewObjectID()
	product.CreatedDate = time.Now()
	raw, _ := json.Marshal(product)

	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: *topic,
			Value: raw,
		},
	)
	if err != nil {
		log.Fatal("failed to write message:", err)
		return err
	}

	return nil
}

func getTopicsOffsetsForAllPartitions(topics []string) (*kafka.ListOffsetsResponse, error) {
	ctx := context.Background()
	client := kafka.Client{Addr: kafka.TCP(os.Getenv("KAFKA_URI"))}

	metadataInput := &kafka.MetadataRequest{Topics: topics}
	metadataOutput, err := client.Metadata(ctx, metadataInput)
	if err != nil {
		return nil, fmt.Errorf("error getting metadata: %d", err)
	}

	offsetsInput := kafka.ListOffsetsRequest{Topics: make(map[string][]kafka.OffsetRequest)}
	for _, topic := range metadataOutput.Topics {
		if topic.Error != nil {
			return nil, fmt.Errorf("error getting topic %s metadata: %d", topic.Name, err)
		}

		var topicPartitions []kafka.OffsetRequest
		for _, partition := range topic.Partitions {
			topicPartitions = append(topicPartitions, kafka.OffsetRequest{
				Partition: partition.ID,
				Timestamp: kafka.LastOffset,
			})
		}
		offsetsInput.Topics[topic.Name] = topicPartitions
	}

	return client.ListOffsets(ctx, &offsetsInput)
}

func getLatestByPartition(ch chan models.Product, topic *string, partition kafka.PartitionOffsets, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{os.Getenv("KAFKA_URI")},
		Topic:     *topic,
		Partition: partition.Partition,
		MaxBytes:  10e6, // 10MB
		MaxWait:   time.Second,
	})
	defer reader.Close()

	product := models.Product{}
	log.Println("partition.LastOffset", partition.LastOffset)
	if partition.LastOffset > 0 {
		reader.SetOffset(partition.LastOffset - 1)
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Println(err.Error())
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		json.Unmarshal(m.Value, &product)
		fmt.Println("xxxxxxx", product)
	}

	ch <- product
}

func (k *KafkaServiceImpl) GetLatest(topic *string) (*models.Product, error) {

	resp, err := getTopicsOffsetsForAllPartitions([]string{*topic})
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	numPartitions := len(resp.Topics[*topic])
	log.Println("num partitions", numPartitions)
	ch := make(chan models.Product, numPartitions)
	var wg sync.WaitGroup
	for _, partition := range resp.Topics[*topic] {
		if partition.Error != nil {
			log.Println(partition.Error)
		}

		fmt.Printf("partition: %d, last offset(%d)\n", partition.Partition, partition.LastOffset)

		go getLatestByPartition(ch, topic, partition, &wg)
		wg.Add(1)
	}

	wg.Wait()
	close(ch)

	var products []models.Product
	for product := range ch {
		if len(product.Name) > 0 {
			products = append(products, product)
		}
	}

	if len(products) == 0 {
		return nil, errors.New("No messages found")
	}

	latest := products[0]
	for i := 1; i < len(products); i++ {
		product := products[i]
		if product.CreatedDate.After(latest.CreatedDate) {
			latest = product
		}
	}

	return &latest, nil
}
