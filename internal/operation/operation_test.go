package operation

import (
	"context"
	"errors"
	"testing"
	"time"
)

type recordingAdapter struct {
	calls int
	fail  map[string]error
}

func (a *recordingAdapter) Apply(_ context.Context, target TargetChange) error {
	a.calls++
	return a.fail[target.Key]
}

func planInput() PlanInput {
	return PlanInput{
		Profile:  "test-profile",
		Kind:     "crm-update",
		Operator: "tester",
		Targets: []TargetChange{{
			Key:      "company:1:website",
			System:   "crm",
			RecordID: "1",
			Field:    "website",
			Before:   "",
			After:    "https://example.test",
			Evidence: "https://example.test/about",
		}},
	}
}

func newService(t *testing.T, adapter Adapter, now time.Time) Service {
	t.Helper()
	store, err := NewFileStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return Service{Store: store, Adapter: adapter, Now: func() time.Time { return now }}
}

func TestPlanIsImmutableAndSideEffectFree(t *testing.T) {
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	adapter := &recordingAdapter{}
	service := newService(t, adapter, now)

	plan, err := service.Plan(planInput())
	if err != nil {
		t.Fatal(err)
	}
	if plan.ID == "" || plan.Digest == "" {
		t.Fatalf("plan identity missing: %#v", plan)
	}
	if adapter.calls != 0 {
		t.Fatalf("adapter calls = %d, want 0", adapter.calls)
	}
	stored, err := service.Store.LoadPlan(plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored.Digest != plan.Digest {
		t.Fatalf("stored digest = %q, want %q", stored.Digest, plan.Digest)
	}
}

func TestApplyRequiresApprovalAndIsIdempotent(t *testing.T) {
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	adapter := &recordingAdapter{}
	service := newService(t, adapter, now)
	plan, err := service.Plan(planInput())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Apply(context.Background(), plan.ID, false); !errors.Is(err, ErrApprovalRequired) {
		t.Fatalf("Apply without approval error = %v", err)
	}
	result, err := service.Apply(context.Background(), plan.ID, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Receipts) != 1 || result.Receipts[0].Status != "applied" {
		t.Fatalf("receipts = %#v", result.Receipts)
	}
	if adapter.calls != 1 {
		t.Fatalf("adapter calls = %d, want 1", adapter.calls)
	}
	if _, err := service.Apply(context.Background(), plan.ID, true); err != nil {
		t.Fatal(err)
	}
	if adapter.calls != 1 {
		t.Fatalf("adapter calls after replay = %d, want 1", adapter.calls)
	}
}

func TestApplyRecordsPartialFailureWithoutReplayingSuccess(t *testing.T) {
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	adapter := &recordingAdapter{fail: map[string]error{"company:1:email": errors.New("remote failure")}}
	service := newService(t, adapter, now)
	input := planInput()
	input.Targets = append(input.Targets, TargetChange{Key: "company:1:email", System: "crm", RecordID: "1", Field: "email", Before: "", After: "info@example.test", Evidence: "https://example.test/contact"})
	plan, err := service.Plan(input)
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.Apply(context.Background(), plan.ID, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Receipts) != 2 || result.Receipts[0].Status != "failed" || result.Receipts[1].Status != "applied" {
		t.Fatalf("receipts = %#v", result.Receipts)
	}
	if adapter.calls != 2 {
		t.Fatalf("adapter calls = %d, want 2", adapter.calls)
	}
	if _, err := service.Apply(context.Background(), plan.ID, true); err != nil {
		t.Fatal(err)
	}
	if adapter.calls != 3 {
		t.Fatalf("adapter calls after retry = %d, want 3", adapter.calls)
	}
}

func TestPlanRejectsMissingEvidenceAndExpiredTimestamp(t *testing.T) {
	now := time.Date(2026, 8, 2, 12, 0, 0, 0, time.UTC)
	service := newService(t, nil, now)
	input := planInput()
	input.Targets[0].Evidence = ""
	if _, err := service.Plan(input); !errors.Is(err, ErrPlanInvalid) {
		t.Fatalf("missing evidence error = %v", err)
	}
	input = planInput()
	input.ExpiresAt = now
	if _, err := service.Plan(input); !errors.Is(err, ErrPlanInvalid) {
		t.Fatalf("expired plan error = %v", err)
	}
}
