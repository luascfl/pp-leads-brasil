package cli

import (
	"os"
	"path/filepath"
)

const ptyWrapperEnv = "PP_LEADS_PTY_WRAPPER"

func detectDefaultPtyWrapper() string {
	if explicit := os.Getenv(ptyWrapperEnv); explicit != "" {
		if info, err := os.Stat(explicit); err == nil && !info.IsDir() {
			return explicit
		}
	}
	const rel = "tmux_qterminal_setup/omp_pty_wrapper.py"
	candidates := []string{}
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates,
			filepath.Join(cwd, "..", rel),
			filepath.Join(cwd, "..", "..", rel),
			filepath.Join(cwd, rel),
		)
	}
	if exe, err := os.Executable(); err == nil {
		base := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(base, "..", rel),
			filepath.Join(base, "..", "..", rel),
		)
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate
		}
	}
	return ""
}
