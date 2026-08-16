package publisher

import (
	"fmt"
	"log"

	"github.com/temuka-api-service/internal/constant"
	"github.com/temuka-api-service/util/queue"
)

type SuggestionEvent struct {
	Operation string                 `json:"operation"`
	Type      string                 `json:"type"`
	EntityID  string                 `json:"entity_id"`
	Data      map[string]interface{} `json:"data"`
}

type SuggestionPublisher interface {
	PublishSuggestionEvent(op, entityType, entityID string, data map[string]interface{}) error
}

type suggestionPublisherImpl struct {
	rmq queue.RabbitMQChannel
}

func NewSearchIndexPublisher(rmq queue.RabbitMQChannel) SuggestionPublisher {
	if err := rmq.RegisterExchange(constant.SuggestionExchange, "direct", true, false); err != nil {
		log.Printf("failed to declare exchange '%s': %v", constant.SuggestionExchange, err)
	}

	if _, err := rmq.InitQueue(constant.SuggestionExchange, constant.SuggestionSyncRoutingKey, true, false); err != nil {
		log.Printf("failed to initialize queue with routing key '%s': %v", constant.SuggestionSyncRoutingKey, err)
	}

	log.Printf("initialized successfully with exchange '%s' and key '%s'", constant.SuggestionExchange, constant.SuggestionSyncRoutingKey)

	return &suggestionPublisherImpl{rmq: rmq}
}

func (p *suggestionPublisherImpl) PublishSuggestionEvent(op, entityType, entityID string, data map[string]interface{}) error {
	event := SuggestionEvent{
		Operation: op,
		Type:      entityType,
		EntityID:  entityID,
		Data:      data,
	}

	log.Printf("publishing event op='%s', type='%s', entityID='%s'", op, entityType, entityID)

	if err := p.rmq.PublishMessage(constant.SuggestionExchange, constant.SuggestionSyncRoutingKey, event); err != nil {
		log.Fatalf("failed to publish event op='%s', entityID='%s': %v", op, entityID, err)
		return fmt.Errorf("failed to publish suggestion event: %w", err)
	}

	log.Printf("successfully published event for entityID='%s'", entityID)

	return nil
}
