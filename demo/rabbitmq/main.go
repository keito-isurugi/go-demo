package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	rabbitmq "github.com/isurugikeito/go-demo/demo/rabbitmq"
)

const (
	rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	queueName   = "hello"
)

func main() {
	// コマンドライン引数で producer/consumer を切り替え
	mode := flag.String("mode", "producer", "mode: producer or consumer")
	flag.Parse()

	switch *mode {
	case "producer":
		runProducer()
	case "consumer":
		runConsumer()
	default:
		log.Fatalf("Invalid mode: %s. Use 'producer' or 'consumer'", *mode)
	}
}

func runProducer() {
	log.Println("Starting Producer...")

	// Producer を作成
	producer, err := rabbitmq.NewProducer(rabbitmqURL, queueName)
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer producer.Close()

	// メッセージを5つ送信
	for i := 1; i <= 5; i++ {
		message := fmt.Sprintf("Hello RabbitMQ! Message #%d", i)
		if err := producer.Publish(message); err != nil {
			log.Printf("Failed to publish message: %v", err)
			continue
		}
		time.Sleep(1 * time.Second)
	}

	log.Println("Producer finished")
}

func runConsumer() {
	log.Println("Starting Consumer...")

	// Consumer を作成
	consumer, err := rabbitmq.NewConsumer(rabbitmqURL, queueName)
	if err != nil {
		log.Fatalf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	// メッセージを受信し続ける (CTRL+C で終了)
	if err := consumer.Consume(); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
}
