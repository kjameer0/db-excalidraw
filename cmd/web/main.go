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
	"path/filepath"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/s3"
)

type awsConfig struct {
	key string
}
type testReader struct {
	dataPath string
}
type dataSaver interface {
	NewReader(nanoid string) (io.ReadCloser, error)
}

func (tC *testReader) NewReader(nanoid string) (io.ReadCloser, error) {
	pathName := tC.dataPath + string(filepath.Separator) + nanoid + ".txt"
	fmt.Println(pathName)
	f, err := os.Open(pathName)
	// TODO: add user who did this failed get
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Define an application struct to hold the application-wide dependencies for the // web application.
type application struct {
	logger    *slog.Logger
	dataSaver dataSaver
}

// for local writing during testing
type fileWriter struct {
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
	if writeErr != nil {
		log.Printf("Error writing to file %s: %v\n", logFile, writeErr)
		return 0, writeErr
	}

	return bytesWritten, nil
}

// TODO: add production logging location
// TODO: write database connection function
// TODO: connect to database
// TODO: add compression for json
// TODO: add writer for files and S3
// TODO: add reader for S3
// TODO: add shell scripts for prod and local
func main() {
	environment := flag.String("env", "development", "indicates production, testing, or development version of application")
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()
	fmt.Printf("Application running in %s mode\n", *environment)
	// mux is the part of the app that guides requests
	// to the url that matches their path
	writer := &fileWriter{outputPath: "log.txt"}
	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout, writer), nil))
	var dataSaver dataSaver = &testReader{dataPath: "test-drawings"}
	app := application{logger: logger, dataSaver: dataSaver}
	slog.SetDefault(logger)

	app.logger.Info("starting server", slog.String("addr", *addr))
	mux := app.routes()

	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
