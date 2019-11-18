package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Configuration struct {
	variables map[string]string `yaml:"variables"`
	topics    map[string]Topic  `yaml:"topics"`
}

type Topic struct {
	subscriptions []Subscription `yaml:"subscriptions"`
}

type Subscription struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint,omitempty"`
	Pull     bool   `yaml:"pull,omitempty"`
}

func Load(path string) (Configuration, error) {
	configuration := Configuration{}

	if path == "" {
		return configuration, errors.New("Couldn't load config YAML")
	}

	yamlBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return configuration, err
	}

	err = yaml.UnmarshalStrict(yamlBytes, &configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}

func (config Configuration) Variables(projectId string) map[string]string {
	variables := make(map[string]string)
	for key, value := range config.variables {
		variables[key] = strings.Replace(value, "${projectId}", projectId, -1)
	}
	return variables
}

func (config Configuration) Topics() map[string]Topic {
	return config.topics
}

func (topic Topic) Subscriptions() []Subscription {
	return topic.subscriptions
}
