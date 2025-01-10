package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	signStr := `amount=1&chargeNo=23344@qq.com&mctNo=748762514052&outOrderNo=240704343597453722238&prodName=腾讯视频超级影视 SVIP1 个月&prodNo=1&signType=rsa`
	hash := md5.Sum([]byte(signStr))
	fmt.Printf("%x\n", hash)
}
