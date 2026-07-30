package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func doctorEnvPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "leads-brasil-pp-cli", "doctor.env"), nil
}

func loadDoctorEnvFile() error {
	path, err := doctorEnvPath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		value = strings.Trim(value, "\"")
		value = strings.Trim(value, "'")
		if key == "" {
			continue
		}
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
	return nil
}

func persistDoctorEnvVar(key, value string) error {
	if key == "" {
		return fmt.Errorf("doctor env key is required")
	}
	path, err := doctorEnvPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	lines := []string{}
	if data, err := os.ReadFile(path); err == nil {
		s := bufio.NewScanner(strings.NewReader(string(data)))
		for s.Scan() {
			lines = append(lines, s.Text())
		}
		if err := s.Err(); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	prefixes := []string{key + "=", "export " + key + "="}
	replaced := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		for _, prefix := range prefixes {
			if strings.HasPrefix(trimmed, prefix) {
				lines[i] = fmt.Sprintf("export %s=%q", key, value)
				replaced = true
				break
			}
		}
	}
	if !replaced {
		lines = append(lines, fmt.Sprintf("export %s=%q", key, value))
	}
	content := strings.Join(lines, "\n")
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o600)
}
