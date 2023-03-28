module bwdemo

go 1.14

require (
	github.com/bytedance/sonic v1.8.6 // indirect
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/raven-go v0.2.0
	github.com/gin-gonic/gin v1.9.0
	github.com/go-playground/validator/v10 v10.12.0 // indirect
	github.com/go-redis/redis/v8 v8.11.0
	github.com/go-zookeeper/zk v1.0.3
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/olivere/elastic/v7 v7.0.15
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/pelletier/go-toml/v2 v2.0.7 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/client_model v0.3.0
	github.com/prometheus/common v0.37.0
	github.com/sirupsen/logrus v1.8.1
	github.com/tidwall/gjson v1.9.3
	github.com/ugorji/go v1.1.7 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	google.golang.org/protobuf v1.30.0
)

//replace (
//	github.com/getsentry/raven-go => ../raven-go
//	github.com/gin-gonic/gin => ../gin
//	github.com/go-redis/redis/v8 => ../redis
//)

// replace github.com/go-redis/redis/v8 => /home/xuyundong/Github/Golang/redis-1
