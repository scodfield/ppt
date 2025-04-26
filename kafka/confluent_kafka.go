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
	go func() {
		var cnt int
		for ev := range ConsumerClient.Events() {
			switch e := ev.(type) {
			case *kafka.Message:
				if e.TopicPartition.Error != nil {
					logger.Error("StartConsumerLoop Delivery failed", zap.Error(e.TopicPartition.Error))
					continue
				}
				cnt++
				_ = pool.GetConsumerPool().Submit(func() {
					processKafkaMessage(e, cnt >= ProcessBatchSize)
				})
			case *kafka.Error:
				logger.Error("StartConsumerLoop Delivery failed: %v\n", zap.Error(e))
				break
			}
		}
		logger.Info("StartConsumerLoop kafka consumer loop exit")
	}()
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

func StartKafkaProducerLoop(producer *kafka.Producer) {
	go func() {
		for ev := range producer.Events() {
			switch e := ev.(type) {
			case *kafka.Message:
				if e.TopicPartition.Error != nil {
					logger.Error("StartKafkaProducerLoop delivery err", zap.Error(e.TopicPartition.Error))
				} else {
					logger.Info("StartKafkaProducerLoop success delivered message", zap.Any("topic_partition", e.TopicPartition))
				}
			case kafka.Error:
				logger.Error("StartKafkaProducerLoop kafka err", zap.Error(e))
			}
		}
		logger.Info("StartKafkaProducerLoop kafka producer loop exit")
	}()
}

func GetKafkaClientID(serviceName string) string {
	return fmt.Sprintf("%s-%s-%s-%s", config.AppName, config.Env, config.HostName, serviceName)
}

func GetKafkaTopic(serviceName string) string {
	return fmt.Sprintf("%s-%s-%s", config.AppName, config.Env, serviceName)
}

func CloseKafkaClient() {
	if ConsumerClient != nil {
		_ = ConsumerClient.Close()
	}
	if ProducerClient != nil {
		ProducerClient.Close()
	}
}

func SendKafkaMessage(topic string, message []byte) {
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: message,
	}
	err := ProducerClient.Produce(msg, nil)
	if err != nil {
		logger.Error("SendKafkaMessage produce error", zap.Error(err))
	}
	ProducerClient.Flush(0)
}
