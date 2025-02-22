package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()

	PRIVATEKEY_PEM_FILE = "private_key.pem"
	PUBLICKEY_PEM_FILE  = "public_key.pem"

	PRIVATEKEY_DER_FILE = "private_key.der"
	PUBLICKEY_DER_FILE  = "public_key.der"
)

func init() {
	dir, _ := os.Getwd()

	PRIVATEKEY_PEM_FILE = filepath.Join(dir, "rsakey", PRIVATEKEY_PEM_FILE)
	PUBLICKEY_PEM_FILE = filepath.Join(dir, "rsakey", PUBLICKEY_PEM_FILE)

	PRIVATEKEY_DER_FILE = filepath.Join(dir, "rsakey", PRIVATEKEY_DER_FILE)
	PUBLICKEY_DER_FILE = filepath.Join(dir, "rsakey", PUBLICKEY_DER_FILE)
}

func ParseRSADerPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	// rsa 私钥的二进制格式有两种
	// PKCS #8 和 PKCS #1
	// PKCS #8 是通用的私钥存储格式，除了存储 rsa 私钥，还可以存储多种形式的私钥，它的 pem 文件的开头是 -----BEGIN PRIVATE KEY-----
	// PKCS #1 是 RSA 私钥的专属格式，只能存储 RSA 私钥，它的 pem 文件的开头是 -----BEGIN RSA PRIVATE KEY-----
	var privateKey interface{}
	privateKey, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		logger.Warnf("parse pkcs7 private key failed err=%v", err)
		return nil, err
	}
	return privateKey.(*rsa.PrivateKey), nil
}

// ParseRSAPemPrivateKey
// 解析 pem 格式的 rsa 私钥
func ParseRSAPemPrivateKey(data []byte) (*rsa.PrivateKey, error) {
	// 解码 PEM 数据
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 数据")
	}
	logger.Infof("private pem key block type is %v\n", block.Type)
	return ParseRSADerPrivateKey(block.Bytes)
}

func ParseRSADerPublicKey(data []byte) (*rsa.PublicKey, error) {
	// RSA 公钥的 PKIX 格式是 X.509 公钥信息的一种编码格式，通常以 DER 或 PEM 格式存储。
	// 它是标准的公钥信息格式，符合 PKIX（Public Key Infrastructure X.509）规范。
	publicKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		logger.Warnf("parse  pkix public key failed, err=%v", err)
		return nil, fmt.Errorf("无法解析证书: %v", err)
	}

	return publicKey.(*rsa.PublicKey), nil
}

// ParseRSAPemPublicKey
// 解析 PEM 格式的 RSA 公钥
func ParseRSAPemPublicKey(data []byte) (*rsa.PublicKey, error) {
	// 解码 PEM 数据
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 数据")
	}

	logger.Infof("public key pem format %v", block.Type)
	return ParseRSADerPublicKey(block.Bytes)
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
func publicKeyIntoFile(pubKey *rsa.PublicKey) error {
	// 将公钥转换成 pkix 格式
	pubASN1, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		logger.Warnf("无法将公钥转换为 ASN.1 格式: %v", err)
		return err
	}

	err = os.WriteFile(PUBLICKEY_DER_FILE, pubASN1, 0600)
	if err != nil {
		logger.Warnf("write der file failed, file=%v", PUBLICKEY_DER_FILE)
		return err
	}
	logger.Infof("generate public key into %v", PUBLICKEY_DER_FILE)

	// 创建一个文件来保存 PEM 编码的私钥
	fd, err := os.Create(PUBLICKEY_PEM_FILE)
	if err != nil {
		logger.Warnf("create file failed, name=%v err=%v", PUBLICKEY_PEM_FILE, err)
		return err
	}
	defer fd.Close()

	// 将 RSA 公钥以 pem 格式存储
	// PEM 编码（Base64 编码）格式：
	// 通常以 -----BEGIN PUBLIC KEY----- 开头，-----END PUBLIC KEY----- 结尾。
	err = pem.Encode(fd, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	if err != nil {
		logger.Warnf("publick key pem encode failed, err=%v", err)
		return err
	}
	logger.Infof("generate public key into %v", PUBLICKEY_PEM_FILE)
	return nil
}

