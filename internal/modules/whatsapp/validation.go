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

const providerOperationTimeout = 30 * time.Second

var ErrInvalidDataURL = errors.New("invalid media data URL")

func providerContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithTimeout(ctx, providerOperationTimeout)
}

// decodeMediaDataURL validates a base64 data URL and its MIME type before decoding.
func decodeMediaDataURL(value string, allowedMIMEs ...string) (*dataurl.DataURL, error) {
	if strings.TrimSpace(value) == "" || !strings.HasPrefix(value, "data:") {
		return nil, ErrInvalidDataURL
	}

	decoded, err := dataurl.DecodeString(value)
	if err != nil || decoded.Encoding != dataurl.EncodingBase64 || len(decoded.Data) == 0 {
		return nil, ErrInvalidDataURL
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
