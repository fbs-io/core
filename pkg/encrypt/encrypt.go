/*
 * @Author: reel
 * @Date: 2023-05-16 21:15:13
 * @LastEditors: reel
 * @LastEditTime: 2023-06-03 20:38:54
 * @Description: 字符串加密
 */
package encrypt

import (
	"bytes"
	"core/pkg/errorx"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

var (
	InternalPrivKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDkGO3/7Xh2fEpgk1XiQ++6UAkhAjedka+zYstuFdC6N+UAi5ri
zH/j61BfghiW7Lh4vWIgD7L0CwWFogm6GM2rcFHEAJoyrMfTif4HTZTGRsbaWKM5
L2zgzIfEcyYc8r7NYwHIYtX9S5geAglf3u3daEtf3uePbr6C+3C83uY5awIDAQAB
AoGAReo+SUwIuIiwY5nFoW0hdgJCShPb6DhXmqyVnlChgfLQBrLD7vCv4rSmFiBS
WTCF+rxR73HgcF1Qe/2G7bvcjjOnF75mYnN9rOMEc3LO5Fnm+OjK6mf1U0wTrTTS
Jjt1tmK0HYmRToOqJ10ySDWW0+6XYWpImedwhLSmGZGRK4ECQQDkVXcGRFgUpXAv
saABBjNuNr204FBrhSrvrGPuu0Bg7yN9tI2wstc66R5N47wFcjnbB6eGy8rsNg8c
B0OiSONTAkEA/7whRc2izY4lW4E74I+AfAz3bjApD1pDJkZo3vVf6djQ1foixEyk
UdJweGHtJ3zx4iXZJr0AWtfR6HMugrjmiQJANVuvg9xmxPRgQhT9MiAT1raeIG2m
/WlSPk4H0Fsb0Usw/Qg7cEZqu46MkWEdqBwoXwHr6Tkog4iigUdFS+BClwJBAOkW
OqlZrRp3hbsqRj397ZijZN4MjVAN8AgxwqH8ucf1Mxrkms2aIWbmTFacwr/sFLcP
0iWJvIoQDaU1Xl4NUykCQEqL5SpaYS7Z4aTaQk9Meeg5uf9jAE2kP5KZmpGPbhsp
/EOkkbtM/mEfR/JeYx3Ukw8E+2Fq0FiqdBzXmC+LyWc=
-----END RSA PRIVATE KEY-----
`)
	InternalPubKey = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDkGO3/7Xh2fEpgk1XiQ++6UAkh
Ajedka+zYstuFdC6N+UAi5rizH/j61BfghiW7Lh4vWIgD7L0CwWFogm6GM2rcFHE
AJoyrMfTif4HTZTGRsbaWKM5L2zgzIfEcyYc8r7NYwHIYtX9S5geAglf3u3daEtf
3uePbr6C+3C83uY5awIDAQAB
-----END PUBLIC KEY-----
`)
)

// 加密
func RsaEncrypt(data, key []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	partLen := pub.N.BitLen()/8 - 11
	chunks := split([]byte(data), partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		bytes, err := rsa.EncryptPKCS1v15(rand.Reader, pub, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(bytes)
	}
	return buffer.Bytes(), nil
}

// 解密
func RsaDecrypt(ciphertext, key []byte) ([]byte, error) {
	//获取私钥
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errorx.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	partLen := priv.N.BitLen() / 8

	chunks := split(ciphertext, partLen)
	buffer := bytes.NewBuffer([]byte{})
	for _, chunk := range chunks {
		decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, priv, chunk)
		if err != nil {
			return nil, err
		}
		buffer.Write(decrypted)
	}

	return buffer.Bytes(), nil
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf)
	}
	return chunks
}

func S2B(str string) []byte {
	res := make([]byte, 0, 256)
	strList := strings.Split(str, "")
	var s0 string
	for i, s := range strList {
		// i, _ := strconv.Atoi(s)
		s0 += s
		if i%2 == 1 {
			num := Decode(s0)
			res = append(res, byte(num))
			s0 = ""
		}

	}
	return res
}

func B2S(b []byte) string {
	var reslist = make([]string, 0, 256)
	for _, item := range b {
		str := Encode(int(item))
		if len(str) == 1 {
			str = "0" + str
		}
		reslist = append(reslist, str)
	}

	return strings.Join(reslist, "")
}

func GenerMd5(str string) (md5str string) {
	data := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func InternalEncode(body []byte) (str string, err error) {
	b, err := RsaEncrypt(body, InternalPubKey)
	if err != nil {
		return "", err
	}
	str = B2S(b)
	return
}

func InternalDecode(str string) (body []byte, err error) {
	rsaB := S2B(str)
	body, err = RsaDecrypt(rsaB, InternalPrivKey)
	return
}
