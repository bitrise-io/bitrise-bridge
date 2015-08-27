package bridge

import (
	"encoding/json"
	"fmt"

	bitriseModels "github.com/bitrise-io/bitrise/models/models_1_0_0"
	envmanModels "github.com/bitrise-io/envman/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"gopkg.in/yaml.v2"
)

// ----------------------------
// --- Bitrise Config Handling

// BitriseDataModelFromYAMLBytes ...
func BitriseDataModelFromYAMLBytes(bytes []byte) (config bitriseModels.BitriseDataModel, err error) {
	if err = yaml.Unmarshal(bytes, &config); err != nil {
		return
	}
	if err = config.Normalize(); err != nil {
		return
	}
	if err = config.Validate(); err != nil {
		return
	}
	if err = config.FillMissingDefaults(); err != nil {
		return
	}
	return
}

// BitriseDataModelFromJSONBytes ...
func BitriseDataModelFromJSONBytes(bytes []byte) (config bitriseModels.BitriseDataModel, err error) {
	if err = json.Unmarshal(bytes, &config); err != nil {
		return
	}
	if err = config.Normalize(); err != nil {
		return
	}
	if err = config.Validate(); err != nil {
		return
	}
	if err = config.FillMissingDefaults(); err != nil {
		return
	}
	return
}

// WriteBitriseDataModelToFile ...
func WriteBitriseDataModelToFile(pth string, config bitriseModels.BitriseDataModel) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	if err := fileutil.WriteBytesToFile(pth, bytes); err != nil {
		return err
	}
	return nil
}

// ----------------------------
// --- Bitrise Inventory Handling

// InventoryModelFromYAMLBytes ...
func InventoryModelFromYAMLBytes(bytes []byte) (envList envmanModels.EnvsYMLModel, err error) {
	if err = yaml.Unmarshal(bytes, &envList); err != nil {
		fmt.Printf("inventory bytes: %s ", bytes)
		return
	}
	return
}

// WriteInventoryModelToFile ...
func WriteInventoryModelToFile(pth string, envList envmanModels.EnvsYMLModel) error {
	bytes, err := yaml.Marshal(envList)
	if err != nil {
		return err
	}
	if err := fileutil.WriteBytesToFile(pth, bytes); err != nil {
		return err
	}
	return nil
}
