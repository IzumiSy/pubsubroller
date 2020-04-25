package main

import (
	"context"
	"testing"
)

func TestDelete(t *testing.T) {
	ctx := context.Background()
	opts := Options{}

	t.Run("Subscription", func(t *testing.T) {
		t.Parallel()

		cb := fakeSubscriptionCallbacks{}
		deleteSubscriptions(fakeClient{}, &cb, ctx, mockConfig, opts)

		if !cb.IsInitialized {
			t.Error("It must be initialized")
		}

		if !cb.IsFinazlied {
			t.Error("It must be finalized")
		}
	})

	t.Run("Topic", func(t *testing.T) {
		t.Parallel()

		cb := fakeTopicCallbacks{}
		deleteTopics(fakeClient{}, &cb, ctx, mockConfig, opts)

		if !cb.IsInitialized {
			t.Error("It must be initialized")
		}

		if !cb.IsFinazlied {
			t.Error("It must be finalized")
		}
	})
}
