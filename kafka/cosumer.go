package kafka

import (
	"actor2/pool"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var ConsumerClient *kafka.Consumer

const (
	ProcessBatchSize = 200
)

func InitKafka(bootstrapUrls []string, groupID, topic string) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        bootstrapUrls,
		"group.id":                 groupID,
		"auto.offset.reset":        "latest",
		"fetch.min.bytes":          1,        // 最小拉取字节数
		"fetch.max.bytes":          52428800, // 最大拉取字节数
		"fetch.wait.max.ms":        "500",    // 如果没有最新消费消息默认等待500ms
		"enable.auto.commit":       false,    // 手动提交
		"go.events.channel.enable": true,
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s\n", err)
		return err
	}
	ConsumerClient = consumer
	consumer.SubscribeTopics([]string{topic}, nil)
	return nil
}

func StartConsumerLoop() {
	var cnt int
	for ev := range ConsumerClient.Events() {
		switch e := ev.(type) {
		case *kafka.Message:
			if e.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", e.TopicPartition.Error)
				continue
			}
			cnt++
			_ = pool.GetConsumerPool().Submit(func() {
				processKafkaMessage(e, cnt >= ProcessBatchSize)
			})
		case *kafka.Error:
			fmt.Printf("Delivery failed: %v\n", e)
			break
		}
	}
	fmt.Println("kafka consumer loop exit")
	err := ConsumerClient.Close()
	if err != nil {
		fmt.Printf("Consumer close failed: %v\n", err)
	}
}

func processKafkaMessage(msg *kafka.Message, isCommit bool) {
	var kafkaMsg map[string]interface{}
	err := json.Unmarshal(msg.Value, &kafkaMsg)
	if err != nil {
		fmt.Printf("Error unmarshalling message: %v\n", err)
		return
	}
	fmt.Printf("Kafka message: %+v\n", kafkaMsg)
	if isCommit {
		ConsumerClient.Commit()
	}
}
