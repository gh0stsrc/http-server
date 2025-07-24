package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
)

func getLinesChannel(rc io.ReadCloser) <-chan string {
	ch := make(chan string)
	buffer := bytes.NewBuffer(make([]byte, 0, 8))
	var currentLine string

	go func() {
		defer close(ch)
		defer rc.Close()
		for {
			_, err := io.CopyN(buffer, rc, 8)
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

	l, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		slog.Error("listener setup", "error", err)
		os.Exit(1)
	}
	defer l.Close()

	for {
		connection, err := l.Accept()
		if err != nil {
			slog.Error("connection", "error", err)
		}
		slog.Info("connection established", "address", connection.LocalAddr().String())
		ch := getLinesChannel(connection)

		for {
			line, ok := <-ch
			if !ok {
				ch = nil
				break
			}
			fmt.Printf("read: %s\n", line)
		}

	}

}
