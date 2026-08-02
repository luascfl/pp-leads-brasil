package cli

import (
	"strings"
	"testing"
)

func TestOperationApplyDoesNotTreatAgentAsApproval(t *testing.T) {
	flags := &rootFlags{}
	cmd := newRootCmd(flags)
	cmd.SetArgs([]string{"operation", "apply", "plan-123", "--agent"})
	err := cmd.Execute()
	if err == nil || !strings.Contains(err.Error(), "--agent is not approval") {
		t.Fatalf("apply error = %v", err)
	}
}
