package main

import (
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Host        string
	Port        uint64
	NamespaceId string
	Group       string
	ConfigName  string
}

type Nacos struct {
	client config_client.IConfigClient
	config Config
}

// validateConfig
// @description: validate config
// @param config
// @return error
func validateConfig(config *Config) error {
	if config.Host == "" {
		return errors.New("host is empty")
	}
	if config.Port <= 0 {
		return errors.New("port is empty")
	}
	if config.NamespaceId == "" {
		config.NamespaceId = "public"
	}
	if config.Group == "" {
		config.Group = "DEFAULT_GROUP"
	}
	if config.ConfigName == "" {
		return errors.New("config name is empty")
	}
	return nil
}

// Provider
// @description: create nacos provider
// @param config
// @return *Nacos
func Provider(config Config) *Nacos {

	// validate config
	if err := validateConfig(&config); err != nil {
		return nil
	}

	// create nacos config client
	clientConf := constant.ClientConfig{
		NamespaceId:         config.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogLevel:            "error",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      config.Host,
			ContextPath: "/nacos",
			Port:        config.Port,
		},
	}

	nclient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  &clientConf,
		ServerConfigs: serverConfigs,
	})

	// create client failed
	if err != nil {
		return nil
	}

	return &Nacos{nclient, config}
}

// Read
// @description: read config from remote
// @receiver n
// @return map[string]interface{}
// @return error
func (n Nacos) Read() (result map[string]interface{}, err error) {
	if n.client == nil {
		return nil, errors.New("nacos client is nil")
	}

	data, err := n.client.GetConfig(vo.ConfigParam{
		DataId: n.config.ConfigName,
		Group:  n.config.Group,
	})

	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ReadBytes
// @description: read config from remote
// @receiver n
// @return result
// @return err
func (n Nacos) ReadBytes() (result []byte, err error) {
	if n.client == nil {
		return nil, errors.New("nacos client is nil")
	}

	data, err := n.client.GetConfig(vo.ConfigParam{
		DataId: n.config.ConfigName,
		Group:  n.config.Group,
	})

	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}

// Watch
// @description: watch config change
// @receiver n
// @param cb
// @return error
func (n Nacos) Watch(cb func(event interface{}, err error)) error {
	if n.client == nil {
		return errors.New("nacos client is nil")
	}

	err := n.client.ListenConfig(vo.ConfigParam{
		DataId: n.config.ConfigName,
		Group:  n.config.Group,
		OnChange: func(namespace, group, dataId, data string) {
			cb(data, nil)
		},
	})

	return err
}
