package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/streadway/amqp"
)

func failOnError2(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:123456@localhost:5672/")
	failOnError2(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError2(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError2(err, "Failed to declare a queue")
	//让程序通过命令行将任意个消息参数传递到队列，姑且将新文件命名为new_task.go:


	for i:=1;i<10;i++ {
		body := bodyFrom2(os.Args)
		body = strconv.Itoa(i)+" "+body
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError2(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)
	}
}

func bodyFrom2(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "simple"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}