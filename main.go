package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	// open the file
	f, err := os.Open("message.txt")
	if err != nil {
		log.Fatalln(err)
	}

	// Read the file 8 bytes at a time
	buffer := make([]byte, 8)

	currentLineStr := ""
	for {
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatalln(err)
		}
		// n == 0 means that we reach the end of the file
		if n == 0 {
			break
		}
		parts := strings.Split(string(buffer[:n]), "\n")
		for i := range len(parts) - 1 {
			fmt.Printf("read: %s%s\n", currentLineStr, parts[i])
			currentLineStr = ""
		}
		currentLineStr += parts[len(parts)-1]
	}
}
