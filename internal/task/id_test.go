package task

import "testing"

func TestParseID(t *testing.T) {
	tests := []struct {
		name    string
		idStr   string
		wantID  int
		wantErr bool
	}{
		{
			name:    "valid numeric ID",
			idStr:   "42",
			wantID:  42,
			wantErr: false,
		},
		{
			name:    "valid large ID",
			idStr:   "999999",
			wantID:  999999,
			wantErr: false,
		},
		{
			name:    "invalid format - letters",
			idStr:   "abc",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "invalid format - mixed",
			idStr:   "12abc",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "whitespace",
			idStr:   "  42  ",
			wantID:  42,
			wantErr: false, // strconv.Atoi trims whitespace
		},
		{
			name:    "empty string",
			idStr:   "",
			wantID:  0,
			wantErr: true,
		},
		{
			name:    "negative ID",
			idStr:   "-5",
			wantID:  -5,
			wantErr: false, // Sscanf parses negatives successfully
		},
		{
			name:    "zero ID",
			idStr:   "0",
			wantID:  0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, err := ParseID(tt.idStr)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseID(%q) error = %v, wantErr %v", tt.idStr, err, tt.wantErr)
				return
			}

			if !tt.wantErr && gotID != tt.wantID {
				t.Errorf("ParseID(%q) = %d, want %d", tt.idStr, gotID, tt.wantID)
			}
		})
	}
}
