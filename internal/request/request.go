package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"httpfromtcp.haonguyen.tech/internal/headers"
)

type ParseState int

const (
	requestStateInitilized ParseState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	RequestLine    RequestLine
	Headers        headers.Headers
	Body           []byte
	bodyReadLength int
	ParseState     ParseState
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
		ParseState: requestStateInitilized,
		Headers:    headers.NewHeaders(),
		Body:       make([]byte, 0),
	}
	buffer := make([]byte, bufferSize)
	for r.ParseState != requestStateDone {
		if readToIndex >= len(buffer) {
			newBuf := make([]byte, len(buffer)*2)
			copy(newBuf, buffer)
			buffer = newBuf
		}
		numBytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				// if we "Read" to EOF and the parse state is not Done, an actual error happens
				if r.ParseState != requestStateDone {
					return nil, fmt.Errorf("incomplete request, in state: %d, read n bytes on EOF: %d", r.ParseState, numBytesRead)
				}
				break
			}
			return nil, err
		}
		readToIndex += numBytesRead

		numBytesParsed, err := r.parse(buffer[:readToIndex])
		if err != nil {
			return r, err
		}
		copy(buffer, buffer[numBytesParsed:])
		readToIndex -= numBytesParsed
	}
	return r, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.ParseState != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.ParseState {
	case requestStateInitilized:
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
		r.ParseState = requestStateParsingHeaders // Once request line is parse change state to start parse header
		return n, nil
	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			// just need more data
			return 0, nil
		}
		if done {
			r.ParseState = requestStateParsingBody
		}
		return n, nil
	case requestStateParsingBody:
		// If header doesn't contain Content-Length, don't parse the body
		contentLengthVal, ok := r.Headers.Get("Content-Length")
		if !ok {
			r.ParseState = requestStateDone
			return 0, nil
		}
		contentLengthValInt, err := strconv.Atoi(contentLengthVal)
		if err != nil {
			return 0, err
		}

		// Append the data to the r.Body
		r.Body = append(r.Body, data...)
		r.bodyReadLength += len(data)
		if r.bodyReadLength > contentLengthValInt {
			return 0, fmt.Errorf("Content-Length too large")
		}
		if r.bodyReadLength == contentLengthValInt {
			r.ParseState = requestStateDone
		}

		return len(data), nil

	case requestStateDone:
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
