package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
)

// 解密函数，使用公钥进行解密
func decryptByPublicKey(data string, sourceKey string) (string, error) {
	// 解码 Base64 公钥
	pubKeyBytes, err := base64.StdEncoding.DecodeString(sourceKey)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 public key: %v", err)
	}

	// 解析公钥
	pubKey, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %v", err)
	}

	// 将接口类型转换为 rsa.PublicKey 类型
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("invalid public key type")
	}
	fmt.Println("===")
	fmt.Println(rsaPubKey.N)
	fmt.Println("===")

	// 解密数据
	return decryptWithRSA(data, rsaPubKey)
}

// 使用 RSA 公钥解密数据
func decryptWithRSA(data string, pubKey *rsa.PublicKey) (string, error) {
	// 解码加密数据
	//encryptedData, err := base64.StdEncoding.DecodeString(data)
	//if err != nil {
	//	return "", fmt.Errorf("failed to decode base64 encrypted data: %v", err)
	//}

	// 执行解密
	//decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, pubKey, encryptedData)
	//if err != nil {
	//	return "", fmt.Errorf("failed to decrypt data: %v", err)
	//}

	//return string(decryptedData), nil
	return "", nil
}

func demo() {
	// 假设输入的 Base64 公钥和加密数据
	sourceKey := "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC0VkbMhQNerL6YcyBU9Femu8ZSdgtW0r/3FoBzOErWb1pBsGPXKZaXS3dOjbAn5BWPD7YbBU2EqNmpkTtiYsQzFz85wQGiHY46CIaTYkFA6Wo0VGyyKDnxQLwXpYNS2/xMoEUhBu2iztt59lq/eX8INRt3Il08ZMIxyl1KS4Q/gQIDAQAB" // 这里替换成你的实际 Base64 公钥
	encryptedData := "k/SlzoUNqln3/CeXhl38Pnzx6CcP/sINSxRSEiEXCHzjBTn9uG1iUIUwB+uBjQ2KDFoh/dHk4w1B5BgdBvQKt9tpRmF86WnPyVAFdCIwGukbRKymBw/gDhnzLeTG/eoDiRbrE+BWwipY97i3150VWe8/EsWryfSNXrLnh77/WBQ="                                         // 这里替换成你的实际 Base64 加密数据

	// 解密操作
	decryptedData, err := decryptByPublicKey(encryptedData, sourceKey)
	if err != nil {
		log.Fatalf("Error decrypting data: %v", err)
	}

	// 输出解密后的数据
	fmt.Println("Decrypted Data:", decryptedData)
}
