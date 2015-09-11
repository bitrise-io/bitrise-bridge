package bridge

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/bitrise-io/go-utils/cmdex"
)

// // BitriseRunOrTrigger ...
// func BitriseRunOrTrigger(inventoryPth, configPth, workflowNameOrTriggerPattern string, isUseTrigger bool) error {
// 	logLevel := log.GetLevel().String()
// bitriseCommandToUse := "run"
// if isUseTrigger {
// 	bitriseCommandToUse = "trigger"
// }
// 	args := []string{"--loglevel", logLevel, bitriseCommandToUse, workflowNameOrTriggerPattern, "--path", configPth}
// 	if inventoryPth != "" {
// 		args = append(args, "--inventory", inventoryPth)
// 	}
// 	return cmdex.RunCommand("bitrise", args...)
// }

// CMDBridgeDoBitriseRunOrTrigger ...
func CMDBridgeDoBitriseRunOrTrigger(inventoryPth, configPth, workflowNameOrTriggerPattern string, isUseTrigger bool, workdirPath string) error {
	logLevel := log.GetLevel().String()

	bitriseCommandToUse := "run"
	if isUseTrigger {
		bitriseCommandToUse = "trigger"
	}

	params := fmt.Sprintf("bitrise --loglevel %s %s %s --path %s", logLevel, bitriseCommandToUse, workflowNameOrTriggerPattern, configPth)
	if inventoryPth != "" {
		params = params + fmt.Sprintf(" --inventory %s", inventoryPth)
	}

	args := []string{"-workdir", workdirPath, "-do", params}
	log.Infof("args: %#v", args)

	return cmdex.RunCommand("cmd-bridge", args...)
}
