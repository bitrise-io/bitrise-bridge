package bridge

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
)

const (
	// CommandHostIDCmdBridge ...
	CommandHostIDCmdBridge = "cmd-bridge"
	// CommandHostIDNone ...
	CommandHostIDNone = "none"
	// CommandHostIDDocker ...
	CommandHostIDDocker = "docker"
)

// PerformRunOrTrigger ...
func PerformRunOrTrigger(commandHostID string, hostSpecificArgs map[string]string, inventoryBase64Str, configBase64Str, workflowNameOrTriggerPattern string, isUseTrigger bool, workdirPath string) error {
	// Decode inventory & write to file
	if inventoryBase64Str != "" {
		_, err := base64.StdEncoding.DecodeString(inventoryBase64Str)
		if err != nil {
			return fmt.Errorf("Failed to decode base 64 string, error: %s", err.Error())
		}
	}

	// Decode bitrise config & write to file
	_, err := base64.StdEncoding.DecodeString(configBase64Str)
	if err != nil {
		return fmt.Errorf("Failed to decode base64 string, error: %s", err)
	}

	// run or trigger ?
	bitriseCommandToUse := "run"
	if isUseTrigger {
		bitriseCommandToUse = "trigger"
	}

	// Call bitrise
	switch commandHostID {
	case CommandHostIDCmdBridge:
		if err := performRunOrTriggerWithCmdBridge(bitriseCommandToUse, inventoryBase64Str, configBase64Str, workflowNameOrTriggerPattern, workdirPath); err != nil {
			return err
		}
	case CommandHostIDNone:
		if err := performRunOrTriggerWithoutCommandHost(bitriseCommandToUse, inventoryBase64Str, configBase64Str, workflowNameOrTriggerPattern, workdirPath); err != nil {
			return err
		}
	case CommandHostIDDocker:
		if err := performRunOrTriggerWithDocker(hostSpecificArgs, bitriseCommandToUse, inventoryBase64Str, configBase64Str, workflowNameOrTriggerPattern, workdirPath); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid / not supported commandHostID: %s", commandHostID)
	}

	return nil
}

func createBitriseCallArgs(bitriseCommandToUse, inventoryBase64, configBase64, workflowNameOrTriggerPattern string) []string {
	logLevel := log.GetLevel().String()

	retArgs := []string{
		"--loglevel", logLevel,
		bitriseCommandToUse, workflowNameOrTriggerPattern,
		"--config-base64", configBase64,
	}

	if inventoryBase64 != "" {
		retArgs = append(retArgs, "--inventory-base64", inventoryBase64)
	}

	return retArgs
}

func performRunOrTriggerWithoutCommandHost(bitriseCommandToUse, inventoryBase64, configBase64, workflowNameOrTriggerPattern, workdirPath string) error {
	bitriseCallArgs := createBitriseCallArgs(bitriseCommandToUse, inventoryBase64, configBase64, workflowNameOrTriggerPattern)
	log.Debugf("=> (debug) bitriseCallArgs: %s", bitriseCallArgs)

	if err := cmdex.RunCommandInDir(workdirPath, "bitrise", bitriseCallArgs...); err != nil {
		log.Debugf("cmd: `bitrise %s` failed, error: %s", bitriseCallArgs)
		return err
	}

	return nil
}

func performRunOrTriggerWithDocker(hostSpecificArgs map[string]string, bitriseCommandToUse, inventoryBase64, configBase64, workflowNameOrTriggerPattern, workdirPath string) error {
	dockerImageToUse := hostSpecificArgs["docker-image-id"]
	if dockerImageToUse == "" {
		return errors.New("No docker-image-id specified")
	}

	bitriseCallArgs := createBitriseCallArgs(bitriseCommandToUse, inventoryBase64, configBase64, workflowNameOrTriggerPattern)
	log.Debugf("=> (debug) bitriseCallArgs: %s", bitriseCallArgs)

	fullDockerArgs := []string{
		"run", "--rm", dockerImageToUse, "bitrise",
	}
	fullDockerArgs = append(fullDockerArgs, bitriseCallArgs...)

	if err := cmdex.RunCommand("docker", fullDockerArgs...); err != nil {
		log.Debugf("cmd: `docker %s` failed, error: %s", fullDockerArgs)
		return err
	}

	return nil
}

func performRunOrTriggerWithCmdBridge(bitriseCommandToUse, inventoryBase64, configBase64, workflowNameOrTriggerPattern, workdirPath string) error {
	bitriseCallArgs := createBitriseCallArgs(bitriseCommandToUse, inventoryBase64, configBase64, workflowNameOrTriggerPattern)
	log.Debugf("=> (debug) bitriseCallArgs: %s", bitriseCallArgs)

	bitriseCmdStr := fmt.Sprintf("bitrise %s", strings.Join(bitriseCallArgs, " "))
	args := []string{"-workdir", workdirPath, "-do", bitriseCmdStr}

	if err := cmdex.RunCommand("cmd-bridge", args...); err != nil {
		log.Debugf("cmd: `%s` failed, error: %s", bitriseCmdStr)
		return err
	}

	return nil
}
