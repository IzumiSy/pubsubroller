package main

import (
	"context"
	"flag"
	"fmt"
	config "pubsubroller/config"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

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

	configuration, err := config.Load(configFilePath)
	if err != nil {
		panic(err)
	}

	if projectId == "" {
		fmt.Println("Error: GCP project ID required with `-projectId` option.")
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
