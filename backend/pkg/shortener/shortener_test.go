package shortener

import "testing"

func TestBase62Shortener(t *testing.T) {
	s := NewBase62Shortener()

	tests := []struct {
		name    string
		id      int64
		want    string
		wantErr bool
	}{
		{"zero case", 0, "0", false},
		{"id one", 1, "1", false},
		{"base-1", 61, "Z", false},
		{"base", 62, "10", false},
		{"random large id", 12345, "3d7", false},
		{"max int64", 9223372036854775807, "aZl8N0y58M7", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test Encode
			got := s.Encode(tt.id)
			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}

			// Test Decode (Round-trip)
			back, err := s.Decode(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && back != tt.id {
				t.Errorf("Decode() = %v, want %v", back, tt.id)
			}
		})
	}
}

func TestBase62Shortener_DecodeErrors(t *testing.T) {
	s := NewBase62Shortener()

	tests := []struct {
		name      string
		shortCode string
	}{
		{"invalid symbol", "abc#123"},
		{"invalid emoji", "abc😊"},
		{"space in string", "abc 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Decode(tt.shortCode)
			if err == nil {
				t.Errorf("Decode(%v) expected error, got nil", tt.shortCode)
			}
		})
	}
}
