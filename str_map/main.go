/*
测试：map 取不存在的值的时候，不会 panic
*/

package main

import "fmt"

func main() {
	urlPath := "/api"
	var oAuthPathMapper = map[string]bool{
		"/interstitial/":            true,
		"/api/frodo/feed_ad":        true,
		"/api/frodo/group_topic_ad": true,
	}
	val := oAuthPathMapper[urlPath]
	fmt.Println(val)
}
