package main

import (
	"fmt"
	"github.com/VictoriaMetrics/metricsql"
	"strings"
)

func main() {
	rawExpr := `(topk(50,(sum(dae_app_nworkers{}) by (app,role,task))) - topk(50,(sum(dae_app_nworkers{} offset 1d) by (app,role,task))))/topk(50,(sum(dae_app_nworkers{} offset 1d) by (app,role,task)))`
	expr, err := metricsql.Parse(rawExpr)
	if err != nil {
		fmt.Println("parsed error", err)
		return
	}
	fmt.Printf("parsed expr: %s\n", expr.AppendString(nil))
	if strings.Contains(rawExpr, "/") {
		ae := expr.(*metricsql.BinaryOpExpr)
		fmt.Println(ae.Op, ae.Left, ae.Right)
	}
}
