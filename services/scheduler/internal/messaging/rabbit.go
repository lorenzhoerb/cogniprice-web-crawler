package messaging

import "github.com/rabbitmq/amqp091-go"

type Rabbit struct {
	Conn *amqp091.Connection
}

func NewRabbit(url string) (*Rabbit, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Rabbit{Conn: conn}, nil
}

func (r *Rabbit) Close() error {
	return r.Conn.Close()
}
