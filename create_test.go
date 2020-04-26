package main

import (
	"context"
	"testing"
)

func TestCreate(t *testing.T) {
	ctx := context.Background()
	opts := appOptions{}

	t.Run("Subscription", func(t *testing.T) {
		t.Parallel()

		cb := fakeSubscriptionCallbacks{}
		createSubscriptions(fakeClient{}, &cb, ctx, mockConfig, opts)

		if !cb.IsInitialized {
			t.Error("It must be initialized")
		}

		if !cb.IsFinazlied {
			t.Error("It must be finalized")
		}

		if cb.Calls != 9 {
			t.Errorf("Calls must be 9, but it is %d", cb.Calls)
		}
	})

	t.Run("Topic", func(t *testing.T) {
		t.Parallel()

		cb := fakeTopicCallbacks{}
		createTopics(fakeClient{}, &cb, ctx, mockConfig, opts)

		if !cb.IsInitialized {
			t.Error("It must be initialized")
		}

		if !cb.IsFinazlied {
			t.Error("It must be finalized")
		}

		if cb.Calls != 3 {
			t.Errorf("Calls must be 3, but it is %d", cb.Calls)
		}
	})
}
