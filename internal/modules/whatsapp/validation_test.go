package whatsapp

import (
	"encoding/base64"
	"testing"
)

func TestDecodeMediaDataURL(t *testing.T) {
	encoded := base64.StdEncoding.EncodeToString([]byte("media"))
	tests := []struct {
		name      string
		value     string
		allowed   []string
		wantError bool
	}{
		{
			name:    "valid image",
			value:   "data:image/png;base64," + encoded,
			allowed: []string{"image/*"},
		},
		{
			name:      "missing data URL header",
			value:     encoded,
			allowed:   []string{"image/*"},
			wantError: true,
		},
		{
			name:      "not base64",
			value:     "data:image/png,media",
			allowed:   []string{"image/*"},
			wantError: true,
		},
		{
			name:      "wrong MIME type",
			value:     "data:audio/ogg;base64," + encoded,
			allowed:   []string{"image/*"},
			wantError: true,
		},
		{
			name:      "empty payload",
			value:     "data:image/png;base64,",
			allowed:   []string{"image/*"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := decodeMediaDataURL(tt.value, tt.allowed...)
			if tt.wantError {
				if err == nil {
					t.Fatal("expected validation error")
				}
				return
			}
			if err != nil {
				t.Fatalf("decodeMediaDataURL() error = %v", err)
			}
			if string(decoded.Data) != "media" {
				t.Fatalf("decoded data = %q, want media", decoded.Data)
			}
		})
	}
}

func TestValidDocumentFilename(t *testing.T) {
	if !validDocumentFilename("report.pdf") {
		t.Fatal("expected non-empty filename to be valid")
	}
	if validDocumentFilename(" \t") {
		t.Fatal("expected blank filename to be invalid")
	}
}
