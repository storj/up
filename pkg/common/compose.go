package common

import (
	"github.com/goccy/go-yaml"
	"github.com/zeebo/errs/v2"
	"io/ioutil"
	"strconv"
	"strings"
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

// hasIndexedPrefix is true if s has prefix and remaining is a number
func hasIndexedPrefix(s string, prefix string) bool {
	if strings.HasPrefix(s, prefix) {
		_, err := strconv.Atoi(s[len(prefix):])
		if err == nil {
			return true
		}
	}
	return false
}

// FilterPrefix returns with the list of services which matches 'service' or starts with 'service' and has a number at the end.
func (c *SimplifiedCompose) FilterPrefix(service string) map[string]*ServiceConfig {
	filtered := map[string]*ServiceConfig{}
	for k, v := range c.Services {
		if k == service || hasIndexedPrefix(k, service) {
			filtered[k] = v
		}
	}

	return filtered
}

// FilterPrefixAndGroup returns with the list of services which either indexed/exactly the same or part of group definition
func (c *SimplifiedCompose) FilterPrefixAndGroup(service string, groups map[string][]string) map[string]*ServiceConfig {
	filtered := map[string]*ServiceConfig{}

	for _, part := range strings.Split(service, ",") {
		selector := strings.TrimSpace(part)
		for k, v := range c.Services {
			if selector == "all" || selector == k || hasIndexedPrefix(k, selector) {
				filtered[k] = v
			} else if group, found := Presets[selector]; found {
				for _, s := range group {
					if s == k || hasIndexedPrefix(k, s) {
						filtered[k] = v
					}
				}
			}
		}
	}
	return filtered
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
	for _, v := range content.FilterPrefixAndGroup(selector, Presets) {
		err := updater(v)
		if err != nil {
			return err
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
