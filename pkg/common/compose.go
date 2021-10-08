package common

import (
	"github.com/goccy/go-yaml"
	"github.com/zeebo/errs/v2"
	"io/ioutil"
)

type SimplifiedCompose struct {
	Version  string
	Services map[string]*ServiceConfig
}

func (c *SimplifiedCompose) GetService(service string) (*ServiceConfig, error) {
	for k, v := range c.Services {
		if k == service {
			return v, nil
		}
	}
	return nil, errs.Errorf("Service %s couldn't be found in the compose file", service)
}

func ReadCompose(file string) (SimplifiedCompose, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return SimplifiedCompose{}, err
	}
	return ParseCompose(content)
}

func ParseCompose(raw []byte) (SimplifiedCompose, error) {
	content := SimplifiedCompose{}
	err := yaml.Unmarshal(raw, &content)
	return content, err
}

func UpdateEach(selector string, updater func(service *ServiceConfig) error) error {
	in, err := ioutil.ReadFile("docker-compose.yaml")
	if err != nil {
		return err
	}
	content := SimplifiedCompose{}
	if err = yaml.Unmarshal(in, &content); err != nil {
		return err
	}
	for k, v := range content.Services {
		if Selected(selector, k) {
			err := updater(v)
			if err != nil {
				return err
			}
		}
	}
	out, err := yaml.Marshal(content)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile("docker-compose.yaml", out, 0644); err != nil {
		return err
	}
	return nil
}

func Update(selector string, updater func(compose *SimplifiedCompose) error) error {
	in, err := ioutil.ReadFile("docker-compose.yaml")
	if err != nil {
		return err
	}
	content := SimplifiedCompose{}
	if err = yaml.Unmarshal(in, &content); err != nil {
		return err
	}
	if err = updater(&content); err != nil {
		return err
	}
	out, err := yaml.Marshal(content)
	if err != nil {
		return err
	}
	if err = ioutil.WriteFile("docker-compose.yaml", out, 0644); err != nil {
		return err
	}
	return nil
}
