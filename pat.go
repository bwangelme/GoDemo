package main

import (
	"fmt"
	"regexp"
)

func main() {
	var pattern *regexp.Regexp
	pattern, _ = regexp.Compile("dae\\..*\\.waylife\\..*")
	res := pattern.MatchString("dae.hello_dae_go.waylife.dae.api.beanstalk.web.default.put")
	fmt.Println(res)
}
