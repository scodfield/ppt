package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"ppt/logger"
)

type ListenConfig interface {
	OnChange(space, group, dataID, data string)
}

type ClientConfig struct {
	host    string
	port    int
	spaceID string
	client  config_client.IConfigClient
}

func (c *ClientConfig) Close() {
	c.client.CloseClient()
}

func (c *ClientConfig) GetConfig(dataID, group string) (string, error) {
	return c.client.GetConfig(vo.ConfigParam{DataId: dataID, Group: group})
}

func (c *ClientConfig) PublishConfig(dataID, group string, content string) error {
	_, err := c.client.PublishConfig(vo.ConfigParam{DataId: dataID, Group: group, Content: content})
	return err
}

func (c *ClientConfig) ListenConfig(dataID, group string, listenCfgI ListenConfig) error {
	return c.client.ListenConfig(vo.ConfigParam{DataId: dataID, Group: group, OnChange: listenCfgI.OnChange})
}

func (c *ClientConfig) CancelConfigListen(dataID, group string) error {
	return c.client.CancelListenConfig(vo.ConfigParam{DataId: dataID, Group: group})
}

func (c *ClientConfig) DeleteConfig(dataID, group string) error {
	_, err := c.client.DeleteConfig(vo.ConfigParam{DataId: dataID, Group: group})
	return err
}

func NewConfigClient(nacosServe []constant.ServerConfig, spaceID string) (*ClientConfig, error) {
	clientConfig := constant.NewClientConfig(
		constant.WithNamespaceId(spaceID),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogLevel("info"),
		constant.WithDisableUseSnapShot(true))
	client, err := clients.NewConfigClient(vo.NacosClientParam{ClientConfig: clientConfig, ServerConfigs: nacosServe})
	if err != nil {
		logger.Error("Nacos NewConfigClient clients.NewConfigClient error", zap.Any("nacos_server", nacosServe), zap.String("space_id", spaceID), zap.Error(err))
		return nil, err
	}
	return &ClientConfig{
		host:    "",
		port:    0,
		spaceID: spaceID,
		client:  client,
	}, nil
}
