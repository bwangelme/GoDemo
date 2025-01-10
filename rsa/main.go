package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

// 生成 RSA 公私钥对
func generateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("生成 RSA 密钥对失败: %v", err)
	}
	return privKey, &privKey.PublicKey, nil
}

// 将公钥转为 PEM 格式的字符串
func publicKeyToPEM(pubKey *rsa.PublicKey) string {
	pubASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		log.Fatalf("无法将公钥转换为 ASN.1 格式: %v", err)
	}

	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	return string(pubPEM)
}

// 将私钥转为 PEM 格式的字符串
func privateKeyToPEM(privKey *rsa.PrivateKey) string {
	privASN1 := x509.MarshalPKCS1PrivateKey(privKey)
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privASN1,
	})
	return string(privPEM)
}

// 使用公钥进行加密
func encryptWithPublicKey(plainText []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	cipherText, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, plainText, nil)
	if err != nil {
		return nil, fmt.Errorf("加密失败: %v", err)
	}
	return cipherText, nil
}

// 使用私钥进行解密
func decryptWithPrivateKey(cipherText []byte, privKey *rsa.PrivateKey) ([]byte, error) {
	plainText, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, cipherText, nil)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v", err)
	}
	return plainText, nil
}

func demo() {
	// 生成 RSA 公私钥对
	privKey, pubKey, err := generateRSAKeyPair()
	if err != nil {
		log.Fatalf("生成 RSA 密钥对失败: %v", err)
	}

	// 将公私钥转换为 PEM 格式的字符串
	pubKeyPEM := publicKeyToPEM(pubKey)
	privKeyPEM := privateKeyToPEM(privKey)

	// 打印 PEM 格式的公钥和私钥
	fmt.Println("公钥 PEM 格式:")
	fmt.Println(pubKeyPEM)

	fmt.Println("私钥 PEM 格式:")
	fmt.Println(privKeyPEM)

	// 要加密的明文
	plainText := []byte("Hello, this is a secret message using RSA encryption!")

	// 使用公钥进行加密
	encryptedText, err := encryptWithPublicKey(plainText, pubKey)
	if err != nil {
		log.Fatalf("加密失败: %v", err)
	}
	fmt.Printf("加密后的数据: %x\n", encryptedText)

	// 使用私钥进行解密
	decryptedText, err := decryptWithPrivateKey(encryptedText, privKey)
	if err != nil {
		log.Fatalf("解密失败: %v", err)
	}
	fmt.Printf("解密后的数据: %s\n", decryptedText)
}
