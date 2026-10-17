package newsletter_service

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func decodeNewsletterPicture(input string) ([]byte, error) {
	if input == "" {
		return nil, nil
	}
	encoded := input
	if strings.HasPrefix(input, "data:") {
		prefix, value, ok := strings.Cut(input, ",")
		if !ok || !strings.HasSuffix(prefix, ";base64") {
			return nil, errors.New("invalid image data URL")
		}
		mime := strings.TrimPrefix(strings.TrimSuffix(prefix, ";base64"), "data:")
		if mime != "image/jpeg" && mime != "image/png" {
			return nil, errors.New("unsupported image format")
		}
		encoded = value
	}
	if len(encoded) > 7*1024*1024 {
		return nil, errors.New("image exceeds maximum size")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, errors.New("invalid image base64")
	}
	if len(data) > 5*1024*1024 {
		return nil, errors.New("image exceeds 5MB")
	}
	detected := http.DetectContentType(data)
	if detected != "image/jpeg" && detected != "image/png" {
		return nil, errors.New("invalid image content")
	}
	return data, nil
}
