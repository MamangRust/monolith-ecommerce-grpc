package outbox

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type fakeExecutor struct {
	calls int
	query string
	args  []any
}

func (f *fakeExecutor) Exec(_ context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	f.calls++
	f.query = query
	f.args = args
	return pgconn.NewCommandTag("UPDATE 1"), nil
}

func (*fakeExecutor) Query(context.Context, string, ...any) (pgx.Rows, error) {
	return nil, errors.New("query not used by inbox tests")
}

func (f *fakeExecutor) QueryRow(_ context.Context, query string, args ...any) pgx.Row {
	f.calls++
	f.query = query
	f.args = args
	return stubRow{}
}

type stubRow struct{}

func (stubRow) Scan(dest ...any) error {
	for _, d := range dest {
		switch v := d.(type) {
		case *bool:
			*v = true
		case *int64:
			*v = 1
		}
	}
	return nil
}

func TestReserveValidatesKeys(t *testing.T) {
	for name, tc := range map[string]struct {
		tx           InboxExecutor
		consumerName string
		eventKey     string
	}{
		"nil executor": {tx: nil, consumerName: "email-service-group", eventKey: "topic:evt-1"},
		"empty consumer": {tx: &fakeExecutor{}, consumerName: "", eventKey: "topic:evt-1"},
		"empty event key": {tx: &fakeExecutor{}, consumerName: "email-service-group", eventKey: ""},
	} {
		t.Run(name, func(t *testing.T) {
			if _, _, _, err := Reserve(context.Background(), tc.tx, tc.consumerName, tc.eventKey, "topic", 0, 1); !errors.Is(err, ErrInvalidInboxKey) {
				t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
			}
		})
	}
}

func TestMarkProcessedValidatesKeys(t *testing.T) {
	if err := MarkProcessed(context.Background(), nil, "email-service-group", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
	if err := MarkProcessed(context.Background(), &fakeExecutor{}, "", "topic:evt-1", 1); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}
}

func TestReleaseValidatesKeysAndRecordsError(t *testing.T) {
	if err := Release(context.Background(), nil, "email-service-group", "topic:evt-1", 1, nil); !errors.Is(err, ErrInvalidInboxKey) {
		t.Fatalf("expected ErrInvalidInboxKey, got %v", err)
	}

	exec := &fakeExecutor{}
	if err := Release(context.Background(), exec, "email-service-group", "topic:evt-1", 1, errors.New("smtp down")); err != nil {
		t.Fatalf("Release returned error: %v", err)
	}
	if exec.calls != 1 || len(exec.args) != 4 {
		t.Fatalf("expected one four-argument update, calls=%d args=%d", exec.calls, len(exec.args))
	}
	if !strings.Contains(exec.args[2].(string), "smtp down") {
		t.Fatalf("expected last_error to carry the processing error, got %q", exec.args[2])
	}
}

func TestReserveUsesLeaseFencedUpsert(t *testing.T) {
	exec := &fakeExecutor{}
	if _, _, _, err := Reserve(context.Background(), exec, "email-service-group", "topic:evt-1", "topic", 0, 1); err != nil {
		t.Fatalf("Reserve returned error: %v", err)
	}
	if exec.calls != 1 {
		t.Fatalf("expected one query call, got %d", exec.calls)
	}
	for _, needle := range []string{"consumer_inbox", "ON CONFLICT", "lease_until <= current_timestamp"} {
		if !strings.Contains(exec.query, needle) {
			t.Fatalf("reserve query does not contain %q: %s", needle, exec.query)
		}
	}
}
