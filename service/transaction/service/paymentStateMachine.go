package service

// Payment status values used across the transaction domain.
// The set is deliberately small and canonical; any status outside this set is
// rejected by validation instead of being persisted silently.
const (
	PaymentStatusPending  = "pending"
	PaymentStatusSuccess  = "success"
	PaymentStatusFailed   = "failed"
	PaymentStatusRefunded = "refunded"
)

// allowedPaymentTransitions defines which status may follow another status.
// Statuses not listed as a source cannot transition at all.
var allowedPaymentTransitions = map[string]map[string]bool{
	PaymentStatusPending: {
		PaymentStatusSuccess: true,
		PaymentStatusFailed:  true,
	},
	PaymentStatusSuccess: {
		PaymentStatusRefunded: true,
	},
	PaymentStatusFailed:   {},
	PaymentStatusRefunded: {},
}

// normalizePaymentStatus maps legacy status values onto the canonical set. The
// pre-existing schema used 'completed' as the default payment status; it is
// semantically identical to 'success', so rows created before the state machine
// must keep working instead of becoming permanently un-updatable.
func normalizePaymentStatus(status string) string {
	if status == "completed" {
		return PaymentStatusSuccess
	}
	return status
}

// IsValidPaymentStatus reports whether status is a known canonical payment status.
// The legacy 'completed' value is accepted as an alias for 'success'.
func IsValidPaymentStatus(status string) bool {
	_, ok := allowedPaymentTransitions[normalizePaymentStatus(status)]
	return ok
}

// CanTransitionPaymentStatus reports whether from may legally move to to.
// A transition to the same status is always allowed so idempotent updates and
// partial updates (where the caller echoes the current status) do not fail.
// Legacy 'completed' is normalized to 'success' on both sides.
func CanTransitionPaymentStatus(from, to string) bool {
	from = normalizePaymentStatus(from)
	to = normalizePaymentStatus(to)
	if from == to {
		return true
	}
	transitions, ok := allowedPaymentTransitions[from]
	if !ok {
		return false
	}
	return transitions[to]
}
