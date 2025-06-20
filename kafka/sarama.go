package kafka

import (
	"errors"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"ppt/logger"
	"time"
)

type SaramaAsyncClient struct {
	producer sarama.AsyncProducer
	topic    string
}

func InitSaramaAsyncClient(bootstrapServer, clientID, topic string) (*SaramaAsyncClient, error) {
	config := sarama.NewConfig()
	// 检查topic是否存在
	admin, err := sarama.NewClusterAdmin([]string{bootstrapServer}, config)
	if err != nil {
		logger.Error("InitSaramaAsyncClient NewClusterAdmin error", zap.String("kafka_bootstrap_server", bootstrapServer), zap.Error(err))
		return nil, err
	}
	defer admin.Close()
	topics, err := admin.ListTopics()
	if err != nil {
		logger.Error("InitSaramaAsyncClient ListTopics error", zap.String("kafka_bootstrap_server", bootstrapServer), zap.String("topic", topic), zap.Error(err))
		return nil, err
	}
	if _, exists := topics[topic]; !exists {
		logger.Error("InitSaramaAsyncClient kafka server has no topic", zap.String("kafka_bootstrap_server", bootstrapServer), zap.String("topic", topic))
		return nil, errors.New("InitSaramaAsyncClient kafka server has no topic")
	}

	config.ClientID = clientID

	//config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 100 * time.Millisecond
	config.Producer.RequiredAcks = sarama.NoResponse // 确认机制
	config.Producer.Compression = sarama.CompressionGZIP
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	config.Producer.MaxMessageBytes = 1024 * 16
	config.Producer.Flush.Bytes = 1024 * 1024
	config.Producer.Flush.Messages = 1024
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true

	config.Net.SASL.Enable = false
	config.Net.TLS.Enable = false

	producer, err := sarama.NewAsyncProducer([]string{bootstrapServer}, config)
	if err != nil {
		logger.Error("InitSaramaAsyncClient NewAsyncProducer error", zap.Error(err))
		return nil, err
	}

	return &SaramaAsyncClient{producer, topic}, nil
}

// SendMessage 异步发送消息
func (sara *SaramaAsyncClient) SendMessage(key, value []byte) {
	msg := &sarama.ProducerMessage{
		Topic: sara.topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(value),
	}
	sara.producer.Input() <- msg
	return
}

func (sara *SaramaAsyncClient) Close() {
	if sara.producer != nil {
		sara.producer.Close()
		sara.producer = nil
	}
}

func ConsumeAsyncProducer(asyncClient *SaramaAsyncClient) {
	if asyncClient == nil {
		logger.Info("SaramaAsyncClient nil producer")
		return
	}
	producer := asyncClient.producer
	go func() {
		for {
			select {
			case <-producer.Successes():
			case err := <-producer.Errors():
				if err != nil {
					if errors.Is(err.Err, sarama.ErrUnknownTopicOrPartition) {
						logger.Fatal("ConsumeAsyncProducer topic does not exists", zap.Error(err))
					}
					logger.Error("SaramaAsyncClient GetAsyncProducer error", zap.Error(err))
				}
			}
		}
	}()
}
