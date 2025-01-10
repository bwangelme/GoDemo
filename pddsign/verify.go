package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
)

// 公钥（用于验证签名）
var publicKeyStr = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0VkbMhQNerL6YcyBU9Femu8ZSdgtW0r/3FoBzOErWb1pBsGPXKZaXS3dOjbAn5BWPD7YbBU2EqNmpkTtiYsQzFz85wQGiHY46CIaTYkFA6Wo0VGyyKDnxQLwXpYNS2/xMoEUhBu2iztt59lq/eX8INRt3Il08ZMIxyl1KS4Q/gQIDAQAB
-----END PUBLIC KEY-----`

// 签名（base64编码后的签名，假设是来自某个 RSA 签名过程）
var signatureBase64 = "k/SlzoUNqln3/CeXhl38Pnzx6CcP/sINSxRSEiEXCHzjBTn9uG1iUIUwB+uBjQ2KDFoh/dHk4w1B5BgdBvQKt9tpRmF86WnPyVAFdCIwGukbRKymBw/gDhnzLeTG/eoDiRbrE+BWwipY97i3150VWe8/EsWryfSNXrLnh77/WBQ="

// 原始消息（待验证的消息）
var message = []byte("amount=1&chargeNo=13971074920&chargeTag=1&mctNo=795601066&notifyUrl=http://gw-api.pinduoduo.com/api/router&outOrderNo=3012-1736232078987&prodName=作业帮vip-1月卡&prodNo=101321033&signType=rsa")

func main() {
	// Step 1: 解析公钥
	pub, err := parseRSAPublicKey(publicKeyStr)
	if err != nil {
		log.Fatalf("Error parsing public key: %v", err)
	}
	fmt.Println("=======")
	fmt.Println(pub.N, pub.E)
	fmt.Println("=======")

	// Step 2: 计算消息的哈希值（这里使用MD5）
	hashed := md5.Sum(message)

	fmt.Println(hex.EncodeToString(hashed[:]), len(hashed))

	// Step 3: 解码签名
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		log.Fatalf("Error decoding signature: %v", err)
	}

	// Step 4: 使用公钥验证签名
	err = rsa.VerifyPKCS1v15(pub, crypto.MD5, hashed[:], signature)
	if err != nil {
		log.Fatalf("Signature verification failed: %v", err)
	}

	fmt.Println("Signature verified successfully")
}

// 解析RSA公钥
func parseRSAPublicKey(pubKeyStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubKeyStr))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}

	// 将解析的公钥断言为 rsa.PublicKey 类型
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not of type RSA")
	}
	return rsaPub, nil
}
