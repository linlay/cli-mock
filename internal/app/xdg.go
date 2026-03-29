package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
)

type xdgManifest struct {
	Entries []xdgEntrySpec `json:"entries"`
}

type xdgEntrySpec struct {
	Path    string `json:"path"`
	Type    string `json:"type"`
	Format  string `json:"format"`
	Content any    `json:"content"`
	Mode    string `json:"mode"`
}

type xdgPreparedEntry struct {
	RelativePath string
	AbsolutePath string
	Type         string
	Content      []byte
	Mode         fs.FileMode
}

type xdgApplySummary struct {
	Root    string
	Created []string
	Updated []string
}

type xdgInspectResponse struct {
	Root      string            `json:"root"`
	Managed   map[string]string `json:"managed"`
	Reveal    bool              `json:"reveal"`
	Entries   []xdgInspectEntry `json:"entries"`
}

type xdgInspectEntry struct {
	Path    string `json:"path"`
	Type    string `json:"type"`
	Mode    string `json:"mode"`
	Bytes   int64  `json:"bytes"`
	Format  string `json:"format,omitempty"`
	Content any    `json:"content,omitempty"`
}

func newXDGCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:         "xdg",
		Short:       "Create and inspect mock XDG-style config trees",
		Description: "Create and inspect mock .config and .local trees under an explicit root.",
		Example: "mock xdg apply --root /tmp/mock-home --manifest ./manifest.json\n" +
			"mock xdg inspect --root /tmp/mock-home\n" +
			"mock xdg inspect --root /tmp/mock-home --reveal",
		Args: cobra.NoArgs,
	}
	cmd.CompletionOptions.DisableDefaultCmd = true

	cmd.AddCommand(
		newXDGApplyCommand(),
		newXDGInspectCommand(),
	)

	return cmd
}

func newXDGApplyCommand() *cobra.Command {
	var root string
	var manifestPath string
	var overwrite bool

	cmd := &cobra.Command{
		Use:         "apply",
		Short:       "Create a mock XDG tree from a JSON manifest",
		Description: "Create .config and .local content under the provided root using a JSON manifest.",
		Example: "mock xdg apply --root /tmp/mock-home --manifest ./manifest.json\n" +
			"mock xdg apply --root /tmp/mock-home --manifest - --overwrite",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rootPath, err := validateXDGRoot(root)
			if err != nil {
				return err
			}

			manifest, err := readXDGManifest(cmd.InOrStdin(), manifestPath)
			if err != nil {
				return err
			}

			prepared, err := prepareXDGEntries(rootPath, manifest)
			if err != nil {
				return err
			}

			summary, err := applyXDGTree(rootPath, prepared, overwrite)
			if err != nil {
				return err
			}

			return writeXDGApplySummary(cmd.OutOrStdout(), summary)
		},
	}

	cmd.Flags().StringVar(&root, "root", "", "Fake home root to manage")
	cmd.Flags().StringVar(&manifestPath, "manifest", "", "JSON manifest path or - for stdin")
	cmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing files")
	return cmd
}

func newXDGInspectCommand() *cobra.Command {
	var root string
	var reveal bool

	cmd := &cobra.Command{
		Use:         "inspect",
		Short:       "Inspect a mock XDG tree as JSON",
		Description: "Inspect .config and .local metadata under the provided root and optionally reveal file content.",
		Example:     "mock xdg inspect --root /tmp/mock-home\nmock xdg inspect --root /tmp/mock-home --reveal",
		Args:        cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			rootPath, err := validateXDGRoot(root)
			if err != nil {
				return err
			}

			response, err := inspectXDGTree(rootPath, reveal)
			if err != nil {
				return err
			}
			enc := json.NewEncoder(cmd.OutOrStdout())
			enc.SetEscapeHTML(false)
			return enc.Encode(response)
		},
	}

	cmd.Flags().StringVar(&root, "root", "", "Fake home root to inspect")
	cmd.Flags().BoolVar(&reveal, "reveal", false, "Include readable text and JSON file content")

	return cmd
}

func validateXDGRoot(root string) (string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return "", fmt.Errorf("--root is required")
	}
	if !filepath.IsAbs(root) {
		return "", fmt.Errorf("--root must be an absolute path")
	}
	return filepath.Clean(root), nil
}

