package mq

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
)

type Callback func(msg string)

func Connect() (*amqp.Connection,error){
	conn,err := amqp.Dial("amqp://guest:guest@127.0.0.1:5672/")
	return conn,err
}

//简单模式和工作模式 发送端
func Publish(exchange string, queueName string, body string) error{
	//1：建立连接
	conn,err := Connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	//2：创建通道
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	//3：创建队列
	q,err := channel.QueueDeclare(queueName,true,false,false,false,nil)
	if err != nil {
		return err
	}

	//4：发送消息
	err = channel.Publish(exchange, q.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	return err
}
//简单模式和工作模式 接收端
func Consumer(exchange string, queueName string, callback Callback){
	//1：建立连接
	conn,err := Connect()
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	//2：创建通道
	channel,err := conn.Channel()
	defer channel.Close()
	if err != nil{
		fmt.Println(err)
		return
	}

	//3：创建队列
	q,err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}


	//4：从队列中获取数据 第二参数是路由 第三个参数自动应答
	msgs,err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//5：进行消息处理
	forever := make(chan bool)
	go func(){
		for d := range msgs{
			s := BytesToString(&(d.Body))
			callback(*s)
			//回调函数执行完 代表业务处理完消息进行手动应答
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for messages")
	<-forever
}

func BytesToString(b *[]byte) *string {
	s := bytes.NewBuffer(*b)
	r := s.String()
	return &r
}

//订阅模式-路由模式-主题模式
func PublishEx(exchange string, types string, routingKey string, body string) error {
	//1：建立连接
	conn,err := Connect()
	defer conn.Close()
	if err != nil {
		return err
	}

	//2：创建通道
	channel, err := conn.Channel()
	defer channel.Close()
	if err != nil{
		return err
	}

	//3：创建交换机
	err = channel.ExchangeDeclare(exchange, types,true, false, false,false,nil)
	if err != nil{
		return err
	}

	//4：发送消息
	err = channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType: "text/plain",
		Body: []byte(body),
	})

	return err
}
func ConsumerEx(exchange string, types string, routingKey string, callback Callback){
	//1：建立连接
	conn,err := Connect()
	defer conn.Close()
	if err != nil{
		fmt.Println(err)
		return
	}

	//2：创建通道
	channel,err := conn.Channel()
	defer channel.Close()
	if err != nil{
		fmt.Println(err)
		return
	}

	//3：创建交换机
	err = channel.ExchangeDeclare(exchange,types,true,false, false, false, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

	//4：创建队列
	q,err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

	//5：绑定
	err = channel.QueueBind(q.Name, routingKey, exchange, false, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

	//6：从队列中获取数据 第二参数是路由 第三个参数自动应答
	msgs,err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//7：进行消息处理
	forever := make(chan bool)
	go func(){
		for d := range msgs{
			s := BytesToString(&(d.Body))
			callback(*s)
			//回调函数执行完 代表业务处理完消息进行手动应答
			d.Ack(false)
		}
	}()
	fmt.Println("Waiting for messages\n")
	<-forever

}

//死信队列
func PublishDlx(exchangeA string, body string) error {
	//1：建立连接
	conn,err := Connect()
	if err != nil{
		return err
	}
	defer conn.Close()

	//2：创建通道
	channel,err := conn.Channel()
	if err != nil{
		return err
	}
	defer channel.Close()

	//3：消息发送到A交换机
	err = channel.Publish(exchangeA, "", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType: "text/plain",
		Body:[]byte(body),
	})
	return err
}
func ConsumerDlx(exchangeA string, queueAName string, exchangeB string, queueBName string, ttl int, callback Callback){
	//1:建立连接
	conn,err := Connect()
	if err != nil{
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//2:创建通道
	channel,err := conn.Channel()
	if err != nil{
		fmt.Println(err)
		return
	}
	defer channel.Close()

	//3: 创建A交换机-创建A队列-A交换机和A队列绑定
	err = channel.ExchangeDeclare(exchangeA, "fanout", true, false, false, false, nil)
	if err != nil{
		fmt.Println(err)
		return
	}
	//3.1 创建一个queue，指定消息过期时间，并且绑定过期以后发送到那个交换机
	queueA,err := channel.QueueDeclare(queueAName, true, false, false, false, amqp.Table{
		// 当消息过期时把消息发送到 exchangeB
		"x-dead-letter-exchange": exchangeB,
		"x-message-ttl":          ttl,
		//"x-dead-letter-queue" : queueBName,
		//"x-dead-letter-routing-key" :
	})
	if err != nil{
		fmt.Println(err)
		return
	}

	err = channel.QueueBind(queueA.Name, "", exchangeA, false, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

	//4: 创建B交换机-创建B队列-B交换机和B队列绑定
	err = channel.ExchangeDeclare(exchangeB, "fanout", true, false, false, false, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

	queueB,err := channel.QueueDeclare(queueBName, true, false, false, false, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

	err = channel.QueueBind(queueB.Name, "", exchangeB, false, nil)
	if err != nil{
		fmt.Println(err)
		return
	}

	//5：从队列中获取数据 第二参数是路由 第三个参数自动应答
	msgs, err := channel.Consume(queueB.Name, "", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	//6：进行消息处理
	forever := make(chan bool)
	go func(){
		for d := range msgs{
			s := BytesToString(&(d.Body))
			callback(*s)
			//回调函数执行完 代表业务处理完消息进行手动应答
			d.Ack(false)
		}
	}()
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
