# koanf-nacos-provider

## introduction
koanf的Nacos服务提供者。

## install
```go
go get github.com/cowardmrx/koanf-nacos-provider
```
## usage
```go
package main

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"
)

var (
	k    = koanf.New(".")
	conf = YourConfig{}
)

type YourConfig struct {
	// Fields
	Username string `yaml:"username"`
}

func main() {
	np,err := Provider(Config{
		Host:        "ip or domain",
		Port:        8848,
		NamespaceId: "",
		Group:       "",
		ConfigName:  "",
	})
	
	if err != nil {
		fmt.Println("get provider failed: ",err.Error())
		return
    }

	if err := k.Load(np, yaml.Parser()); err != nil {
		// error
		fmt.Println(fmt.Errorf("error: %w", err))
		return
	}

	if err := k.Unmarshal("", &conf); err != nil {
		// error
		fmt.Println(fmt.Errorf("error: %w", err))
		return
	}

	// watch
	if err := np.Watch(func(data interface{}, err error) {

		k = koanf.New(".")
		if err = k.Load(np, yaml.Parser()); err != nil {
			// error
			fmt.Println(fmt.Errorf("error: %w", err))
			return
		}
		err = k.Unmarshal("", &conf)
		if err != nil {
			// error
			fmt.Println(fmt.Errorf("error: %w", err))
			return
		}

		fmt.Println("2 ", conf.Username)

		// do something

	}); err != nil {
		// error
		fmt.Println(fmt.Errorf("error: %w", err))
		return
	}

}
```