package main
import(
	"log"
	"github.com/streadway/amqp"
	// "time"
	"database/db"
	// "context"
	// "os/signal"
	// "os"
)
func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}
//AddToQueue will add new users data to queue
func main(){





	//creating connection to rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()



	//creating channel
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()


	//Declaring a queue
	q, err := ch.QueueDeclare(
		"register_queue", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")




	//defining custom exchange
	err = ch.ExchangeDeclare(
		"register_exchange",   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)



	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")



	//declare binding
	err = ch.QueueBind(
		q.Name,       // queue name
		"register_key",      // routing key
		"register_exchange", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")






	//Consumer consuming messages from queue
	queries, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")




	//go routine for displaying values
	go func() {
		for d := range queries {
			done:=db.Adduser(string(d.Body))
			if done == true{
				log.Println("successfully done:"+string(d.Body))
			}else{
				log.Println("failed:"+string(d.Body))
			}
			d.Ack(false)
		}
		// log.Printf("Done")
	}()
	forever := make(chan bool)
	<-forever
}
