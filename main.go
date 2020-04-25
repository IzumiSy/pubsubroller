package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	flags "github.com/jessevdk/go-flags"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	gcp "pubsubroller/adapters/google"
	config "pubsubroller/config"
)

type appOptions struct {
	IsDryRun  bool
	Variables map[string]string
}

var opts struct {
	ProjectID      string `short:"p" long:"projectId" description:"target GCP project ID"`
	ConfigFilePath string `short:"c" long:"config" description:"configuration file path" required:"true"`
	Endpoint       string `short:"e" long:"endpoint" description:"service endpoint"`
	DryRun         bool   `long:"dry" description:"dry run"`
	Delete         bool   `long:"delete" description:"delete all topics and their subscriptions"`
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	configuration, err := config.Load(opts.ConfigFilePath)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		return
	}

	ctx := context.Background()
	projectID := opts.ProjectID

	if projectID == "" {
		credentials, err := google.FindDefaultCredentials(ctx)
		if err != nil {
			panic(err)
		}
		if credentials == nil {
			fmt.Println("Error: invalid credential")
			return
		}
		projectID = credentials.ProjectID

		if projectID == "" {
			fmt.Println("Error: Failed retrieving project ID from gcloud credentials.")
			fmt.Println("Could you manually provide your project ID with --projectID option?")
			return
		}
	}

	var internalClient *pubsub.Client
	var cerr error

	if opts.Endpoint != "" {
		internalClient, cerr = pubsub.NewClient(ctx, projectID, option.WithEndpoint(opts.Endpoint))
	} else {
		internalClient, cerr = pubsub.NewClient(ctx, projectID)
	}

	if cerr != nil {
		fmt.Println("Error on initializing pubsub client:", err.Error())
		return
	}

	client := gcp.PubsubClient{
		Client: internalClient,
		Ctx:    ctx,
	}

	fmt.Printf("Target project ID: %s\n\n", projectID)

	appOpts := appOptions{
		IsDryRun:  opts.DryRun,
		Variables: configuration.Variables(projectID),
	}

	if opts.Delete {
		deleteTopics(client, deleteTopicsLogger{}, ctx, configuration, appOpts)
		deleteSubscriptions(client, deleteSubscriptionLogger{}, ctx, configuration, appOpts)
	} else {
		createTopics(client, createTopicsLogger{}, ctx, configuration, appOpts)
		createSubscriptions(client, createSubscriptionsLogger{}, ctx, configuration, appOpts)
	}

	return
}
