package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()

	PRIVATEKEY_FILE = "private_key.pem"
	PUBLICKEY_FILE  = "public_key.pem"
)

func init() {
	dir, _ := os.Getwd()

	PRIVATEKEY_FILE = filepath.Join(dir, "rsakey", PRIVATEKEY_FILE)
	PUBLICKEY_FILE = filepath.Join(dir, "rsakey", PUBLICKEY_FILE)
}

func ParseRSAPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	// 解码 PEM 数据
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 数据")
	}
	fmt.Println(block.Type)
	var privateKey interface{}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return privateKey.(*rsa.PrivateKey), nil
}

// 解析 PEM 格式的 RSA 公钥
func ParseRSAPublicKey(data []byte) (*rsa.PublicKey, error) {
	// 解码 PEM 数据
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 数据")
	}

	// 解析证书
	fmt.Println(block.Type)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("无法解析证书: %v", err)
	}

	return publicKey.(*rsa.PublicKey), nil
}

// 生成 RSA 公私钥对
func generateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Warnf("generateRSAKeyPair failed err=%v", err)
		return nil, nil, fmt.Errorf("生成 RSA 密钥对失败: %v", err)
	}
	return privKey, &privKey.PublicKey, nil
}

// 将公钥转为 PEM 格式的字符串
func publicKeyToPEM(pubKey *rsa.PublicKey) error {
	// 将公钥转换成 pkix 格式
	pubASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		logger.Warnf("无法将公钥转换为 ASN.1 格式: %v", err)
	}
	// 创建一个文件来保存 PEM 编码的私钥
	fd, err := os.Create(PUBLICKEY_FILE)
	if err != nil {
		logger.Warnf("create file failed, name=%v err=%v", PUBLICKEY_FILE, err)
		return err
	}
	defer fd.Close()

	err = pem.Encode(fd, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	if err != nil {
		logger.Warnf("publick key pem encode failed, err=%v", err)
		return err
	}
	logger.Infof("generate public key into %v", PUBLICKEY_FILE)
	return nil
}

// 将私钥转为 PEM 格式的字符串
func privateKeyToPEM(privKey *rsa.PrivateKey) error {
	// 将私钥转换成 pkcs#8 格式
	privASN1, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		logger.Warnf("MarshalPKCS8PrivateKey failed, err=%v", err)
		return err
	}
	// 创建一个文件来保存 PEM 编码的私钥
	fd, err := os.Create(PRIVATEKEY_FILE)
	if err != nil {
		logger.Warnf("create file failed, name=%v err=%v", PRIVATEKEY_FILE, err)
		return err
	}
	defer fd.Close()

	err = pem.Encode(fd, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privASN1,
	})
	if err != nil {
		logger.Warnf("pem encode failed err=%v", err)
		return err
	}
	logger.Infof("generate private key into %v", PRIVATEKEY_FILE)
	return nil
}

func generateKey() {
	priKey, pubKey, err := generateRSAKeyPair()
	if err != nil {
		return
	}
	err = privateKeyToPEM(priKey)
	if err != nil {
		return
	}
	err = publicKeyToPEM(pubKey)
	if err != nil {
		return
	}
}

func main() {
	generateKey()
}
