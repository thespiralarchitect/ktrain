package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/cmd/repository"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type RabbitMqManager struct {
	*amqp.Connection
	activityLogRepository repository.ActivityLogRepository
}

func ConectRambbitMQ(activityLogRepository repository.ActivityLogRepository) (*RabbitMqManager, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp:%s", viper.GetString("rabbitmq.host")))
	if err != nil {
		return nil, err
	}
	return &RabbitMqManager{
		Connection:            conn,
		activityLogRepository: activityLogRepository,
	}, nil
}
func (m *RabbitMqManager) Consumers(ctx context.Context) error {
	ch, err := m.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.ExchangeDeclare(
		"activityLog", // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}
	q, err := ch.QueueDeclare(
		"activityLog", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		return err
	}
	err = ch.QueueBind(
		q.Name,
		"",
		"activityLog",
		false,
		nil,
	)
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	err = m.CreateAction(ctx, msgs)
	if err != nil {
		return err
	}
	return nil
}
func (m *RabbitMqManager) CreateAction(ctx context.Context, msgs <-chan amqp.Delivery) error {
	forever := make(chan bool)
	for d := range msgs {
		log := dto.UserActivityLogMessage{}
		err := json.Unmarshal(d.Body, &log)
		if err != nil {
			return err
		}
		_, err = m.activityLogRepository.CreateAction(ctx, log.ID, log.Log)
		if err != nil {
			return err
		}
	}
	<-forever
	return nil
}
func (m *RabbitMqManager) Close() {
	err := m.Connection.Close()
	if err != nil {
		log.Fatalf("Could not close server, err: %v", err)
	}
}
