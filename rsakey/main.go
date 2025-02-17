package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

var (
	pubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsD6nqxDi+1bX7feWs+h0thPCbHQ7y3rk6TzXwd2jvCGK17r8LeAX7GhoLDSd2ECOXHZ+Bdu8uRmbox5Z7mWrLJnRge6NmVwBS6mQkNFkA+nEOumQleoPz2PCNbKcHqV/ivZrq2OIwGsKXN6isDL7StutcUyN60ZlbASb8uiP5y5b0rRDhX+o8sndIXx9+yN97T3T4lflDRyKfYpg+jQIXELMsbsXfpdobJhrYoiD+OO3ILheqaYNyJrogsP8kDt0bTXeE3KsOe7ffwz221aQpnfUokuthQai+Z4TVrnlgztC39C/6+d8jUpn1o+S6ClOpPqP9UsW54FVaW8z5jPYoQIDAQAB
-----END PUBLIC KEY-----`

	priKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCwPqerEOL7Vtft95az6HS2E8JsdDvLeuTpPNfB3aO8IYrXuvwt4BfsaGgsNJ3YQI5cdn4F27y5GZujHlnuZassmdGB7o2ZXAFLqZCQ0WQD6cQ66ZCV6g/PY8I1spwepX+K9murY4jAawpc3qKwMvtK261xTI3rRmVsBJvy6I/nLlvStEOFf6jyyd0hfH37I33tPdPiV+UNHIp9imD6NAhcQsyxuxd+l2hsmGtiiIP447cguF6ppg3ImuiCw/yQO3RtNd4Tcqw57t9/DPbbVpCmd9SiS62FBqL5nhNWueWDO0Lf0L/r53yNSmfWj5LoKU6k+o/1SxbngVVpbzPmM9ihAgMBAAECggEATvnsmXAKPpWWUikB66GNdy/Yjk/xoYdy/39HsbR3lCy1smE0cvw5zDKnB7QWTVr0UEO6yjZC1fE/OHO32efOMkDFTMOQyTmczQJxaSujdUtyJvIV28/UgNsTootkgSkpQ6ST7+u80njE3oPXhDx8NfnFuYEUEWtstGBGX63OGkoBSkQB0GFlFCxuRmSiz83OH25WL9qmOXdGaeTxvQd3j3yTPAOQMqFpCHBJs20Pg6bKBviEb6FW7spCbY5Unj7HVM/tw+DqVBslULwA2FBTkhk6ScGLuiqDulx7Wpc2R1bdQ/cs1rGVECcTuBgiekOlamOk54FdNqyGFQTqtE1yvwKBgQDnBmhNam6NpRjn1EQx9vr/M1Slt6pKceaWytnOFTVhZqTX0hzxFxRGWzxWNpfGg9oOkGrS5LH+hgrmBGV7JIbvioUZG1FiASrfjkrVGUZj2oqcdj7BheSz7a5xVNe6rC7Ywidw4cunC/CP4HdfTSu/o6TUlUUL78+cZiKiHn58zwKBgQDDTDWKTDiI15rN2XhtT9Je2rfU6jHWbBKrou9OLprK3ZvgJdTiCeTPs3k5CmzWZjzetlBDHX/cH89gfw0boVjK4zYVgQkBr7g7erwGW33gZ57ZYdu8VF3O6F5O0J0dpoyvnnLXbtuGjESWiJHPIxyPOilI/xh6Rq/LQwVwrwQPjwKBgCtgX5sRfbpoojl8+GTtO4lJCP6ocnfR1PrBEY4JG2GzVQYUtExsCel/3d9OFsc2IG4VnYkFWYoxfsBbWPZ7ED7PolfpcilVkMgyvkgum7HJ6bag2P2a9yr1WIh85phtFcqrAZ7HNmah7kQFYERrh+hOgHdNo44vM6ro3l3UHemvAoGAMIUMgDFzkjvOj/nJe47rOvmn1lPg0d7DvLScM5ZMir4H7eY4P3gpyphSM6Otao637LTqt+HqVCvq/5RRE15Aixdr5mfKbwrTAKP7drDgUxIrWuJ/Dwj+zVrZo0cc4bLxHOiGq5M1IvZSS/veDdIxVDwk6afG0wogvqUGAvrYTW0CgYEAiEY37hFTX/ZzqzgWsa/yq45qIlcQKr/JQtsovJqZ1ZzOwbs4b/pt12ehTpoCp5SrdYYX5TuohzQPyyMo60CkEEkaHQpOhqRcwIA6fQZfBET5OrC7R1BYu9GR5CVLUqBv+7votesZ3aJzsIDUTfrmmk3Q9+UT6P7SvL5t8p898jM=
-----END PRIVATE KEY-----`
)

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

func main() {
	// 解析公钥
	pubKey, err := ParseRSAPublicKey([]byte(pubKey))
	if err != nil {
		log.Fatalf("解析公钥失败: %v", err)
	}

	// 打印公钥信息
	fmt.Println(pubKey.E)
	fmt.Println(pubKey.N)

	privateKey, err := ParseRSAPrivateKey([]byte(priKey))
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}
	fmt.Println(privateKey.E)
	fmt.Println(privateKey.N)
}
