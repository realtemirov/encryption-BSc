package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"

	"github.com/realtemirov/encryption/repo"
)

type Service interface {
	Encryption(u *repo.User, text string) (string, error)
	Decryption(u *repo.User, text string) (string, error)
}

type service struct {
	db repo.Db
}

func NewService(db repo.Db) Service {
	return &service{db}
}

// Decryption implements Service.
func (s *service) Decryption(u *repo.User, text string) (string, error) {

	key := []byte("1234567890123456")
	iv := []byte("1234567890123456")

	if u.Type == repo.AES {
		if u.Decryption {
			aesDecrypted, _ := decryptAES(key, []byte(text), iv)
			return fmt.Sprintf("AES Decrypted (Go): `%s`", string(aesDecrypted)), nil
		} else {
			return text, nil
		}
	} else if u.Type == repo.DES {
		if u.Encryption {
			desDecrypted, _ := decryptDES(key[:8], []byte(text)[des.BlockSize:], iv[:8])
			return fmt.Sprintf("DES Decrypted (Go): `%s`", string(desDecrypted)), nil
		} else {
			return text, nil
		}
	} else {
		return text, nil
	}

}

// Encryption implements Service.
func (s *service) Encryption(u *repo.User, text string) (string, error) {
	key := []byte("1234567890123456")
	iv := []byte("1234567890123456")

	if u.Type == repo.AES {
		if u.Encryption {
			aesEncrypted, _ := encryptAES(key, []byte(text), iv)
			return fmt.Sprintf("AES Encrypted (Go): `%s`", string(aesEncrypted)), nil
		} else {
			return text, nil
		}
	} else if u.Type == repo.DES {
		if u.Encryption {
			desEncrypted, _ := encryptDES(key[:8], []byte(text), iv[:8])
			return fmt.Sprintf("DES Encrypted (Go): `%s`", string(desEncrypted)), nil
		} else {
			return text, nil
		}
	}

	return text, nil
}

// PKCS7 padding
func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := make([]byte, padding)
	for i := 0; i < padding; i++ {
		padText[i] = byte(padding)
	}
	return append(src, padText...)
}

// PKCS7 unpadding
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("unpad error: input is empty")
	}
	padding := int(src[length-1])
	if padding > length {
		return nil, fmt.Errorf("unpad error: invalid padding size")
	}
	return src[:length-padding], nil
}

// AES encryption function
func encryptAES(key, text, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedText := pad(text, aes.BlockSize)
	ciphertext := make([]byte, len(paddedText))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES decryption function
func decryptAES(key, ciphertext, iv []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	decrypted := make([]byte, len(decodedCiphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, decodedCiphertext)

	unpaddedText, err := unpad(decrypted)
	if err != nil {
		return "", err
	}

	return string(unpaddedText), nil
}

// DES encryption function
func encryptDES(key, text, iv []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	paddedText := pad(text, des.BlockSize)
	ciphertext := make([]byte, len(paddedText))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DES decryption function
func decryptDES(key, ciphertext, iv []byte) (string, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}

	decodedCiphertext, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < des.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	decrypted := make([]byte, len(decodedCiphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, decodedCiphertext)

	unpaddedText, err := unpad(decrypted)
	if err != nil {
		return "", err
	}

	return string(unpaddedText), nil
}
