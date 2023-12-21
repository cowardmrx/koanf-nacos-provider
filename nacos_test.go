package koanf_nacos_provider

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"testing"
)

func ExampleProvider() {
	_ = Provider(Config{
		Host:        "",
		Port:        8848,
		NamespaceId: "",
		Group:       "",
		ConfigName:  "",
	})

	// output nacos provider
}

func ExampleNacos_Read() {
	np := Provider(Config{
		Host:        "ip or domain",
		Port:        8848,
		NamespaceId: "",
		Group:       "",
		ConfigName:  "",
	})

	result, err := np.Read()
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
	// output map[email:map[host:123 password:stmp123 port:123 username:hjahaha] password:hass123 username:user1331]
}

func ExampleNacos_ReadBytes() {

	np := Provider(Config{
		Host:        "ip or domain",
		Port:        8848,
		NamespaceId: "",
		Group:       "",
		ConfigName:  "",
	})

	result, err := np.ReadBytes()
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

	// output
	//username: "user1331"
	//password: "hass123"
	//
	//email:
	//  host: "123"
	//  port: 123
	//  username: "hjahaha"
	//  password: "stmp123"
}

func ExampleNacos_Watch() {
	k := koanf.New(".")

	np := Provider(Config{
		Host:        "ip or domain",
		Port:        8848,
		NamespaceId: "",
		Group:       "",
		ConfigName:  "",
	})

	_ = k.Load(np, yaml.Parser())

	err := np.Watch(func(event interface{}, err error) {
		// reload config
		k = koanf.New(".")
		k.Load(np, yaml.Parser())

		// unmarshal config
		// k.Unmarshal("email", &email)

		// TODO do something
	})
	if err != nil {
		panic(err)
	}

}

func TestNacos_Read(t *testing.T) {
	type fields struct {
		client config_client.IConfigClient
		config Config
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult map[string]interface{}
		wantErr    bool
	}{
		{
			fields: fields{
				config: Config{
					Host:        "ip or domain",
					Port:        8848,
					NamespaceId: "namespace id",
					Group:       "", // default value is DEFAULT_GROUP
					ConfigName:  "config name",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			n := Provider(tt.fields.config)

			gotResult, err := n.Read()
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("gotResult: %v", gotResult)
		})
	}
}

func TestNacos_ReadBytes(t *testing.T) {
	type fields struct {
		client config_client.IConfigClient
		config Config
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult []byte
		wantErr    bool
	}{
		{
			fields: fields{
				config: Config{
					Host:        "ip or domain",
					Port:        8848,
					NamespaceId: "namespace id",
					Group:       "", // default value is DEFAULT_GROUP
					ConfigName:  "config name",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			n := Provider(tt.fields.config)

			gotResult, err := n.ReadBytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("gotResult: %v", string(gotResult))
		})
	}
}

func TestNacos_Watch(t *testing.T) {
	type fields struct {
		client config_client.IConfigClient
		config Config
	}
	type args struct {
		cb func(event interface{}, err error)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			fields: fields{
				config: Config{
					Host:        "ip or domain",
					Port:        8848,
					NamespaceId: "namespace id",
					Group:       "", // default value is DEFAULT_GROUP
					ConfigName:  "config name",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Provider(tt.fields.config)

			if err := n.Watch(tt.args.cb); (err != nil) != tt.wantErr {
				t.Errorf("Watch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
