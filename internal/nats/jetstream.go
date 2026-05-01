package nats

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	js      jetstream.JetStream
	subject string
}

func New(url, stream, subject string) (*Client, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     stream,
		Subjects: []string{subject},
		Storage:  jetstream.FileStorage,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		js:      js,
		subject: subject,
	}, nil
}

func (c *Client) Publish(ctx context.Context, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = c.js.Publish(ctx, c.subject, payload)
	if err != nil {
		return err
	}
	return nil
}
