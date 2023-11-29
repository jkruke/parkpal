package kafka

import (
	"context"
	"encoding/json"
	"parkpal-web-server/internal/business"
	"parkpal-web-server/internal/entity"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	business business.Business
}

func NewKafkaConsumer(brokers []string, groupID string, topics []string, business business.Business) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: c, business: business}, nil
}

func (kc *KafkaConsumer) ConsumeMessages(ctx context.Context) {
	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			// do something
		}

		var prod entity.ParkingLot
		err = json.Unmarshal(msg.Value, &prod)

		// save msg to memStore
		kc.business.UpdateParkingLot(ctx, (*business.UpdateParkingLotRequest)(&prod))
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
