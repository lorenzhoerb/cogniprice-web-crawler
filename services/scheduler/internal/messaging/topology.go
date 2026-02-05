package messaging

import (
	"github.com/lorenzhoerb/cogniprice/shared/messaging"
	"github.com/rabbitmq/amqp091-go"
)

func DeclareTopology(ch *amqp091.Channel) error {
	return ch.ExchangeDeclare(
		messaging.CrawlCommandExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
}
