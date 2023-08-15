/*
本程序为了测试 yaml 结构体如何设置默认值

重写结构体的 UnmarshalYAML 接口

既可以在结构体创建的时候指定默认值
也可以在结构体 Decode 完成之后，设置默认值

## 参考链接

- https://github.com/go-yaml/yaml/issues/165#issuecomment-727092641
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
	App     string `yaml:"app,omitempty"`
	Type    string `yaml:"type,omitempty"`
	Version string `yaml:"version"`
}

// UnmarshalYAML
//
//	此程序演示了，在 Decode 之前设置 Type 字段的默认值为 thrift
func (u *UseService) UnmarshalYAML(value *yaml.Node) error {
	type tmp UseService
	ru := tmp{
		Type: "thrift",
	}
	err := value.Decode(&ru)
	if err != nil {
		return errors.Wrap(err, "UseService: failed to unmarshal")
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

// UnmarshalYAML
//
//	此程序演示了更复杂的情况, Name 字段的默认值是根据 Interface 的值来设置的
func (s *Service) UnmarshalYAML(value *yaml.Node) error {
	type tmp Service
	rs := tmp{}
	err := value.Decode(&rs)
	if err != nil {
		return errors.Wrap(err, "Service: failed to unmarshal")
	}

	if !strings.Contains(rs.Interface, ".") {
		return errors.New("Service: invalid Interface value")
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

	// 输出
	//linaewen linaewen.fdoulist
	for _, svc := range config.Services {
		fmt.Println(svc.Name, svc.Interface)
	}

	// 输出
	//fm grpc
	//music thrift
	//pony thrift
	for _, svc := range config.UseServices {
		fmt.Println(svc.App, svc.Type)
	}
}
