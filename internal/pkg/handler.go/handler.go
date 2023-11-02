package handler

import (
	"log"
	"rec/internal/pkg/postgresql"
	Publisher "rec/internal/pkg/publisher"
	Subscriber "rec/internal/pkg/subscriber"
	"time"

	"github.com/nats-io/nats.go"
	stan "github.com/nats-io/stan.go"
)

type StreamingHandler struct {
	conn  *stan.Conn
	sub   *Subscriber.Subscriber
	pub   *Publisher.Publisher
	name  string
	isErr bool
}

func NewStreamingHandler(db *postgresql.DBService) *StreamingHandler {
	sh := StreamingHandler{}
	sh.Init(db)
	return &sh
}

func (sh *StreamingHandler) Init(db *postgresql.DBService) {
	sh.name = "StreamingHandler"
	err := sh.Connect()

	if err != nil {
		sh.isErr = true
		log.Printf("%s: StreamingHandler error: %s", sh.name, err)
	} else {
		sh.sub = Subscriber.NewSubscriber(db, sh.conn)
		sh.sub.Subscribe()

		sh.pub = Publisher.NewPublisher(sh.conn)
		sh.pub.Publish()
	}
}

func (sh *StreamingHandler) Connect() error {
	conn, err := stan.Connect(
		"test-cluster",
		"test-client",
		stan.NatsURL("localhost"),
		stan.NatsOptions(
			nats.ReconnectWait(time.Second*4),
			nats.Timeout(time.Second*4),
		),
		stan.Pings(5, 3),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Printf("%s: connection lost, reason: %v", sh.name, reason)
		}),
	)
	if err != nil {
		log.Printf("%s: can't connect: %v.\n", sh.name, err)
		return err
	}
	sh.conn = &conn

	log.Printf("%s: connected!", sh.name)
	return nil
}

func (sh *StreamingHandler) Finish() {
	if !sh.isErr {
		log.Printf("%s: Finish...", sh.name)
		sh.sub.Unsubscribe()
		(*sh.conn).Close()
		log.Printf("%s: Finished!", sh.name)
	}
}
