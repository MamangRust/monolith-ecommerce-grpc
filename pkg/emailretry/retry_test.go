package emailretry

import (
	"errors"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

func TestHeadersRoundTrip(t *testing.T) {
	meta := RetryMeta{
		SrcTopic: "email-service-topic-auth-register", SrcPartition: 2, SrcOffset: 42,
		Attempt: 3, RetryAt: time.Date(2026, 8, 16, 1, 0, 0, 0, time.UTC), Reason: "smtp: 421 busy",
	}
	// Consumer messages carry pointer headers — convert the producer value
	// headers to match sarama's ConsumerMessage shape.
	valueHeaders := BuildHeaders(meta)
	msg := &sarama.ConsumerMessage{Headers: toPtrHeaders(valueHeaders)}

	got, err := ParseMeta(msg)
	if err != nil {
		t.Fatalf("ParseMeta returned error: %v", err)
	}
	if got != meta {
		t.Fatalf("round-trip mismatch:\n got %+v\nwant %+v", got, meta)
	}
}

func toPtrHeaders(in []sarama.RecordHeader) []*sarama.RecordHeader {
	out := make([]*sarama.RecordHeader, 0, len(in))
	for i := range in {
		h := in[i]
		out = append(out, &h)
	}
	return out
}

func TestBuildDLQHeadersAddsAttempts(t *testing.T) {
	meta := RetryMeta{
		SrcTopic: "t", SrcPartition: 0, SrcOffset: 1, Attempt: 5,
		RetryAt: time.Now().UTC(), Reason: "boom", Attempts: 5,
	}
	headers := BuildDLQHeaders(meta)

	found := false
	for _, h := range headers {
		if string(h.Key) == HdrAttempts && string(h.Value) == "5" {
			found = true
		}
	}
	if !found {
		t.Fatal("DLQ headers must carry the final attempt count")
	}
	if len(headers) != len(BuildHeaders(meta))+1 {
		t.Fatalf("DLQ headers must extend retry headers, got %d", len(headers))
	}
}

func TestParseMetaRejectsMissingHeaders(t *testing.T) {
	cases := map[string]*sarama.ConsumerMessage{
		"nil message": nil,
		"no headers":  {Headers: nil},
		"missing src": {Headers: []*sarama.RecordHeader{{Key: []byte(HdrAttempt), Value: []byte("1")}}},
	}
	for name, msg := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseMeta(msg); !errors.Is(err, ErrMissingRetryMeta) {
				t.Fatalf("expected ErrMissingRetryMeta, got %v", err)
			}
		})
	}
}

func TestBackoffIsLinearAndCapped(t *testing.T) {
	if got := Backoff(1, 30*time.Second); got != 30*time.Second {
		t.Fatalf("attempt 1 backoff = %v, want 30s", got)
	}
	if got := Backoff(3, 30*time.Second); got != 90*time.Second {
		t.Fatalf("attempt 3 backoff = %v, want 90s", got)
	}
	// 30s * 1000 attempts must be capped at MaxBackoff.
	if got := Backoff(1000, 30*time.Second); got != MaxBackoff {
		t.Fatalf("capped backoff = %v, want %v", got, MaxBackoff)
	}
	// Non-positive base falls back to the default.
	if got := Backoff(2, 0); got != 2*DefaultBackoff {
		t.Fatalf("default backoff = %v, want %v", got, 2*DefaultBackoff)
	}
}

func TestNextRetryAtIsInFuture(t *testing.T) {
	before := time.Now().UTC()
	at := NextRetryAt(1, time.Second)
	if at.Before(before) {
		t.Fatalf("NextRetryAt returned past time %v", at)
	}
}
