package bridge

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-tools/bitrise-bridge/config"
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
func PerformRunOrTrigger(commandHostID string, bridgeConfig config.Model, inventoryBase64Str, configBase64Str, runParamJSONBase64, workflowNameOrTriggerPattern string, isUseTrigger bool, workdirPath string) error {
	// Decode inventory - for error/encoding check
	if inventoryBase64Str != "" {
		_, err := base64.StdEncoding.DecodeString(inventoryBase64Str)
		if err != nil {
			return fmt.Errorf("Failed to decode base 64 string, error: %s", err.Error())
		}
	}

	// Decode bitrise config - for error/encoding check
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
		if err := performRunOrTriggerWithCmdBridge(bitriseCommandToUse, inventoryBase64Str, configBase64Str, runParamJSONBase64, workflowNameOrTriggerPattern, workdirPath); err != nil {
			return err
		}
	case CommandHostIDNone:
		if err := performRunOrTriggerWithoutCommandHost(bitriseCommandToUse, inventoryBase64Str, configBase64Str, runParamJSONBase64, workflowNameOrTriggerPattern, workdirPath); err != nil {
			return err
		}
	case CommandHostIDDocker:
		if err := performRunOrTriggerWithDocker(bridgeConfig, bitriseCommandToUse, inventoryBase64Str, configBase64Str, runParamJSONBase64, workflowNameOrTriggerPattern, workdirPath); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid / not supported commandHostID: %s", commandHostID)
	}

	return nil
}

func createBitriseCallArgs(bitriseCommandToUse, inventoryBase64, configBase64, runParamJSONBase64, workflowNameOrTriggerPattern string) []string {
	logLevel := log.GetLevel().String()

	retArgs := []string{
		"--loglevel", logLevel,
	}

	if len(runParamJSONBase64) > 0 {
		// new style, all params in one (Base64 encoded) JSON
		retArgs = append(retArgs, bitriseCommandToUse, "--json-params-base64", runParamJSONBase64)
	} else {
		// old style, separate params
		retArgs = append(retArgs, bitriseCommandToUse, workflowNameOrTriggerPattern)
	}

	// config / bitrise.yml
	retArgs = append(retArgs, "--config-base64", configBase64)

	// inventory / secrets
	if inventoryBase64 != "" {
		retArgs = append(retArgs, "--inventory-base64", inventoryBase64)
	}

	return retArgs
}

func performRunOrTriggerWithoutCommandHost(bitriseCommandToUse, inventoryBase64, configBase64, runParamJSONBase64, workflowNameOrTriggerPattern, workdirPath string) error {
	bitriseCallArgs := createBitriseCallArgs(bitriseCommandToUse, inventoryBase64, configBase64, runParamJSONBase64, workflowNameOrTriggerPattern)
	log.Debugf("=> (debug) bitriseCallArgs: %s", bitriseCallArgs)

	if err := cmdex.RunCommandInDir(workdirPath, "bitrise", bitriseCallArgs...); err != nil {
		log.Debugf("cmd: `bitrise %s` failed, error: %s", bitriseCallArgs)
		return err
	}

	return nil
}

func performRunOrTriggerWithDocker(bridgeConfig config.Model, bitriseCommandToUse, inventoryBase64, configBase64, runParamJSONBase64, workflowNameOrTriggerPattern, workdirPath string) error {
	dockerParamImageToUse := bridgeConfig.Docker.Image
	if dockerParamImageToUse == "" {
		return errors.New("No docker-image-id specified")
	}
	dockerParamIsAllowAccessToDockerInContainer := bridgeConfig.Docker.IsAllowAccessToDockerInContainer

	bitriseCallArgs := createBitriseCallArgs(bitriseCommandToUse, inventoryBase64, configBase64, runParamJSONBase64, workflowNameOrTriggerPattern)
	log.Debugf("=> (debug) bitriseCallArgs: %s", bitriseCallArgs)

	fullDockerArgs := []string{
		"run", "--rm",
	}
	if dockerParamIsAllowAccessToDockerInContainer {
		// mount the docker.sock socker & the docker binary as volumes, to make it
		//  accessible inside the container
		dockerPth, err := cmdex.RunCommandAndReturnStdout("which", "docker")
		if err != nil || dockerPth == "" {
			return errors.New("Failed to determin docker binary path; required for the 'docker-allow-access-to-docker-in-container' option")
		}
		fullDockerArgs = append(fullDockerArgs,
			"-v", "/var/run/docker.sock:/var/run/docker.sock",
			"-v", fmt.Sprintf("%s:%s", dockerPth, "/bin/docker"),
		)
	}
	if len(bridgeConfig.Docker.Volumes) > 0 {
		for _, aVolDef := range bridgeConfig.Docker.Volumes {
			fullDockerArgs = append(fullDockerArgs, "-v", aVolDef)
		}
	}
	if len(bridgeConfig.Docker.AdditionalRunArguments) > 0 {
		fullDockerArgs = append(fullDockerArgs, bridgeConfig.Docker.AdditionalRunArguments...)
	}
	// these are the docker specific params
	fullDockerArgs = append(fullDockerArgs, dockerParamImageToUse)
	// append Bitrise specific params
	fullDockerArgs = append(fullDockerArgs, "bitrise")
	fullDockerArgs = append(fullDockerArgs, bitriseCallArgs...)

	log.Debugf("fullDockerArgs: %#v", fullDockerArgs)

	if err := cmdex.RunCommand("docker", fullDockerArgs...); err != nil {
		log.Debugf("cmd: `docker %s` failed, error: %s", fullDockerArgs)
		return err
	}

	return nil
}

func performRunOrTriggerWithCmdBridge(bitriseCommandToUse, inventoryBase64, configBase64, runParamJSONBase64, workflowNameOrTriggerPattern, workdirPath string) error {
	bitriseCallArgs := createBitriseCallArgs(bitriseCommandToUse, inventoryBase64, configBase64, runParamJSONBase64, workflowNameOrTriggerPattern)
	log.Debugf("=> (debug) bitriseCallArgs: %s", bitriseCallArgs)

	bitriseCmdStr := fmt.Sprintf("bitrise %s", strings.Join(bitriseCallArgs, " "))
	args := []string{"-workdir", workdirPath, "-do", bitriseCmdStr}

	if err := cmdex.RunCommand("cmd-bridge", args...); err != nil {
		log.Debugf("cmd: `%s` failed, error: %s", bitriseCmdStr)
		return err
	}

	return nil
}
