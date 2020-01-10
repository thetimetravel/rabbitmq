package main

import (
	"bytes"
	"log"
	"time"

	"github.com/streadway/amqp"
)
//公平调度
func failOnErrorw22(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorw22(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnErrorw22(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnErrorw22(err, "Failed to declare a queue")

	//关键代码，1代表一次只能放松一条，而且要等消息返回才能再发下一条99
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnErrorw22(err, "Failed to set QoS")

	//根据消息体中"."的个数来模拟任务的耗时长度。该文件的任务还是从队列中取出一个任务并执行，我们姑且称之为work.go
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnErrorw22(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			time.Sleep(2* time.Second)
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}