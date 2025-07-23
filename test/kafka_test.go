package test

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"ppt/kafka"
	"ppt/logger"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup() {

}

func teardown() {

}

func TestKafka(t *testing.T) {
	bootstrapServerAddr := "localhost:9092"
	//groupID := "ppt-test-group"
	clientID := "ppt-test-client"
	topic := "ppt-test-statics"
	asyncClient, err := kafka.InitSaramaAsyncClient(bootstrapServerAddr, clientID, topic)
	if err != nil {
		logger.Error("TestKafka InitSaramaAsyncClient error", zap.Error(err))
		return
	}
	kafka.ConsumeAsyncProducer(asyncClient)

	userID := "101000851"
	kafkaMsg := map[string]interface{}{
		"user_id": userID,
		"gender":  true,
		"age":     24,
		"level":   95,
		"exp":     234523,
	}
	kafkaMsgBytes, err := json.Marshal(kafkaMsg)
	if err != nil {
		logger.Error("TestKafka Marshal error", zap.Error(err))
		return
	}
	asyncClient.SendMessage([]byte(userID), kafkaMsgBytes)
	time.Sleep(5 * time.Second)
}

func TestSaramaConsumer(t *testing.T) {
	bootstrapServerAddr := "localhost:9092"
	clientID := "ppt-test-client"
	groupID := "ppt-test-group"
	topics := []string{"ppt-test-statics"}
	handler := &kafka.SaramaKafkaHandler{}
	consumer, err := kafka.InitSaramaConsumerClient(bootstrapServerAddr, clientID, groupID, topics, handler)
	if err != nil {
		logger.Error("TestKafka InitSaramaConsumerClient error", zap.Error(err))
		return
	}
	consumer.Start()
}
