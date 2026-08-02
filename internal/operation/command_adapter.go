package operation

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// CommandAdapter invokes a profile-owned executable once for each stored target.
// The target is provided as JSON on standard input; no credentials are passed by
// the public platform. The configured command belongs to the external profile.
type CommandAdapter struct {
	Command []string
}

func (a CommandAdapter) Apply(ctx context.Context, target TargetChange) error {
	if len(a.Command) == 0 || strings.TrimSpace(a.Command[0]) == "" {
		return fmt.Errorf("operation adapter command is required")
	}
	payload, err := json.Marshal(target)
	if err != nil {
		return fmt.Errorf("serializing operation target: %w", err)
	}
	cmd := exec.CommandContext(ctx, a.Command[0], a.Command[1:]...)
	cmd.Stdin = bytes.NewReader(payload)
	output, err := cmd.CombinedOutput()
	if err != nil {
		message := strings.TrimSpace(string(output))
		if message == "" {
			return fmt.Errorf("running operation adapter: %w", err)
		}
		return fmt.Errorf("running operation adapter: %w: %s", err, message)
	}
	return nil
}
