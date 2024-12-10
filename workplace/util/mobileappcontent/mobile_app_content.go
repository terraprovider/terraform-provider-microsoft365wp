package mobileappcontent

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"
)

type ContentFileInfo struct {
	Name          string
	Size          int64
	SizeEncrypted int64
	Manifest      any // string, but may be nil
}

type FileEncryptionInfo struct {
	FileDigest           string `json:"fileDigest"`
	FileDigestAlgorithm  string `json:"fileDigestAlgorithm"`
	EncryptionKey        string `json:"encryptionKey"`
	InitializationVector string `json:"initializationVector"`
	Mac                  string `json:"mac"`
	MacKey               string `json:"macKey"`
	ProfileIdentifier    string `json:"profileIdentifier"`
}

const kEncryptedContentTempName = "IntuneMobileAppContentEncrypted_*.bin"

func runOnZipFile(archivePath string, filePathInArchive string, impl func(io.Reader) (any, error)) (result any, err error) {

	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return
	}
	defer zipReader.Close()

	for _, zipFile := range zipReader.File {
		if strings.EqualFold(zipFile.Name, filePathInArchive) {

			readCloser, err := zipFile.Open()
			if err != nil {
				return nil, err
			}
			defer readCloser.Close()

			return impl(readCloser)

		}
	}

	return nil, fmt.Errorf("file %s not found in archive %s", filePathInArchive, archivePath)
}
