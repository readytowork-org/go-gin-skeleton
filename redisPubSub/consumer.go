package redisPubSub

import (
	"boilerplate-api/infrastructure"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

// Consumer is a generic consumer for different message type
type MessageConsumer struct {
	redisClient  infrastructure.Redis
	subscription *redis.PubSub
}

// NewMessageConsumer creates a new instance of MessageConsumer.
func NewMessageConsumer(redis infrastructure.Redis) *MessageConsumer {
	return &MessageConsumer{
		redisClient: redis,
	}
}

// This Function takes queue names in an array and uses a switch statement to perform required logic for the queues
func (c *MessageConsumer) ConsumerMessages(ctx context.Context, queueNames []string) {
	for _, queueName := range queueNames {
		switch queueName {
		case "Test":
			// We will handle the go routines in the custom function
			go c.handleCustomType1Logic(ctx, queueName)
		case "Test2":
			// handleCustomType2Logic(c)
		default:
			log.Printf("[%s] Unsupported message type: %+v\n", queueName, queueName)
		}
	}
}

// handleCustomType1Logic initiates a goroutine to handle messages from the specified queue.
func (c *MessageConsumer) handleCustomType1Logic(ctx context.Context, queueName string) {

	// Create a cancellation context to gracefully stop the goroutine
	consumerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("[%s] Consumer started listening...\n", queueName)

	// Subscribe to the specified Redis channel
	c.subscription = c.redisClient.RedisClient.Subscribe(queueName)
	defer c.subscription.Close()

	// Obtain the channel for receiving messages
	channel := c.subscription.Channel()

	for {
		select {
		// Check if the main context is canceled to stop the goroutine
		case <-consumerCtx.Done():
			log.Printf("[%s] Consumer stopped listening...\n", queueName)
			return
			// Listen for incoming messages on the channel
		case msg := <-channel:
			var messageObj interface{}
			// Deserialize the message payload
			err := json.Unmarshal([]byte(msg.Payload), &messageObj)
			if err != nil {
				log.Printf("[%s] Failed to deserialize message: %v", queueName, err)
				continue
			}

			// Continue with your logic here:

			fmt.Printf("[%s] Received message: %+v\n", queueName, messageObj)
		}
	}
}
