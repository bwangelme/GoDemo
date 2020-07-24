module bwdemo

go 1.14

require (
	github.com/certifi/gocertifi v0.0.0-20200211180108-c7c1fbc02894 // indirect
	github.com/getsentry/raven-go v0.2.0
	github.com/gin-gonic/gin v1.6.3 // indirect
	github.com/olivere/elastic/v7 v7.0.15
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
)

replace (
    github.com/getsentry/raven-go => /Users/michaeltsui/Github/Golang/raven-go
	github.com/gin-gonic/gin => /Users/michaeltsui/Github/Golang/gin
)
