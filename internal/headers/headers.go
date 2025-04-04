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

	h.Set(key, string(value))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	h[key] = value
}
