package convert

import "testing"

func TestNullableString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want *string
	}{
		{name: "empty becomes nil", in: "", want: nil},
		{name: "non-empty keeps value", in: "hello", want: ptr("hello")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NullableString(tt.in)
			if tt.want == nil {
				if got != nil {
					t.Fatalf("NullableString(%q) = %v, want nil", tt.in, *got)
				}
				return
			}
			if got == nil {
				t.Fatalf("NullableString(%q) = nil, want %q", tt.in, *tt.want)
			}
			if *got != *tt.want {
				t.Fatalf("NullableString(%q) = %q, want %q", tt.in, *got, *tt.want)
			}
		})
	}
}

func TestStringPtr(t *testing.T) {
	got := StringPtr("value")
	if got == nil {
		t.Fatal("StringPtr returned nil")
	}
	if *got != "value" {
		t.Fatalf("StringPtr = %q, want %q", *got, "value")
	}
}

func ptr(s string) *string { return &s }
