package main

import (
	"fmt"
	"net/url"
	"path"
)

var fileURL = `https://charge.zuoyebang.cc//album_f26f0ad023facfcb69df655bd08d56a6.mp4?authorization=bce-auth-v1%2F80054b2251d24cb98eda994c93426d7d%2F2021-05-31T03%3A34%3A53Z%2F1576800000%2Fhost%2F63ce7802da49bbce6d3b907a9941de71c5518e39d002a9cf86efcf629c32b497`

func getFileExtension(urlStr string) (string, error) {
	// 解析 URL
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// 获取路径部分
	ext := path.Ext(parsedURL.Path)

	// 返回文件后缀
	return ext, nil
}

func main() {
	fmt.Println(getFileExtension(fileURL))
}
