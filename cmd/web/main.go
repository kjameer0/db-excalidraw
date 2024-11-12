// npx nodemon --exec go run ./cmd/web --signal SIGTERM -e go
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
)

// Define an application struct to hold the application-wide dependencies for the // web application.
type application struct {
	logger *slog.Logger
}

type fileWriter struct{}

func (f *fileWriter) Write(p []byte) (n int, err error) {
	logFile := "log.txt"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failure to open %s", logFile)
		return 0, err
	}
	defer file.Close()

	bytesWritten, writeErr := file.Write(p)
	if bytesWritten < len(p) {
		return 0, fmt.Errorf("failure to write all data to %s", logFile)
	}

	if writeErr != nil {
		log.Printf("Error writing to file %s: %v\n", logFile, writeErr)
		return 0, writeErr
	}

	return bytesWritten, nil
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	// mux is the part of the app that guides requests
	// to the url that matches their path
	mux := http.NewServeMux()
	writer := &fileWriter{}
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, writer), nil))
	slog.SetDefault(logger)
	// the file server with assets comes from a specific folder
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// StripPrefix gets rid of /static from the URL so we aren't searching
	// for /static/static/path-to-asset
	// create a get route for all assets
	mux.Handle("GET /static/", http.StripPrefix("/static", neuter(fileServer)))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)
	mux.HandleFunc("GET /drawing/{name}", getDrawingByName)
	mux.HandleFunc("POST /drawing/{$}", postDrawing)
	mux.HandleFunc("POST /compressed/drawing", postCompressedDrawing)
	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
