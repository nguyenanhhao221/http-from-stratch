package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"httpfromtcp.haonguyen.tech/internal/request"
	"httpfromtcp.haonguyen.tech/internal/response"
	"httpfromtcp.haonguyen.tech/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, testHandler)
	if err != nil {
		log.Fatalf("error starting server: %v\n", err)
	}
	defer func() {
		if err := server.Close(); err != nil {
			log.Fatalf("error closing server: %v\n", err)
		}
	}()
	log.Println("Server started on port:", port)

	// Common pattern to exit the program.
	sigChn := make(chan os.Signal, 1)
	signal.Notify(sigChn, syscall.SIGINT, syscall.SIGTERM)
	<-sigChn
	log.Println("Server gracefully shutdown")
}

func testHandler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handler400(w, req)
		return

	case "/myproblem":
		handler500(w, req)
		return
	default:
		handler200(w, req)
		return
	}
}

func handler400(w *response.Writer, _ *request.Request) {
	err := w.WriteStatusLine(response.StatusBadRequest)
	if err != nil {
		log.Printf("error write status line :%v\n", err)
		return
	}
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	err = w.WriteHeaders(h)
	if err != nil {
		log.Printf("error when write header %v\n", err)
	}
	_, err = w.WriteBody(body)
	if err != nil {
		log.Printf("error when write body %v\n", err)
	}
}

func handler500(w *response.Writer, _ *request.Request) {
	if err := w.WriteStatusLine(response.StatusServerInternalError); err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")

	if err := w.WriteHeaders(h); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	if _, err := w.WriteBody(body); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
}

func handler200(w *response.Writer, _ *request.Request) {
	if err := w.WriteStatusLine(response.StatusOK); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)
	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	if err := w.WriteHeaders(h); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	if _, err := w.WriteBody(body); err != nil {
		log.Printf("error: %v\n", err)
		return
	}
}
