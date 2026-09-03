package event

import (
	"encoding/json"
	"testing"
)

func TestMarshalEmail_BuildsValidEnvelope(t *testing.T) {
	b, err := MarshalEmail("auth.register", "user@example.com", "Welcome", "<b>Hello</b>")
	if err != nil {
		t.Fatalf("MarshalEmail returned error: %v", err)
	}

	var env EmailEnvelope
	if err := json.Unmarshal(b, &env); err != nil {
		t.Fatalf("unmarshal envelope: %v", err)
	}

	if env.EventID == "" {
		t.Error("event_id must be non-empty")
	}
	if env.SchemaVersion != SchemaVersion {
		t.Errorf("schema_version = %d, want %d", env.SchemaVersion, SchemaVersion)
	}
	if env.EventType != "auth.register" {
		t.Errorf("event_type = %q, want auth.register", env.EventType)
	}
	if env.OccurredAt == "" {
		t.Error("occurred_at must be set")
	}
	if env.Email != "user@example.com" || env.Subject != "Welcome" || env.Body != "<b>Hello</b>" {
		t.Error("email/subject/body must be preserved")
	}
	if !env.IsValid() {
		t.Error("fresh envelope must be valid")
	}
}

func TestMarshalEmail_UniqueEventID(t *testing.T) {
	b1, err := MarshalEmail("a.b", "e@x.com", "s", "b")
	if err != nil {
		t.Fatal(err)
	}
	b2, err := MarshalEmail("a.b", "e@x.com", "s", "b")
	if err != nil {
		t.Fatal(err)
	}

	var e1, e2 EmailEnvelope
	_ = json.Unmarshal(b1, &e1)
	_ = json.Unmarshal(b2, &e2)
	if e1.EventID == e2.EventID {
		t.Error("event_id must be unique per event")
	}
}

func TestEmailEnvelope_IsValid(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*EmailEnvelope)
		wantErr bool
	}{
		{name: "valid", mutate: func(*EmailEnvelope) {}},
		{name: "missing event_id", mutate: func(e *EmailEnvelope) { e.EventID = "" }, wantErr: true},
		{name: "wrong schema_version", mutate: func(e *EmailEnvelope) { e.SchemaVersion = 0 }, wantErr: true},
		{name: "missing event_type", mutate: func(e *EmailEnvelope) { e.EventType = "" }, wantErr: true},
		{name: "missing email", mutate: func(e *EmailEnvelope) { e.Email = "" }, wantErr: true},
		{name: "missing subject", mutate: func(e *EmailEnvelope) { e.Subject = "" }, wantErr: true},
		{name: "missing body", mutate: func(e *EmailEnvelope) { e.Body = "" }, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := NewEmail("a.b", "e@x.com", "s", "b")
			tc.mutate(&env)
			if got := env.IsValid(); got == tc.wantErr {
				t.Errorf("IsValid() = %v, wantErr = %v", got, tc.wantErr)
			}
		})
	}
}
