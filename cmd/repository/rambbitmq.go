package repository

import (
	"ktrain/rambbitmq"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ActivityLogRambbitMqRepository interface {
	Publish(id int64, log string) error
}

type activityLogRambbitMqRepository struct {
	manager *rambbitmq.RabbitMqManager
}

func NewActivityLogRambbitMqRepository(db *rambbitmq.RabbitMqManager) ActivityLogRambbitMqRepository {
	return &activityLogRambbitMqRepository{
		manager: db,
	}
}
func (m *activityLogRambbitMqRepository) Publish(id int64, log string) error {
	body := strconv.Itoa(int(id)) + log
	err := m.manager.Channel.Publish(
		"activityLog", // exchange
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	return nil
}
