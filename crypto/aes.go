package main

/*
本程序使用 aes-256-cbc 模式对一个密文进行解密

秘钥的长度需要是 256 位，即 32 个字节

密文的格式是

iv + hmacSum + 密文

iv 长 16 (aes.BlockSize)个字节
hmacSum 长 32 个字节
密文的长度是 16 (aes.BlockSize)的整数倍
*/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

func main() {
	//encrypted := "yFA+imJBPNwkiLVDnT5z98YfnLo8mh7Wv5gp51Vs5zD1EJNmJKLkR/1hqGL4+tQ/ZWYaggPpBjfWigxCTBg5WjMYaylf1iJD7CvJ1l2Urft8MY4b2Ays0vGvAELduef8"
	//key := "qnr8Rp@ICi6&yeoW8B5sG"

	encrypted := "bsfnzwyOHI16vaDXkEqEU0K4UolWDTdEMcYzPJRegefuMy9mGDYBkcIkGO9L0MZHvU+FjZeMA9SbD+cE2FDrRUrsna5kY6YiG/oA2HqV5z+ABH63kkGkM3VjVDn4MPVX"
	key := "2xMtq&#g7RI.hwV0sjwCl2MLGC6SHRVM"

	fmt.Println(len(encrypted), len(key))

	decoded, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		panic(err)
	}

	ivlen := aes.BlockSize
	if len(decoded) < ivlen+32 {
		panic("Ciphertext too short")
	}

	iv := decoded[:ivlen]
	hmacReceived := decoded[ivlen : ivlen+32]
	ciphertextRaw := decoded[ivlen+32:]

	hmacCalculated := hmac.New(sha256.New, []byte(key))
	hmacCalculated.Write(iv)
	hmacCalculated.Write(ciphertextRaw)
	hmacCalculatedSum := hmacCalculated.Sum(nil)

	if !hmac.Equal(hmacCalculatedSum, hmacReceived) {
		panic(errors.New("HMAC verification failed"))
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextRaw, ciphertextRaw)

	// 去除填充
	padding := int(ciphertextRaw[len(ciphertextRaw)-1])
	ciphertextRaw = ciphertextRaw[:len(ciphertextRaw)-padding]

	fmt.Println(string(ciphertextRaw))
}
