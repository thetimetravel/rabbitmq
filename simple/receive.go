package main

import (
	"log"

	"github.com/streadway/amqp"
)

func failOnErrorr(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}


func main() {
	//连接RabbitMQ服务器
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnErrorr(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	//上面代码会建立一个socket连接，处理一些协议转换及版本对接和登录授权的问题。建立连接之后，我们需要创建一个通道channel
	ch, err := conn.Channel()
	failOnErrorr(err, "Failed to open a channel")
	defer ch.Close()

	//需要RabbitMQ服务器让它将消息分发到我们的消费者程序中，消息转发操作是异步执行的，
	//这里使用goroutine来完成从队列中的读取消息操作：
	q, err := ch.QueueDeclare(
		"simple", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnErrorr(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnErrorr(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}