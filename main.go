package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {

	workingDir, err := os.Getwd()
	s := make([]byte, 0, 8)
	buffer := bytes.NewBuffer(s)
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

	for {
		_, err := io.CopyN(buffer, f, 8)
		fmt.Printf("read: %s\n", buffer.String())
		buffer.Reset()
		if err != nil {
			if err == io.EOF {
				break
			}
			slog.Error("reading from file to buffer", "error", err)
		}
	}
}
