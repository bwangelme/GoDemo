RSA 秘钥存储格式
===

## ANS1

- ANS1（Abstract Syntax Notation One）是一种用于描述数据结构的标准，它是一个符号语言，旨在为数据交换协议（如电子邮件、数字证书、网络协议等）定义数据结构。ANS1 主要用于在不同系统和平台之间进行数据交换，尤其是在涉及到加密、身份认证和数字签名等领域。

## PEM, DER

- PEM（Privacy-Enhanced Mail）是将 ANS1 编码数据通过 Base64 编码和 PEM 格式的头尾标识符封装在一起（例如 -----BEGIN CERTIFICATE----- 和 -----END CERTIFICATE-----）。
- DER 是 ANS1 编码数据的二进制表示形式，通常直接存储为文件。

RSA 公钥的 PEM 格式以 `-----BEGIN PUBLIC KEY-----` 开头，`-----END PUBLIC KEY-----` 结尾
RSA 公钥的 PKIX 格式是 X.509 公钥信息的一种编码格式，通常以 DER 或 PEM 格式存储。
它是标准的公钥信息格式，符合 PKIX（Public Key Infrastructure X.509）规范。

RSA 私钥的 PEM 格式有两种

## PKCS #1

PKCS#1 是用于 RSA 密钥的标准，它定义了 RSA 公钥和私钥的格式以及与 RSA 算法相关的加密操作。

主要内容：
RSA 公钥：PKCS#1 定义了如何表示 RSA 公钥。RSA 公钥由两个主要部分组成：模数（n） 和 公钥指数（e）。在 PKCS#1 中，公钥是一个大整数对（n, e），这些值是用来加密数据的。

RSA 私钥：PKCS#1 还定义了如何表示 RSA 私钥。私钥包括多个字段，最重要的是：

模数（n）
私钥指数（d）
与公钥指数相关的私钥参数（例如 p, q, dp, dq 等，表示素数因子以及加速加密过程的参数）
PKCS#1 对 RSA 私钥进行了详细定义，它使用标准的二进制编码格式（如 DER 格式）来存储这些数据。

典型格式：
RSA 私钥（PKCS#1 格式）：通常以 PEM 格式表示，开始于 -----BEGIN RSA PRIVATE KEY----- 和 -----END RSA PRIVATE KEY-----。
RSA 公钥（PKCS#1 格式）：通常以 PEM 格式表示，开始于 -----BEGIN PUBLIC KEY----- 和 -----END PUBLIC KEY-----。
应用：
RSA 密钥生成和使用中，如数字签名、加密数据、密钥交换等操作。

## PKCS #8

PKCS#8 是一个更通用的标准，定义了 私钥信息的格式，不仅仅是针对 RSA 私钥，也支持多种不同类型的私钥（如 DSA、ECDSA 等）。它的设计目标是提供一种更通用的格式，可以用来存储和交换各种类型的私钥信息，而不仅限于 RSA。

主要内容：
私钥信息结构：PKCS#8 定义了一个通用的私钥信息结构，可以存储任意加密算法的私钥。结构包括：

版本：标识 PKCS#8 格式的版本号。
算法标识符：指示所用加密算法（例如 RSA、DSA 或 ECDSA 等）。
私钥：私钥数据本身，通常是一个二进制编码（DER 格式）表示的私钥。
PKCS#8 的优势：与 PKCS#1 不同，PKCS#8 支持多种加密算法，而不仅仅是 RSA。此外，PKCS#8 格式支持通过密码加密私钥数据，从而提高私钥的安全性。

典型格式：
私钥（PKCS#8 格式）：与 PKCS#1 不同，PKCS#8 也使用 PEM 或 DER 编码格式，但文件格式的标识符为 -----BEGIN PRIVATE KEY----- 和 -----END PRIVATE KEY-----。PKCS#8 可以容纳多种类型的私钥，不仅限于 RSA。

示例：

-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT4wggE8AgEAAkEArEsmrxr7XYlNQ==
-----END PRIVATE KEY-----
应用：
PKCS#8 被广泛应用于各种加密库中，尤其在需要支持不同类型密钥的系统中（例如，RSA、DSA、ECDSA 等）。
