package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"os"
	"ppt/config"
	"ppt/log"
	"ppt/nacos/wrapper"
	"strings"
	"sync"
	"time"
)

var (
	KafkaProducerClient *SaramaAsyncClient
)

func InitKafkaSarama(kafkaCfg *wrapper.KafkaConfig) error {
	var err error
	clientID := GetKafkaClientID()
	topic := GetKafkaTopic("statics")
	KafkaProducerClient, err = InitSaramaAsyncClient(kafkaCfg.BootstrapServer, clientID, topic)
	if err != nil {
		return err
	}
	return nil
}

func GetKafkaClientID() string {
	hostName, err := os.Hostname()
	if err != nil {
		log.Error("GetKafkaClientID Hostname error", zap.String("service_name", config.AppName), zap.Error(err))
		return ""
	}
	return fmt.Sprintf("%s-%s-%s", config.AppName, config.Env, hostName)
}

func GetKafkaTopic(moduleName string) string {
	return fmt.Sprintf("%s-%s-%s", config.AppName, config.Env, moduleName)
}

type SaramaAsyncClient struct {
	producer sarama.AsyncProducer
	topic    string
}

func InitSaramaAsyncClient(bootstrapServer, clientID, topic string) (*SaramaAsyncClient, error) {
	config := sarama.NewConfig()
	// 检查topic是否存在
	//admin, err := sarama.NewClusterAdmin([]string{bootstrapServer}, config)
	//if err != nil {
	//	logger.Error("InitSaramaAsyncClient NewClusterAdmin error", zap.String("kafka_bootstrap_server", bootstrapServer), zap.Error(err))
	//	return nil, err
	//}
	//defer admin.Close()
	//topics, err := admin.ListTopics()
	//if err != nil {
	//	logger.Error("InitSaramaAsyncClient ListTopics error", zap.String("kafka_bootstrap_server", bootstrapServer), zap.String("topic", topic), zap.Error(err))
	//	return nil, err
	//}
	//if _, exists := topics[topic]; !exists {
	//	logger.Error("InitSaramaAsyncClient kafka server has no topic", zap.String("kafka_bootstrap_server", bootstrapServer), zap.String("topic", topic))
	//	return nil, errors.New("InitSaramaAsyncClient kafka server has no topic")
	//}
	client, err := sarama.NewClient([]string{bootstrapServer}, config)
	if err != nil {
		log.Error("InitSaramaAsyncClient NewClient error", zap.String("kafka_bootstrap_server", bootstrapServer), zap.Error(err))
		return nil, err
	}
	defer client.Close()

	if len(topic) > 0 {
		topics, err := client.Topics()
		if err != nil {
			log.Error("InitSaramaAsyncClient Topics error", zap.String("kafka_bootstrap_server", bootstrapServer), zap.String("kafka_topic", topic), zap.Error(err))
			return nil, err
		}
		var topicExists bool
		for _, remoteTopic := range topics {
			if remoteTopic == topic {
				topicExists = true
			}
		}
		if !topicExists {
			log.Error("InitSaramaAsyncClient Topic Not Found", zap.String("kafka_bootstrap_server", bootstrapServer), zap.String("kafka_topic", topic))
			return nil, errors.New("topic does not exist")
		}
	}

	config.ClientID = clientID

	//config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Idempotent = true
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
		log.Error("InitSaramaAsyncClient NewAsyncProducer error", zap.Error(err))
		return nil, err
	}

	return &SaramaAsyncClient{producer, topic}, nil
}

func (sara *SaramaAsyncClient) GetProducer() sarama.AsyncProducer {
	return sara.producer
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
		log.Info("SaramaAsyncClient nil producer")
		return
	}
	producer := asyncClient.producer
	go func() {
		for {
			select {
			case msg := <-producer.Successes():
				log.Info("ConsumeAsyncProducer success to send msg", zap.Any("msg", msg))
			case err := <-producer.Errors():
				if err != nil {
					if errors.Is(err.Err, sarama.ErrUnknownTopicOrPartition) {
						log.Fatal("ConsumeAsyncProducer topic does not exists", zap.Error(err))
					}
					log.Error("SaramaAsyncClient GetAsyncProducer error", zap.Error(err))
				}
			}
		}
	}()
}

type SaramaConsumerClient struct {
	consumerGroup sarama.ConsumerGroup
	handler       MessageHandler
	topic         []string
	ready         chan struct{} // 无缓冲空结构体通道
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	readyOnce     sync.Once
}

func InitSaramaConsumerClient(bootstrapServer, clientID, groupID string, topic []string, handler MessageHandler) (*SaramaConsumerClient, error) {
	if bootstrapServer == "" || len(topic) <= 0 || handler == nil {
		log.Error("InitSaramaConsumerClient invalid parameters", zap.String("bootstrap_server", bootstrapServer), zap.Strings("topic", topic), zap.Any("handler_func", handler))
		return nil, errors.New("invalid parameters")
	}
	config := sarama.NewConfig()
	if config == nil {
		return nil, errors.New("sarama NewConfig return nil")
	}
	config.ClientID = clientID
	rangeStrategy := sarama.NewBalanceStrategyRange()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{rangeStrategy} // 工厂函数创建策略
	config.Consumer.Offsets.Initial = sarama.OffsetOldest                                     // 从最早的消息开始消费
	config.Consumer.Offsets.AutoCommit.Enable = true                                          // 启用自动提交偏移量
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second                             // 自动提交间隔

	config.Net.SASL.Enable = false
	config.Net.TLS.Enable = false

	config.Consumer.Return.Errors = true

	brokers := strings.Split(bootstrapServer, ",")
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Error("InitSaramaConsumerClient NewConsumerGroup error", zap.Error(err))
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())

	return &SaramaConsumerClient{
		consumerGroup: consumerGroup,
		topic:         topic,
		handler:       handler,
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

func (sara *SaramaConsumerClient) markReady() {
	sara.readyOnce.Do(func() {
		close(sara.ready)
		log.Info("SaramaConsumerClient ready to mark ready")
	})
}

func (sara *SaramaConsumerClient) Start() error {
	sara.ready = make(chan struct{})
	sara.wg.Add(1)
	go func() {
		defer sara.wg.Done()
		handler := &SaramaHandler{
			handler:   sara.handler,
			readyFunc: sara.markReady,
		}

		for {
			select {
			case <-sara.ctx.Done():
				return
			default:
				err := sara.consumerGroup.Consume(sara.ctx, sara.topic, handler)
				if err != nil {
					log.Error("SaramaConsumerClient ConsumerGroup Consume error", zap.Error(err))
					time.Sleep(5 * time.Second)
				}
			}
		}
	}()

	<-sara.ready
	log.Info("SaramaConsumerClient start success")
	return nil
}

func (sara *SaramaConsumerClient) Close() error {
	sara.cancel()  // 发送取消信号
	sara.wg.Wait() // 等待consumer退出
	if err := sara.consumerGroup.Close(); err != nil {
		log.Error("SaramaConsumerClient Close() error", zap.Error(err))
	}
	return nil
}

func (sara *SaramaConsumerClient) WaitUntilReady(timeout time.Duration) error {
	select {
	case <-sara.ready:
		return nil
	case <-time.After(timeout):
		return errors.New("timeout waiting for consumer ready")
	}
}

func StartSaramaKafka() {
	ConsumeAsyncProducer(KafkaProducerClient)
}

func CloseSaramaKafka() {
	if KafkaProducerClient != nil {
		KafkaProducerClient.Close()
		KafkaProducerClient = nil
	}
}
