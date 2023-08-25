package natsstreaming

import (
	"L0_EVRONE/internal/usecase"
	"L0_EVRONE/pkg/logger"
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// connect to nats server
	nc, _ := nats.Connect(nats.DefaultURL)

	// create jetstream context from nats connection
	js, _ := jetstream.New(nc)

	// Create a stream
	s, _ := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "ORDERS",
		Subjects: []string{"ORDERS.*"},
	})

	c, _ := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "CONS",
		AckPolicy: jetstream.AckExplicitPolicy,
	})

	// Receive messages continuously in a callback
	cons, _ := c.Consume(func(msg jetstream.Msg) {
		msg.Ack()
		data, err := Serialization(msg.Data())
		if err != nil {
			l.Error("")
		}
		err = t.Set(context.Background(), data)
		if err != nil {
			l.Error("")
		}
		//fmt.Printf("Received a JetStream message via callback: %s\n", string(msg.Data()))
	})

	return &Nats{
		cancel:   cancel,
		consStop: cons,
	}
}

func Close() {

}

func GetMessage() {

}

/*func main() {
//	// Создаем клиента
	client, err := nats.NewClient("nats://<your-nats-host-address>:4222")
	if err != nil {
		log.Fatal(err)
	}

	// Регистрируемся для прослушивания канала
	_, err = client.Channel("my-channel").Subscribe(
        func(string) {
		// Обрабатываем сообщение
	}
    )
	if err != nil {
		log.Fatalf("Failed to subscribe: %s", err)
	}
}
*/
