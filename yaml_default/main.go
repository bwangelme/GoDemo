/*
本程序为了测试 yaml 结构体如何设置默认值

重写结构体的 UnmarshalYAML 接口

既可以在结构体创建的时候指定默认值
也可以在结构体 Decode 完成之后，设置默认值
*/
package main

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type UseService struct {
	App        string `yaml:"app,omitempty"`
	TargetPath string `yaml:"target_path"`
	Type       string `yaml:"type,omitempty"`
	Version    string `yaml:"version"`
}

func (u *UseService) UnmarshalYAML(value *yaml.Node) error {
	type rawUseService UseService
	ru := rawUseService{
		Type: "thrift",
	}
	err := value.Decode(&ru)
	if err != nil {
		return errors.Wrap(err, "rawUseService: failed to unmarshal")
	}
	*u = UseService(ru)
	return nil
}

type Service struct {
	Name      string `yaml:"name"`
	Interface string `yaml:"interface"`
	Type      string `yaml:"type"`
	Handler   string `yaml:"handler"`
}

func (s *Service) UnmarshalYAML(value *yaml.Node) error {
	type tmp Service
	rs := tmp{}
	err := value.Decode(&rs)
	if err != nil {
		return errors.Wrap(err, "rawUseService: failed to unmarshal")
	}

	if rs.Name == "" {
		rs.Name = strings.Split(rs.Interface, ".")[0]
	}

	*s = Service(rs)
	return nil
}

type AppConfig struct {
	Application string        `yaml:"application"`
	Runtime     string        `yaml:"runtime"`
	Services    []*Service    `yaml:"services"`
	UseServices []*UseService `yaml:"use_services"`
}

func main() {
	config := &AppConfig{}
	content, _ := os.ReadFile("test.yaml")
	err := yaml.Unmarshal(content, config)
	if err != nil {
		log.Fatalln(err)
	}
	for _, svc := range config.Services {
		fmt.Println(svc.Name, svc.Interface)
	}
	for _, svc := range config.UseServices {
		fmt.Println(svc.App, svc.Type)
	}
}
