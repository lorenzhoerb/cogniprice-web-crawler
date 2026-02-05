package dispatcher

import (
	"encoding/json"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/messaging"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/shared/logger"
	shared "github.com/lorenzhoerb/cogniprice/shared/messaging"
	"github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type rabbitDispatcher struct {
	ch *amqp091.Channel
}

func NewRabbitDispatcher(rabbit *messaging.Rabbit) (*rabbitDispatcher, error) {
	ch, err := rabbit.Conn.Channel()
	if err != nil {
		return nil, err
	}
	return &rabbitDispatcher{ch: ch}, nil
}

func (d *rabbitDispatcher) DispatchJobs(jobs []model.JobDispatched) error {
	for _, job := range jobs {
		logger.Log.Info("dispatching job",
			// structured fields
			zap.Uint64("jobID", uint64(job.ID)),
			zap.String("url", job.URL),
		)
		if err := d.publish(job); err != nil {
			return err
		}
	}
	return nil
}

func (d *rabbitDispatcher) publish(job model.JobDispatched) error {
	jsonPayload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	err = d.ch.Publish(
		shared.CrawlCommandExchange,
		shared.CrawlCommandExecuteKey,
		false,
		false,
		amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			Body:         jsonPayload,
		},
	)

	if err != nil {
		return err
	}

	return nil
}
