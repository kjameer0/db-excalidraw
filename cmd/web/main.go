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

	// "github.com/joho/godotenv"
)

// Define an application struct to hold the application-wide dependencies for the // web application.
type application struct {
	logger *slog.Logger
}

func init() {
	// var drawingPath string
	environment := os.Getenv("ENVIRONMENT")
	// if environment == "" {
	// 	err := godotenv.Load(".env.local")
	// 	environment = os.Getenv("ENVIRONMENT")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }
	fmt.Println("Program running in:", environment, "mode")
}

// for local writing during testing
type fileWriter struct {
	outputPath string
}

type drawingUploader interface {
	Upload(p []byte) (int, error)
}
type localDrawingUploader struct {
	outputPath string
}

func (f *fileWriter) Write(p []byte) (n int, err error) {
	logFile := f.outputPath
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
	writer := &fileWriter{outputPath: "log.txt"}
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, writer), nil))
	app := application{logger: logger}
	slog.SetDefault(logger)
	// the file server with assets comes from a specific folder
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// StripPrefix gets rid of /static from the URL so we aren't searching
	// for /static/static/path-to-asset
	// create a get route for all assets
	mux.Handle("GET /static/", http.StripPrefix("/static", neuter(fileServer)))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /stream/{$}", app.streamHandler)
	mux.HandleFunc("GET /drawing/{name}", app.getDrawingByName)
	mux.HandleFunc("POST /drawing/{$}", app.postDrawing)
	mux.HandleFunc("POST /compressed/drawing", app.postCompressedDrawing)
	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
