package cli

import (
	"os"
	"path/filepath"
)

func detectDefaultUseCaseConfig() string {
	const rel = "organizejr-pp-leads/icp/ejs-comunicacao"
	candidates := []string{}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(cwd, rel),
			filepath.Join(cwd, "..", rel),
			filepath.Join(cwd, "..", "..", rel),
		)
	}
	if exe, err := os.Executable(); err == nil {
		base := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(base, rel),
			filepath.Join(base, "..", rel),
			filepath.Join(base, "..", "..", rel),
		)
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			return candidate
		}
	}
	return ""
}
