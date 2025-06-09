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
