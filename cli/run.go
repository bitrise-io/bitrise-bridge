package cli

import (
	"encoding/base64"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-bridge/bridge"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/codegangsta/cli"
)

func run(c *cli.Context) {
	// Input validation
	configBase64Str := c.String(ConfigDataKey)
	if configBase64Str == "" {
		log.Fatal("Missing config data")
	}

	inventoryBase64Str := c.String(InventoryDataKey)

	// Decode bitrise config
	configBase64Data, err := base64.StdEncoding.DecodeString(configBase64Str)
	if err != nil {
		log.Fatal("Failed to encode base 64 string:", err)
	}

	if inventoryBase64Str != "" {

	}

	config, err := bridge.BitriseDataModelFromYAMLBytes(configBase64Data)
	if err != nil {
		log.Fatal("Failed to parse bitrise config:", err)
	}

	// Write config to file
	tempBitriseConfigDirPath, err := pathutil.NormalizedOSTempDirPath("config")
	if err != nil {
		log.Fatal("Failed to create temp bitrise config dir path:", err)
	}
	tempBitriseConfigFilePath := path.Join(tempBitriseConfigDirPath, "config.yml")

	if err := bridge.WriteBitriseDataModel(tempBitriseConfigFilePath, config); err != nil {
		log.Fatal("Failed to write bitrise config to file:", err)
	}
}
