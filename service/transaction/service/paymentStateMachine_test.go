package service

import "testing"

func TestPaymentStateMachineValidity(t *testing.T) {
	valid := []string{PaymentStatusPending, PaymentStatusSuccess, PaymentStatusFailed, PaymentStatusRefunded}
	for _, status := range valid {
		if !IsValidPaymentStatus(status) {
			t.Errorf("expected %q to be a valid payment status", status)
		}
	}
	if !IsValidPaymentStatus("completed") {
		t.Error("expected legacy 'completed' to be accepted as a valid alias for 'success'")
	}
	for _, status := range []string{"", "canceled", "PENDING", "success "} {
		if IsValidPaymentStatus(status) {
			t.Errorf("expected %q to be an invalid payment status", status)
		}
	}
}

func TestPaymentStateMachineTransitions(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want bool
	}{
		{"pending to success", PaymentStatusPending, PaymentStatusSuccess, true},
		{"pending to failed", PaymentStatusPending, PaymentStatusFailed, true},
		{"pending to refunded", PaymentStatusPending, PaymentStatusRefunded, false},
		{"success to refunded", PaymentStatusSuccess, PaymentStatusRefunded, true},
		{"success to failed", PaymentStatusSuccess, PaymentStatusFailed, false},
		{"success to pending", PaymentStatusSuccess, PaymentStatusPending, false},
		{"failed to success", PaymentStatusFailed, PaymentStatusSuccess, false},
		{"failed to pending", PaymentStatusFailed, PaymentStatusPending, false},
		{"refunded to success", PaymentStatusRefunded, PaymentStatusSuccess, false},
		{"refunded to refunded", PaymentStatusRefunded, PaymentStatusRefunded, true},
		{"success to success", PaymentStatusSuccess, PaymentStatusSuccess, true},
		{"pending to pending", PaymentStatusPending, PaymentStatusPending, true},
		{"legacy completed to refunded", "completed", PaymentStatusRefunded, true},
		{"legacy completed to failed", "completed", PaymentStatusFailed, false},
		{"success to legacy completed", PaymentStatusSuccess, "completed", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanTransitionPaymentStatus(tt.from, tt.to); got != tt.want {
				t.Errorf("CanTransitionPaymentStatus(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}
