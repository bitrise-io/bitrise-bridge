package bridge

import (
	"encoding/base64"
	"testing"

	"github.com/bitrise-io/bitrise-bridge/config"
	"github.com/stretchr/testify/require"
)

const (
	testCommandHostID = CommandHostIDNone
	// testCommandHostID = CommandHostIDCmdBridge
)

func stringBase64Convert(strToConvert string) string {
	return base64.StdEncoding.EncodeToString([]byte(strToConvert))
}

func TestRun_RunMode(t *testing.T) {
	inventoryStr := `
envs:
  - MY_HOME: $HOME
    opts:
      is_expand: true
`

	inventoryBase64Str := stringBase64Convert(inventoryStr)
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

	configBase64Str := stringBase64Convert(configStr)
	t.Log("Config:", configBase64Str)

	t.Log("Perform - run")
	{
		// workflow param
		require.NoError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, inventoryBase64Str, configBase64Str,
			"", "target", false, "/"))
		// JSON param
		require.NoError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, inventoryBase64Str, configBase64Str,
			stringBase64Convert(`{"workflow":"target"}`), "", false, "/"))
	}

	t.Log("Perform - run without inventory")
	{
		// workflow param
		require.NoError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			"", "simple-success", false, ""))
		// JSON param
		require.NoError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			stringBase64Convert(`{"workflow":"simple-success"}`), "", false, ""))
	}

	t.Log("Perform - invalid workflow")
	{
		// workflow param
		require.EqualError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			"", "does-not-exist", false, ""),
			"exit status 1")
		// JSON param
		require.EqualError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			stringBase64Convert(`{"workflow":"does-not-exist"}`), "", false, ""),
			"exit status 1")
	}

	t.Log("Perform - fail-test")
	{
		// workflow param
		require.EqualError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			"", "fail-test", false, ""),
			"exit status 1")
		// JSON param
		require.EqualError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			stringBase64Convert(`{"workflow":"fail-test"}`), "", false, ""),
			"exit status 1")
	}
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
	{
		// pattern param
		require.NoError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, inventoryBase64Str, configBase64Str,
			"", "trig-target", true, ""))
		// JSON param
		require.NoError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, inventoryBase64Str, configBase64Str,
			stringBase64Convert(`{"pattern":"trig-target"}`), "", true, ""))
	}

	t.Log("Perform - no definition")
	{
		// pattern param
		require.EqualError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			"", "no-def", true, ""),
			"exit status 1")
		// JSON param
		require.EqualError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			stringBase64Convert(`{"pattern":"no-def"}`), "", true, ""),
			"exit status 1")
	}

	t.Log("Perform - fail-test")
	{
		// pattern param
		require.EqualError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			"", "trig-fail-test", true, ""),
			"exit status 1")
		// JSON param
		require.EqualError(t, PerformRunOrTrigger(testCommandHostID, config.Model{}, "", configBase64Str,
			stringBase64Convert(`{"pattern":"trig-fail-test"}`), "", true, ""),
			"exit status 1")
	}
}
