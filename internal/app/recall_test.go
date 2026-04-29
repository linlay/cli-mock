package app

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRecallDefaultTextOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		wantQuery string
		wantText  string
	}{
		{
			name:      "knowledge",
			args:      []string{"recall", "knowledge"},
			wantQuery: "query: 办公用品申请流程",
			wantText:  "Office supply requests must include the item list",
		},
		{
			name:      "web search",
			args:      []string{"recall", "web-search"},
			wantQuery: "query: release checklist best practices",
			wantText:  "A release checklist usually confirms scope",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := runCommand(t, nil, tt.args...)

			if result.code != ExitSuccess {
				t.Fatalf("expected exit %d, got %d stderr=%q", ExitSuccess, result.code, result.stderr)
			}
			if result.stderr != "" {
				t.Fatalf("expected empty stderr, got %q", result.stderr)
			}
			if json.Valid([]byte(result.stdout)) {
				t.Fatalf("expected default output to be structured text, got JSON %q", result.stdout)
			}
			for _, want := range []string{
				tt.wantQuery,
				"source_count: 1",
				"chunk_count: 2",
				tt.wantText,
			} {
				if !strings.Contains(result.stdout, want) {
					t.Fatalf("expected stdout to contain %q, got %q", want, result.stdout)
				}
			}
		})
	}
}

func TestRecallJSONOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		args      []string
		wantQuery string
		wantText  string
		wantTitle string
	}{
		{
			name:      "knowledge",
			args:      []string{"recall", "knowledge", "--output", "json"},
			wantQuery: "办公用品申请流程",
			wantText:  "Office supply requests must include the item list, estimated cost, requester department, and business reason before manager approval.",
			wantTitle: "Office Supplies Policy",
		},
		{
			name:      "web search",
			args:      []string{"recall", "web-search", "--output", "json"},
			wantQuery: "release checklist best practices",
			wantText:  "A release checklist usually confirms scope, owners, rollback steps, monitoring dashboards, and customer communication before deployment.",
			wantTitle: "Release Checklist Example",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := runCommand(t, nil, tt.args...)

			if result.code != ExitSuccess {
				t.Fatalf("expected exit %d, got %d stderr=%q", ExitSuccess, result.code, result.stderr)
			}
			var object map[string]json.RawMessage
			if err := json.Unmarshal([]byte(result.stdout), &object); err != nil {
				t.Fatalf("unmarshal recall object: %v", err)
			}
			if len(object) != 4 {
				t.Fatalf("expected only 4 top-level fields, got %#v", object)
			}
			for _, key := range []string{"query", "sourceCount", "chunkCount", "sources"} {
				if _, ok := object[key]; !ok {
					t.Fatalf("expected top-level key %q in %#v", key, object)
				}
			}
			var response recallResponse
			if err := json.Unmarshal([]byte(result.stdout), &response); err != nil {
				t.Fatalf("unmarshal recall response: %v", err)
			}
			if response.Query != tt.wantQuery {
				t.Fatalf("unexpected recall response: %#v", response)
			}
			if response.SourceCount != len(response.Sources) {
				t.Fatalf("expected sourceCount to match sources, got %#v", response)
			}
			if response.ChunkCount != 2 || len(response.Sources) != 1 || len(response.Sources[0].Chunks) != 2 {
				t.Fatalf("unexpected source/chunk shape: %#v", response)
			}
			if response.Sources[0].Title != tt.wantTitle || response.Sources[0].Chunks[0].Content != tt.wantText {
				t.Fatalf("unexpected source content: %#v", response.Sources[0])
			}
		})
	}
}

func TestRecallJSONSupportsQueryFlag(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil,
		"recall", "knowledge",
		"--query", "custom query",
		"--output", "json",
	)

	if result.code != ExitSuccess {
		t.Fatalf("expected exit %d, got %d stderr=%q", ExitSuccess, result.code, result.stderr)
	}
	var response recallResponse
	if err := json.Unmarshal([]byte(result.stdout), &response); err != nil {
		t.Fatalf("unmarshal recall response: %v", err)
	}
	if response.Query != "custom query" {
		t.Fatalf("query flag was not reflected in response: %#v", response)
	}
}

func TestRecallRejectsInvalidOutput(t *testing.T) {
	t.Parallel()

	result := runCommand(t, nil, "recall", "web-search", "--output", "yaml")

	if result.code != ExitUsage {
		t.Fatalf("expected exit %d, got %d", ExitUsage, result.code)
	}
	if !strings.Contains(result.stderr, `invalid --output "yaml": must be one of text, json`) {
		t.Fatalf("unexpected stderr: %q", result.stderr)
	}
}
