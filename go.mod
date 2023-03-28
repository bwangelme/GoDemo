module bwdemo

go 1.14

require (
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/raven-go v0.2.0
	github.com/gin-gonic/gin v1.7.0
	github.com/go-redis/redis/v8 v8.11.0
	github.com/go-zookeeper/zk v1.0.3
	github.com/golang/protobuf v1.5.2
	github.com/olivere/elastic/v7 v7.0.15
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.14.0
	github.com/prometheus/client_model v0.3.0
	github.com/prometheus/common v0.37.0
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tidwall/gjson v1.9.3
	google.golang.org/protobuf v1.28.1
)

//replace (
//	github.com/getsentry/raven-go => ../raven-go
//	github.com/gin-gonic/gin => ../gin
//	github.com/go-redis/redis/v8 => ../redis
//)

// replace github.com/go-redis/redis/v8 => /home/xuyundong/Github/Golang/redis-1
