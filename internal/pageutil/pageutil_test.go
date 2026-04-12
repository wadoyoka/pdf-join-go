package pageutil

import (
	"testing"
)

func TestParsePages(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []int
		wantErr bool
	}{
		{"single page", "3", []int{3}, false},
		{"multiple pages", "2,5,8", []int{2, 5, 8}, false},
		{"unsorted input", "8,2,5", []int{2, 5, 8}, false},
		{"with spaces", " 2 , 5 , 8 ", []int{2, 5, 8}, false},
		{"empty string", "", nil, true},
		{"non-numeric", "abc", nil, true},
		{"zero", "0", nil, true},
		{"negative", "-1", nil, true},
		{"duplicate", "2,2", nil, true},
		{"mixed invalid", "1,abc,3", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePages(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParsePages(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !tt.wantErr {
				if len(got) != len(tt.want) {
					t.Fatalf("ParsePages(%q) = %v, want %v", tt.input, got, tt.want)
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Fatalf("ParsePages(%q) = %v, want %v", tt.input, got, tt.want)
					}
				}
			}
		})
	}
}
