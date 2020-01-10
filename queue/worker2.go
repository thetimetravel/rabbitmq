package main

import (
	"bytes"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnErrorw2(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
//循环调度
func main() {
	conn, err := amqp.Dial("amqp://admin:123456@localhost:5672/")
	failOnErrorw2(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrorw2(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnErrorw2(err, "Failed to declare a queue")



	//根据消息体中"."的个数来模拟任务的耗时长度。该文件的任务还是从队列中取出一个任务并执行，我们姑且称之为work.go
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnErrorw2(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			time.Sleep(2 * time.Second)
			log.Printf("2 Received a message: %s\r", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done\r")

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}