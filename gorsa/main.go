package main

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/wenzhenxi/gorsa"
)

var (
	PubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwIMIQtpub8EqPle5rLImU35EziL26uxYPVVg8jRWm/HhlKx+QGF5m41xJIL6mY5U2C8UjWaVNYGI1mn91mJ/t++PganX3AKrIb7t+V42LsHG8Pz6jYiTdLnSIyM9XbWxzMDUSb1XyqJKXw9fCOCpfKWJgfi0aLh6pteYqrKNH6dY6WcKqGIXEPRuYfXlvlKStYl2bNdpyLD2GmhcH51wvXSdFanjXXZqQEojogj5lfoKVDNkWYrQ9UJz16jS2Cl0FTyQ4vcQEKzqvs5vvSfuemJ8n6wYEXk72VHn4SO+wmPQIYyhulVjQT0F6qScokYZuH9PexA8PgEfyLohhLu+xQIDAQAB
-----END PUBLIC KEY-----`
	PriKey = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDAgwhC2m5vwSo+V7mssiZTfkTOIvbq7Fg9VWDyNFab8eGUrH5AYXmbjXEkgvqZjlTYLxSNZpU1gYjWaf3WYn+374+BqdfcAqshvu35XjYuwcbw/PqNiJN0udIjIz1dtbHMwNRJvVfKokpfD18I4Kl8pYmB+LRouHqm15iqso0fp1jpZwqoYhcQ9G5h9eW+UpK1iXZs12nIsPYaaFwfnXC9dJ0VqeNddmpASiOiCPmV+gpUM2RZitD1QnPXqNLYKXQVPJDi9xAQrOq+zm+9J+56YnyfrBgReTvZUefhI77CY9AhjKG6VWNBPQXqpJyiRhm4f097EDw+AR/IuiGEu77FAgMBAAECggEAAhmxN1qfBQDXfmXNA8Zdsy6PBjsQe8Ms3dcb/Xlk/REwRpICcEmMaZVOx9jkw8X2Wkro++kuaqBkFw6QCpBV1DnHacu4D8bAAp+My/afwqgc93ZMddkhpZ0ELSnMkcukKeJPs/V6zqU/lsEWe/2Bghs5B4qTBuCwKhwdyiel8QVeDFHLYaD0Rgxp9CiquKd5vh6bu+7JdliW4t7P6UGGBhOlcRgIrOYK1p67pxrSLrDFomFdrw1J0bsyC3Aunnaj3fPedV+PCGK8Dmm4luSp60/SIDvPrbvOkFaRXXeiYs5RzCZ4OhMA/BJPOtNt9iFKyW5Wyq5j7xfTu1Z+3EUEGQKBgQDHM6jbnfjgzMT3MoB3LxM2pHAPZMUytuIZy4xQoFdRnv9LBtYvmFxj7oAWKyoUbxlbIXCdHcXgyDWH8UUVvA/FaKocqvBGowijri8uP1nN6kqh/jJkdGjJebMgc5ONGBql8AwH0pTuEwstc5+HUjZCWkp/30+6MSrGNRDx6fDoiQKBgQD3Zw2QcR1Qrg4HZ6VUyTKMHi5KQ+JWK54Go65TTBRTPF9QqL96CEDgGf2risVAw1y5aD8FV0mp3b3ydzXXRfTdu0axeWkCeV7wvL6F0uzbATlDL2+mCD05G1GrqITMsVzjkr2ZdkD9FbQmP22+TPATfxmNCavH3JG2BK0z03HdXQKBgDiOez+3T0UL/lQs23p7PBpEf6hvNOncML+gIft2OrRqzxLPIxqqSHj66xTgNndMv8c27FE/dcIHNeOd5b5xAY1L1RiEk+mKOla2iqC3zdh/z6bElCwfPO0PB4OMLPVfRZmdWN0TtcMOVxsVe9KgzdTy67n4fhtgAEP8Jw54HDT5AoGAL+mRQuvlFX7f0KdN6YGdfG1L4a4L40xHKloApLkTJpuGigRkMKnwhnYCHnqdgNFU38NEkTA6X99FrfNQRgSSS269Xhl5MLV4oX6sFxamMDOod050fN8TnD+iLXBRZ6LhrmD4vwObymetn8qg4j3cMKpotFuvHOGmhm2ZbXQBlTECgYAXACrcbO3snVczwI8pZMkIm4J6GWLnKJkNf0IQqeh9GQrMz8zYHS2cUQCk1a4kxPGph1yEYow4jeyYoX1OT8dofGu15hT8JGF7u+XW36lEvT8Evc9t7YWQ5lDqPmPQGeKJqFVC9g1fwv/pEiW7w22JfqEauocWJ9/U3fgBXaMjEA==
-----END PRIVATE KEY-----`
)

func main() {
	data := `hello, world`

	content, err := PrivateEncrypt(data, PriKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("after encrypt")
	fmt.Println(content)

	dData, err := PublicDecrypt(content, PubKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dData)
}

func PublicDecrypt(data string, publicKey string) (string, error) {
	dataBs, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	gRsa := gorsa.RSASecurity{}
	if err := gRsa.SetPublicKey(publicKey); err != nil {
		return "", err
	}

	rsaData, err := gRsa.PubKeyDECRYPT(dataBs)
	if err != nil {
		return "", err
	}

	return string(rsaData), nil
}

func PrivateEncrypt(data string, privateKey string) (string, error) {
	gRsa := gorsa.RSASecurity{}
	gRsa.SetPrivateKey(privateKey)

	rsaData, err := gRsa.PriKeyENCTYPT([]byte(data))
	if err != nil {
		return "", err
	}

	baseData := base64.StdEncoding.EncodeToString(rsaData)
	return baseData, nil
}
