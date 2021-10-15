package rambbitmq

import (
	"encoding/json"
	"fmt"
	"ktrain/cmd/api/user-api/dto"
	"ktrain/pkg/logger"

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
func (m *RabbitMqManager) Publish(body dto.UserActivityLogMessage) error {
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
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"activityLog", // exchange
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bodyJson,
		})
	if err != nil {
		return err
	}
	return nil
}
func (m *RabbitMqManager) Close() {
	err := m.Connection.Close()
	if err != nil {
		logger.Log().Fatalf("Could not close server, err: %v", err)
	}
}
