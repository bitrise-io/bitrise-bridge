package config

import (
	"encoding/json"
	"path"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// DockerConfigModel ...
type DockerConfigModel struct {
	// Image
	//  docker image to use
	Image string `json:"image"`
	// Volumes
	//  in docker format: /path/on/host:/path/in/container
	Volumes []string `json:"volumes"`
	// IsAllowAccessToDockerInContainer
	//  shares the docker socker & docker binary with the container
	IsAllowAccessToDockerInContainer bool `json:"allow_access_to_docker_in_container"`
	// AdditionalRunArguments
	//  additional arguments for the `docker run .. ` command,
	//  appended (!) after other arguments (e.g. volumes)
	AdditionalRunArguments []string `json:"additional_run_arguments"`
}

// Model ...
type Model struct {
	Docker DockerConfigModel `json:"docker"`
}

// DefaultConfigFilePath ...
func DefaultConfigFilePath() string {
	return path.Join(pathutil.UserHomeDir(), ".bitrise-bridge", "config.json")
}

// LoadConfigFromFile ...
func LoadConfigFromFile(pth string) (Model, error) {
	bytes, err := fileutil.ReadBytesFromFile(pth)
	if err != nil {
		return Model{}, err
	}

	var m Model
	if err := json.Unmarshal(bytes, &m); err != nil {
		return Model{}, err
	}

	return m, nil
}
