package usecase

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const EnvConfigPath = "PP_LEADS_USE_CASE_CONFIG"
const EnvConfigDir = "PP_LEADS_ICP_DIR"

type Config struct {
	Name          string            `json:"name"`
	Label         string            `json:"label"`
	LeadTablePath string            `json:"lead_table_path"`
	OutputDir     string            `json:"output_dir"`
	MessageField  string            `json:"message_field"`
	FieldMap      map[string]string `json:"field_map"`
	Path          string            `json:"-"`
}

func LoadFromEnv() (*Config, error) {
	path := strings.TrimSpace(os.Getenv(EnvConfigPath))
	if path != "" {
		return Load(path)
	}
	dir := strings.TrimSpace(os.Getenv(EnvConfigDir))
	if dir != "" {
		return Load(dir)
	}
	return nil, nil
}

func Load(path string) (*Config, error) {
	resolved, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return nil, fmt.Errorf("reading ICP source %s: %w", resolved, err)
	}
	if info.IsDir() {
		return LoadDir(resolved)
	}
	data, err := os.ReadFile(resolved)
	if err != nil {
		return nil, fmt.Errorf("reading ICP config %s: %w", resolved, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing ICP config %s: %w", resolved, err)
	}
	cfg.Path = resolved
	if cfg.FieldMap == nil {
		cfg.FieldMap = map[string]string{}
	}
	return &cfg, nil
}

func LoadDir(dir string) (*Config, error) {
	resolved, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return nil, fmt.Errorf("reading ICP directory %s: %w", resolved, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("ICP directory %s is not a directory", resolved)
	}
	leadTablePath, err := detectLeadTablePath(resolved)
	if err != nil {
		return nil, err
	}
	cfg := &Config{
		Name:          strings.TrimSpace(filepath.Base(resolved)),
		Label:         strings.ReplaceAll(strings.TrimSpace(filepath.Base(resolved)), "-", " "),
		LeadTablePath: leadTablePath,
		OutputDir:     filepath.Join("outputs", "enrichment"),
		FieldMap:      map[string]string{},
		Path:          resolved,
	}
	return cfg, nil
}

func detectLeadTablePath(dir string) (string, error) {
	patterns := []string{"lead-table-*.csv", "lead-table-*.xlsx", "lead-table-*.xlsm"}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return "", err
		}
		if len(matches) > 0 {
			sort.Strings(matches)
			return matches[0], nil
		}
	}
	return "", fmt.Errorf("no lead-table file detected under ICP directory %s", dir)
}

func (c *Config) baseDir() string {
	if c == nil || strings.TrimSpace(c.Path) == "" {
		return ""
	}
	if info, err := os.Stat(c.Path); err == nil && info.IsDir() {
		return c.Path
	}
	return filepath.Dir(c.Path)
}

func (c *Config) ResolvePath(path string) string {
	if c == nil || strings.TrimSpace(path) == "" {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	base := c.baseDir()
	if base == "" {
		return path
	}
	return filepath.Join(base, path)
}

func (c *Config) Metadata() map[string]any {
	if c == nil {
		return nil
	}
	meta := map[string]any{}
	if value := strings.TrimSpace(c.Name); value != "" {
		meta["name"] = value
	}
	if value := strings.TrimSpace(c.Label); value != "" {
		meta["label"] = value
	}
	if len(meta) == 0 {
		return nil
	}
	return meta
}

func (c *Config) field(alias string, fallbacks ...string) string {
	if c != nil {
		if value := strings.TrimSpace(c.FieldMap[alias]); value != "" {
			return value
		}
	}
	for _, candidate := range fallbacks {
		if strings.TrimSpace(candidate) != "" {
			return candidate
		}
	}
	return ""
}

func Name(record map[string]string, cfg *Config) string {
	return firstValue(record, cfg.field("name", "lead", "name", "empresa", "company"))
}

func CNPJ(record map[string]string, cfg *Config) string {
	return firstValue(record, cfg.field("cnpj", "CNPJ", "cnpj"))
}

func CityUF(record map[string]string, cfg *Config) string {
	return firstValue(record, cfg.field("city_uf", "cidade/UF", "cidade_uf", "city"))
}

func Contact(record map[string]string, cfg *Config) string {
	return firstValue(record, cfg.field("contact", "contato", "email", "contact"))
}

func Site(record map[string]string, cfg *Config) string {
	return firstValue(record, cfg.field("site", "site", "website", "domain"))
}

func Message(record map[string]string, cfg *Config) string {
	if cfg != nil && strings.TrimSpace(cfg.MessageField) != "" {
		if value := strings.TrimSpace(record[cfg.MessageField]); value != "" {
			return value
		}
	}
	for _, key := range []string{"primeira mensagem", "first_message", "message", "mensagem", "outreach_message"} {
		if value := strings.TrimSpace(record[key]); value != "" {
			return value
		}
	}
	return ""
}

func SearchText(record map[string]string) string {
	if len(record) == 0 {
		return ""
	}
	keys := make([]string, 0, len(record))
	for key := range record {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		value := strings.TrimSpace(record[key])
		if value != "" {
			parts = append(parts, value)
		}
	}
	return strings.Join(parts, " ")
}

func LeadContext(record map[string]string) map[string]any {
	context := make(map[string]any, len(record))
	for key, value := range record {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		normalized := normalizeKey(key)
		if strings.HasPrefix(normalized, "pp_") {
			continue
		}
		context[normalized] = value
	}
	return context
}

func firstValue(record map[string]string, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(record[key]); value != "" {
			return value
		}
	}
	return ""
}

func normalizeKey(value string) string {
	var b strings.Builder
	b.Grow(len(value))
	lastUnderscore := false
	for _, r := range strings.TrimSpace(value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
			lastUnderscore = false
			continue
		}
		if !lastUnderscore {
			b.WriteByte('_')
			lastUnderscore = true
		}
	}
	return strings.Trim(b.String(), "_")
}