func readXDGManifest(stdin io.Reader, path string) (*xdgManifest, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, fmt.Errorf("--manifest is required")
	}

	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(stdin)
		if err != nil {
			return nil, failuref("read manifest stdin: %v", err)
		}
	} else {
		data, err = os.ReadFile(path)
		if err != nil {
			return nil, failuref("read manifest %q: %v", path, err)
		}
	}

	var manifest xdgManifest
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&manifest); err != nil {
		return nil, fmt.Errorf("invalid manifest JSON: %w", err)
	}
	if dec.More() {
		return nil, fmt.Errorf("invalid manifest JSON: multiple JSON values")
	}
	if len(manifest.Entries) == 0 {
		return nil, fmt.Errorf("manifest.entries must contain at least one entry")
	}
	return &manifest, nil
}

func prepareXDGEntries(root string, manifest *xdgManifest) ([]xdgPreparedEntry, error) {
	seen := make(map[string]struct{}, len(manifest.Entries))
	prepared := make([]xdgPreparedEntry, 0, len(manifest.Entries))

	for _, entry := range manifest.Entries {
		next, err := prepareXDGEntry(root, entry)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[next.RelativePath]; ok {
			return nil, fmt.Errorf("manifest contains duplicate path %q", next.RelativePath)
		}
		seen[next.RelativePath] = struct{}{}
		prepared = append(prepared, next)
	}

	sort.Slice(prepared, func(i, j int) bool {
		return prepared[i].RelativePath < prepared[j].RelativePath
	})
	return prepared, nil
}

func prepareXDGEntry(root string, entry xdgEntrySpec) (xdgPreparedEntry, error) {
	rel, err := validateXDGRelativePath(entry.Path)
	if err != nil {
		return xdgPreparedEntry{}, err
	}

	mode, err := parseXDGMode(entry.Mode, rel, entry.Type)
	if err != nil {
		return xdgPreparedEntry{}, err
	}

	switch entry.Type {
	case "dir":
		if entry.Format != "" {
			return xdgPreparedEntry{}, fmt.Errorf("directory entry %q cannot set format", rel)
		}
		if entry.Content != nil {
			return xdgPreparedEntry{}, fmt.Errorf("directory entry %q cannot include content", rel)
		}
		return xdgPreparedEntry{
			RelativePath: rel,
			AbsolutePath: filepath.Join(root, filepath.FromSlash(rel)),
			Type:         entry.Type,
			Mode:         mode,
		}, nil
	case "file":
		if entry.Content == nil {
			return xdgPreparedEntry{}, fmt.Errorf("file entry %q requires content", rel)
		}
		content, err := encodeXDGContent(rel, entry.Format, entry.Content)
		if err != nil {
			return xdgPreparedEntry{}, err
		}
		return xdgPreparedEntry{
			RelativePath: rel,
			AbsolutePath: filepath.Join(root, filepath.FromSlash(rel)),
			Type:         entry.Type,
			Content:      content,
			Mode:         mode,
		}, nil
	default:
		return xdgPreparedEntry{}, fmt.Errorf("entry %q has unsupported type %q", rel, entry.Type)
	}
}

func validateXDGRelativePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("entry path is required")
	}
	if filepath.IsAbs(path) {
		return "", fmt.Errorf("entry path %q must be relative", path)
	}

	cleaned := filepath.ToSlash(filepath.Clean(path))
	if strings.HasPrefix(cleaned, "../") || cleaned == ".." {
		return "", fmt.Errorf("entry path %q must not escape the root", path)
	}
	if strings.Contains(cleaned, "/../") {
		return "", fmt.Errorf("entry path %q must not contain ..", path)
	}
	if cleaned == ".config" || cleaned == ".local" {
		return "", fmt.Errorf("entry path %q must target a child under .config/ or .local/", path)
	}
	if !strings.HasPrefix(cleaned, ".config/") && !strings.HasPrefix(cleaned, ".local/") {
		return "", fmt.Errorf("entry path %q must start with .config/ or .local/", path)
	}
	return cleaned, nil
}

