package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type RedisConfig struct {
	Dbs           map[string]string `json:"dbs"`
	Host          string            `json:"host,omitempty"`
	Port          uint16            `json:"port"`
	SilentTime    uint16            `json:"silent_time"`
	Slave         string            `json:"slave,omitempty"`
	SocketTimeout uint16            `json:"socket_timeout"`
}

type Config struct {
	redisCfgs map[string]*RedisConfig
	Note      string `json:"note"`
	Options   struct {
		SentryDSN string `json:"sentry_dsn"`
	} `json:"options"`
}

var config = []byte(`{
	"dba": {
		"dbs": {
			"1": "DEFAULT"
		},
		"host": "dba-redis-m",
		"port": 6403,
		"silent_time": 5,
		"socket_timeout": 3
	},
	"sa": {
		"dbs": {
			"1": "DEFAULT-sa"
		},
		"host": "sa-redis-m",
		"port": 6403,
		"silent_time": 5,
		"socket_timeout": 3
	},
	"note": "not used, just as doc",
	"options": {
		"sentry_dsn": "xxxx"
	}
}`)

func Loads(data []byte) (*Config, error) {
	var cfg = &Config{}
	var rdsCfgs = make(map[string]*RedisConfig)
	err := json.Unmarshal(config, &cfg)
	if err != nil {
		return nil, err
	}

	var productionConfig = struct {
		X map[string]json.RawMessage `json:"-"`
	}{}

	err = json.Unmarshal(config, &productionConfig.X)
	if err != nil {
		return nil, err
	}

	for k, v := range productionConfig.X {
		var cfg = &RedisConfig{}
		err = json.Unmarshal(v, &cfg)
		if err != nil {
			continue
		}
		if cfg.Host == "" {
			continue
		}
		rdsCfgs[k] = cfg
	}

	cfg.redisCfgs = rdsCfgs
	return cfg, nil
}

func Dumps(c *Config) ([]byte, error) {
	var res = make(map[string]interface{})

	res["options"] = c.Options
	res["note"] = c.Note
	for key, cfg := range c.redisCfgs {
		res[key] = cfg
	}

	return json.Marshal(res)
}

func main() {
	cfg, err := Loads(config)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Note:", cfg.Note)
	fmt.Println("DSN:", cfg.Options.SentryDSN)
	for key, rds := range cfg.redisCfgs {
		fmt.Println(key, rds.Host)
	}

	cfg.redisCfgs["sa"].Host = "localhsot"
	data, err := Dumps(cfg)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("==================")
	fmt.Println(string(data))
}
