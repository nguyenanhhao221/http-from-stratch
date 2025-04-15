package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// Find of data have the crlf
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, done, nil
	}

	// If it's at the start of the data, you've found the end of the headers, so return the proper values immediately.
	if idx == 0 {
		return idx + 2, true, nil
	}

	parts := bytes.SplitN(data[:idx], []byte(":"), 2)

	key := string(parts[0])
	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	// Trim leading white spaces
	key = strings.TrimSpace(key)
	value := bytes.TrimSpace(parts[1])

	// Validate field name
	if err := validateFieldName(key); err != nil {
		return 0, false, err
	}

	h.Set(key, string(value))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	// Convert key to lowercase
	key = strings.ToLower(key)

	if existVal, ok := h[key]; ok {
		h[key] = strings.Join([]string{existVal, value}, ", ")
	} else {
		h[key] = value
	}
}

func (h Headers) Get(key string) (string, bool) {
	// Convert to lower case because http header is case-insensitive
	lowerCaseKey := strings.ToLower(key)
	v, ok := h[lowerCaseKey]
	return v, ok
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Delete(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}

func validateFieldName(key string) error {
	if len(key) < 1 {
		return fmt.Errorf("field name length must be at least 1")
	}
	allowedSpecials := "!#$%&'*+-.^_`|~"
	for _, c := range key {
		if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && !strings.ContainsRune(allowedSpecials, c) {
			return fmt.Errorf("field name contains invalid character: %q", c)
		}
	}
	return nil
}
