package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func printPKCS1RSAKey() {
	// 生成 RSA 密钥对
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key:", err)
		return
	}

	// 提取公钥
	pubKey := &privKey.PublicKey

	// 将公钥编码为 PKCS#1 格式
	// 将公钥数据写入文件
	pubKeyBytes := x509.MarshalPKCS1PublicKey(pubKey)
	data := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	})
	fmt.Println("Public key pem")
	fmt.Println(string(data))

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	data = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	fmt.Println("Private key pem")
	fmt.Println(string(data))
}
