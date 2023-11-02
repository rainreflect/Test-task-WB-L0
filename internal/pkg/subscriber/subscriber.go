package Subscriber

import (
	"encoding/json"
	"log"
	"os"
	"rec/internal/pkg/postgresql"
	"rec/internal/pkg/postgresql/model"
	"time"

	stan "github.com/nats-io/stan.go"
)

type Subscriber struct {
	sub      stan.Subscription
	dbObject *postgresql.DBService
	sc       *stan.Conn
	name     string
}

func NewSubscriber(db *postgresql.DBService, conn *stan.Conn) *Subscriber {
	return &Subscriber{
		name:     "Subscriber",
		dbObject: db,
		sc:       conn,
	}
}

func (s *Subscriber) Subscribe() {

	var err error

	s.sub, err = (*s.sc).Subscribe(
		"testch",
		func(m *stan.Msg) {
			log.Printf("%s: received a message!\n", s.name)
			if s.messageHandler(m.Data) {
				err := m.Ack()
				if err != nil {
					log.Printf("%s ack() err: %s", s.name, err)
				}
			}
		},
		stan.AckWait(time.Duration(30)*time.Second),
		stan.DeliverAllAvailable(),

		stan.MaxInflight(10))

	if err != nil {
		log.Printf("%s: error: %v\n", s.name, err)
	}
	log.Printf("%s: subscribed to subject %s\n", s.name, os.Getenv("NATS_SUBJECT"))
}

func (s *Subscriber) messageHandler(data []byte) bool {
	recievedOrder := model.OrderItem{}
	err := json.Unmarshal(data, &recievedOrder)
	if err != nil {
		log.Printf("%s: messageHandler() error, %v\n", s.name, err)
		return true
	}
	log.Printf("%s: unmarshal Order to struct: %v\n", s.name, recievedOrder)

	_, err = s.dbObject.CreateOrder(&recievedOrder)
	if err != nil {
		log.Printf("%s: unable to add order: %v\n", s.name, err)
		return false
	}
	return true
}

func (s *Subscriber) Unsubscribe() {
	if s.sub != nil {
		s.sub.Unsubscribe()
	}
}
