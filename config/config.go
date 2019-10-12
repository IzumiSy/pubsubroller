package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/pkg/errors"
)

type Configuration struct {
	Variables map[string]string `yaml:"variables"`
	Topics    map[string]Topic  `yaml:"topics"`
}

type Topic struct {
	Subsciptions []Subscription `yaml:"subscriptions"`
}

type Subscription struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint,omitempty"`
	Pull     bool   `yaml:"pull,omitempty"`
}

func Load(string path) (*Configuration, error) {
	if path == "" {
		return nil, errors.New("Couldn't load config YAML")
	}

	yamlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configuration := Configuration{}
	err = yaml.Unmarshal(yamlBytes, &configuration)
	if err != nil {
		return nil, err
	}

	return &configuration, nil
}

func (config Configuration) Variables(projectId string) map[string]string {
	variables := make(map[string]string)
	for key, value := range configuration.Variables {
		variables[key] = strings.Replace(value, "${projectId}", projectId, -1)
	}
	return variables
}