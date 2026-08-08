package publisher

import (
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
	_ = rmq.RegisterExchange(constant.SuggestionExchange, "direct", true, false)
	_, _ = rmq.InitQueue(constant.SuggestionExchange, constant.SuggestionSyncRoutingKey, true, false)

	return &suggestionPublisherImpl{rmq: rmq}
}

func (p *suggestionPublisherImpl) PublishSuggestionEvent(op, entityType, entityID string, data map[string]interface{}) error {
	event := SuggestionEvent{
		Operation: op,
		Type:      entityType,
		EntityID:  entityID,
		Data:      data,
	}

	return p.rmq.PublishMessage(constant.SuggestionExchange, constant.SuggestionSyncRoutingKey, event)
}
