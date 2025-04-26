package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"ppt/config"
	"ppt/logger"
	"ppt/pool"
)

var ConsumerClient *kafka.Consumer
var ProducerClient *kafka.Producer

const (
	ProcessBatchSize       = 200
	ProducerDefaultRetries = 3
)

func InitKafkaConsumerClient(bootstrapUrls []string, groupID, clientID, topic string) error {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":        bootstrapUrls,
		"group.id":                 groupID,
		"client.id":                clientID,
		"auto.offset.reset":        "latest",
		"fetch.min.bytes":          1,        // 最小拉取字节数
		"fetch.max.bytes":          52428800, // 最大拉取字节数
		"fetch.wait.max.ms":        "500",    // 如果没有最新消费消息默认等待500ms
		"enable.auto.commit":       false,    // 手动提交
		"go.events.channel.enable": true,
	})
	if err != nil {
		logger.Error("InitKafkaConsumerClient kafka new consumer error", zap.Error(err))
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

func InitKafkaProducerClient(bootstrapUrls []string, groupID, clientID, topic string) error {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":        bootstrapUrls,
		"group.id":                 groupID,
		"client.id":                clientID,
		"retries":                  ProducerDefaultRetries,
		"auto.offset.reset":        "latest",
		"enable.auto.commit":       true,
		"go.events.channel.enable": true,
	})
	if err != nil {
		logger.Error("InitKafkaProducerClient init kafka producer err", zap.Error(err))
		return err
	}
	ProducerClient = producer
	return nil
}

func GetKafkaClientID(serviceName string) string {
	return fmt.Sprintf("%s-%s-%s-%s", config.AppName, config.Env, config.HostName, serviceName)
}

func GetKafkaTopic(serviceName string) string {
	return fmt.Sprintf("%s-%s-%s", config.AppName, config.Env, serviceName)
}

func CloseKafka() {
	if ConsumerClient != nil {
		_ = ConsumerClient.Close()
	}
	if ProducerClient != nil {
		ProducerClient.Close()
	}
}
