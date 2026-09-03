// Package emailretry provides the shared retry/DLQ contract for the email
// service (Phase 4 — Unified Retry dan DLQ): topic names, per-message retry
// metadata headers and backoff computation. Producers and the retry processor
// must agree on these constants so retries and dead letters can be operated
// and replayed safely.
package emailretry

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/IBM/sarama"
)

const (
	// RetryTopic is the shared retry topic for transient SMTP failures.
	RetryTopic = "email-service-topic-email-retry"
	// DLQTopic is the shared dead-letter topic for events that exhausted retries.
	DLQTopic = "email-service-topic-email-dlq"
	// RetryGroup is the consumer group of the retry processor.
	RetryGroup = "email-service-retry-group"
)

// Retry metadata header keys. Kept on the message headers (never in the
// payload) so the payload remains the original envelope: replaying a DLQ event
// only requires republishing the payload to the source topic, and the
// consumer-side inbox guarantees no duplicate email if the event was
// already sent.
const (
	HdrSrcTopic     = "_srcTopic"
	HdrSrcPartition = "_srcPartition"
	HdrSrcOffset    = "_srcOffset"
	HdrAttempt      = "_attempt"
	HdrRetryAt      = "_retryAt"
	HdrReason       = "_reason"
	HdrAttempts     = "_attempts"
)

// Defaults for retry escalation.
const (
	DefaultMaxAttempts = 5
	DefaultBackoff     = 30 * time.Second
	MaxBackoff         = 10 * time.Minute
)

// ErrMissingRetryMeta indicates a retry-topic message without the required
// metadata headers. Such messages are treated as poison.
var ErrMissingRetryMeta = errors.New("message is missing retry metadata headers")

// RetryMeta carries the retry/DLQ metadata of one message.
type RetryMeta struct {
	SrcTopic     string
	SrcPartition int32
	SrcOffset    int64
	Attempt      int
	RetryAt      time.Time
	Reason       string
	Attempts     int // only set on DLQ messages (final attempt count)
}

// BuildHeaders serializes meta into sarama record headers (Phase 4 contract).
// Producer messages carry value headers ([]sarama.RecordHeader); consumer
// messages carry pointers — see sarama's ProducerMessage vs ConsumerMessage.
func BuildHeaders(meta RetryMeta) []sarama.RecordHeader {
	return []sarama.RecordHeader{
		{Key: []byte(HdrSrcTopic), Value: []byte(meta.SrcTopic)},
		{Key: []byte(HdrSrcPartition), Value: []byte(strconv.FormatInt(int64(meta.SrcPartition), 10))},
		{Key: []byte(HdrSrcOffset), Value: []byte(strconv.FormatInt(meta.SrcOffset, 10))},
		{Key: []byte(HdrAttempt), Value: []byte(strconv.Itoa(meta.Attempt))},
		{Key: []byte(HdrRetryAt), Value: []byte(meta.RetryAt.UTC().Format(time.RFC3339))},
		{Key: []byte(HdrReason), Value: []byte(meta.Reason)},
	}
}

// BuildDLQHeaders serializes meta plus the final attempt count into headers.
func BuildDLQHeaders(meta RetryMeta) []sarama.RecordHeader {
	headers := BuildHeaders(meta)
	return append(headers, sarama.RecordHeader{Key: []byte(HdrAttempts), Value: []byte(strconv.Itoa(meta.Attempts))})
}

// ParseMeta reads retry metadata from a consumed message's headers.
func ParseMeta(msg *sarama.ConsumerMessage) (RetryMeta, error) {
	if msg == nil {
		return RetryMeta{}, ErrMissingRetryMeta
	}
	values := make(map[string]string, len(msg.Headers))
	for _, h := range msg.Headers {
		if h != nil {
			values[string(h.Key)] = string(h.Value)
		}
	}

	meta := RetryMeta{SrcTopic: values[HdrSrcTopic]}
	if meta.SrcTopic == "" {
		return RetryMeta{}, fmt.Errorf("%w: missing %s", ErrMissingRetryMeta, HdrSrcTopic)
	}
	partition, err := strconv.ParseInt(values[HdrSrcPartition], 10, 32)
	if err != nil {
		return RetryMeta{}, fmt.Errorf("%w: invalid %s: %v", ErrMissingRetryMeta, HdrSrcPartition, err)
	}
	meta.SrcPartition = int32(partition)

	offset, err := strconv.ParseInt(values[HdrSrcOffset], 10, 64)
	if err != nil {
		return RetryMeta{}, fmt.Errorf("%w: invalid %s: %v", ErrMissingRetryMeta, HdrSrcOffset, err)
	}
	meta.SrcOffset = offset

	attempt, err := strconv.Atoi(values[HdrAttempt])
	if err != nil {
		return RetryMeta{}, fmt.Errorf("%w: invalid %s: %v", ErrMissingRetryMeta, HdrAttempt, err)
	}
	meta.Attempt = attempt

	retryAt, err := time.Parse(time.RFC3339, values[HdrRetryAt])
	if err != nil {
		return RetryMeta{}, fmt.Errorf("%w: invalid %s: %v", ErrMissingRetryMeta, HdrRetryAt, err)
	}
	meta.RetryAt = retryAt
	meta.Reason = values[HdrReason]

	if attempts, err := strconv.Atoi(values[HdrAttempts]); err == nil {
		meta.Attempts = attempts
	}
	return meta, nil
}

// Backoff returns the delay before the given attempt: base * attempt, capped at
// MaxBackoff so a retry processor never sleeps for an unbounded period.
func Backoff(attempt int, base time.Duration) time.Duration {
	if base <= 0 {
		base = DefaultBackoff
	}
	if attempt < 1 {
		attempt = 1
	}
	d := base * time.Duration(attempt)
	if d > MaxBackoff {
		return MaxBackoff
	}
	return d
}

// NextRetryAt computes when the next attempt should run.
func NextRetryAt(attempt int, base time.Duration) time.Time {
	return time.Now().UTC().Add(Backoff(attempt, base))
}
