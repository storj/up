package cmd

import (
	"io/ioutil"
	"path"

	composeloader "github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
)

const (
	mainConfigDefaultPath = "docker-compose.yml"
)

type LoadParams struct {
	WorkDir string
	// if empty docker-compose.yml used
	MainConfigFilePath string
}

type Loader interface {
	Load(params LoadParams) (*types.Project, error)
}

func NewLoader() Loader {
	return &loader{}
}

type loader struct {
}

func (l *loader) Load(params LoadParams) (*types.Project, error) {
	mainConfigPath := mainConfigDefaultPath
	if params.MainConfigFilePath != "" {
		mainConfigPath = params.MainConfigFilePath
	}

	b, err := ioutil.ReadFile(path.Join(params.WorkDir, mainConfigPath))
	if err != nil {
		return nil, err
	}
	config, err := composeloader.ParseYAML(b)
	if err != nil {
		return nil, err
	}

	return composeloader.Load(types.ConfigDetails{
		WorkingDir: params.WorkDir,
		ConfigFiles: []types.ConfigFile{
			{
				Filename: mainConfigPath,
				Config:   config,
			},
		},
		Environment: nil,
	})
}
