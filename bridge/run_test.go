package bridge

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testCommandHostID = CommandHostIDNone
	// testCommandHostID = CommandHostIDCmdBridge
)

func TestRun_RunMode(t *testing.T) {
	inventoryStr := `
envs:
  - MY_HOME: $HOME
    opts:
      is_expand: true
`

	inventoryBytes := []byte(inventoryStr)
	inventoryBase64Str := base64.StdEncoding.EncodeToString(inventoryBytes)
	t.Log("Inventory:", inventoryBase64Str)

	configStr := `
format_version: 0.9.10
default_step_lib_source: "https://github.com/bitrise-io/bitrise-steplib.git"

workflows:
  target:
    title: target
    steps:
    - script:
        title: Should success
        inputs:
        - content: |
            #!/bin/bash
            set -v
            if [[ "$MY_HOME" != "$HOME" ]] ; then
              exit 1
            fi
  simple-success:
    title: "Simple success"
    steps:
    - script:
        title: Should success
        inputs:
        - content: exit 0
  fail-test:
    title: "Fail test"
    steps:
    - script:
        title: Should fail
        inputs:
        - content: exit 1
`

	configBytes := []byte(configStr)
	configBase64Str := base64.StdEncoding.EncodeToString(configBytes)
	t.Log("Config:", configBase64Str)

	t.Log("Perform - run")
	err := PerformRunOrTrigger(testCommandHostID, inventoryBase64Str, configBase64Str, "target", false, "/")
	require.NoError(t, err)

	t.Log("Perform - run without inventory")
	err = PerformRunOrTrigger(testCommandHostID, "", configBase64Str, "simple-success", false, "")
	require.NoError(t, err)

	t.Log("Perform - invalid workflow")
	err = PerformRunOrTrigger(testCommandHostID, "", configBase64Str, "does-not-exist", false, "")
	require.Error(t, err)

	t.Log("Perform - fail-test")
	err = PerformRunOrTrigger(testCommandHostID, "", configBase64Str, "fail-test", false, "")
	require.Error(t, err)
}

func TestRun_TriggerMode(t *testing.T) {
	inventoryStr := `
envs:
  - MY_HOME: $HOME
    opts:
      is_expand: true
`

	inventoryBytes := []byte(inventoryStr)
	inventoryBase64Str := base64.StdEncoding.EncodeToString(inventoryBytes)
	t.Log("Inventory:", inventoryBase64Str)

	configStr := `
format_version: 0.9.10
default_step_lib_source: "https://github.com/bitrise-io/bitrise-steplib.git"

trigger_map:
- pattern: trig-target
  is_pull_request_allowed: true
  workflow: target-wf
- pattern: trig-fail-test
  is_pull_request_allowed: true
  workflow: fail-test-wf

workflows:
  target-wf:
    title: target-wf
    steps:
    - script:
        title: Should success
        inputs:
        - content: |
            #!/bin/bash
            set -v
            if [[ "$MY_HOME" != "$HOME" ]] ; then
              exit 1
            fi
  fail-test-wf:
    title: "Fail test"
    steps:
    - script:
        title: Should fail
        inputs:
        - content: exit 1
`

	configBytes := []byte(configStr)
	configBase64Str := base64.StdEncoding.EncodeToString(configBytes)
	t.Log("Config:", configBase64Str)

	t.Log("Perform - simple OK")
	err := PerformRunOrTrigger(testCommandHostID, inventoryBase64Str, configBase64Str, "trig-target", true, "")
	require.NoError(t, err)

	t.Log("Perform - no definition")
	err = PerformRunOrTrigger(testCommandHostID, "", configBase64Str, "no-def", true, "")
	require.Error(t, err)

	t.Log("Perform - fail-test")
	err = PerformRunOrTrigger(testCommandHostID, "", configBase64Str, "trig-fail-test", true, "")
	require.Error(t, err)
}
