package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)

func Connect() *amqp.Channel {
	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")

	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	return channel
}

func Notify(payload []byte, exchange string, routingKey string, ch *amqp.Channel) {

	err := ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Message sent")
}

// consumir a fila no sistema de mensageria 
// ultimo parametro é um chanel (pode ser lido em outros lugares )
func StartConsuming(queue string, ch *amqp.Channel, in chan []byte) {

	// para consumir, precisa declarar a fila, se nao existir, é criada 
	// passar propriedades 
	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err.Error())
	}

	// pegar as mensagens da fila para consumir 
	msgs, err := ch.Consume(
		q.Name,
		"checkout",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		panic(err.Error())
	}

	// assincronismo no go lang
	go func() {
		// só continua executado quando o chanel for lido em outro lugar 
		for m := range msgs {
			// mensagens inseridas no chanel 
			in <- []byte(m.Body)
		}
		close(in)
	}()
}
