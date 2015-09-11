package cli

import (
	"encoding/base64"
	"testing"
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
	_, _, err := performRunOrTrigger(inventoryBase64Str, configBase64Str, "target", false, "/")
	if err != nil {
		t.Fatalf("Perform - run: %s", err)
	}

	t.Log("Perform - run without inventory")
	_, _, err = performRunOrTrigger("", configBase64Str, "simple-success", false, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Perform - invalid workflow")
	_, _, err = performRunOrTrigger("", configBase64Str, "does-not-exist", false, "")
	if err == nil {
		t.Fatal("Should fail for invalid workflow!")
	}

	t.Log("Perform - fail-test")
	_, _, err = performRunOrTrigger("", configBase64Str, "fail-test", false, "")
	if err == nil {
		t.Fatal("Should fail for failing (fail-test) workflow!")
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
	_, _, err := performRunOrTrigger(inventoryBase64Str, configBase64Str, "trig-target", true, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Perform - no definition")
	_, _, err = performRunOrTrigger("", configBase64Str, "no-def", true, "")
	if err == nil {
		t.Fatal("Should fail for failing (no-def) trigger pattern!")
	}

	t.Log("Perform - fail-test")
	_, _, err = performRunOrTrigger("", configBase64Str, "trig-fail-test", true, "")
	if err == nil {
		t.Fatal("Should fail for failing (trig-fail-test) build!")
	}
}
