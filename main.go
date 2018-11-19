package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
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

func main() {
	projectIdPtr := flag.String("projectId", "", "target GCP project ID")
	configFilePathPtr := flag.String("config", "", "configuration file path")
	endpointPtr := flag.String("endpoint", "", "service endpoint")
	isDryRunPtr := flag.Bool("dry", false, "dry run")
	flag.Parse()

	projectId := *projectIdPtr
	configFilePath := *configFilePathPtr
	endpoint := *endpointPtr
	isDryRun := *isDryRunPtr

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

	configration := Configuration{}
	err = yaml.Unmarshal(yamlBytes, &configration)
	if err != nil {
		panic(err)
	}

	// variablesの取得

	variables := make(map[string]string)
	for key, value := range configration.Variables {
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

	// topicを並列処理で作成

	egTopics := errgroup.Group{}
	topicSkippedCount := 0
	topicCreatedCount := 0

	fmt.Printf("\nStart creating topics...\n\n")
	for topicName, _ := range configration.Topics {
		topicName := topicName

		egTopics.Go(func() error {
			isCreated, err := createTopic(client, ctx, topicName, isDryRun)
			if isCreated {
				topicCreatedCount += 1
			} else {
				topicSkippedCount += 1
			}
			return err
		})
	}

	if err := egTopics.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("\nTopics Created: %d, Skipped: %d\n", topicCreatedCount, topicSkippedCount)

	// subscriptionを並列処理で作成

	egSubscriptions := errgroup.Group{}
	subscriptionSkippedCount := 0
	subscriptionCreatedCount := 0

	fmt.Printf("\nStart creating subscriptions...\n\n")
	for topicName, topic := range configration.Topics {
		topicName := topicName
		topic := topic

		for _, subscription := range topic.Subsciptions {
			subscription := subscription
			endpoint := subscription.Endpoint
			for key, value := range variables {
				endpoint = strings.Replace(endpoint, "${"+key+"}", value, -1)
			}

			subscription.Endpoint = endpoint
			egSubscriptions.Go(func() error {
				isCreated, err := createSubscription(client, ctx, subscription, topicName, isDryRun)
				if isCreated {
					subscriptionCreatedCount += 1
				} else {
					subscriptionSkippedCount += 1
				}
				return err
			})
		}
	}

	if err := egSubscriptions.Wait(); err != nil {
		panic(err)
	}

	fmt.Printf("\nSubscriptions Created: %d, Skipped: %d\n", subscriptionCreatedCount, subscriptionSkippedCount)

	return
}

func createTopic(client *pubsub.Client, ctx context.Context, topicId string, isDryRun bool) (bool, error) {
	exists, err := client.Topic(topicId).Exists(ctx)
	if err != nil {
		return false, err
	}

	if exists {
		fmt.Println("Skip:", topicId)
		return false, nil
	}

	// dryrunであれば必ず成功扱いとするため作成系の実行はスキップする
	if !isDryRun {
		_, err = client.CreateTopic(ctx, topicId)
		if err != nil {
			return false, err
		}
	}

	fmt.Println("Created:", topicId)
	return true, nil
}

func createSubscription(client *pubsub.Client, ctx context.Context, subscription Subscription, topicName string, isDryRun bool) (bool, error) {
	name := subscription.Name

	endpoint := subscription.Endpoint

	s := client.Subscription(name)
	exists, err := s.Exists(ctx)
	if err != nil {
		return false, err
	}

	if exists {
		fmt.Println("Skip:", name)
		return false, nil
	}

	// dryrunであれば必ず成功扱いとするため作成系の実行はスキップする
	if !isDryRun {
		pushConfig := pubsub.PushConfig{
			Endpoint: endpoint,
		}

		// 空のpubsub.PushConfigを指定してpullなsubscriptionにする
		if subscription.Pull {
			pushConfig = pubsub.PushConfig{}
		} else {
			if subscription.Endpoint == "" {
				return false, fmt.Errorf("Failed because no endpoint specified to subscription: %s", name)
			}
		}

		topic := client.Topic(topicName)
		_, err = client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
			Topic:      topic,
			PushConfig: pushConfig,
		})
		if err != nil {
			return false, err
		}
	}

	fmt.Println("Created:", name)
	return true, nil
}
