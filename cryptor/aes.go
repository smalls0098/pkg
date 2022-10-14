package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

func pkcs7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	var padText []byte
	if padding == 0 {
		// 已对齐，填充一整块数据，每个数据为 blockSize
		padText = bytes.Repeat([]byte{byte(blockSize)}, blockSize)
	} else {
		// 未对齐 填充 padding 个数据，每个数据为 padding
		padText = bytes.Repeat([]byte{byte(padding)}, padding)
	}
	return append(src, padText...)
}

func pkcs7UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}

// AesCbcEncrypt encrypt data with key use AES CBC algorithm
// len(key) should be 16, 24 or 32
func AesCbcEncrypt(data, key []byte, iv []byte) []byte {
	block, _ := aes.NewCipher(key)
	data = pkcs7Padding(data, block.BlockSize())

	mode := cipher.NewCBCEncrypter(block, iv)

	encrypted := make([]byte, len(data))
	mode.CryptBlocks(encrypted, data)
	return encrypted
}

// AesCbcDecrypt decrypt data with key use AES CBC algorithm
// len(key) should be 16, 24 or 32
func AesCbcDecrypt(encrypted, key []byte, iv []byte) []byte {
	block, _ := aes.NewCipher(key)

	mode := cipher.NewCBCDecrypter(block, iv)

	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	decrypted = pkcs7UnPadding(decrypted)
	return decrypted
}
