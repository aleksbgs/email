package main

import (
	"encoding/json"
	"fmt"
	"github.com/aleksbgs/email/utils"
	"net/smtp"
)
import "github.com/confluentinc/confluent-kafka-go/kafka"

func main() {

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": utils.ViperEnvVariable("bootstrap.servers"),
		"security.protocol": "SASL_SSL",
		"sasl.username":     utils.ViperEnvVariable("sasl.username"),
		"sasl.password":     utils.ViperEnvVariable("sasl.password"),
		"sasl.mechanism":    "PLAIN",
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

		ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link #%s", message["ambassador_revenue"], message["code"]))

		smtp.SendMail("host.docker.internal:1025", nil, "no-reply@email.com", []string{message["ambassador_email"].(string)}, ambassadorMessage)

		adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed", message["id"], message["admin_revenue"]))

		smtp.SendMail("host.docker.internal:1025", nil, "no-reply@email.com", []string{"admin@admin.com"}, adminMessage)
	}

	consumer.Close()

}
