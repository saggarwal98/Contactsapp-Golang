package rabbit
import(
	"log"
	"github.com/streadway/amqp"
	// "time"
)
func failOnError(err error, msg string) {
	if err != nil {
	  log.Fatalf("%s: %s", msg, err)
	}
}
//AddToQueue will add new users data to queue
func AddToQueue(query string){

	
	//creating connection to rabbitmq
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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




	//declare binding
	err = ch.QueueBind(
		q.Name,       // queue name
		"register_key",      // routing key
		"register_exchange", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")





	//Sending data to Queue
	err = ch.Publish(
		"register_exchange",     // exchange
		"register_key", // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			DeliveryMode: amqp.Persistent,
		    ContentType: "text/plain",
		    Body:        []byte(query),
		})
	failOnError(err, "Failed to publish a message")

}