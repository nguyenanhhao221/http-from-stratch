package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// Read all bytes from the reader
	rawBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	// Convert bytes to string and split on CRLF to get request lines
	rawRequest := string(rawBytes)
	requestLines := strings.Split(rawRequest, "\r\n")

	// HTTP request must have at least request-line and 2 CRLF
	if len(requestLines) < 3 {
		return nil, errors.New("invalid request")
	}

	// Parse the request-line which is the first line
	// Request-line format: Method SP Request-Target SP HTTP-Version CRLF
	requestLineStr := requestLines[0]
	requestLineParts := strings.Split(requestLineStr, " ")
	if len(requestLineParts) != 3 {
		return nil, errors.New("invalid request line")
	}
	methodPart, requestTargetPart, httpVersion := requestLineParts[0], requestLineParts[1], requestLineParts[2]

	for _, c := range methodPart {
		if !unicode.IsUpper(c) {
			return nil, errors.New("invalid method: must be uppercase")
		}
	}

	if httpVersion != "HTTP/1.1" {
		return nil, errors.New("unsupport http version, only support HTTP/1.1")
	}

	versionParts := strings.Split(httpVersion, "/")
	if len(versionParts) != 2 {
		return nil, errors.New("invalid http version part")
	}

	// Create and return Request object from parsed components
	request := &Request{
		RequestLine: RequestLine{
			Method:        methodPart,
			RequestTarget: requestTargetPart,
			HttpVersion:   versionParts[1],
		},
	}

	return request, nil
}
