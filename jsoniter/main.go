package main

/*
* 针对 JSON  body 生成签名

JSON body 是一个多层嵌套的结构体

生成的签名格式是 key1=v1&key2=v2，key1 和 key2 需要按字母序排序

如果 v 是一个 json 对象，那么将这个 json 对象也解析成 key=v 的格式
注意：如果 v 的值也是一个对象，此时不用再解析了，直接贴上原始内容即可(要保留原始字符串的顺序)

测试用例详见 main 函数
*/

import (
	"fmt"
	"sort"
	"strings"

	"github.com/json-iterator/go"
)

func getSortedMapKeys(v map[string]jsoniter.Any) []string {
	keys := make([]string, 0)
	for key, _ := range v {
		if key == "sign" {
			continue
		}

		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func strDropWhite(strV string) string {
	strV = strings.Replace(strV, "\n", "", -1)
	strV = strings.Replace(strV, "\r", "", -1)
	strV = strings.Replace(strV, " ", "", -1)
	strV = strings.Replace(strV, "\t", "", -1)
	return strV
}

func BuildSign(jsonData string) (string, error) {
	// 使用 map[string]interface{} 和 jsoniter.Any 解析 JSON
	var req map[string]jsoniter.Any
	err := jsoniter.Unmarshal([]byte(jsonData), &req)
	if err != nil {
		return "", err
	}

	signItems := make([]string, 0)
	sortedKeys := getSortedMapKeys(req)
	for _, key := range sortedKeys {
		v := req[key]
		switch v.ValueType() {
		case jsoniter.StringValue, jsoniter.NumberValue:
			signItems = append(signItems, fmt.Sprintf("%s=%s", key, v.ToString()))
		case jsoniter.ObjectValue:
			var req map[string]jsoniter.Any
			err := jsoniter.Unmarshal([]byte(v.ToString()), &req)
			if err != nil {
				return "", err
			}
			innerSignItems := make([]string, 0)
			sortedKeys := getSortedMapKeys(req)
			for _, key := range sortedKeys {
				v := req[key]
				var strV string
				switch v.ValueType() {
				case jsoniter.ArrayValue, jsoniter.ObjectValue:
					strV = strDropWhite(v.ToString())
				default:
					strV = v.ToString()
				}
				innerSignItems = append(innerSignItems, fmt.Sprintf("%s=%s", key, strV))
			}
			innerSign := strings.Join(innerSignItems, "&")
			signItems = append(signItems, fmt.Sprintf("%s=%s", key, innerSign))
		}
	}

	sign := strings.Join(signItems, "&")
	return sign, nil
}

func main() {
	jsonData := `{
		"req_biz_content": {
			"brand_code": "homework",
			"total_price": "10.00",
			"other_info": {},
			"prods_info": [
				{
					"price": "14.00",
					"rcg_type_code": "month",
					"service_code": "101381029"
				}
			],
			"mobile": "13364001234",
			"order_id": "1742286194",
			"timestamp": 1742286194
		},
		"timestamp": 1742286198,
		"sign": "MEYCIQDO3xDxpsylHl0rWrKYZgUlufTytTGZzlUXiNDKVMnKigIhAJuCX2gq5+OJOHwOGNERqhSBwqKdmfJq/lFk35NXk7Yd"
	}`

	sign, err := BuildSign(jsonData)
	fmt.Println(sign)
	// 期望的 sign 结果是
	// req_biz_content=brand_code=homework&mobile=13364001234&order_id=1742286194&other_info={}&prods_info=[{"price":"14.00","rcg_type_code":"month","service_code":"101381029"}]&timestamp=1742286194&total_price=10.00&timestamp=1742286198
	fmt.Println(err)
}
