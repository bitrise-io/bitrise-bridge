package cli

import (
	"encoding/base64"
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/bitrise-bridge/bridge"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/codegangsta/cli"
)

func runWorkflow(inventoryBase64Str, configBase64Str, workflowName string) (string, string, error) {
	// Paths
	tempBitriseWorkDirPath, err := pathutil.NormalizedOSTempDirPath("bitrise")
	if err != nil {
		return "", "", err
	}
	inventoryFilePath := ""
	configFilePath := ""

	// Decode inventory & write to file
	if inventoryBase64Str != "" {
		inventoryBase64Bytes, err := base64.StdEncoding.DecodeString(inventoryBase64Str)
		if err != nil {
			return inventoryFilePath, configFilePath, fmt.Errorf("Failed to decode base 64 string, error: %s", err.Error())
		}

		inventory, err := bridge.InventoryModelFromYAMLBytes(inventoryBase64Bytes)
		if err != nil {
			return inventoryFilePath, configFilePath, fmt.Errorf("Failed to parse bitrise inventory, error: %s", err.Error())
		}

		inventoryFilePath = path.Join(tempBitriseWorkDirPath, "inventory.yml")
		if err := bridge.WriteInventoryModelToFile(inventoryFilePath, inventory); err != nil {
			return inventoryFilePath, configFilePath, fmt.Errorf("Failed to write bitrise inventory to file, error: %s", err.Error())
		}
	}

	// Decode bitrise config & write to file
	configBase64Bytes, err := base64.StdEncoding.DecodeString(configBase64Str)
	if err != nil {
		return inventoryFilePath, configFilePath, fmt.Errorf("Failed to decode base 64 string, error: %s", err)
	}

	config, err := bridge.BitriseDataModelFromYAMLBytes(configBase64Bytes)
	if err != nil {
		return inventoryFilePath, configFilePath, fmt.Errorf("Failed to parse bitrise config, error: %s", err.Error())
	}

	configFilePath = path.Join(tempBitriseWorkDirPath, "config.yml")
	if err := bridge.WriteBitriseDataModelToFile(configFilePath, config); err != nil {
		return inventoryFilePath, configFilePath, fmt.Errorf("Failed to write bitrise config to file, error: %s", err.Error())
	}

	// Call bitrise
	if err := bridge.CMDBridgeDoBitriseRun(inventoryFilePath, configFilePath, workflowName); err != nil {
		return inventoryFilePath, configFilePath, fmt.Errorf("cmd: `bitrise run %s --inventory %s --path %s` failed, error: %s", workflowName, inventoryFilePath, configFilePath, err.Error())
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

	if inventoryFilePath, configFilePath, err := runWorkflow(inventoryBase64Str, configBase64Str, workflowName); err != nil {
		log.Warn("Failed to run workflow, params:")
		log.Warnf("Inventory path (%s)", inventoryFilePath)
		log.Warnf("Config path (%s)", configFilePath)
		log.Fatalln("Error:", err)
	}
}
