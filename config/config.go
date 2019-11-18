package config

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Configuration struct {
	Internal_Variables_ map[string]string `yaml:"variables"`
	Internal_Topics_    map[string]Topic  `yaml:"topics"`
}

type Topic struct {
	Internal_Subscriptions_ []Subscription `yaml:"subscriptions"`
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
	for key, value := range config.Internal_Variables_ {
		variables[key] = strings.Replace(value, "${projectId}", projectId, -1)
	}
	return variables
}

func (config Configuration) Topics() map[string]Topic {
	return config.Internal_Topics_
}

func (topic Topic) Subscriptions() []Subscription {
	return topic.Internal_Subscriptions_
}
