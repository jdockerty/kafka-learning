package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	brokerAddress *string = flag.String("brokers", "", "Address of the kafka broker")
	apiKey        *string = flag.String("api-key", "", "Confluent Cloud API Key")
	strimzi       *bool   = flag.Bool("strimzi", true, "Using Strimzi locally, no keys required")
	secretKey     *string = flag.String("secret-key", "", "Confluent Cloud Secret Key")
	topic         *string = flag.String("topic", "", "Kafka topic name to send messages to")
)

func main() {

	flag.Parse()

	if *brokerAddress == "" {
		fmt.Println("Broker address must be set!")
		os.Exit(2)
	}

	var producer *kafka.Producer
	if !*strimzi {
		if *apiKey == "" || *secretKey == "" {
			fmt.Println("Must provide Confluent Cloud keys!")
			os.Exit(2)
		}
		p, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": *brokerAddress,
			"security.protocol": "SASL_SSL",
			"sasl.mechanisms":   "PLAIN",
			"sasl.username":     *apiKey,
			"sasl.password":     *secretKey,
			"client.id":         "1",
			"acks":              "all",
		})
		if err != nil {
			fmt.Printf("Failed to create producer: %s\n", err)
			os.Exit(1)
		}

		producer = p
	} else {
		p, err := kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": *brokerAddress,
			"client.id":         "1",
			"acks":              "all",
		})
		if err != nil {
			fmt.Printf("Failed to create producer: %s\n", err)
			os.Exit(1)
		}

		producer = p
	}
	// Go-routine to handle message delivery reports and
	// possibly other event types (errors, stats, etc)
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Failed to deliver message: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Produced event to topic %s: key = %-10s value = %s\n",
						*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
				}
			}
		}
	}()

	users := [...]string{"eabara", "jsmith", "sgarcia", "jbernard", "htanaka", "awalther"}
	items := [...]string{"book", "alarm clock", "t-shirts", "gift card", "batteries"}

	for n := 0; n < 10; n++ {
		key := users[rand.Intn(len(users))]
		data := items[rand.Intn(len(items))]
		producer.Produce(&kafka.Message{

			TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
			Key:            []byte(key),
			Value:          []byte(data),
		}, nil)
	}

	// Wait for all messages to be delivered
	producer.Flush(15 * 1000)
	producer.Close()
}
