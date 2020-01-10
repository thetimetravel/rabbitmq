package main
import (
	"fmt"
	"github.com/streadway/amqp"
	"log"

)
func main() {
	//连接RabbitMQ服务器：
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect to RabbitMQ")
	fmt.Println("eeeeeed")
	defer conn.Close()
	//上面代码会建立一个socket连接，处理一些协议转换及版本对接和登录授权的问题。建立连接之后，
	//我们需要创建一个通道channel，之后我们的大多数API操作都是围绕通道来实现的：
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel ")
	defer ch.Close()
	//需要定义一个队列用来存储、转发消息，然后我们的sender只需要将消息发送到这个队列中，就完成了消息的publish操作

}


	func failOnError(err error, msg string){
	if err != nil {
		log.Fatalf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}
