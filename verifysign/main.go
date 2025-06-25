package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// 生成RSA密钥对
func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatal(err)
	}
	return privateKey, &privateKey.PublicKey
}

// 将私钥保存为PEM文件
func savePrivateKeyToFile(privateKey *rsa.PrivateKey, filename string) {
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = pem.Encode(file, privPem)
	if err != nil {
		log.Fatal(err)
	}
}

// 将公钥保存为PEM文件
func savePublicKeyToFile(publicKey *rsa.PublicKey, filename string) {
	pubBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.Fatal(err)
	}

	pubPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	}

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	err = pem.Encode(file, pubPem)
	if err != nil {
		log.Fatal(err)
	}
}

// 使用RSA私钥对消息进行签名
func signMessage(privateKey *rsa.PrivateKey, message []byte) ([]byte, error) {
	hashed := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// 使用RSA公钥验证签名
func verifySignature(publicKey *rsa.PublicKey, message []byte, signature []byte) error {
	hashed := sha256.Sum256(message)
	err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	// 生成2048位的RSA密钥对
	privateKey, publicKey := generateKeyPair(2048)

	// 保存密钥到文件（可选）
	//savePrivateKeyToFile(privateKey, "private.pem")
	//savePublicKeyToFile(publicKey, "public.pem")

	// 要签名的消息
	message := []byte("这是要签名的消息")

	// 签名
	signature, err := signMessage(privateKey, message)
	if err != nil {
		log.Fatal("签名失败:", err)
	}
	fmt.Printf("签名结果: %x\n", signature)

	// 验证签名
	err = verifySignature(publicKey, message, signature)
	if err != nil {
		log.Fatal("验签失败:", err)
	} else {
		fmt.Println("验签成功!")
	}

	// 测试篡改后的消息
	tamperedMessage := []byte("这是被篡改的消息")
	err = verifySignature(publicKey, tamperedMessage, signature)
	if err != nil {
		fmt.Println("对篡改消息验签失败（符合预期）:", err)
	} else {
		log.Fatal("对篡改消息验签成功（不符合预期）")
	}
}
