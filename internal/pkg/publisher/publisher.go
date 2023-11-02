package Publisher

import (
	"encoding/json"
	"log"
	"rec/internal/pkg/postgresql/model"
	"time"

	stan "github.com/nats-io/stan.go"
)

type Publisher struct {
	sc   *stan.Conn
	name string
}

func NewPublisher(conn *stan.Conn) *Publisher {
	return &Publisher{
		name: "Publisher",
		sc:   conn,
	}
}

func (p *Publisher) Publish() {

	dlvr := model.Delivery{Name: "Test Testov", Phone: "+9720000000", Zip: "2639809", City: "Kiryat Mozkin", Address: "Ploshad Mira 15", Region: "Kraiot", Email: "test@gmail.com"}
	pmnt := model.Payment{Transaction: "b563feb7b2b84b6test", RequestId: "", Currency: "USD", Provider: "wbpay", Amount: 1817, PaymentDt: 1637907727, Bank: "alpha", DeliveryCost: 1500, GoodsTotal: 317, CustomFee: 0}
	itm1 := model.Items{ChrtId: 1, TrackNumber: "WBILMTESTTRACK", Price: 453, Rid: "ab4219087a764ae0btest", Name: "Mascaras", Sale: 30, Size: "0", TotalPrice: 317, NmId: 2389212, Brand: "Vivienne Sabo", Status: 202}
	itm2 := model.Items{ChrtId: 2, TrackNumber: "WBILMTESTTRACK", Price: 678, Rid: "ab4219087a764ae0btest", Name: "Boots", Sale: 15, Size: "42", TotalPrice: 1322, NmId: 2389212, Brand: "Nike", Status: 202}
	tT := time.Date(2021, 11, 26, 18, 22, 19, 0, time.UTC)
	od := model.OrderData{OrderUid: "b563feb7b2b84b6test", TrackNumber: "WBILMTESTTRACK", Entry: "WBIL", Delivery: dlvr, Payment: pmnt, Items: []model.Items{itm1, itm2}, Locale: "en", InternalSignature: "", CustomerId: "test", DeliveryService: "meest", Shardkey: "9", SmId: 99, DateCreated: tT, OofShard: "1"}
	oi := model.OrderItem{"2", od}
	orderData, err := json.Marshal(oi)
	if err != nil {
		log.Printf("%s: json.Marshal error: %v\n", p.name, err)
	}

	AHandler := func(idd string, err error) {
		if err != nil {
			log.Printf("%s: error publishing msg id %s: %v\n", p.name, idd, err.Error())
		} else {
			log.Printf("%s: received ack for msg id: %s\n", p.name, idd)
		}
	}

	log.Printf("%s: publishing data ...\n", p.name)
	nuid, err := (*p.sc).PublishAsync("testch", orderData, AHandler)
	if err != nil {
		log.Printf("%s: error publishing msg %s: %v\n", p.name, nuid, err.Error())
	}
}
