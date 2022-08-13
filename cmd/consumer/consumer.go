package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	brokerAddress *string = flag.String("brokers", "localhost:9092", "address of the kafka broker")
	topic         *string = flag.String("topic", "", "the kafka topic name to send messages to")
	apiKey        *string = flag.String("api-key", "", "Confluent Cloud API Key")
	secretKey     *string = flag.String("secret-key", "", "Confluent Cloud Secret Key")
)

func main() {
	flag.Parse()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": *brokerAddress,
		"group.id":          "foo",
		"auto.offset.reset": "earliest",
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     *apiKey,
		"sasl.password":     *secretKey,
		"client.id":         "1",
		"acks":              "all",
	})
	if err != nil {
		fmt.Errorf("Unable to create consumer: %s\n", err)
		return
	}

	err = consumer.SubscribeTopics([]string{*topic}, nil)

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}

	consumer.Close()

}
