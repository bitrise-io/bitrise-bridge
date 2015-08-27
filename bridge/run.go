package bridge

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
)

// BitriseRun ...
func BitriseRun(inventoryPth, configPth, workflowName string) error {
	logLevel := log.GetLevel().String()
	args := []string{"--loglevel", logLevel, "run", workflowName, "--path", configPth}
	if inventoryPth != "" {
		args = append(args, "--inventory", inventoryPth)
	}
	return cmdex.RunCommand("bitrise", args...)
}

// CMDBridgeDoBitriseRun ...
func CMDBridgeDoBitriseRun(inventoryPth, configPth, workflowName string) error {
	logLevel := log.GetLevel().String()

	params := fmt.Sprintf("bitrise --loglevel %s run %s --path %s", logLevel, workflowName, configPth)
	if inventoryPth != "" {
		params = params + fmt.Sprintf(" --inventory %s", inventoryPth)
	}

	args := []string{"-do", params}

	return cmdex.RunCommand("cmd-bridge", args...)
}
