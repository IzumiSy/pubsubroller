package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
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
	flag.Parse()

	projectId := *projectIdPtr
	configFilePath := *configFilePathPtr

	if projectId == "" {
		fmt.Println("Error: GCP project ID required with `-projectId` option.")
		return
	}

	if configFilePath == "" {
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

	// pubsubクライアントの生成

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		fmt.Println("Error on initializing pubsub client:", err.Error())
		return
	}

	// topicを並列処理で作成

	egTopics := errgroup.Group{}

	fmt.Printf("\nStart creating topics...\n\n")
	for topicName, _ := range configration.Topics {
		topicName := topicName

		egTopics.Go(func() error {
			return createTopic(client, ctx, topicName)
		})
	}

	if err := egTopics.Wait(); err != nil {
		fmt.Println("Error on creating topics:", err.Error())
	}

	// subscriptionを並列処理で作成

	egSubscriptions := errgroup.Group{}

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
				return CreateSubscription(client, ctx, subscription, topicName)
			})
		}
	}

	if err := egSubscriptions.Wait(); err != nil {
		fmt.Println("Error on creating subscriptions:", err.Error())
	}

	return
}

func createTopic(client *pubsub.Client, ctx context.Context, topicId string) error {
	exists, err := client.Topic(topicId).Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		fmt.Println("Skip:", topicId)
		return nil
	}

	_, err = client.CreateTopic(ctx, topicId)
	if err != nil {
		return err
	}

	fmt.Println("Created:", topicId)
	return nil
}

func CreateSubscription(client *pubsub.Client, ctx context.Context, subscription Subscription, topicName string) error {
	name := subscription.Name
	endpoint := subscription.Endpoint

	s := client.Subscription(name)
	exists, err := s.Exists(ctx)
	if err != nil {
		return err
	}

	if exists {
		fmt.Println("Skip:", name)
		return nil
	}

	pushConfig := pubsub.PushConfig{
		Endpoint: endpoint,
	}

	// 空のpubsub.PushConfigを指定してpullなsubscriptionにする
	if subscription.Pull {
		pushConfig = pubsub.PushConfig{}
	} else {
		if subscription.Endpoint == "" {
			return fmt.Errorf("Failed because no endpoint specified to subscription: %s", name)
		}
	}

	topic := client.Topic(topicName)
	_, err = client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:      topic,
		PushConfig: pushConfig,
	})
	if err != nil {
		return err
	}

	return nil
}
