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

func (c *MessageConsumer) ConsumerMessages(ctx context.Context, queueNames []string) {
	for _, queueName := range queueNames {
		switch queueName {
		case "Test":
			go c.handleCustomType1Logic(ctx, queueName)
		case "Test2":
			// handleCustomType2Logic(c)
		default:
			log.Printf("[%s] Unsupported message type: %+v\n", queueName, queueName)
		}
	}
}

func (c *MessageConsumer) handleCustomType1Logic(ctx context.Context, queueName string) {
	consumerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("[%s] Consumer started listening...\n", queueName)
	c.subscription = c.redisClient.RedisClient.Subscribe(queueName)
	defer c.subscription.Close()

	channel := c.subscription.Channel()

	for {
		select {
		case <-consumerCtx.Done():
			log.Printf("[%s] Consumer stopped listening...\n", queueName)
			return
		case msg := <-channel:
			var messageObj interface{}
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