// 将私钥转为 PEM 格式的字符串
func privateKeyIntoFile(privKey *rsa.PrivateKey) error {
	// 将私钥转换成 pkcs#8 格式
	privASN1, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		logger.Warnf("MarshalPKCS8PrivateKey failed, err=%v", err)
		return err
	}
	err = os.WriteFile(PRIVATEKEY_DER_FILE, privASN1, 0600)
	if err != nil {
		logger.Warnf("write der file failed, file=%v", PRIVATEKEY_DER_FILE)
		return err
	}
	logger.Infof("generate private key into %v", PRIVATEKEY_DER_FILE)

	fd, err := os.Create(PRIVATEKEY_PEM_FILE)
	if err != nil {
		logger.Warnf("create file failed, name=%v err=%v", PRIVATEKEY_PEM_FILE, err)
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
	logger.Infof("generate private key into %v", PRIVATEKEY_PEM_FILE)
	return nil
}

func generateRSAKey() {
	priKey, pubKey, err := generateRSAKeyPair()
	if err != nil {
		return
	}
	err = privateKeyIntoFile(priKey)
	if err != nil {
		return
	}
	err = publicKeyIntoFile(pubKey)
	if err != nil {
		return
	}
}

/*
readRSAPrivateKey

解析 pem 和 der 格式的私钥

可以看到解析出来的私钥的 N, E, D 都是一样的
*/
func readRSAPrivateKey(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Warnf("read file failed, file=%v", filename)
		return
	}
	if strings.HasSuffix(filename, ".pem") {
		priKey, err := ParseRSAPemPrivateKey(data)
		if err != nil {
			logger.Warnf("parse pem key failed, file=%v", filename)
			return
		}
		logger.Infof("read pem private key from %v\n primes=%v\n N=%v\n E=%v\n D=%v\n", filename, priKey.Primes, priKey.N, priKey.E, priKey.D)
	} else if strings.HasSuffix(filename, ".der") {
		priKey, err := ParseRSADerPrivateKey(data)
		if err != nil {
			logger.Warnf("parse der key failed, file=%v", filename)
			return
		}
		logger.Infof("read der private key from %v\n primes=%v\n N=%v\n E=%v\n D=%v\n", filename, priKey.Primes, priKey.N, priKey.E, priKey.D)
	} else {
		logger.Warnf("file format error file=%v", filename)
	}
}

/*
readRSAPublicKey

解析 pem 和 der 格式的公钥
*/
func readRSAPublicKey(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Warnf("read file failed, file=%v", filename)
		return
	}
	if strings.HasSuffix(filename, ".pem") {
		pubKey, err := ParseRSAPemPublicKey(data)
		if err != nil {
			logger.Warnf("parse pem key failed, file=%v", filename)
			return
		}
		logger.Infof("read public key from %v\n N=%v\n E=%v\n", filename, pubKey.N, pubKey.E)
	} else if strings.HasSuffix(filename, ".der") {
		pubKey, err := ParseRSADerPublicKey(data)
		if err != nil {
			logger.Warnf("parse pem key failed, file=%v", filename)
			return
		}
		logger.Infof("read public key from %v\n N=%v\n E=%v\n", filename, pubKey.N, pubKey.E)
	} else {
		logger.Warnf("file format error file=%v", filename)
	}
}

func parseRSAKey() {
	for _, file := range []string{
		PRIVATEKEY_PEM_FILE, PRIVATEKEY_DER_FILE,
	} {
		readRSAPrivateKey(file)
	}
	logger.Infof("=====================")
	for _, file := range []string{
		PUBLICKEY_PEM_FILE, PUBLICKEY_DER_FILE,
	} {
		readRSAPublicKey(file)
	}
}

func main() {
	//generateKey()
	//parseRSAKey()
	//printPKCS1RSAKey()
	//printPKCS8PrivateKey()
	printPKIXPublicKey()
}
