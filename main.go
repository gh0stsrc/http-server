package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	var currentLine string

	go func() {
		defer close(ch)
		defer f.Close()
		for {
			_, err := io.CopyN(buffer, f, 8)
			currentLine = currentLine + buffer.String()
			parts := strings.Split(currentLine, "\n")
			buffer.Reset()
			if len(parts) > 1 {
				msg := ""
				for i := range len(parts) - 1 {
					msg = msg + parts[i]
				}
				ch <- msg
				currentLine = parts[len(parts)-1]
			}
			if err != nil {
				if err == io.EOF {
					break
				}
				slog.Error("reading from file to buffer", "error", err)
			}
		}
	}()

	return ch
}

func main() {

	workingDir, err := os.Getwd()
	if err != nil {
		slog.Error("get working dir", "error", err)
		os.Exit(1)
	}
	f, err := os.Open(filepath.Join(workingDir, "messages.txt"))
	if err != nil {
		slog.Error("opening file", "error", err)
		os.Exit(1)
	}
	defer f.Close()

	ch := getLinesChannel(f)

	for {
		line, ok := <-ch
		if !ok {
			ch = nil
			break
		}
		fmt.Printf("read: %s\n", line)
	}

}
