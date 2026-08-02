package operation

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	ErrApprovalRequired = errors.New("explicit approval is required")
	ErrPlanExpired      = errors.New("operation plan has expired")
	ErrPlanInvalid      = errors.New("invalid operation plan")
)

type TargetChange struct {
	Key      string `json:"key"`
	System   string `json:"system"`
	RecordID string `json:"record_id"`
	Field    string `json:"field"`
	Before   string `json:"before"`
	After    string `json:"after"`
	Evidence string `json:"evidence"`
}

type PlanInput struct {
	Profile   string         `json:"profile"`
	Kind      string         `json:"kind"`
	Operator  string         `json:"operator"`
	Targets   []TargetChange `json:"targets"`
	ExpiresAt time.Time      `json:"expires_at"`
}

type Plan struct {
	ID        string         `json:"id"`
	Profile   string         `json:"profile"`
	Kind      string         `json:"kind"`
	Operator  string         `json:"operator"`
	Targets   []TargetChange `json:"targets"`
	CreatedAt time.Time      `json:"created_at"`
	ExpiresAt time.Time      `json:"expires_at"`
	Digest    string         `json:"digest"`
}

type Receipt struct {
	PlanID      string       `json:"plan_id"`
	Target      TargetChange `json:"target"`
	Operator    string       `json:"operator"`
	AttemptedAt time.Time    `json:"attempted_at"`
	Status      string       `json:"status"`
	Error       string       `json:"error,omitempty"`
}

type Result struct {
	PlanID   string    `json:"plan_id"`
	Digest   string    `json:"digest"`
	Receipts []Receipt `json:"receipts"`
}

type Store interface {
	SavePlan(Plan) error
	LoadPlan(string) (Plan, error)
	SaveReceipt(Receipt) error
	LoadReceipt(planID, targetKey string) (Receipt, bool, error)
}

type Adapter interface {
	Apply(context.Context, TargetChange) error
}

type Service struct {
	Store   Store
	Adapter Adapter
	Now     func() time.Time
}

func (s Service) Plan(input PlanInput) (Plan, error) {
	if s.Store == nil {
		return Plan{}, fmt.Errorf("operation plan store is required")
	}
	if err := validateInput(input); err != nil {
		return Plan{}, err
	}
	now := s.now().UTC()
	expiresAt := input.ExpiresAt.UTC()
	if expiresAt.IsZero() {
		expiresAt = now.Add(30 * time.Minute)
	}
	if !expiresAt.After(now) {
		return Plan{}, fmt.Errorf("%w: expiration must be in the future", ErrPlanInvalid)
	}
	targets := append([]TargetChange(nil), input.Targets...)
	sort.Slice(targets, func(i, j int) bool { return targets[i].Key < targets[j].Key })
	plan := Plan{
		Profile:   strings.TrimSpace(input.Profile),
		Kind:      strings.TrimSpace(input.Kind),
		Operator:  strings.TrimSpace(input.Operator),
		Targets:   targets,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}
	digest, err := digestPlan(plan)
	if err != nil {
		return Plan{}, err
	}
	plan.Digest = digest
	plan.ID = digest[:24]
	if err := s.Store.SavePlan(plan); err != nil {
		return Plan{}, err
	}
	return plan, nil
}

func (s Service) Apply(ctx context.Context, planID string, approved bool) (Result, error) {
	if !approved {
		return Result{}, ErrApprovalRequired
	}
	if s.Store == nil {
		return Result{}, fmt.Errorf("operation plan store is required")
	}
	plan, err := s.Store.LoadPlan(strings.TrimSpace(planID))
	if err != nil {
		return Result{}, err
	}
	digest, err := digestPlan(plan)
	if err != nil || digest != plan.Digest || plan.ID != plan.Digest[:24] {
		return Result{}, fmt.Errorf("%w: stored plan digest mismatch", ErrPlanInvalid)
	}
	if !plan.ExpiresAt.After(s.now()) {
		return Result{}, ErrPlanExpired
	}
	result := Result{PlanID: plan.ID, Digest: plan.Digest, Receipts: make([]Receipt, 0, len(plan.Targets))}
	for _, target := range plan.Targets {
		if receipt, found, err := s.Store.LoadReceipt(plan.ID, target.Key); err != nil {
			return Result{}, err
		} else if found && receipt.Status == "applied" {
			result.Receipts = append(result.Receipts, receipt)
			continue
		}
		receipt := Receipt{PlanID: plan.ID, Target: target, Operator: plan.Operator, AttemptedAt: s.now().UTC(), Status: "applied"}
		if s.Adapter == nil {
			receipt.Status = "failed"
			receipt.Error = "no external operation adapter is configured for this profile"
		} else if err := s.Adapter.Apply(ctx, target); err != nil {
			receipt.Status = "failed"
			receipt.Error = err.Error()
		}
		if err := s.Store.SaveReceipt(receipt); err != nil {
			return Result{}, err
		}
		result.Receipts = append(result.Receipts, receipt)
	}
	return result, nil
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now().UTC()
	}
	return time.Now().UTC()
}

