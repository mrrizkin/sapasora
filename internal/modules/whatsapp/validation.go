package whatsapp

import (
	"context"
	"errors"
	"fmt"
	"mime"
	"strings"
	"time"

	"github.com/vincent-petithory/dataurl"
)

const (
	providerOperationTimeout = 30 * time.Second
	defaultMediaUploadLimit  = 8 * 1024 * 1024
)

var (
	ErrInvalidDataURL = errors.New("invalid media data URL")
	ErrMediaTooLarge  = errors.New("media exceeds configured upload limit")
)

func providerContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, providerOperationTimeout)
}

// decodeMediaDataURL validates a base64 data URL, its MIME type, and decoded
// size before handing the bytes to a provider upload operation.
func decodeMediaDataURL(value string, maxBytes int, allowedMIMEs ...string) (*dataurl.DataURL, error) {
	if strings.TrimSpace(value) == "" || !strings.HasPrefix(value, "data:") {
		return nil, ErrInvalidDataURL
	}
	if maxBytes <= 0 {
		return nil, ErrMediaTooLarge
	}

	// Reject clearly oversized base64 payloads before DecodeString allocates the
	// decoded byte slice. The decoded length check below remains authoritative.
	if comma := strings.IndexByte(value, ','); comma >= 0 {
		maxEncodedBytes := ((maxBytes + 2) / 3) * 4
		if len(value)-comma-1 > maxEncodedBytes {
			return nil, fmt.Errorf("%w: limit is %d bytes", ErrMediaTooLarge, maxBytes)
		}
	}

	decoded, err := dataurl.DecodeString(value)
	if err != nil || decoded.Encoding != dataurl.EncodingBase64 || len(decoded.Data) == 0 {
		return nil, ErrInvalidDataURL
	}
	if len(decoded.Data) > maxBytes {
		return nil, fmt.Errorf("%w: limit is %d bytes", ErrMediaTooLarge, maxBytes)
	}

	mediaType := decoded.MediaType.Type + "/" + decoded.MediaType.Subtype
	if mediaType == "/" {
		return nil, ErrInvalidDataURL
	}
	if _, _, err := mime.ParseMediaType(mediaType); err != nil {
		return nil, fmt.Errorf("%w: invalid MIME type", ErrInvalidDataURL)
	}

	if len(allowedMIMEs) > 0 && !matchesAllowedMIME(mediaType, allowedMIMEs) {
		return nil, fmt.Errorf("%w: unsupported MIME type %q", ErrInvalidDataURL, mediaType)
	}

	return decoded, nil
}

func matchesAllowedMIME(mediaType string, allowedMIMEs []string) bool {
	for _, allowed := range allowedMIMEs {
		if strings.HasSuffix(allowed, "/*") {
			if strings.HasPrefix(mediaType, strings.TrimSuffix(allowed, "*")) {
				return true
			}
			continue
		}
		if mediaType == allowed {
			return true
		}
	}
	return false
}

func validDocumentFilename(filename string) bool {
	return strings.TrimSpace(filename) != ""
}

func mediaExtension(mediaType string) string {
	extensions, err := mime.ExtensionsByType(mediaType)
	if err != nil || len(extensions) == 0 {
		return ""
	}
	return extensions[0]
}
