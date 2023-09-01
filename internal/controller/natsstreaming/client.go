package natsstreaming

import (
	"L0_WB/internal/usecase"
	"L0_WB/pkg/logger"
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Nats struct {
	cancel   context.CancelFunc
	consStop jetstream.ConsumeContext
}

func New(l logger.Interface, t usecase.OrderUseCase) *Nats {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// connect to nats server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		l.Error("nats.Connect: %v", err)
		return &Nats{cancel: cancel, consStop: nil}
	}
	// create jetstream context from nats connection
	js, err := jetstream.New(nc)
	if err != nil {
		l.Error("jetstream.New: %v", err)
		return &Nats{cancel: cancel, consStop: nil}
	}
	// Create a stream
	s, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})
	if err != nil {
		l.Error("js.CreateStream: %v", err)
		return &Nats{cancel: cancel, consStop: nil}
	}
	c, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		l.Error("s.CreateOrUpdateConsumer: %v", err)
		return &Nats{cancel: cancel, consStop: nil}
	}
	// Receive messages continuously in a callback
	cons, err := c.Consume(func(msg jetstream.Msg) {
		msg.Ack()
		err := Validation(msg.Data())
		if err != nil {
			l.Error("Validation: %v", err)
			return
		}
		data, err := Serialization(msg.Data())
		if err != nil {
			l.Error("Serialization: %v", err)
			return
		}
		err = t.Set(context.Background(), data)
		if err != nil {
			l.Error("Seve data: %v", err)
			return
		}
		//fmt.Printf("Received a JetStream message via callback: %s\n", string(msg.Data()))
	})
	if err != nil {
		l.Error("c.Consume: %v", err)
		return &Nats{cancel: cancel, consStop: nil}
	}

	return &Nats{
		cancel:   cancel,
		consStop: cons,
	}
}

func (n *Nats) Shutdown() error {

	n.consStop.Stop()
	n.cancel()
	return nil
}
