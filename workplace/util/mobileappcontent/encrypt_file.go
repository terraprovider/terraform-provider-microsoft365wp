package mobileappcontent

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func EnryptFile(plainPath string) (fileEncryptionInfo FileEncryptionInfo, encryptedContentPath string, encryptedContentSize int64, err error) {

	const kAesKeySizeBytes = 256 / 8 // 32
	const kAesBlockSize = 16
	const kEncryptionBlockSize = 1024 * 1024 * 4 // 4 MB

	// open source file
	plainFile, err := os.Open(plainPath)
	if err != nil {
		return
	}
	defer plainFile.Close()

	// get file size
	plainFileInfo, err := plainFile.Stat()
	if err != nil {
		return
	}
	plainFileSize := plainFileInfo.Size()

	// calculate source sha256
	plainSha256er := sha256.New()
	if _, err = io.Copy(plainSha256er, plainFile); err != nil {
		return
	}
	if _, err = plainFile.Seek(0, 0); err != nil {
		return
	}
	plainSha256 := plainSha256er.Sum(nil)

	// create/open target file
	encContentFile, err := os.CreateTemp("", kEncryptedContentTempName)
	if err != nil {
		return
	}
	defer encContentFile.Close()

	// prepare HMAC SHA256
	hmacKey := make([]byte, kAesKeySizeBytes)
	_, err = rand.Read(hmacKey)
	if err != nil {
		return
	}
	sha256Hmacer := hmac.New(sha256.New, hmacKey)

	// ensure space for MAC in target file
	macPlaceHolder := make([]byte, sha256Hmacer.Size())
	encContentFile.Write(macPlaceHolder)

	// create AES key
	aesKey := make([]byte, kAesKeySizeBytes)
	_, err = rand.Read(aesKey)
	if err != nil {
		return
	}
	// the key is private and does not get written to the encrypted file!

	// create IV & save to target file
	aesIv := make([]byte, kAesBlockSize)
	_, err = rand.Read(aesIv)
	if err != nil {
		return
	}
	encContentFile.Write(aesIv)

	// prepare AES encryption
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return
	}
	cbcMode := cipher.NewCBCEncrypter(aesBlock, aesIv)

	// prepare buffers
	bufferPlain := make([]byte, kEncryptionBlockSize)
	bufferEnc := make([]byte, cap(bufferPlain))
	plainBytesReadSum := int64(0)

	for {
		// reset buffer to full size
		bufferPlain = bufferPlain[:cap(bufferPlain)]

		// read block from source
		var plainBytesRead int
		plainBytesRead, err = plainFile.Read(bufferPlain)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			break
		}
		// trim buffer to size read
		bufferPlain = bufferPlain[:plainBytesRead]
		plainBytesReadSum += int64(plainBytesRead)

		// pad on last read
		var bufferPlainPadded []byte
		if plainBytesReadSum >= plainFileSize {
			// copy bufferPlain since cap(bufferPlain) must not be changed by append() in PKCS5Padding()
			bufferPlainToPad := append([]byte{}, bufferPlain...)
			bufferPlainPadded = PKCS5Padding(bufferPlainToPad, kAesBlockSize)
		} else if plainBytesRead == kEncryptionBlockSize {
			// kEncryptionBlockSize should be divisable by kAesBlockSize, so no padding needed
			bufferPlainPadded = bufferPlain
		} else {
			err = fmt.Errorf("plainBytesRead (%d) is not equal to kEncryptionBlockSize (%d) but plainBytesReadSum (%d) is not equal or more than plainFileSize (%d)",
				plainBytesRead, kEncryptionBlockSize, plainBytesReadSum, plainFileSize)
			return
		}

		// resize encrypted buffer to plain buffer size
		bufferEnc = bufferEnc[:len(bufferPlainPadded)]
		// encrypt block
		cbcMode.CryptBlocks(bufferEnc, bufferPlainPadded)

		// save encrypted block to target file
		_, err = encContentFile.Write(bufferEnc)
		if err != nil {
			break
		}
	}
	if err != nil {
		return
	}

	// seek to ecnryted content starting at IV / after MAC space
	if _, err = encContentFile.Seek(int64(sha256Hmacer.Size()), 0); err != nil {
		return
	}
	// create HMAC from there
	if _, err = io.Copy(sha256Hmacer, encContentFile); err != nil {
		return
	}
	encContentHmac := sha256Hmacer.Sum(nil)
	// save MAC to beginning of file (before IV)
	if _, err = encContentFile.WriteAt(encContentHmac, 0); err != nil {
		return
	}

	// get encrypted file size
	encFileInfo, err := encContentFile.Stat()
	if err != nil {
		return
	}

	// results
	fileEncryptionInfo.FileDigest = base64.StdEncoding.EncodeToString(plainSha256)
	fileEncryptionInfo.FileDigestAlgorithm = "SHA256"
	fileEncryptionInfo.EncryptionKey = base64.StdEncoding.EncodeToString(aesKey)
	fileEncryptionInfo.InitializationVector = base64.StdEncoding.EncodeToString(aesIv)
	fileEncryptionInfo.Mac = base64.StdEncoding.EncodeToString(encContentHmac)
	fileEncryptionInfo.MacKey = base64.StdEncoding.EncodeToString(hmacKey)
	fileEncryptionInfo.ProfileIdentifier = "ProfileVersion1"
	encryptedContentPath = encContentFile.Name()
	encryptedContentSize = encFileInfo.Size()

	return
}

func PKCS5Padding(bufferIn []byte, blockSize int) []byte {

	padLength := blockSize - len(bufferIn)%blockSize
	// we need to always pad so there is always some padding that can be removed after decryption
	if padLength == 0 {
		padLength = blockSize
	}

	padBytes := bytes.Repeat([]byte{byte(padLength)}, padLength)
	bufferOut := append(bufferIn, padBytes...)

	return bufferOut
}