func validateInput(input PlanInput) error {
	if strings.TrimSpace(input.Profile) == "" || strings.TrimSpace(input.Kind) == "" || strings.TrimSpace(input.Operator) == "" {
		return fmt.Errorf("%w: profile, kind and operator are required", ErrPlanInvalid)
	}
	if len(input.Targets) == 0 {
		return fmt.Errorf("%w: at least one target is required", ErrPlanInvalid)
	}
	seen := make(map[string]struct{}, len(input.Targets))
	for _, target := range input.Targets {
		if strings.TrimSpace(target.Key) == "" || strings.TrimSpace(target.System) == "" || strings.TrimSpace(target.RecordID) == "" || strings.TrimSpace(target.Field) == "" || strings.TrimSpace(target.Evidence) == "" {
			return fmt.Errorf("%w: target key, system, record_id, field and evidence are required", ErrPlanInvalid)
		}
		if _, exists := seen[target.Key]; exists {
			return fmt.Errorf("%w: duplicate target %q", ErrPlanInvalid, target.Key)
		}
		seen[target.Key] = struct{}{}
	}
	return nil
}

func digestPlan(plan Plan) (string, error) {
	copy := plan
	copy.ID = ""
	copy.Digest = ""
	payload, err := json.Marshal(copy)
	if err != nil {
		return "", fmt.Errorf("serializing operation plan: %w", err)
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}

type FileStore struct {
	Root string
}

func NewFileStore(root string) (*FileStore, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return nil, fmt.Errorf("operation store root is required")
	}
	resolved, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	return &FileStore{Root: resolved}, nil
}

func (s *FileStore) SavePlan(plan Plan) error {
	return s.writeJSON(filepath.Join("plans", plan.ID+".json"), plan)
}

func (s *FileStore) LoadPlan(id string) (Plan, error) {
	var plan Plan
	if err := s.readJSON(filepath.Join("plans", safeName(id)+".json"), &plan); err != nil {
		return Plan{}, err
	}
	return plan, nil
}

func (s *FileStore) SaveReceipt(receipt Receipt) error {
	return s.writeJSON(filepath.Join("receipts", receipt.PlanID, safeName(receipt.Target.Key)+".json"), receipt)
}

func (s *FileStore) LoadReceipt(planID, targetKey string) (Receipt, bool, error) {
	var receipt Receipt
	err := s.readJSON(filepath.Join("receipts", safeName(planID), safeName(targetKey)+".json"), &receipt)
	if errors.Is(err, os.ErrNotExist) {
		return Receipt{}, false, nil
	}
	if err != nil {
		return Receipt{}, false, err
	}
	return receipt, true, nil
}

func (s *FileStore) writeJSON(relative string, value any) error {
	path, err := s.path(relative)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(path), ".operation-*")
	if err != nil {
		return err
	}
	name := temp.Name()
	defer os.Remove(name)
	if _, err := temp.Write(data); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Chmod(0o600); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(name, path)
}

func (s *FileStore) readJSON(relative string, target any) error {
	path, err := s.path(relative)
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (s *FileStore) path(relative string) (string, error) {
	if s == nil || strings.TrimSpace(s.Root) == "" {
		return "", fmt.Errorf("operation store root is required")
	}
	path := filepath.Join(s.Root, relative)
	rel, err := filepath.Rel(s.Root, path)
	if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return "", fmt.Errorf("invalid operation store path")
	}
	return path, nil
}

func safeName(value string) string {
	return strings.NewReplacer("/", "_", "\\", "_", "..", "_").Replace(value)
}
