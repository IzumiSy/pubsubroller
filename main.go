package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	gcp "pubsubroller/adapters/google"
	config "pubsubroller/config"
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

	var internalClient *pubsub.Client
	var cerr error

	if endpoint != "" {
		internalClient, cerr = pubsub.NewClient(ctx, projectId, option.WithEndpoint(endpoint))
	} else {
		internalClient, cerr = pubsub.NewClient(ctx, projectId)
	}

	if cerr != nil {
		fmt.Println("Error on initializing pubsub client:", err.Error())
		return
	}

	client := gcp.PubsubClient{
		Client: internalClient,
		Ctx:    ctx,
	}

	// 実行オプションを作成して実行

	opts := Options{
		IsDryRun:  isDryRun,
		Variables: variables,
	}

	if isDeleteMode {
		deleteTopics(client, deleteTopicsLogger{}, ctx, configuration, opts)
		deleteSubscriptions(client, deleteSubscriptionLogger{}, ctx, configuration, opts)
	} else {
		createTopics(client, createTopicsLogger{}, ctx, configuration, opts)
		createSubscriptions(client, createSubscriptionsLogger{}, ctx, configuration, opts)
	}

	return
}
