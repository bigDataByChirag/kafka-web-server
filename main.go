package main

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"sync"
)

// Define Kafka topic name and Redis address
const (
	KafkaTopic = "data"
	RedisAddr  = "localhost:6379"
)

// Create a Redis client
var redisClient *redis.Client
var ctx = context.Background()

// Create a Sarama consumer group
var (
	brokers       = []string{"localhost:9092"}
	kafkaProducer sarama.SyncProducer
	kafkaConsumer sarama.Consumer
	//userData      = make(map[string]string)
	//mu            sync.Mutex
)

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr: RedisAddr,
	})

	// Initialize Kafka producer
	initKafkaProducer()

	// Initialize Kafka consumer
	initKafkaConsumer()

	// Start background goroutine to consume messages from Kafka and store in Redis
	go consumeMessages()

	// Create Gin router
	r := gin.Default()

	// Define POST and GET endpoints
	r.POST("/api/post", handlePostRequest)
	r.GET("/api/get/:userId", handleGetRequest)

	// Start the server
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Initialize Kafka producer
func initKafkaProducer() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}

	kafkaProducer = producer
}

// Initialize Kafka consumer
func initKafkaConsumer() {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatal("Failed to create Kafka consumer:", err)
	}

	kafkaConsumer = consumer
}

// Handle POST request and publish to Kafka
func handlePostRequest(c *gin.Context) {
	var requestData map[string]interface{}

	// Parse request body as JSON
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID and data from request
	userId, exists := requestData["userId"].(string)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId not found"})
		return
	}
	data, exists := requestData["data"].(string)
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "data not found"})
		return
	}
	var wg = new(sync.WaitGroup)
	// Use a goroutine to publish data to Kafka
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Publish data to Kafka
		message := &sarama.ProducerMessage{
			Topic: KafkaTopic,
			Key:   sarama.StringEncoder(userId),
			Value: sarama.StringEncoder(data),
		}

		partition, offset, err := kafkaProducer.SendMessage(message)
		if err != nil {
			log.Println("Failed to publish data to Kafka:", err)
		} else {
			log.Printf("Data published to Kafka - Partition: %d, Offset: %d\n", partition, offset)
		}
	}()
	wg.Wait()
	// Respond immediately with success
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// Handle GET request and retrieve data from Redis
func handleGetRequest(c *gin.Context) {
	userId := c.Param("userId")

	ch := make(chan string)
	// Retrieve data from Redis using the provided user ID
	go func() {

		data, err := redisClient.Get(ctx, userId).Result()
		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user ID or data not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ch <- data
	}()

	data := <-ch
	// Respond with the data
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func consumeMessages() {
	partitionConsumer, err := kafkaConsumer.ConsumePartition(KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatal("Failed to start consuming messages:", err)
	}
	defer partitionConsumer.Close()

	for message := range partitionConsumer.Messages() {
		// Extract user ID and data from message
		userId := string(message.Key)
		data := string(message.Value)

		// Store data in Redis under userId key
		err = redisClient.Set(ctx, userId, data, 0).Err()
		if err != nil {
			log.Println("Failed to store data in Redis:", err)
		} else {
			log.Println("Stored data in Redis under userId:", userId)
		}
	}
}
