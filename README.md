# http-from-scratch

This project is a hands-on exploration of how the HTTP protocol works under the hood, built from scratch in Go. The main goal is to demystify the abstractions provided by high-level frameworks and libraries, and to gain a deep understanding of how HTTP servers operate at the TCP level.

## Why?

As a developer, I've always worked with HTTP through high-level abstractions (like `net/http` in Go, Express in Node.js, etc.). This project is my attempt to peel back those layers and learn how HTTP really works—by implementing it myself, starting from raw TCP sockets.

## What's Inside

- **Custom HTTP Server:**  
  Implements an HTTP/1.1 server directly on top of TCP sockets, handling request parsing, header management, and response formatting manually.
- **TCP Listener:**  
  A simple TCP server that accepts connections and prints out raw HTTP requests, for learning and debugging.
- **UDP Listener:**  
  A basic UDP client for sending messages, included for protocol comparison and experimentation.
- **Manual Request/Response Parsing:**  
  All HTTP request and response parsing, header management, and state handling are implemented from scratch—no use of Go's `net/http` abstractions for these parts.
- **Learning-Oriented Handlers:**  
  Example handlers for different HTTP status codes, proxying, and chunked transfer encoding, to demonstrate protocol features.

## Project Structure

```
cmd/
  httpserver/      # Main HTTP server entry point
  tcplistener/     # Simple TCP listener for raw request inspection
  udplistener/     # UDP client for protocol comparison
internal/
  server/          # Core TCP server logic
  request/         # HTTP request parsing and state machine
  response/        # HTTP response formatting and writing
  headers/         # HTTP header parsing and validation
```

## How to Run

### HTTP Server

```sh
go run ./cmd/httpserver
```
- Listens on port `42069` by default.
- Handles basic HTTP requests, returns different responses based on the path (see `testHandler` in `cmd/httpserver/main.go`).

### TCP Listener

```sh
go run ./cmd/tcplistener
```
- Listens for raw TCP connections on port `42069`.
- Prints out the raw HTTP request line, headers, and body for inspection.

### UDP Listener

```sh
go run ./cmd/udplistener
```
- Simple UDP client for sending messages to `localhost:42069`.

## What You'll Learn

- How to accept and manage TCP connections in Go.
- How to parse HTTP/1.1 requests and headers manually.
- How to construct and send valid HTTP responses, including status lines, headers, and bodies.
- How chunked transfer encoding and proxying work at the protocol level.
- The differences between TCP and UDP for network communication.

## Notes

- This project is for educational purposes and is **not production-ready**.
- Only a subset of HTTP/1.1 is implemented (enough to understand the basics and experiment).
- Error handling and edge cases are handled in a way that prioritizes learning and clarity.
