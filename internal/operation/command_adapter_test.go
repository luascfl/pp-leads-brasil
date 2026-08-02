package operation

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCommandAdapterPassesStoredTargetToProfileCommand(t *testing.T) {
	output := filepath.Join(t.TempDir(), "target.json")
	script := filepath.Join(t.TempDir(), "adapter.sh")
	content := "#!/bin/sh\ncat > \"$1\"\n"
	if err := os.WriteFile(script, []byte(content), 0o700); err != nil {
		t.Fatal(err)
	}
	adapter := CommandAdapter{Command: []string{script, output}}
	target := TargetChange{Key: "lead:1:site", System: "ploomes", RecordID: "1", Field: "Website", After: "https://example.test"}
	if err := adapter.Apply(context.Background(), target); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(data), `{"key":"lead:1:site","system":"ploomes","record_id":"1","field":"Website","before":"","after":"https://example.test","evidence":""}`; got != want {
		t.Fatalf("adapter input = %s, want %s", got, want)
	}
}

func TestCommandAdapterRejectsMissingCommand(t *testing.T) {
	if err := (CommandAdapter{}).Apply(context.Background(), TargetChange{}); err == nil {
		t.Fatal("missing command did not fail")
	}
}