func parseXDGMode(raw, relPath, entryType string) (fs.FileMode, error) {
	if strings.TrimSpace(raw) == "" {
		if entryType == "dir" {
			return 0o755, nil
		}
		if strings.HasPrefix(relPath, ".config/") {
			return 0o644, nil
		}
		return 0o600, nil
	}

	value, err := strconv.ParseUint(raw, 8, 32)
	if err != nil {
		return 0, fmt.Errorf("entry %q has invalid mode %q", relPath, raw)
	}
	return fs.FileMode(value), nil
}

func encodeXDGContent(relPath, format string, content any) ([]byte, error) {
	switch format {
	case "", "text":
		asString, ok := content.(string)
		if !ok {
			return nil, fmt.Errorf("file entry %q with format %q requires string content", relPath, defaultTextFormat(format))
		}
		return []byte(asString), nil
	case "json":
		data, err := json.Marshal(content)
		if err != nil {
			return nil, fmt.Errorf("file entry %q marshal JSON: %w", relPath, err)
		}
		return append(data, '\n'), nil
	default:
		return nil, fmt.Errorf("file entry %q has unsupported format %q", relPath, format)
	}
}

func defaultTextFormat(format string) string {
	if format == "" {
		return "text"
	}
	return format
}

func applyXDGTree(root string, entries []xdgPreparedEntry, overwrite bool) (*xdgApplySummary, error) {
	configDir := filepath.Join(root, ".config")
	localDir := filepath.Join(root, ".local")

	for _, dir := range []string{configDir, localDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, failuref("create managed directory %q: %v", dir, err)
		}
	}

	summary := &xdgApplySummary{Root: root}

	for _, entry := range entries {
		switch entry.Type {
		case "dir":
			existed := pathExists(entry.AbsolutePath)
			if err := os.MkdirAll(entry.AbsolutePath, entry.Mode); err != nil {
				return nil, failuref("create directory %q: %v", entry.RelativePath, err)
			}
			if err := os.Chmod(entry.AbsolutePath, entry.Mode); err != nil {
				return nil, failuref("set mode on %q: %v", entry.RelativePath, err)
			}
			if existed {
				summary.Updated = append(summary.Updated, entry.RelativePath)
			} else {
				summary.Created = append(summary.Created, entry.RelativePath)
			}
		case "file":
			existed := pathExists(entry.AbsolutePath)
			if err := os.MkdirAll(filepath.Dir(entry.AbsolutePath), 0o755); err != nil {
				return nil, failuref("create parent directory for %q: %v", entry.RelativePath, err)
			}
			if err := writeXDGFile(entry, overwrite); err != nil {
				return nil, err
			}

			if existed {
				summary.Updated = append(summary.Updated, entry.RelativePath)
			} else {
				summary.Created = append(summary.Created, entry.RelativePath)
			}
		}
	}

	sort.Strings(summary.Created)
	sort.Strings(summary.Updated)
	return summary, nil
}

