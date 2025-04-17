package config

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type NacosConfig struct {
	ServerAddr string
	Port       uint64
	Namespace  string
	Group      string
	DataId     string
}

func NewNacosClient(config NacosConfig) (clients.IConfigClient, error) {
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         config.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}

	// 创建serverConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: config.ServerAddr,
			Port:   config.Port,
		},
	}

	// 创建动态配置客户端
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create nacos client: %v", err)
	}

	return configClient, nil
}

func GetConfig(configClient clients.IConfigClient, group, dataId string) (string, error) {
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get config: %v", err)
	}
	return content, nil
}
