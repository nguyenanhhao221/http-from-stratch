package response

import (
	"fmt"
	"io"
	"log"

	"httpfromtcp.haonguyen.tech/internal/headers"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

type Writer struct {
	writer      io.Writer
	writerState writerState
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state: %d", w.writerState)
	}
	// Set the next state after write status line
	defer func() { w.writerState = writerStateHeaders }()

	statusLine, ok := statusCodeMap[statusCode]
	if !ok {
		log.Println("cannot find default reason phrase status")
		statusLine = fmt.Sprintf("HTTP/1.1 %d \r\n", statusCode)
	}

	if _, err := w.writer.Write([]byte(statusLine)); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("cannot write header in state: %d", w.writerState)
	}
	defer func() { w.writerState = writerStateBody }()

	for k, v := range headers {
		if _, err := fmt.Fprintf(w.writer, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}
	// write empty line by the end of headers
	if _, err := w.writer.Write([]byte("\r\n")); err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state: %d", w.writerState)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state: %d", w.writerState)
	}
	chunkSize := len(p)
	// Write the chunk size with \r\n
	str := fmt.Sprintf("%x\r\n", chunkSize)
	_, err := w.writer.Write([]byte(str))
	if err != nil {
		return 0, err
	}

	return w.writer.Write(append(p, []byte("\r\n")...))
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	b := []byte("0\r\n\r\n")
	return w.writer.Write(b)
}
