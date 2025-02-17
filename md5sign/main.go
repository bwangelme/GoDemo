package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"
	"time"
)

const SignSecret = "secretkeyfortest"

func generateSignStr(params map[string]string) (string, string) {
	var keys []string
	time.Parse()
	for key, _ := range params {
		if "sign" == key {
			continue
		}

		keys = append(keys, key)
	}

	sort.Strings(keys)
	// 构造k1=v1&k2=v2...格式的字符串
	var buf bytes.Buffer
	buf.WriteString(SignSecret)
	for _, key := range keys {
		value := params[key]
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(value)
		buf.WriteString("&")
	}
	signArgStr := strings.TrimSuffix(buf.String(), "&")
	signArgStr = fmt.Sprintf("%v%v", signArgStr, SignSecret)
	hash := md5.Sum([]byte(signArgStr))
	md5String := hex.EncodeToString(hash[:])

	return signArgStr, md5String
}

func getReqParameterSignStr(params url.Values) (string, string) {
	d := make(map[string]string)
	for key, val := range params {
		if len(val) > 0 {
			d[key] = val[0]
		}
	}
	return generateSignStr(d)
}

func main() {
	urlStr := "https://apivip.zuoyebang.com/vippartner/pdd/dingxin_queryorder?pddOrderNo=217-89276&ts=2025-02-17 14:10:23"
	// 解析 URL
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		log.Fatal(err)
	}

	signStr, sign := getReqParameterSignStr(parsedUrl.Query())
	fmt.Println("计算签名的字符串是: ", signStr)
	fmt.Println("签名内容:", sign)

}
