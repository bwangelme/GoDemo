package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

// RSASigner 结构体用于管理RSA签名操作
type RSASigner struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewRSASigner 创建新的RSA签名器
func NewRSASigner() (*RSASigner, error) {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("生成RSA密钥失败: %v", err)
	}

	return &RSASigner{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}, nil
}

// NewRSASignerFromKeys 从现有密钥创建RSA签名器
func NewRSASignerFromKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *RSASigner {
	return &RSASigner{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

// Sign 使用私钥对数据进行签名
func (r *RSASigner) Sign(data []byte) ([]byte, error) {
	// 计算数据的SHA256哈希
	hash := sha256.Sum256(data)

	// 使用私钥对哈希进行签名
	signature, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, fmt.Errorf("签名失败: %v", err)
	}

	return signature, nil
}

// Verify 使用公钥验证签名
func (r *RSASigner) Verify(data []byte, signature []byte) error {
	// 计算数据的SHA256哈希
	hash := sha256.Sum256(data)

	// 使用公钥验证签名
	err := rsa.VerifyPKCS1v15(r.publicKey, crypto.SHA256, hash[:], signature)
	if err != nil {
		return fmt.Errorf("签名验证失败: %v", err)
	}

	return nil
}

// GetPublicKeyPEM 获取PEM格式的公钥
func (r *RSASigner) GetPublicKeyPEM() string {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(r.publicKey)
	if err != nil {
		return ""
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM)
}

// GetPrivateKeyPEM 获取PEM格式的私钥
func (r *RSASigner) GetPrivateKeyPEM() string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(r.privateKey)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return string(privateKeyPEM)
}

// SaveKeysToFile 将密钥保存到文件
func (r *RSASigner) SaveKeysToFile(privateKeyFile, publicKeyFile string) error {
	// 保存私钥
	err := os.WriteFile(privateKeyFile, []byte(r.GetPrivateKeyPEM()), 0600)
	if err != nil {
		return fmt.Errorf("保存私钥失败: %v", err)
	}

	// 保存公钥
	err = os.WriteFile(publicKeyFile, []byte(r.GetPublicKeyPEM()), 0644)
	if err != nil {
		return fmt.Errorf("保存公钥失败: %v", err)
	}

	return nil
}

// LoadKeysFromFile 从文件加载密钥
func LoadKeysFromFile(privateKeyFile, publicKeyFile string) (*RSASigner, error) {
	// 读取私钥
	privateKeyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("读取私钥失败: %v", err)
	}

	// 解析私钥
	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("解析私钥PEM失败")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析私钥失败: %v", err)
	}

	// 读取公钥
	publicKeyBytes, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("读取公钥失败: %v", err)
	}

	// 解析公钥
	block, _ = pem.Decode(publicKeyBytes)
	if block == nil {
		return nil, fmt.Errorf("解析公钥PEM失败")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %v", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("公钥类型错误")
	}

	return NewRSASignerFromKeys(privateKey, rsaPublicKey), nil
}

// SignString 对字符串进行签名并返回base64编码的签名
func (r *RSASigner) SignString(data string) (string, error) {
	signature, err := r.Sign([]byte(data))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifyString 验证base64编码的签名
func (r *RSASigner) VerifyString(data string, signatureBase64 string) error {
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return fmt.Errorf("解码签名失败: %v", err)
	}

	return r.Verify([]byte(data), signature)
}

// 示例用法
func main() {
	fmt.Println("=== RSA 签名验签示例 ===\n")

	// 1. 创建新的RSA签名器
	signer, err := NewRSASigner()
	if err != nil {
		log.Fatalf("创建RSA签名器失败: %v", err)
	}

	// 2. 显示密钥信息
	fmt.Println("生成的RSA密钥对:")
	fmt.Printf("私钥长度: %d bits\n", signer.privateKey.Size()*8)
	fmt.Printf("公钥长度: %d bits\n", signer.publicKey.Size()*8)

	// 3. 保存密钥到文件
	err = signer.SaveKeysToFile("private_key.pem", "public_key.pem")
	if err != nil {
		log.Printf("保存密钥失败: %v", err)
	} else {
		fmt.Println("\n密钥已保存到文件:")
		fmt.Println("- private_key.pem (私钥)")
		fmt.Println("- public_key.pem (公钥)")
	}

	// 4. 测试签名和验证
	testData := "Hello, RSA签名测试!"
	fmt.Printf("\n测试数据: %s\n", testData)

	// 签名
	signature, err := signer.SignString(testData)
	if err != nil {
		log.Fatalf("签名失败: %v", err)
	}
	fmt.Printf("签名结果 (Base64): %s\n", signature)

	// 验证签名
	err = signer.VerifyString(testData, signature)
	if err != nil {
		log.Fatalf("签名验证失败: %v", err)
	}
	fmt.Println("✓ 签名验证成功!")

	// 5. 测试错误数据验证
	wrongData := "错误的测试数据"
	err = signer.VerifyString(wrongData, signature)
	if err != nil {
		fmt.Printf("✓ 错误数据验证失败 (预期行为): %v\n", err)
	} else {
		fmt.Println("✗ 错误数据验证成功 (异常行为)")
	}

	// 6. 从文件加载密钥测试
	fmt.Println("\n=== 从文件加载密钥测试 ===")
	loadedSigner, err := LoadKeysFromFile("private_key.pem", "public_key.pem")
	if err != nil {
		log.Printf("加载密钥失败: %v", err)
	} else {
		// 使用加载的密钥进行验证
		err = loadedSigner.VerifyString(testData, signature)
		if err != nil {
			log.Printf("使用加载的密钥验证失败: %v", err)
		} else {
			fmt.Println("✓ 使用加载的密钥验证成功!")
		}
	}

	fmt.Println("\n=== 示例完成 ===")
}
