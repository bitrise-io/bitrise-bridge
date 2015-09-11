package cli

import (
	"encoding/base64"
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-bridge/bridge"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/codegangsta/cli"
)

func performRunOrTrigger(inventoryBase64Str, configBase64Str, workflowNameOrTriggerPattern string, isUseTrigger bool, workdirPath string) (string, string, error) {
	// Paths
	tempBitriseWorkDirPath, err := pathutil.NormalizedOSTempDirPath("bitrise-bridge")
	if err != nil {
		return "", "", err
	}
	inventoryFilePath := ""
	configFilePath := ""

	// Decode inventory & write to file
	if inventoryBase64Str != "" {
		inventoryDecodedBytes, err := base64.StdEncoding.DecodeString(inventoryBase64Str)
		if err != nil {
			return inventoryFilePath, configFilePath, fmt.Errorf("Failed to decode base 64 string, error: %s", err.Error())
		}

		inventoryFilePath = path.Join(tempBitriseWorkDirPath, "inventory.yml")
		if err := fileutil.WriteBytesToFile(inventoryFilePath, inventoryDecodedBytes); err != nil {
			return inventoryFilePath, configFilePath, fmt.Errorf("Failed to write bitrise inventory to file, error: %s", err.Error())
		}
	}

	// Decode bitrise config & write to file
	configDecodedBytes, err := base64.StdEncoding.DecodeString(configBase64Str)
	if err != nil {
		return inventoryFilePath, configFilePath, fmt.Errorf("Failed to decode base 64 string, error: %s", err)
	}

	configFilePath = path.Join(tempBitriseWorkDirPath, "config.yml")
	if err := fileutil.WriteBytesToFile(configFilePath, configDecodedBytes); err != nil {
		return inventoryFilePath, configFilePath, fmt.Errorf("Failed to write bitrise config to file, error: %s", err.Error())
	}

	// Call bitrise
	if err := bridge.CMDBridgeDoBitriseRunOrTrigger(inventoryFilePath, configFilePath, workflowNameOrTriggerPattern, isUseTrigger, workdirPath); err != nil {
		bitriseCommandToUse := "run"
		if isUseTrigger {
			bitriseCommandToUse = "trigger"
		}
		log.Debugf("cmd: `bitrise %s %s --workdir %s --inventory %s --path %s` failed, error: %s",
			bitriseCommandToUse, workflowNameOrTriggerPattern,
			workdirPath, inventoryFilePath, configFilePath, err)
		return inventoryFilePath, configFilePath, err
	}

	return inventoryFilePath, configFilePath, nil
}

func run(c *cli.Context) {
	// Input validation
	inventoryBase64Str := c.String(InventoryDataKey)

	configBase64Str := c.String(ConfigDataKey)
	if configBase64Str == "" {
		log.Fatal("Missing required config data")
	}

	workflowName := c.String(WorkflowNameKey)
	if workflowName == "" {
		log.Fatal("Missing required workflow name")
	}

	if _, _, err := performRunOrTrigger(inventoryBase64Str, configBase64Str, workflowName, false, c.String(WorkdirPathKey)); err != nil {
		log.Fatalf("Failed to run, error: %s", err)
	}
}
