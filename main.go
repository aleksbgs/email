package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/smtp"
	"os"
)

func main() {

	fmt.Println(os.Getenv("KAFKA_TOPIC"))

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
		"sasl.username":     os.Getenv("SASL_USERNAME"),
		"sasl.password":     os.Getenv("SASL_PASSWORD"),
		"sasl.mechanism":    os.Getenv("SASL_MECHANISM"),
		"security.protocol": os.Getenv("SECURITY_PROTOCOL"),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	consumer.SubscribeTopics([]string{"default"}, nil)

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {

			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			return
		}
		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

		var message map[string]interface{}

		json.Unmarshal(msg.Value, &message)

		host := os.Getenv("EMAIL_HOST")

		port := os.Getenv("EMAIL_PORT")

		auth := smtp.PlainAuth("", os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"), host)

		ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link #%s", message["ambassador_revenue"], message["code"]))

		smtp.SendMail(host+":"+port, auth, "no-reply@email.com", []string{message["ambassador_email"].(string)}, ambassadorMessage)

		adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed", message["id"], message["admin_revenue"]))

		smtp.SendMail(host+":"+port, auth, "no-reply@email.com", []string{"admin@admin.com"}, adminMessage)
	}

	consumer.Close()

}
