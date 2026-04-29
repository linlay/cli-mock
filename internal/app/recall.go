package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"
)

const (
	defaultRecallOutput = defaultBusinessOutput
)

type recallOptions struct {
	Query  string
	Output string
}

type recallResponse struct {
	Query       string         `json:"query,omitempty"`
	SourceCount int            `json:"sourceCount"`
	ChunkCount  int            `json:"chunkCount"`
	Sources     []recallSource `json:"sources"`
}

type recallSource struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Title          string        `json:"title,omitempty"`
	Icon           string        `json:"icon,omitempty"`
	URL            string        `json:"url,omitempty"`
	Link           string        `json:"link,omitempty"`
	CollectionID   string        `json:"collectionId,omitempty"`
	CollectionName string        `json:"collectionName,omitempty"`
	ChunkIndexes   []int         `json:"chunkIndexes"`
	MinIndex       int           `json:"minIndex"`
	Chunks         []recallChunk `json:"chunks"`
}

type recallChunk struct {
	ChunkID string  `json:"chunkId"`
	Index   int     `json:"index"`
	Content string  `json:"content"`
	Score   float64 `json:"score,omitempty"`
}

func newRecallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "recall",
		Short:       "Mock source recall results",
		Description: "Group mock source recall commands that emit recall result bodies.",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cmd.AddCommand(
		newRecallKnowledgeCommand(),
		newRecallWebSearchCommand(),
	)
	return cmd
}

func newRecallKnowledgeCommand() *cobra.Command {
	defaultQuery := "办公用品申请流程"
	opts := recallOptions{
		Query:  defaultQuery,
		Output: defaultRecallOutput,
	}

	cmd := &cobra.Command{
		Use:         "knowledge",
		Short:       "Publish mock knowledge-base recall sources",
		Description: "Return a fixed knowledge-base recall result body.",
		Example: "mock recall knowledge\n" +
			"mock recall knowledge --query 办公用品申请流程 --output json",
		Args:        cobra.NoArgs,
		ParamFields: recallParamFields(defaultQuery),
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, err := parseOutputFormat(opts.Output)
			if err != nil {
				return err
			}
			return writeRecallResponse(cmd.OutOrStdout(), outputFormat, buildKnowledgeRecall(opts))
		},
	}
	addRecallFlags(cmd, &opts)
	return cmd
}

func newRecallWebSearchCommand() *cobra.Command {
	defaultQuery := "release checklist best practices"
	opts := recallOptions{
		Query:  defaultQuery,
		Output: defaultRecallOutput,
	}

	cmd := &cobra.Command{
		Use:         "web-search",
		Short:       "Publish mock web search recall sources",
		Description: "Return a fixed web search recall result body.",
		Example: "mock recall web-search\n" +
			"mock recall web-search --query \"release checklist best practices\" --output json",
		Args:        cobra.NoArgs,
		ParamFields: recallParamFields(defaultQuery),
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, err := parseOutputFormat(opts.Output)
			if err != nil {
				return err
			}
			return writeRecallResponse(cmd.OutOrStdout(), outputFormat, buildWebSearchRecall(opts))
		},
	}
	addRecallFlags(cmd, &opts)
	return cmd
}

func addRecallFlags(cmd *cobra.Command, opts *recallOptions) {
	cmd.Flags().StringVar(&opts.Query, "query", opts.Query, "Recall query text")
	cmd.Flags().StringVar(&opts.Output, "output", opts.Output, outputUsage())
}

func recallParamFields(defaultQuery string) []cobra.HelpField {
	return []cobra.HelpField{
		optionalField("query", "string", defaultQuery, "Recall query text"),
		outputField(defaultRecallOutput),
	}
}

func buildKnowledgeRecall(opts recallOptions) recallResponse {
	sources := []recallSource{
		{
			ID:             "doc_policy_001",
			Name:           "office-supplies-policy.txt",
			Title:          "Office Supplies Policy",
			Icon:           "ragflow",
			URL:            "https://example.com/docs/office-supplies-policy.txt",
			CollectionID:   "col_policy",
			CollectionName: "Company Policies",
			ChunkIndexes:   []int{1, 2},
			MinIndex:       1,
			Chunks: []recallChunk{
				{
					ChunkID: "chunk_policy_001",
					Index:   1,
					Content: "Office supply requests must include the item list, estimated cost, requester department, and business reason before manager approval.",
					Score:   0.92,
				},
				{
					ChunkID: "chunk_policy_002",
					Index:   2,
					Content: "Requests above the monthly department allowance require finance review after the direct manager approves the purchase.",
					Score:   0.86,
				},
			},
		},
	}
	return newRecallResponse(opts, sources)
}

func buildWebSearchRecall(opts recallOptions) recallResponse {
	sources := []recallSource{
		{
			ID:           "web_result_001",
			Name:         "example-release-checklist",
			Title:        "Release Checklist Example",
			Icon:         "websearch",
			URL:          "https://example.com/release-checklist",
			Link:         "https://example.com/release-checklist",
			ChunkIndexes: []int{1, 2},
			MinIndex:     1,
			Chunks: []recallChunk{
				{
					ChunkID: "chunk_web_001",
					Index:   1,
					Content: "A release checklist usually confirms scope, owners, rollback steps, monitoring dashboards, and customer communication before deployment.",
					Score:   0.9,
				},
				{
					ChunkID: "chunk_web_002",
					Index:   2,
					Content: "Post-release validation should compare expected metrics with live telemetry and record any follow-up work in the release notes.",
					Score:   0.84,
				},
			},
		},
	}
	return newRecallResponse(opts, sources)
}

func newRecallResponse(opts recallOptions, sources []recallSource) recallResponse {
	chunkCount := 0
	for _, source := range sources {
		chunkCount += len(source.Chunks)
	}
	return recallResponse{
		Query:       strings.TrimSpace(opts.Query),
		SourceCount: len(sources),
		ChunkCount:  chunkCount,
		Sources:     sources,
	}
}

func writeRecallResponse(out io.Writer, outputFormat string, response recallResponse) error {
	if outputFormat == "json" {
		return writeJSON(out, response)
	}

	var b strings.Builder
	if strings.TrimSpace(response.Query) != "" {
		writeStructuredField(&b, "query", response.Query, 0)
	}
	writeStructuredField(&b, "source_count", response.SourceCount, 0)
	writeStructuredField(&b, "chunk_count", response.ChunkCount, 0)
	writeStructuredField(&b, "sources", response.Sources, 0)

	_, err := fmt.Fprint(out, b.String())
	return err
}