func writeXDGFile(entry xdgPreparedEntry, overwrite bool) error {
	if info, err := os.Stat(entry.AbsolutePath); err == nil {
		if info.IsDir() {
			return failuref("cannot replace directory %q with a file", entry.RelativePath)
		}
		if !overwrite {
			return fmt.Errorf("file %q already exists; use --overwrite to replace it", entry.RelativePath)
		}
		if err := os.WriteFile(entry.AbsolutePath, entry.Content, entry.Mode); err != nil {
			return failuref("write file %q: %v", entry.RelativePath, err)
		}
		if err := os.Chmod(entry.AbsolutePath, entry.Mode); err != nil {
			return failuref("set mode on %q: %v", entry.RelativePath, err)
		}
		return nil
	} else if !os.IsNotExist(err) {
		return failuref("stat file %q: %v", entry.RelativePath, err)
	}

	if err := os.WriteFile(entry.AbsolutePath, entry.Content, entry.Mode); err != nil {
		return failuref("write file %q: %v", entry.RelativePath, err)
	}
	if err := os.Chmod(entry.AbsolutePath, entry.Mode); err != nil {
		return failuref("set mode on %q: %v", entry.RelativePath, err)
	}
	return nil
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func writeXDGApplySummary(w io.Writer, summary *xdgApplySummary) error {
	if _, err := fmt.Fprintf(w, "root: %s\n", summary.Root); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "managed_dirs:\n- %s\n- %s\n", filepath.Join(summary.Root, ".config"), filepath.Join(summary.Root, ".local")); err != nil {
		return err
	}
	if _, err := io.WriteString(w, "created:\n"); err != nil {
		return err
	}
	if len(summary.Created) == 0 {
		if _, err := io.WriteString(w, "-\n"); err != nil {
			return err
		}
	} else {
		for _, path := range summary.Created {
			if _, err := fmt.Fprintf(w, "- %s\n", path); err != nil {
				return err
			}
		}
	}

	if _, err := io.WriteString(w, "updated:\n"); err != nil {
		return err
	}
	if len(summary.Updated) == 0 {
		if _, err := io.WriteString(w, "-\n"); err != nil {
			return err
		}
	} else {
		for _, path := range summary.Updated {
			if _, err := fmt.Fprintf(w, "- %s\n", path); err != nil {
				return err
			}
		}
	}

	exports := []string{
		"HOME=" + summary.Root,
		"XDG_CONFIG_HOME=" + filepath.Join(summary.Root, ".config"),
		"XDG_DATA_HOME=" + filepath.Join(summary.Root, ".local", "share"),
		"XDG_STATE_HOME=" + filepath.Join(summary.Root, ".local", "state"),
	}
	if _, err := io.WriteString(w, "exports:\n"); err != nil {
		return err
	}
	for _, line := range exports {
		if _, err := fmt.Fprintf(w, "- %s\n", line); err != nil {
			return err
		}
	}
	return nil
}

func inspectXDGTree(root string, reveal bool) (*xdgInspectResponse, error) {
	response := &xdgInspectResponse{
		Root: root,
		Managed: map[string]string{
			"config": filepath.Join(root, ".config"),
			"local":  filepath.Join(root, ".local"),
		},
		Reveal:  reveal,
		Entries: []xdgInspectEntry{},
	}

	for _, rel := range []string{".config", ".local"} {
		entries, err := collectXDGEntries(root, rel, reveal)
		if err != nil {
			return nil, err
		}
		response.Entries = append(response.Entries, entries...)
	}
	sort.Slice(response.Entries, func(i, j int) bool {
		return response.Entries[i].Path < response.Entries[j].Path
	})
	return response, nil
}

func collectXDGEntries(root, relativeBase string, reveal bool) ([]xdgInspectEntry, error) {
	basePath := filepath.Join(root, filepath.FromSlash(relativeBase))
	info, err := os.Stat(basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []xdgInspectEntry{}, nil
		}
		return nil, failuref("inspect %q: %v", relativeBase, err)
	}
	if !info.IsDir() {
		return nil, failuref("inspect %q: expected a directory", relativeBase)
	}

	entries := []xdgInspectEntry{}
	err = filepath.WalkDir(basePath, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == basePath {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		item := xdgInspectEntry{
			Path:  filepath.ToSlash(rel),
			Type:  "file",
			Mode:  fmt.Sprintf("%04o", info.Mode().Perm()),
			Bytes: info.Size(),
		}
		if d.IsDir() {
			item.Type = "dir"
			item.Bytes = 0
		} else if reveal {
			format, content, ok, err := revealXDGContent(path)
			if err != nil {
				return err
			}
			if ok {
				item.Format = format
				item.Content = content
			}
		}
		entries = append(entries, item)
		return nil
	})
	if err != nil {
		return nil, failuref("walk %q: %v", relativeBase, err)
	}
	return entries, nil
}

func revealXDGContent(path string) (string, any, bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", nil, false, err
	}
	if json.Valid(data) {
		var value any
		if err := json.Unmarshal(data, &value); err != nil {
			return "", nil, false, err
		}
		return "json", value, true, nil
	}
	if isTextContent(data) {
		return "text", string(data), true, nil
	}
	return "", nil, false, nil
}

func isTextContent(data []byte) bool {
	if !utf8.Valid(data) {
		return false
	}
	for _, r := range string(data) {
		if r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		if r < 0x20 {
			return false
		}
	}
	return true
}
