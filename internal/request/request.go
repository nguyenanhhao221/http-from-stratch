package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type ParseState int

const (
	stateInitilized ParseState = iota
	stateDone
)

type Request struct {
	RequestLine RequestLine
	ParseState  ParseState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const crlf = "\r\n"

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	readToIndex := 0
	r := &Request{
		ParseState: stateInitilized,
	}
	buffer := make([]byte, bufferSize)
	for r.ParseState != stateDone {
		if readToIndex >= len(buffer) {
			newBuf := make([]byte, len(buffer)*2)
			copy(newBuf, buffer)
			buffer = newBuf
		}
		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				r.ParseState = stateDone
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := r.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}
		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.ParseState {
	case stateInitilized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			// something actually went wrong
			return 0, err
		}
		if n == 0 {
			// just need more data
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.ParseState = stateDone
		return n, nil
	case stateDone:
		return 0, fmt.Errorf("error: trying to read data in a done state")
	default:
		return 0, fmt.Errorf("unknown state: %q", r.ParseState)
	}
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {
	idx := bytes.Index(data, []byte(crlf))
	// If it cannot find the \r\n, it means we need more data before we process
	if idx == -1 {
		return nil, 0, nil
	}

	requestLineStr := string(data[:idx])
	requestLineParts := strings.Split(requestLineStr, " ")
	if len(requestLineParts) != 3 {
		return nil, 0, errors.New("invalid request line")
	}
	methodPart, requestTargetPart, httpVersion := requestLineParts[0], requestLineParts[1], requestLineParts[2]

	for _, c := range methodPart {
		if !unicode.IsUpper(c) {
			return nil, 0, errors.New("invalid method: must be uppercase")
		}
	}

	if httpVersion != "HTTP/1.1" {
		return nil, 0, errors.New("unsupport http version, only support HTTP/1.1")
	}

	versionParts := strings.Split(httpVersion, "/")
	if len(versionParts) != 2 {
		return nil, 0, errors.New("invalid http version part")
	}

	// Update the Request struct with the Parsed RequestLine
	return &RequestLine{
		Method:        methodPart,
		RequestTarget: requestTargetPart,
		HttpVersion:   versionParts[1],
	}, idx + 2, nil
}
