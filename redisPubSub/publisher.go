package redisPubSub

import (
	"boilerplate-api/infrastructure"
	"context"
	"encoding/json"
	"log"
)

// MessagePublisher is a generic publisher for different message types.
type MessagePublisher struct {
	redisClient infrastructure.Redis
}

func NewMessagePublisher(redisClient infrastructure.Redis) *MessagePublisher {
	return &MessagePublisher{redisClient}
}

// PublishMessages publishes messages to  channels.
func (p *MessagePublisher) PublishMessages(ctx context.Context, message interface{}, queueName string) {

	serializedMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("[%s] Failed to serialize message: %v", queueName, err)
	}

	// Use the context for the publishing operation
	err = p.redisClient.RedisClient.Publish(queueName, serializedMessage).Err()
	if err != nil {
		log.Printf("[%s] Failed to publish message: %v", queueName, err)
	}

}
