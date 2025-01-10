package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
)

func custom() {

	pubKeyPem := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoZ67dtUTLxoXnNEzRBFB
mwukEJGC+y69cGgpNbtElQj3m4Aft/7cu9qYbTNguTSnCDt7uovZNb21u1vpZwKH
yVgFEGO4SA8RNnjhJt2D7z8RDMWX3saody7jo9TKlrPABLZGo2o8vadW8Dly/v+I
d0YDheCkVCoCEeUjQ8koXZhTwhYkGPu+vkdiqX5cUaiVTu1uzt591aO5Vw/hV4DI
hFKnOTnYXnpXiwRwtPyYoGTa64yWfi2t0bv99qz0BgDjQjD0civCe8LRXGGhyB1U
1aHjDDGEnulTYJyEqCzNGwBpzEHUjqIOXElFjt55AFGpCHAuyuoXoP3gQvoSj6RC
sQIDAQAB
-----END PUBLIC KEY-----`

	// Import public key
	pubKey := ImportSPKIPublicKeyPEM(pubKeyPem)
	// Base64 decode ciphertext
	ciphertextBytes, _ := base64.StdEncoding.DecodeString("ajQbkszbZ97YZaPSRBab9vj0DDLm9tTrQwSZ+ucPj+cYSmw06KLCtRH3SPn3b2DqSd1revLXqxMtSzFmjRvZ5F8y3nzdP8NJaRplOigbPFhKZTv7xBVK5ATEmLukgtI7f+d3KdmGUG+cyTkfxIrMBvB3BIS5oTiMNmC9pqLaWcDVF9qpuxnwEMQJbeO9nTklpdv+F8BrchHmeUkKRrMJBoPbbcfq9Hi4bHiFyxPWhwB66d/AryCKsFRhaX6hSkTL+0NvuhVhv98wdo3juv2Il50XKOCbfc8kUG628TcSK6n31piLF9cntSVTB/L/pVfcAxEwx4hcUhLuqmk6EZIJvGo0G5LM22fe2GWj0kQWm/b49Awy5vbU60MEmfrnD4/nGEpsNOiiwrUR90j5929g6knda3ry16sTLUsxZo0b2eRfMt583T/DSWkaZTooGzxYSmU7+8QVSuQExJi7pILSO3/ndynZhlBvnMk5H8SKzAbwdwSEuaE4jDZgvaai2lnA1RfaqbsZ8BDECW3jvZ05JaXb/hfAa3IR5nlJCkazCQaD223H6vR4uGx4hcsT1ocAeunfwK8girBUYWl+oUpEy/tDb7oVYb/fMHaN47r9iJedFyjgm33PJFButvE3Eiup99aYixfXJ7UlUwfy/6VX3AMRMMeIXFIS7qppOhGSCbxqNBuSzNtn3thlo9JEFpv2+PQMMub21OtDBJn65w+P5xhKbDToosK1EfdI+fdvYOpJ3Wt68terEy1LMWaNG9nkXzLefN0/w0lpGmU6KBs8WEplO/vEFUrkBMSYu6SC0jt/53cp2YZQb5zJOR/EiswG8HcEhLmhOIw2YL2motpZwNUX2qm7GfAQxAlt472dOSWl2/4XwGtyEeZ5SQpGswkGg9ttx+r0eLhseIXLE9aHAHrp38CvIIqwVGFpfqFKRMv7Q2+6FWG/3zB2jeO6/YiXnRco4Jt9zyRQbrbxNxIrqffWmIsX1ye1JVMH8v+lV9wDETDHiFxSEu6qaToRkgm8")
	// Split ciphertext into signature chunks a 2048/8 bytes and decrypt each chunk
	reader := bytes.NewReader(ciphertextBytes)
	var writer bytes.Buffer
	ciphertextBytesChunk := make([]byte, 2048/8)
	for {
		n, _ := io.ReadFull(reader, ciphertextBytesChunk)
		if n == 0 {
			break
		}
		decryptChunk(ciphertextBytesChunk, &writer, pubKey)
	}
	// Concatenate decrypted signature chunks
	decryptedData := writer.String()
	fmt.Println(decryptedData)
}

func ImportSPKIPublicKeyPEM(spkiPEM string) *rsa.PublicKey {
	body, _ := pem.Decode([]byte(spkiPEM))
	publicKey, _ := x509.ParsePKIXPublicKey(body.Bytes)
	if publicKey, ok := publicKey.(*rsa.PublicKey); ok {
		return publicKey
	} else {
		return nil
	}
}

func decryptChunk(ciphertextBytesChunk []byte, writer *bytes.Buffer, pubKey *rsa.PublicKey) {
	// Decrypt each signature chunk
	ciphertextInt := new(big.Int)
	ciphertextInt.SetBytes(ciphertextBytesChunk)
	decryptedPaddedInt := decrypt(new(big.Int), pubKey, ciphertextInt)
	// Remove padding
	decryptedPaddedBytes := make([]byte, pubKey.Size())
	decryptedPaddedInt.FillBytes(decryptedPaddedBytes)
	start := bytes.Index(decryptedPaddedBytes[1:], []byte{0}) + 1 // // 0001FF...FF00<data>: Find index after 2nd 0x00
	decryptedBytes := decryptedPaddedBytes[start:]
	// Write decrypted signature chunk
	writer.Write(decryptedBytes)
}

func decrypt(c *big.Int, pub *rsa.PublicKey, m *big.Int) *big.Int {
	// Textbook RSA
	e := big.NewInt(int64(pub.E))
	c.Exp(m, e, pub.N)
	return c
}
