package cli

import (
	"encoding/base64"
	"testing"
)

func TestRun(t *testing.T) {
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
format_version: 0.9.11
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
`

	configBytes := []byte(configStr)
	configBase64Str := base64.StdEncoding.EncodeToString(configBytes)
	t.Log("Config:", configBase64Str)

	workflowName := "target"

	// Test
	_, _, err := runWorkflow(inventoryBase64Str, configBase64Str, workflowName)
	if err != nil {
		t.Fatal(err)
	}
}
