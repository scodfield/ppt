package wrapper

import (
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"go.uber.org/zap"
	"ppt/logger"
	"ppt/nacos"
)

var ServerConf constant.ServerConfig

func init() {
	ServerConf = *constant.NewServerConfig(nacos.NacosUrl, nacos.NacosPort)
}

func GetNacosConfig(spaceID, group, dataID string) (string, error) {
	nacosClient, err := nacos.NewConfigClient([]constant.ServerConfig{ServerConf}, spaceID)
	if err != nil {
		return "", err
	}
	defer nacosClient.Close()

	data, err := nacosClient.GetConfig(dataID, group)
	if err != nil {
		logger.Error("GetNacosConfig nacos client GetConfig error", zap.String("space_id", spaceID), zap.String("group", group), zap.String("data_id", dataID), zap.Error(err))
		return "", err
	}

	return data, nil
}

func GetNacosDBConfig() (*DBConfig, error) {
	data, err := GetNacosConfig(nacos.NacosRegion, nacos.NacosDefaultGroup, nacos.NacosDataIDDBConfig)
	if err != nil {
		logger.Error("GetNacosDBConfig GetNacosDBConfig error", zap.Error(err))
		return nil, err
	}
	nacosDBConfig := &DBConfig{}
	err = json.Unmarshal([]byte(data), nacosDBConfig)
	if err != nil {
		logger.Error("GetNacosDBConfig Unmarshal NacosDBConfig error", zap.Error(err))
		return nil, err
	}
	return nacosDBConfig, nil
}
