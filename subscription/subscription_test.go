package subscription

import (
	"pubsubroller/config"
	"testing"
)

func TestSubscription(t *testing.T) {
	t.Run("FromConfig", func(t *testing.T) {
		mockConfig := config.Configuration{
			Internal_Topics_: map[string]config.Topic{
				"topic1": {
					Internal_Subscriptions_: []config.Subscription{
						{
							Name:     "subscription${suffix1}",
							Endpoint: "https://example.com/path1",
						},
						{
							Name:     "subscription2",
							Endpoint: "https://example.com/${path2}",
						},
						{
							Name:     "subscription${suffix3}",
							Endpoint: "https://example.com/${path3}",
						},
					},
				},
			},
		}

		vars := map[string]string{
			"path1":   "replaced1",
			"path2":   "replaced2",
			"path3":   "replaced3",
			"suffix1": "1",
			"suffix3": "3",
		}

		subs := FromConfig(mockConfig, vars)

		for _, sub := range subs {
			switch sub.Name {
			case "subscription1":
				if sub.Endpoint != "https://example.com/path1" {
					t.Error("Expected to be replaced")
				}
			case "subscription2":
				if sub.Endpoint != "https://example.com/replaced2" {
					t.Error("Expected to be replaced")
				}
			case "subscription3":
				if sub.Endpoint != "https://example.com/replaced3" {
					t.Error("Expected to be replaced")
				}
			default:
				t.Error("Replacing malfunctioning")
			}
		}
	})
}
