package main

import (
	"context"
	"flag"
	"fmt"
	config "pubsubroller/config"

	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type Options struct {
	IsDryRun  bool
	Variables map[string]string
}

func main() {
	projectIdPtr := flag.String("projectId", "", "target GCP project ID")
	configFilePathPtr := flag.String("config", "", "configuration file path (Required)")
	endpointPtr := flag.String("endpoint", "", "service endpoint")
	isDryRunPtr := flag.Bool("dry", false, "dry run")
	isDeleteModePtr := flag.Bool("delete", false, "delete all topics and their subscriptions")
	flag.Parse()

	configFilePath := *configFilePathPtr
	endpoint := *endpointPtr
	isDryRun := *isDryRunPtr
	isDeleteMode := *isDeleteModePtr

	ctx := context.Background()

	// configFilePathは必須パラメータ

	if len(configFilePath) == 0 {
		fmt.Println("Make sure you give -config flag which is required.")
		fmt.Println("You can see more with -help option")
		return
	}

	configuration, err := config.Load(configFilePath)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	var projectId string

	if projectIdPtr != nil {
		projectId = *projectIdPtr
	} else {
		credentials, err := google.FindDefaultCredentials(ctx)
		if err != nil {
			panic(err)
		}
		if credentials == nil {
			fmt.Println("Error: invalid credential")
			return
		}
		projectId = credentials.ProjectID
	}

	if len(projectId) == 0 {
		fmt.Println("Error: Project ID must not be empty")
		return
	}

	fmt.Printf("Target project ID: %s\n\n", projectId)
	variables := configuration.Variables(projectId)

	// クライアント生成
	// endpointが指定されていればここでクライアントに設定する

	var opt option.ClientOption
	if endpoint != "" {
		opt = option.WithEndpoint(endpoint)
	}

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
