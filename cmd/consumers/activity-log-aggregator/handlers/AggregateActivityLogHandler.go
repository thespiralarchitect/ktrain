package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/pkg/logger"
	"ktrain/proto/pb"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type RabbitMqManager struct {
	*amqp.Connection
}

func ConectRambbitMQ() (*RabbitMqManager, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp:%s", viper.GetString("rabbitmq.host")))
	if err != nil {
		return nil, err
	}
	return &RabbitMqManager{
		Connection: conn,
	}, nil
}
func (m *RabbitMqManager) Consumers(ctx context.Context, activityLogClient pb.ActivityLogDMSServiceClient) error {
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
	err = m.CreateAction(ctx, msgs, activityLogClient)
	if err != nil {
		return err
	}
	return nil
}
func (m *RabbitMqManager) CreateAction(ctx context.Context, msgs <-chan amqp.Delivery, activityLogClient pb.ActivityLogDMSServiceClient) error {
	forever := make(chan bool)
	for d := range msgs {
		log := dto.UserActivityLogMessage{}
		err := json.Unmarshal(d.Body, &log)
		if err != nil {
			return err
		}
		req := &pb.CreateActionRequest{
			Id:  log.ID,
			Log: log.Log,
		}
		_, err = activityLogClient.CreateAction(ctx, req)
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
		logger.InitLogger().Fatalf("Could not close server, err: %v", err)
	}
}
