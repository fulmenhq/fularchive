package cmd

import (
	"strings"
	"testing"

	"github.com/fulmenhq/refbolt/internal/config"
	"github.com/fulmenhq/refbolt/internal/provider"
)

func boolPtr(b bool) *bool { return &b }

var testTopics = []config.Topic{
	{
		Slug: "llm-api",
		Providers: []provider.ProviderConfig{
			{Slug: "openai", Name: "OpenAI"},
			{Slug: "anthropic", Name: "Anthropic"},
			{Slug: "disabled-llm", Name: "Disabled", Enabled: boolPtr(false)},
		},
	},
	{
		Slug: "data-platform",
		Providers: []provider.ProviderConfig{
			{Slug: "trino", Name: "Trino"},
			{Slug: "aws-glue-dg", Name: "AWS Glue"},
		},
	},
}

func slugs(selected []selectedProvider) []string {
	var s []string
	for _, sp := range selected {
		s = append(s, sp.cfg.Slug)
	}
	return s
}

func TestResolveProviders_NoSelector(t *testing.T) {
	_, err := resolveProviders(testTopics, false, nil, nil, nil)
	if err == nil {
		t.Fatal("expected error with no selectors")
	}
	if !strings.Contains(err.Error(), "no providers selected") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveProviders_All(t *testing.T) {
	result, err := resolveProviders(testTopics, true, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := slugs(result)
	// disabled-llm should be skipped
	expected := []string{"openai", "anthropic", "trino", "aws-glue-dg"}
	if len(got) != len(expected) {
		t.Fatalf("got %v, want %v", got, expected)
	}
	for i, s := range expected {
		if got[i] != s {
			t.Errorf("got[%d] = %q, want %q", i, got[i], s)
		}
	}
}

func TestResolveProviders_AllSkipsDisabled(t *testing.T) {
	result, err := resolveProviders(testTopics, true, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, sp := range result {
		if sp.cfg.Slug == "disabled-llm" {
			t.Error("disabled-llm should not be selected by --all")
		}
	}
}

func TestResolveProviders_ProviderFlag(t *testing.T) {
	result, err := resolveProviders(testTopics, false, []string{"openai"}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := slugs(result)
	if len(got) != 1 || got[0] != "openai" {
		t.Errorf("got %v, want [openai]", got)
	}
}

func TestResolveProviders_ProviderOverridesDisabled(t *testing.T) {
	result, err := resolveProviders(testTopics, false, []string{"disabled-llm"}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := slugs(result)
	if len(got) != 1 || got[0] != "disabled-llm" {
		t.Errorf("got %v, want [disabled-llm]", got)
	}
}

func TestResolveProviders_TopicFlag(t *testing.T) {
	result, err := resolveProviders(testTopics, false, nil, []string{"llm-api"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := slugs(result)
	// disabled-llm should be skipped even with --topic
	expected := []string{"openai", "anthropic"}
	if len(got) != len(expected) {
		t.Fatalf("got %v, want %v", got, expected)
	}
}

func TestResolveProviders_TopicSkipsDisabled(t *testing.T) {
	result, err := resolveProviders(testTopics, false, nil, []string{"llm-api"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, sp := range result {
		if sp.cfg.Slug == "disabled-llm" {
			t.Error("disabled-llm should not be selected by --topic")
		}
	}
}

func TestResolveProviders_Union(t *testing.T) {
	result, err := resolveProviders(testTopics, false, []string{"trino"}, []string{"llm-api"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := slugs(result)
	// Union: llm-api providers (enabled) + explicit trino
	expected := []string{"openai", "anthropic", "trino"}
	if len(got) != len(expected) {
		t.Fatalf("got %v, want %v", got, expected)
	}
}

func TestResolveProviders_ExcludeProvider(t *testing.T) {
	result, err := resolveProviders(testTopics, true, nil, nil, []string{"trino"})
	if err != nil {
		t.Fatal(err)
	}
	for _, sp := range result {
		if sp.cfg.Slug == "trino" {
			t.Error("trino should be excluded")
		}
	}
}

func TestResolveProviders_UnknownProvider(t *testing.T) {
	_, err := resolveProviders(testTopics, false, []string{"nonexistent"}, nil, nil)
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
	if !strings.Contains(err.Error(), "unknown provider slug") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveProviders_UnknownTopic(t *testing.T) {
	_, err := resolveProviders(testTopics, false, nil, []string{"nonexistent"}, nil)
	if err == nil {
		t.Fatal("expected error for unknown topic")
	}
	if !strings.Contains(err.Error(), "unknown topic slug") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveProviders_UnknownExclude(t *testing.T) {
	_, err := resolveProviders(testTopics, true, nil, nil, []string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown exclude provider")
	}
	if !strings.Contains(err.Error(), "unknown provider slug in --exclude-provider") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveProviders_ConflictProviderAndExclude(t *testing.T) {
	_, err := resolveProviders(testTopics, false, []string{"openai"}, nil, []string{"openai"})
	if err == nil {
		t.Fatal("expected error for conflicting provider and exclude")
	}
	if !strings.Contains(err.Error(), "both --provider and --exclude-provider") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestResolveProviders_AllExcluded(t *testing.T) {
	_, err := resolveProviders(testTopics, true, nil, nil, []string{"openai", "anthropic", "trino", "aws-glue-dg"})
	if err == nil {
		t.Fatal("expected error when all providers excluded")
	}
	if !strings.Contains(err.Error(), "no providers matched") {
		t.Errorf("unexpected error: %v", err)
	}
}
