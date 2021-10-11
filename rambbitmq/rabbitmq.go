package rambbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type RabbitMqManager struct {
	*amqp.Connection
	*amqp.Channel
	amqp.Queue
}

func ConectRambbitMQ() (*RabbitMqManager, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp:%s ", viper.GetString("rabbitmq.host")))
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
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
		return nil, err
	}
	return &RabbitMqManager{
		Connection: conn,
		Channel:    ch,
		Queue:      amqp.Queue{},
	}, nil
}
func (m *RabbitMqManager) Close() {
	err := m.Connection.Close()
	if err != nil {
		log.Fatalf("Could not close server, err: %v", err)
	}
	err = m.Channel.Close()
	if err != nil {
		log.Fatalf("Could not close channel, err: %v", err)
	}
}
