package lightsocks

import (
	"time"
)

type cipher struct {
	// 编码用的密码
	encodePassword *password
	// 解码用的密码
	decodePassword *password
	// xor
	xorEncryption [xorEncryptionLength]byte
	// current index for xor
	xorIndex int
}

// 加密原数据
func (cipher *cipher) encode(bs []byte) {
	for i, v := range bs {
		bs[i] = cipher.encodePassword[v] ^ cipher.xorEncryption[cipher.xorIndex]
		cipher.xorIndex++
		if cipher.xorIndex >= xorEncryptionLength {
			cipher.xorIndex = 0
		}
	}
}

// 解码加密后的数据到原数据
func (cipher *cipher) decode(bs []byte) {
	for i, v := range bs {
		bs[i] = v ^ cipher.xorEncryption[cipher.xorIndex]
		bs[i] = cipher.decodePassword[bs[i]]
		cipher.xorIndex++
		if cipher.xorIndex >= xorEncryptionLength {
			cipher.xorIndex = 0
		}
	}
}

// 新建一个编码解码器
func newCipher(encodePassword *password) *cipher {
	decodePassword := &password{}
	for i, v := range encodePassword {
		encodePassword[i] = v
		decodePassword[v] = byte(i)
	}
	var xorEncryption [xorEncryptionLength]byte
	for i, v := range encodePassword {
		xorEncryption[i] = v
	}
	for i, v := range decodePassword {
		xorEncryption[i+passwordLength] = v
	}
	for i, v := range encodePassword {
		xorEncryption[3*passwordLength-i-1] = v
	}
	for i, v := range decodePassword {
		xorEncryption[4*passwordLength-i-1] = v
	}
	seed := time.Now().UTC().Hour()

	return &cipher{
		encodePassword: encodePassword,
		decodePassword: decodePassword,
		xorEncryption:  xorEncryption,
		xorIndex:       seed,
	}
}
