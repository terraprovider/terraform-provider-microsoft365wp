package azurestorage

import (
	"bytes"
	"context"
	"terraform-provider-microsoft365wp/workplace/util/retryablehttputil"

	// nosemgrep
	"crypto/md5"

	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type blockBlobBlockList struct {
	XMLName xml.Name `xml:"BlockList"`
	BlockId []string `xml:"Latest"`
}

func UploadBlobInBlocks(ctx context.Context, fileName string, sasUri string) (fileMd5Sum string, err error) {

	const kBlobBlockSize = 1024 * 1024 * 4 // 4 MB

	const kLogPref = "azurestorage.UploadBlobInBlocks: "

	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	md5Hash := md5.New()
	if _, err = io.Copy(md5Hash, file); err != nil {
		return
	}
	if _, err = file.Seek(0, 0); err != nil {
		return
	}
	fileMd5SumRaw := md5Hash.Sum(nil)
	tflog.Trace(ctx, kLogPref+fmt.Sprintf("file md5 sum: %x", fileMd5SumRaw))

	client := retryablehttp.NewClient()
	retryablehttputil.ConfigureClientRetryLimitsAndBackoff(client)
	client.CheckRetry = retryablehttp.ErrorPropagatedRetryPolicy
	client.ErrorHandler = retryablehttp.PassthroughErrorHandler // return http error and response body

	blockBuffer := make([]byte, kBlobBlockSize)
	blockList := blockBlobBlockList{}
	for blockNum := uint32(0); true; blockNum++ {
		blockBuffer = blockBuffer[:cap(blockBuffer)]
		bytesRead, err := file.Read(blockBuffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		blockBuffer = blockBuffer[:bytesRead]

		blockNumBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(blockNumBytes, blockNum)
		blockIdBase64 := base64.StdEncoding.EncodeToString(blockNumBytes)
		blockList.BlockId = append(blockList.BlockId, blockIdBase64)

		blockUri := fmt.Sprintf("%s&comp=block&blockid=%s", sasUri, blockIdBase64)
		req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPut, blockUri, bytes.NewBuffer(blockBuffer))
		if err != nil {
			return "", err
		}
		req.Header.Add("x-ms-blob-type", "BlockBlob")

		err = blobHttpDoWithErrorProcessing(client, req)
		if err != nil {
			return "", err
		}
	}

	blockListUri := fmt.Sprintf("%s&comp=blocklist", sasUri)
	blockListBytes, err := xml.Marshal(blockList)
	if err != nil {
		return
	}
	tflog.Trace(ctx, kLogPref+fmt.Sprintf("blockList: %s", string(blockListBytes)))
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodPut, blockListUri, bytes.NewBuffer(blockListBytes))
	if err != nil {
		return
	}
	fileMd5SumB64 := base64.StdEncoding.EncodeToString(fileMd5SumRaw)
	req.Header.Add("x-ms-blob-content-md5", fileMd5SumB64)

	err = blobHttpDoWithErrorProcessing(client, req)
	if err != nil {
		return
	}

	return fileMd5SumB64, nil
}

func blobHttpDoWithErrorProcessing(client *retryablehttp.Client, req *retryablehttp.Request) error {
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode >= 300 {

		errorMessage := fmt.Sprintf("azurestorage.blobClientDo: http server returned %d on %s", resp.StatusCode, req.URL.String())

		var respBodyString string
		if resp.Body != nil {
			respBodyBytes, err2 := io.ReadAll(resp.Body)
			if err2 == nil {
				respBodyString = string(respBodyBytes)
			} else {
				respBodyString += "(error reading response body: " + err2.Error() + ")"
			}
		} else {
			respBodyString = "(nil)"
		}
		errorMessage += fmt.Sprintf("\nResponse body: %s", respBodyString)

		if err != nil {
			errorMessage += fmt.Sprintf("\nretryablehttp.Client error: %v", err)
		}

		return errors.New(errorMessage)
	}

	return nil
}
