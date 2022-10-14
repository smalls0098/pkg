package cryptor

import (
	"testing"
)

func Test_AesCbc(t *testing.T) {
	key := []byte("smalls0098")
	content := `id=1`
	encrypt := AesCbcEncrypt([]byte(content), key, key)
	decrypt := AesCbcDecrypt(encrypt, key, key)
	t.Log(content == string(decrypt))
}
