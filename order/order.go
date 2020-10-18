package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
	"order/db"
	"order/queue"
	"os"
	"time"
)

type Product struct {
	Uuid    string  `json:"uuid"`
	Product string  `json:"product"`
	Price   float32 `json:"price,string"`
}

type Order struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	ProductId string    `json:"product_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at,string"`
}

var productsUrl string

func init() {
	productsUrl = os.Getenv("PRODUCT_URL")
}

func main() {
	var param string

	// passar um paramentro na execução da aplicação para indicar qual fila ler 
	flag.StringVar(&param, "opt","","Usage")
	flag.Parse()

	in := make(chan []byte)
	connection := queue.Connect()

	switch param {
	case "checkout":
		// chanel está recebendo o valores 
		queue.StartConsuming("checkout_queue", connection, in)
		// leitura do conteudo do chanel, para poder esvaziar o chanel 
		for payload := range in {
			// criação de uma order e gerar uma nitificação
			notifyOrderCreated(createOrder(payload), connection)
			fmt.Println(string(payload))
		}
	case "payment":
		// aqui é a fila para realizar e salvar o processo de compra com seu status de aprovedo/negado 
		queue.StartConsuming("payment_queue", connection, in)
		var order Order
		for payload := range in {
			json.Unmarshal(payload, &order)
			saveOrder(order)
			fmt.Println("Payment: ",string(payload))
		}
	}
}

// criar uma order e salvar no BS em outro metodo 
func createOrder(payload []byte) Order {
	var order Order
	json.Unmarshal(payload, &order)

	uuid, _ := uuid.NewV4()
	order.Uuid = uuid.String()
	order.Status = "pendente"
	order.CreatedAt = time.Now()
	saveOrder(order)
	return order
}

// serviço do redis está no docker compose 
func saveOrder(order Order) {
	// converter em json
	json, _ := json.Marshal(order)
	// abrir a conexão 
	connection := db.Connect()

// inserir no red 
	err := connection.Set(order.Uuid, string(json), 0).Err()
	if err != nil {
		panic(err.Error())
	}
}

// notificar quando criar order 
func notifyOrderCreated(order Order, ch *amqp.Channel)  {
	json, _ := json.Marshal(order)

	// enviar a notificação na fila 
	queue.Notify(json, "order_ex", "", ch)
}