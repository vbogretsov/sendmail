package client

import (
	"encoding/json"

	"github.com/streadway/amqp"

	"github.com/vbogretsov/sendmail/model"
)

type Client struct {
	channel *amqp.Channel
	topic   string
}

func New(conn *amqp.Connection, topic string) (*Client, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	sender := Client{
		channel: ch,
		topic:   topic,
	}
	return &sender, nil
}

func (s *Client) Send(req model.Request) error {
	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{Body: buf}

	return s.channel.Publish(s.topic, s.topic, false, false, msg)
}

func (s *Client) Close() error {
	return s.channel.Close()
}
