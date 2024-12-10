package mobileappcontent

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/text/encoding/unicode"
)

type AppxMsixMetadata struct {
	XMLName  xml.Name `xml:"Package" tfsdk:"-"`
	Identity struct {
		Name                  string `xml:",attr" tfsdk:"name"`
		ProcessorArchitecture string `xml:",attr" tfsdk:"processor_architecture"`
		Publisher             string `xml:",attr" tfsdk:"publisher"`
		PublisherHash         string `xml:"-" tfsdk:"publisher_hash"`
		Version               string `xml:",attr" tfsdk:"version"`
	} `tfsdk:"identity"`
	Properties struct {
		Framework            bool   `tfsdk:"framework"`
		DisplayName          string `tfsdk:"display_name"`
		PublisherDisplayName string `tfsdk:"publisher_display_name"`
		Description          string `tfsdk:"description"`
		Logo                 string `tfsdk:"logo"`
	} `tfsdk:"properties"`
	Logo struct {
		Type        string `tfsdk:"type"`
		ValueBase64 string `tfsdk:"value_base64"`
	} `xml:"-" tfsdk:"logo"`
}

func ParseAppxMsixMetadata(ctx context.Context, appxMsixPath string) (appxMsixMetadata AppxMsixMetadata, err error) {

	const kLogPref = "mobileappcontent.ParseAppxMsixMetadata: "

	// read manifest
	appxMsixMetadataAny, err := runOnZipFile(appxMsixPath, "AppxManifest.xml",
		func(reader io.Reader) (result any, err error) {
			xmlRaw, err := io.ReadAll(reader)
			if err != nil {
				return
			}
			var metadata AppxMsixMetadata
			if err = xml.Unmarshal(xmlRaw, &metadata); err != nil {
				return
			}
			result = metadata
			return
		})
	if err != nil {
		return
	}
	appxMsixMetadata = appxMsixMetadataAny.(AppxMsixMetadata)
	tflog.Trace(ctx, kLogPref, map[string]any{"appxMsixMetadata": appxMsixMetadata})

	appxMsixMetadata.Identity.PublisherHash, err = AppxMsixPublisherHashFromPublisher(ctx, appxMsixMetadata.Identity.Publisher)
	if err != nil {
		return
	}

	// check if logo has been declared and exit if not
	if appxMsixMetadata.Properties.Logo == "" {
		return
	}

	// try to read logo
	logoPath := strings.ReplaceAll(appxMsixMetadata.Properties.Logo, "\\", "/")
	tflog.Trace(ctx, kLogPref+"logoPath: "+logoPath)
	logoBytesAny, err := runOnZipFile(appxMsixPath, logoPath,
		func(reader io.Reader) (result any, err error) {
			return io.ReadAll(reader)
		})
	if err != nil {
		// convert error to warning to still allow retrieving other data from erroneous appx/msix
		tflog.Warn(ctx, kLogPref+fmt.Sprintf("Unable to read logo from path %s: %v", logoPath, err))
		err = nil
		return
	}
	logoBase64 := base64.StdEncoding.EncodeToString(logoBytesAny.([]byte))
	tflog.Trace(ctx, kLogPref+"logoBase64: "+logoBase64)
	logoType := mime.TypeByExtension(filepath.Ext(logoPath))
	tflog.Trace(ctx, kLogPref+"logoType: "+logoType)

	appxMsixMetadata.Logo.ValueBase64 = logoBase64
	appxMsixMetadata.Logo.Type = logoType

	return
}

func AppxMsixPublisherHashFromPublisher(ctx context.Context, publisher string) (publisherHash string, err error) {
	const kLogPref = "mobileappcontent.AppxMsixPublisherIdFromPublisher: "

	utf16Enc := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	publisherBytes, err := utf16Enc.Bytes([]byte(publisher))
	if err != nil {
		return
	}
	tflog.Trace(ctx, kLogPref+fmt.Sprintf("publisherBytes: %v", publisherBytes))

	sha256 := sha256.New()
	sha256.Write(publisherBytes)
	publisherSha256 := sha256.Sum(nil)[0:8]
	tflog.Trace(ctx, kLogPref+fmt.Sprintf("publisherSha256: %v", publisherSha256))

	const crockfordChars = "0123456789abcdefghjkmnpqrstvwxyz"
	var crockfordEnc = base32.NewEncoding(crockfordChars).WithPadding(base32.NoPadding)
	publisherHash = crockfordEnc.EncodeToString(publisherSha256)
	tflog.Trace(ctx, kLogPref+fmt.Sprintf("publisherId: %v", publisherHash))

	return
}
