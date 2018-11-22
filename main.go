package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
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

type Options struct {
	IsDryRun  bool
	Variables map[string]string
}

func main() {
	projectIdPtr := flag.String("projectId", "", "target GCP project ID")
	configFilePathPtr := flag.String("config", "", "configuration file path")
	endpointPtr := flag.String("endpoint", "", "service endpoint")
	isDryRunPtr := flag.Bool("dry", false, "dry run")
	isDeleteModePtr := flag.Bool("delete", false, "delete all topics and their subscriptions")
	flag.Parse()

	projectId := *projectIdPtr
	configFilePath := *configFilePathPtr
	endpoint := *endpointPtr
	isDryRun := *isDryRunPtr
	isDeleteMode := *isDeleteModePtr

	// projectIdとconfigFilePathは必須パラメータ

	if projectId == "" {
		fmt.Println("Error: GCP project ID required with `-projectId` option.")
		return
	} else if configFilePath == "" {
		fmt.Println("Error: no configuration file specified.")
		return
	}

	yamlBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Target project ID: %s\n\n", projectId)

	configuration := Configuration{}
	err = yaml.Unmarshal(yamlBytes, &configuration)
	if err != nil {
		panic(err)
	}

	// variablesの取得

	variables := make(map[string]string)
	for key, value := range configuration.Variables {
		_value := strings.Replace(value, "${projectId}", projectId, -1)
		variables[key] = _value
		fmt.Println(key, "=", _value)
	}

	// クライアント生成
	// endpointが指定されていればここでクライアントに設定する

	var opt option.ClientOption
	if endpoint != "" {
		opt = option.WithEndpoint(endpoint)
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if opt != nil {
		client, err = pubsub.NewClient(ctx, projectId, opt)
	}
	if err != nil {
		fmt.Println("Error on initializing pubsub client:", err.Error())
		return
	}

	// 実行オプションを作成して実行

	opts := Options{
		IsDryRun:  isDryRun,
		Variables: variables,
	}

	if isDeleteMode {
		deleteTopics(client, ctx, configuration, opts)
		deleteSubscriptions(client, ctx, configuration, opts)
	} else {
		createTopics(client, ctx, configuration, opts)
		createSubscriptions(client, ctx, configuration, opts)
	}

	return
}
