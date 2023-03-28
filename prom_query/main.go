package main

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"log"
	"time"
)

func main() {
	config := api.Config{
		Address: "http://localhost:9090",
	}
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalln("init client error", err)
	}

	promql := `
sum by (le) (
	rate(
		demo_api_request_duration_seconds_bucket{method="GET", path="/api/bar", status="200"}[5m]
	)
)
`

	v1api := v1.NewAPI(client)
	value, warn, err := v1api.Query(context.Background(), promql, time.Now())
	if err != nil {
		log.Fatalln("query failed", err)
	}

	if len(warn) != 0 {
		log.Print("query warning", warn)
	}

	switch value.Type() {
	case model.ValVector:
		v, _ := value.(model.Vector)
		fmt.Println(v[0].Metric, v[0].Value, v[0].Timestamp)
	}
}
