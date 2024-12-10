package mobileappcontent

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type IntunewinMetadata struct {
	XMLName                xml.Name `xml:"ApplicationInfo" tfsdk:"-"`
	Name                   string   `tfsdk:"name"`
	UnencryptedContentSize int64    `tfsdk:"unencrypted_content_size"`
	FileName               string   `tfsdk:"file_name"`
	SetupFile              string   `tfsdk:"setup_file"`
	EncryptionInfo         struct {
		XMLName              xml.Name `tfsdk:"-"`
		EncryptionKey        string   `tfsdk:"encryption_key"`
		MacKey               string   `tfsdk:"mac_key"`
		InitializationVector string   `tfsdk:"initialization_vector"`
		Mac                  string   `tfsdk:"mac"`
		ProfileIdentifier    string   `tfsdk:"profile_identifier"`
		FileDigest           string   `tfsdk:"file_digest"`
		FileDigestAlgorithm  string   `tfsdk:"file_digest_algorithm"`
	} `tfsdk:"encryption_info"`
	MsiInfo struct {
		XMLName                       xml.Name `tfsdk:"-"`
		MsiProductCode                string   `tfsdk:"msi_product_code"`
		MsiProductVersion             string   `tfsdk:"msi_product_version"`
		MsiPackageCode                string   `tfsdk:"msi_package_code"`
		MsiUpgradeCode                string   `tfsdk:"msi_upgrade_code"`
		MsiExecutionContext           string   `tfsdk:"msi_execution_context"`
		MsiRequiresLogon              bool     `tfsdk:"msi_requires_logon"`
		MsiRequiresReboot             bool     `tfsdk:"msi_requires_reboot"`
		MsiIsMachineInstall           bool     `tfsdk:"msi_is_machine_install"`
		MsiIsUserInstall              bool     `tfsdk:"msi_is_user_install"`
		MsiIncludesServices           bool     `tfsdk:"msi_includes_services"`
		MsiIncludesODBCDataSource     bool     `tfsdk:"msi_includes_odbc_data_source"`
		MsiContainsSystemRegistryKeys bool     `tfsdk:"msi_contains_system_registry_keys"`
		MsiContainsSystemFolders      bool     `tfsdk:"msi_contains_system_folders"`
		MsiPublisher                  string   `tfsdk:"msi_publisher"`
	} `tfsdk:"msi_info"`
	EncryptedContentSize int64 `tfsdk:"-"`
}

func ParseAndExtractIntunewinOuter(ctx context.Context, intuneWinPath string) (intunewinMetadata IntunewinMetadata, fileEncryptionInfo FileEncryptionInfo, encryptedContentPath string, err error) {

	const kLogPref = "mobileappcontent.ParseAndExtractIntunewin: "

	// get encryption info first

	intunewinMetadata, err = ParseIntunewinMetadata(ctx, intuneWinPath)
	if err != nil {
		return
	}
	if intunewinMetadata.FileName == "" || intunewinMetadata.UnencryptedContentSize == 0 || intunewinMetadata.EncryptionInfo.EncryptionKey == "" {
		err = fmt.Errorf("file Detection.xml not valid")
		return
	}

	fileEncryptionInfo.FileDigest = intunewinMetadata.EncryptionInfo.FileDigest
	fileEncryptionInfo.FileDigestAlgorithm = intunewinMetadata.EncryptionInfo.FileDigestAlgorithm
	fileEncryptionInfo.EncryptionKey = intunewinMetadata.EncryptionInfo.EncryptionKey
	fileEncryptionInfo.InitializationVector = intunewinMetadata.EncryptionInfo.InitializationVector
	fileEncryptionInfo.Mac = intunewinMetadata.EncryptionInfo.Mac
	fileEncryptionInfo.MacKey = intunewinMetadata.EncryptionInfo.MacKey
	fileEncryptionInfo.ProfileIdentifier = intunewinMetadata.EncryptionInfo.ProfileIdentifier

	// now get encrypted inner content

	contentPathInArchive := fmt.Sprintf("IntuneWinPackage/Contents/%s", intunewinMetadata.FileName)
	enrcyptedContentPathRaw, err := runOnZipFile(intuneWinPath, contentPathInArchive,
		func(reader io.Reader) (result any, err error) {
			mobileAppContentFile, err := os.CreateTemp("", kEncryptedContentTempName)
			if err != nil {
				return
			}
			defer mobileAppContentFile.Close()

			if _, err = io.Copy(mobileAppContentFile, reader); err != nil {
				return
			}

			return mobileAppContentFile.Name(), nil
		})
	if err != nil {
		return
	}

	encryptedContentPath = enrcyptedContentPathRaw.(string)
	tflog.Trace(ctx, kLogPref, map[string]any{"encryptedContentPath": encryptedContentPath})

	fi, err := os.Stat(encryptedContentPath)
	if err != nil {
		return
	}
	intunewinMetadata.EncryptedContentSize = fi.Size()

	return
}

func ParseIntunewinMetadata(ctx context.Context, intuneWinPath string) (intunewinMetadata IntunewinMetadata, err error) {

	const kLogPref = "mobileappcontent.ReadIntunewinDetectionXml: "

	intunewinMetadataAny, err := runOnZipFile(intuneWinPath, "IntuneWinPackage/Metadata/Detection.xml",
		func(reader io.Reader) (result any, err error) {
			xmlRaw, err := io.ReadAll(reader)
			if err != nil {
				return
			}
			var metadata IntunewinMetadata
			if err = xml.Unmarshal(xmlRaw, &metadata); err != nil {
				return
			}
			result = metadata
			return
		})
	if err != nil {
		return
	}

	intunewinMetadata = intunewinMetadataAny.(IntunewinMetadata)
	tflog.Trace(ctx, kLogPref, map[string]any{"intunewinMetadata": intunewinMetadata})

	return
}
