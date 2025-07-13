package test

import (
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"ppt/kafka"
	"ppt/logger"
	"testing"
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

	userID := "101000001"
	kafkaMsg := map[string]interface{}{
		"user_id": userID,
		"gender":  true,
		"age":     23,
		"level":   99,
		"exp":     234523,
	}
	kafkaMsgBytes, err := json.Marshal(kafkaMsg)
	if err != nil {
		logger.Error("TestKafka Marshal error", zap.Error(err))
		return
	}
	asyncClient.SendMessage([]byte(userID), kafkaMsgBytes)
}

func TestSaramaConsumer(t *testing.T) {

}
