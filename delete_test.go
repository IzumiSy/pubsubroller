package main

import (
	"context"
	"testing"
)

func TestDeleteSubscription(t *testing.T) {
	ctx := context.Background()
	opts := Options{}
	cb := fakeCallbacks{}

	deleteSubscriptions(fakeClient{}, &cb, ctx, mockConfig, opts)

	if !cb.IsInitialized {
		t.Error("It must be initialized")
	}

	if !cb.IsFinazlied {
		t.Error("It must be finalized")
	}
}
