package kafka

import (
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"ppt/logger"
)

type MessageHandler interface {
	Handle(message *sarama.ConsumerMessage) error
}

type SaramaHandler struct {
	handler   MessageHandler
	readyFunc func() // 就绪回调函数
}

func NewSaramaHandler(handler MessageHandler) *SaramaHandler {
	return &SaramaHandler{handler: handler}
}

func (h *SaramaHandler) Setup(session sarama.ConsumerGroupSession) error {
	logger.Info("Consumer group setup", zap.Any("Initialized", session.MemberID()), zap.Any("Claims", session.Claims()))
	if h.readyFunc != nil {
		h.readyFunc()
	}
	return nil
}

func (h *SaramaHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	logger.Info("Consumer group cleanup", zap.Any("Initialized", session.MemberID()))
	return nil
}

func (h *SaramaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logger.Info("Decoded message", zap.Any("consume_message", message))
		if err := h.handler.Handle(message); err != nil {
			logger.Error("Failed to handle message", zap.Error(err))
			continue
		}
		session.MarkMessage(message, "")
	}
	return nil
}
