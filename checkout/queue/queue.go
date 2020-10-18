package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
)
// importar pacote do protocolo do rabbitmq 

// ponteiro para modificar no mesmo lugar 
// conectar com o sistema de mensageria 
func Connect() *amqp.Channel {

	// link de conexão 
	dsn := "amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":" + os.Getenv("RABBITMQ_DEFAULT_PORT") + os.Getenv("RABBITMQ_DEFAULT_VHOST")

	// fazer conexão 
	conn, err := amqp.Dial(dsn)
	if err != nil {
		panic(err.Error())
	}

	// criar canal 
	channel, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}
	return channel
}

// envio de mensanges 
// dados(mensangem),
func Notify(payload []byte, exchange string, routingKey string, ch *amqp.Channel) {

	// nao precisa retorno de algo nessa variavel 
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
