package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"payment/queue"
	"time"
)

type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}

func main() {
	// abertura do canal 
	in := make(chan []byte)

	// abrir conexão
	connection := queue.Connect()
	// fila para consumir 
	queue.StartConsuming("order_queue", connection, in)

	var order Order
	for payload := range in {
		json.Unmarshal(payload, &order)
		// mudar o status da order 
		order.Status = "aprovado"
		notifyPaymentProcessed(order, connection)
	}
}

func notifyPaymentProcessed(order Order, ch *amqp.Channel) {
	json, _ := json.Marshal(order)
	// fila para informar a aprovação 
	queue.Notify(json, "payment_ex", "", ch)
	fmt.Println(string(json))
}
