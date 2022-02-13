module bwdemo

go 1.14

require (
	github.com/certifi/gocertifi v0.0.0-20210507211836-431795d63e8d // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/raven-go v0.2.0
	github.com/gin-gonic/gin v1.7.0
	github.com/go-playground/validator/v10 v10.3.0 // indirect
	github.com/go-redis/redis/v8 v8.11.0
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/olivere/elastic/v7 v7.0.15
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/tidwall/gjson v1.9.3
)

replace (
	github.com/getsentry/raven-go => ../raven-go
	github.com/gin-gonic/gin => ../gin
	github.com/go-redis/redis/v8 => ../redis
)

// replace github.com/go-redis/redis/v8 => /home/xuyundong/Github/Golang/redis-1
