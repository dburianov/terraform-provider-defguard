package provider

import (
	"context"
	"testing"
)

func TestNewProvider(t *testing.T) {
	// Basic test to verify the provider is initialized
	version := "test"
	providerFunc := New(version)
	p := providerFunc()

	if p == nil {
		t.Error("expected provider to be not nil")
	}
}

func TestProviderResources(t *testing.T) {
	version := "test"
	providerFunc := New(version)
	p := providerFunc()

	// Test that all resources are registered
	resources := p.Resources(context.Background())

	if resources == nil {
		t.Fatal("expected Resources to be not nil")
	}

	expectedCount := 5
	if len(resources) != expectedCount {
		t.Errorf("expected %d resources, got %d", expectedCount, len(resources))
	}
}

func TestProviderDataSources(t *testing.T) {
	version := "test"
	providerFunc := New(version)
	p := providerFunc()

	// Test that no data sources are registered (for now)
	dataSources := p.DataSources(context.Background())

	if dataSources == nil {
		t.Error("expected DataSources to be not nil")
	}
}
