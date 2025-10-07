package rabbitmq

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer は RabbitMQ からメッセージを受信する構造体
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

// NewConsumer は新しい Consumer を作成する
func NewConsumer(url, queueName string) (*Consumer, error) {
	// RabbitMQ に接続
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// チャネルを開く
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// キューを宣言
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

// Consume はメッセージを受信し続ける
func (c *Consumer) Consume() error {
	// メッセージを受信するためのチャネルを取得
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		"",           // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	// ブロッキングチャネル
	forever := make(chan bool)

	// ゴルーチンでメッセージを処理
	go func() {
		for d := range msgs {
			log.Printf(" [x] Received: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	return nil
}

// Close はコネクションとチャネルをクローズする
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
